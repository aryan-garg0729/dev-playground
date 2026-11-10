# usage: wrong.py <lesson-dir> <name> <i>=<exit>:<stdout> ...   (http: <i>=h<status>:<body>)
# Copies solution.json, replaces the given steps (0-based), writes wrong/<name>.json.
import json, sys, os
d, name, *edits = sys.argv[1:]
sol = json.load(open(os.path.join(d, "solution.json")))
for e in edits:
    i, rest = e.split("=", 1)
    code, out = rest.split(":", 1)
    out = out.encode().decode("unicode_escape")
    if code.startswith("h"):
        sol["steps"][int(i)] = {"http": {"status": int(code[1:]), "headers": {"Content-Type": "text/plain"}, "body": out}}
    else:
        sol["steps"][int(i)] = {"cli": {"exit_code": int(code), "stdout": out}}
os.makedirs(os.path.join(d, "wrong"), exist_ok=True)
json.dump(sol, open(os.path.join(d, "wrong", name + ".json"), "w"))
