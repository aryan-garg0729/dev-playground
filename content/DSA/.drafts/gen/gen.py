"""Generate content/DSA from the lesson data in data_*.py. Run: python3 gen.py <content/DSA>"""
import hashlib
import json
import os
import sys

sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))
from data_a import CHAPTERS as A  # noqa: E402
from data_b import CHAPTERS as B  # noqa: E402

CHAPTERS = A + B
SITES = {
    "cses.fi": "CSES",
    "atcoder.jp": "AtCoder",
    "leetcode.com": "LeetCode",
    "www.spoj.com": "SPOJ",
}


def q(s):
    return json.dumps(s, ensure_ascii=False)


def site_of(url):
    host = url.split("/")[2]
    return SITES[host]


def write(path, text):
    os.makedirs(os.path.dirname(path), exist_ok=True)
    with open(path, "w") as f:
        f.write(text)


def lesson_files(root, ch, idx, les):
    folder = f"{idx:02d}-{les['slug']}"
    d = os.path.join(root, ch["folder"], folder)
    # rotate the correct option to a stable pseudo-random position
    opts = les["opts"]
    correct, wrong = opts[0], opts[1:]
    pos = int(hashlib.sha1(les["slug"].encode()).hexdigest(), 16) % len(opts)
    ordered = wrong[:pos] + [correct] + wrong[pos:]
    ids = "abcdefgh"
    answer = ids[pos]
    is_problem = les["kind"] == "problem"
    lines = [
        f"id: {ch['id']}.{les['slug']}",
        "type: mcq",
        "version: 1",
        f"title: {q(les['title'])}",
        f"xp: {10 if is_problem else 5}",
        f"difficulty: {les['diff']}",
        f"tags: [{', '.join(les['tags'])}]",
        "free: true",
        "spec:",
        "  question: prompt.md",
        "  multiple: false",
        "  options:",
    ]
    for i, text in enumerate(ordered):
        lines.append(f"    - {{ id: {ids[i]}, text: {q(text)} }}")
    lines.append(f"  answer: [{answer}]")
    lines.append(f"  explain: {q(les['exp'])}")
    write(os.path.join(d, "exercise.yaml"), "\n".join(lines) + "\n")

    md = [f"# {les['title']}", ""]
    if is_problem:
        site = site_of(les["url"])
        md.append(f"**Practice:** [{les['title']} on {site}]({les['url']})")
        if les.get("read"):
            name, url = les["read"]
            md.append("")
            md.append(f"**Read:** [{name}]({url})")
        md.append("")
    md.append(les["body"].strip())
    md.append("")
    if is_problem:
        md.append("## Check your approach")
        md.append("")
        md.append(
            f"Solve it on {site} first, or give it a real try (about 30 minutes) before reading on. "
            "The explanation after you answer walks through the full approach."
        )
        md.append("")
    else:
        md.append("## Check yourself")
        md.append("")
    md.append(les["q"].strip())
    write(os.path.join(d, "prompt.md"), "\n".join(md) + "\n")
    return folder


def main(root):
    ch_folders = []
    for ci, ch in enumerate(CHAPTERS, 1):
        ch["folder"] = f"{ci:02d}-{ch['slug']}"
        folders = [lesson_files(root, ch, i, les) for i, les in enumerate(ch["lessons"], 1)]
        y = [f"id: {ch['id']}", f"title: {q(ch['title'])}", "exercises:"] + [f"  - {f}" for f in folders]
        write(os.path.join(root, ch["folder"], "chapter.yaml"), "\n".join(y) + "\n")
        ch_folders.append(ch["folder"])
    course = [
        "id: dsa",
        "title: DSA Interview Prep",
        "description: >",
        "  Revise data structures and algorithms before interviews. Each lesson recaps a pattern or",
        "  links a classic problem on CSES, AtCoder or LeetCode, then checks you know the key idea.",
        "  Every lesson is open: pick any topic, in any order.",
        "order: 6",
        "published: true",
        "chapters:",
    ] + [f"  - {f}" for f in ch_folders] + [
        "icon: dsa",
        "level: intermediate",
        "tags: [DSA, Algorithms, Interviews]",
        "topics: Graphs, trees, DP, strings, range queries",
    ]
    write(os.path.join(root, "course.yaml"), "\n".join(course) + "\n")
    n = sum(len(c["lessons"]) for c in CHAPTERS)
    print(f"{len(CHAPTERS)} chapters, {n} lessons")


def P(slug, title, diff, tags, body, q, opts, exp):
    return dict(kind="pattern", slug=slug, title=title, diff=diff, tags=tags, body=body, q=q, opts=opts, exp=exp)


def X(slug, title, url, diff, tags, body, q, opts, exp, read=None):
    return dict(kind="problem", slug=slug, title=title, url=url, diff=diff, tags=tags, body=body, q=q,
                opts=opts, exp=exp, read=read)


if __name__ == "__main__":
    main(sys.argv[1])
