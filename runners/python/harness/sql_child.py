"""Runs one learner SQL statement (and, when asked, the reference solution) on
fresh in-memory SQLite databases. Reads its config on stdin, writes one JSON
object on stdout. The parent harness enforces the time limit.

Result comparison happens in Go (api/internal/checks/sql_result_equals.go) and,
for Run in the browser, in sandbox/src/sqlcompare.ts.
"""

import json
import sqlite3
import sys

ROW_CAP = 1000
DENY = {sqlite3.SQLITE_ATTACH, sqlite3.SQLITE_DETACH, sqlite3.SQLITE_PRAGMA}


def authorizer(action, *_):
    return sqlite3.SQLITE_DENY if action in DENY else sqlite3.SQLITE_OK


def cell(v):
    if isinstance(v, bytes):
        return v.hex()
    if isinstance(v, float) and v != v:  # NaN
        return None
    return v


def run_query(setup, query):
    if not query.strip():
        return {"error": "Write one SQL statement"}
    con = sqlite3.connect(":memory:")
    try:
        con.executescript(setup)  # author data, before the authorizer
        con.set_authorizer(authorizer)
        try:
            cur = con.execute(query)
        except sqlite3.ProgrammingError as e:
            if "one statement" in str(e):
                return {"error": "Write a single SQL statement"}
            return {"error": str(e)}
        except sqlite3.Error as e:
            return {"error": str(e)}
        if cur.description is None:
            return {"error": "Your statement returned no rows; write a SELECT"}
        out = {"columns": [d[0] for d in cur.description], "rows": [], "truncated": False}
        for row in cur:
            if len(out["rows"]) >= ROW_CAP:
                out["truncated"] = True
                break
            out["rows"].append([cell(v) for v in row])
        return out
    except sqlite3.Error as e:
        return {"error": str(e)}
    finally:
        con.close()


def main():
    cfg = json.loads(sys.stdin.read())
    results = []
    for ds in cfg["datasets"]:
        item = {"learner": run_query(ds["setup"], cfg["query"]), "solution": None}
        if ds.get("solution"):
            item["solution"] = run_query(ds["setup"], cfg["solution"])
        results.append(item)
    sys.stdout.write(json.dumps({"results": results}))


if __name__ == "__main__":
    main()
