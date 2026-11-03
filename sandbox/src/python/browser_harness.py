"""Browser (Pyodide) harness for Run: the learner's entry plus the visible tests.

Run is ungraded and nothing here is trusted by the server. Rows mirror the
server's visible unit_tests rows (api/internal/checks/unit_tests.go).
"""

import ast
import contextlib
import io
import json
import os
import runpy
import sys
import time
import traceback

WORK = "/work"
LIMIT = 16 * 1024


class Capped(io.StringIO):
    def write(self, s):
        room = LIMIT - self.tell()
        if room > 0:
            super().write(s[:room])
        return len(s)


def reset(files):
    os.makedirs(WORK, exist_ok=True)
    for f in os.listdir(WORK):
        os.remove(os.path.join(WORK, f))
    for name, mod in list(sys.modules.items()):
        if (getattr(mod, "__file__", None) or "").startswith(WORK + "/"):
            del sys.modules[name]
    for name, src in files.items():
        with open(os.path.join(WORK, name), "w", encoding="utf-8") as f:
            f.write(src)
    if WORK not in sys.path:
        sys.path.insert(0, WORK)


def learner_line(exc, entry_path):
    line = 0
    for fs in traceback.extract_tb(exc.__traceback__):
        if fs.filename == entry_path:
            line = fs.lineno or 0
    return line


def describe(exc):
    msg = str(exc)
    return f"{type(exc).__name__}: {msg}"[:500] if msg else type(exc).__name__


def run(job_json):
    from cq import Mismatch

    job = json.loads(job_json)
    files, entry = job["files"], job["entry"]
    entry_path = os.path.join(WORK, entry)
    t0 = time.perf_counter()
    reset(files)
    out, err = Capped(), Capped()
    rows = []

    try:
        ast.parse(files[entry], filename=entry_path)
        syntax = None
    except SyntaxError as e:
        syntax = (e.lineno or 0, f"SyntaxError: {e.msg}")

    with contextlib.redirect_stdout(out), contextlib.redirect_stderr(err):
        if syntax is None:
            try:
                runpy.run_path(entry_path, run_name="__main__")
            except SystemExit:
                pass
            except BaseException:
                traceback.print_exc()
    # Tests import the entry again; that output (and prints made during tests) is
    # not shown, matching the server, so Run shows the program's output once.
    with contextlib.redirect_stdout(Capped()), contextlib.redirect_stderr(Capped()):
        for tf in job["tests"]:
            names = []
            try:
                tree = ast.parse(tf["source"], filename=tf["file"])
                for node in tree.body:
                    if isinstance(node, (ast.FunctionDef, ast.AsyncFunctionDef)) and node.name.startswith("test_"):
                        doc = (ast.get_docstring(node) or "").strip().splitlines()
                        names.append((node.name, doc[0] if doc else node.name[len("test_"):]))
            except SyntaxError:
                continue
            if syntax is not None:
                for _, label in names:
                    rows.append({"name": label, "passed": False, "visible": True,
                                 "feedback": f"{syntax[1]} (line {syntax[0]})", "line": syntax[0]})
                continue
            ns = {"__name__": "cq_" + tf["file"][:-3]}
            try:
                exec(compile(tf["source"], tf["file"], "exec"), ns)
            except BaseException as e:
                for _, label in names:
                    rows.append({"name": label, "passed": False, "visible": True,
                                 "feedback": describe(e), "line": learner_line(e, entry_path)})
                continue
            for fn, label in names:
                row = {"name": label, "passed": False, "visible": True}
                try:
                    ns[fn]()
                    row["passed"] = True
                except Mismatch as m:
                    row.update(got=m.got, expected=m.expected)
                    if m.message:
                        row["feedback"] = m.message
                except AssertionError as e:
                    row["feedback"] = str(e)[:500] or "Assertion failed"
                    row["line"] = learner_line(e, entry_path)
                except BaseException as e:
                    row["feedback"] = describe(e)
                    row["line"] = learner_line(e, entry_path)
                rows.append(row)


    return json.dumps({
        "stdout": out.getvalue(), "stderr": err.getvalue(), "checks": rows,
        "runtimeMs": int((time.perf_counter() - t0) * 1000),
    })
