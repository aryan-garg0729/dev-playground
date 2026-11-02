"""Runs test functions against learner code. Learner code runs in this process.

Reads its config (nonce, result fd, test sources) from stdin before importing any
learner code, then reports one JSON line per test on the result fd. Test sources
are compiled from memory and never written to disk.
"""

import json
import os
import sys
import time
import traceback


def main():
    cfg = json.loads(sys.stdin.read())
    devnull = os.open(os.devnull, os.O_RDONLY)
    os.dup2(devnull, 0)
    os.close(devnull)
    sys.stdin = open(0, closefd=False)

    out = os.fdopen(cfg["result_fd"], "w", buffering=1)
    nonce = cfg["nonce"]
    entry_path = os.path.join("/work", cfg["entry"])
    sys.path.insert(0, "/work")

    from cq import Mismatch  # imported before learner code can patch it

    def emit(rec):
        rec["nonce"] = nonce
        out.write(json.dumps(rec) + "\n")

    def learner_line(exc):
        line = 0
        for fs in traceback.extract_tb(exc.__traceback__):
            if fs.filename == entry_path:
                line = fs.lineno or 0
        return line

    def describe(exc):
        msg = str(exc)
        return f"{type(exc).__name__}: {msg}"[:500] if msg else type(exc).__name__

    for tf in cfg["tests"]:
        ns = {"__name__": "cq_" + tf["file"][:-3]}
        try:
            exec(compile(tf["source"], tf["file"], "exec"), ns)
        except BaseException as e:  # usually the learner module failing to import
            for name in tf["names"]:
                emit({"file": tf["file"], "name": name, "passed": False,
                      "error": describe(e), "line": learner_line(e)})
            continue
        for name in tf["names"]:
            rec = {"file": tf["file"], "name": name, "passed": False}
            t0 = time.perf_counter()
            try:
                ns[name]()
                rec["passed"] = True
            except Mismatch as m:
                rec.update(message=m.message, got=m.got, expected=m.expected)
            except AssertionError as e:
                rec["message"] = str(e)[:500]
                rec["line"] = learner_line(e)
            except BaseException as e:  # SystemExit and friends from learner code count as errors
                rec["error"] = describe(e)
                rec["line"] = learner_line(e)
            rec["ms"] = int((time.perf_counter() - t0) * 1000)
            emit(rec)
    emit({"done": True})
    out.flush()


if __name__ == "__main__":
    main()
