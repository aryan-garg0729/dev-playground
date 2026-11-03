"""CodeQuest Python grading harness (trusted parent; PID 1 in the runner container).

Protocol
  stdin   one JSON job (from the Go runner):
            files        {name: source} learner files, written to /work
            entry        learner file the tests import (e.g. "main.py")
            tests        [{file, source}] test files; never written to disk
            run_entry    run the entry file as __main__ first (for stdout checks)
            forbidden    names whose calls are reported by the static scan
            timeout_ms   per phase (program, tests)
            output_limit bytes kept of each stdout/stderr stream
  stdout  exactly one JSON report line, written only by this process.

Learner code never runs in this process. It runs in child processes whose
stdout/stderr are pipes read here, so learner `print` calls land inside the
report's "stdout" field and can never produce or alter the report itself.
Test results come back on a separate pipe (the dedicated result fd), tagged with
a per-run nonce; records for unknown tests or with a wrong nonce are dropped.
"""

import ast
import ctypes
import json
import os
import re
import secrets
import signal
import subprocess
import sys
import threading
import time

WORK = "/work"
CHILD = os.path.join(os.path.dirname(os.path.abspath(__file__)), "child.py")
NAME_RE = re.compile(r"^[A-Za-z_][A-Za-z0-9_]*\.py$")
RESULT_LIMIT = 1 << 20
JOIN_GRACE_S = 0.5


def make_undumpable():
    """Stop same-uid learner processes from ptrace-ing or reading /proc/1/mem."""
    try:
        ctypes.CDLL(None, use_errno=True).prctl(4, 0, 0, 0, 0)  # PR_SET_DUMPABLE, 0
    except Exception:
        pass


class Capped:
    """Drains a pipe in a thread, keeping at most `limit` bytes."""

    def __init__(self, stream, limit):
        self.buf = bytearray()
        self.truncated = False
        self._stream, self._limit = stream, limit
        self._t = threading.Thread(target=self._run, daemon=True)
        self._t.start()

    def _run(self):
        try:
            while True:
                chunk = self._stream.read(65536)
                if not chunk:
                    return
                room = self._limit - len(self.buf)
                if room > 0:
                    self.buf += chunk[:room]
                if len(chunk) > max(room, 0):
                    self.truncated = True
        except (OSError, ValueError):
            return

    def text(self, wait=JOIN_GRACE_S):
        self._t.join(wait)
        s = bytes(self.buf).decode("utf-8", "replace")
        return s + "\n[output truncated]" if self.truncated else s


def scan(files, forbidden):
    """Static scan of learner files: syntax errors and calls to forbidden names."""
    hits, syntax = [], None
    for name, src in files.items():
        try:
            tree = ast.parse(src, filename=name)
        except SyntaxError as e:
            syntax = syntax or {"file": name, "line": e.lineno or 0, "message": f"SyntaxError: {e.msg}"}
            continue
        for node in ast.walk(tree):
            if not isinstance(node, ast.Call):
                continue
            f = node.func
            called = f.id if isinstance(f, ast.Name) else f.attr if isinstance(f, ast.Attribute) else None
            if called in forbidden:
                hits.append({"name": called, "file": name, "line": node.lineno})
    return syntax, hits


def discover(tests):
    """Top-level test_* functions of each test file, in source order, with docstring labels.

    A docstring line "input: <call>" declares the test's input (for Peek, which
    may reveal it for a failed hidden test; never its expected output).
    """
    out = []
    for t in tests:
        tree = ast.parse(t["source"], filename=t["file"])  # validated by cq build
        cases = []
        for node in tree.body:
            if isinstance(node, (ast.FunctionDef, ast.AsyncFunctionDef)) and node.name.startswith("test_"):
                doc = (ast.get_docstring(node) or "").strip().splitlines()
                inp = next((l.split(":", 1)[1].strip() for l in doc[1:] if l.strip().lower().startswith("input:")), "")
                cases.append({"name": node.name, "label": doc[0] if doc else "", "input": inp[:200]})
        out.append({"file": t["file"], "cases": cases})
    return out


def child_env():
    return {
        "PATH": "/usr/local/bin:/usr/bin:/bin",
        "HOME": WORK,
        "LANG": "C.UTF-8",
        "PYTHONDONTWRITEBYTECODE": "1",
        "PYTHONUNBUFFERED": "1",
        "PYTHONHASHSEED": "0",
        "LANGSMITH_TRACING": "false",
        "LANGCHAIN_TRACING_V2": "false",
    }


def run_proc(argv, stdin_bytes, timeout_s, limit, pass_fds=()):
    """Run a child in its own process group; kill the whole group on timeout."""
    t0 = time.monotonic()
    p = subprocess.Popen(
        argv, cwd=WORK, env=child_env(), stdin=subprocess.PIPE, stdout=subprocess.PIPE,
        stderr=subprocess.PIPE, pass_fds=pass_fds, start_new_session=True,
    )
    out, err = Capped(p.stdout, limit), Capped(p.stderr, limit)
    try:
        if stdin_bytes:
            p.stdin.write(stdin_bytes)
        p.stdin.close()
    except BrokenPipeError:
        pass
    timed_out = False
    try:
        p.wait(timeout=timeout_s)
    except subprocess.TimeoutExpired:
        timed_out = True
    try:
        os.killpg(p.pid, signal.SIGKILL)  # also reaps anything the learner forked
    except ProcessLookupError:
        pass
    p.wait()
    ms = int((time.monotonic() - t0) * 1000)
    return {
        "stdout": out.text(), "stderr": err.text(), "exit_code": p.returncode,
        "timed_out": timed_out, "runtime_ms": ms,
    }


def run_program(entry, timeout_s, limit):
    code = (
        "import runpy, sys; sys.path.insert(0, '/work'); "
        f"runpy.run_path('/work/{entry}', run_name='__main__')"
    )
    return run_proc([sys.executable, "-I", "-c", code], b"", timeout_s, limit)


def run_tests(entry, tests, plan, timeout_s, limit):
    nonce = secrets.token_hex(16)
    r, w = os.pipe()
    cfg = {
        "nonce": nonce, "result_fd": w, "entry": entry,
        "tests": [{"file": t["file"], "source": t["source"], "names": [c["name"] for c in p["cases"]]}
                  for t, p in zip(tests, plan)],
    }
    results = Capped(os.fdopen(r, "rb"), RESULT_LIMIT)
    try:
        proc = run_proc([sys.executable, "-I", CHILD], json.dumps(cfg).encode(), timeout_s, limit, pass_fds=(w,))
    finally:
        os.close(w)
    raw = results.text()

    known = {(p["file"], c["name"]) for p in plan for c in p["cases"]}
    got, done = {}, False
    for line in raw.splitlines():
        try:
            rec = json.loads(line)
        except ValueError:
            continue
        if not isinstance(rec, dict) or rec.get("nonce") != nonce:
            continue
        if rec.get("done") is True:
            done = True
            continue
        key = (rec.get("file"), rec.get("name"))
        if key in known and key not in got:
            got[key] = rec

    files = []
    for p in plan:
        cases = []
        for c in p["cases"]:
            rec = got.get((p["file"], c["name"]))
            case = {"name": c["name"], "label": c["label"], "input": c["input"], "ran": rec is not None, "passed": False,
                    "message": "", "error": "", "got": "", "expected": "", "line": 0}
            if rec is not None:
                case["passed"] = rec.get("passed") is True
                for k in ("message", "error", "got", "expected"):
                    v = rec.get(k)
                    case[k] = v[:1000] if isinstance(v, str) else ""
                line = rec.get("line")
                case["line"] = line if isinstance(line, int) and line > 0 else 0
            cases.append(case)
        files.append({"file": p["file"], "cases": cases})
    proc["files"] = files
    proc["completed"] = done
    return proc


SQL_CHILD = os.path.join(os.path.dirname(os.path.abspath(__file__)), "sql_child.py")


def run_sql(job, t0):
    """mode "sql": one learner statement over each dataset (see sql_child.py)."""
    timeout_s = int(job["timeout_ms"]) / 1000
    cfg = {"query": job["query"], "solution": job.get("solution", ""), "datasets": job["datasets"]}
    proc = run_proc([sys.executable, "-I", SQL_CHILD], json.dumps(cfg).encode(), timeout_s, 4 << 20)
    report = {"version": 1, "sql": None, "timed_out": proc["timed_out"], "error": ""}
    if not proc["timed_out"]:
        try:
            report["sql"] = json.loads(proc["stdout"])
        except ValueError:
            report["error"] = (proc["stderr"].strip().splitlines() or ["sql runner failed"])[-1][:500]
    report["runtime_ms"] = int((time.monotonic() - t0) * 1000)
    sys.stdout.write(json.dumps(report) + "\n")
    sys.stdout.flush()


def main():
    make_undumpable()
    t0 = time.monotonic()
    job = json.loads(sys.stdin.buffer.read())
    if job.get("mode") == "sql":
        run_sql(job, t0)
        return
    files, entry = job["files"], job["entry"]
    tests = job.get("tests", [])
    limit = int(job.get("output_limit", 16384))
    timeout_s = int(job["timeout_ms"]) / 1000

    for name in list(files) + [t["file"] for t in tests]:
        if not NAME_RE.match(name):
            raise SystemExit(f"harness: bad file name {name!r}")
    if entry not in files:
        raise SystemExit("harness: entry is not a learner file")

    for name, src in files.items():
        with open(os.path.join(WORK, name), "w", encoding="utf-8") as f:
            f.write(src)

    syntax, hits = scan(files, set(job.get("forbidden", [])))
    plan = discover(tests)
    report = {"version": 1, "syntax_error": syntax, "forbidden": hits, "program": None,
              "tests": None, "timed_out": False}

    if syntax is None:
        if job.get("run_entry"):
            report["program"] = run_program(entry, timeout_s, limit)
            report["timed_out"] = report["program"]["timed_out"]
        if not report["timed_out"] and tests:
            report["tests"] = run_tests(entry, tests, plan, timeout_s, limit)
            report["timed_out"] = report["tests"]["timed_out"]
    if report["tests"] is None:
        report["tests"] = {"stdout": "", "stderr": "", "exit_code": 0, "timed_out": False, "runtime_ms": 0,
                           "completed": False,
                           "files": [{"file": p["file"], "cases": [
                               {"name": c["name"], "label": c["label"], "input": c["input"], "ran": False, "passed": False,
                                "message": "", "error": "", "got": "", "expected": "", "line": 0}
                               for c in p["cases"]]} for p in plan]}
    report["runtime_ms"] = int((time.monotonic() - t0) * 1000)
    sys.stdout.write(json.dumps(report) + "\n")
    sys.stdout.flush()


if __name__ == "__main__":
    main()
