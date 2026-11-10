"""Lesson constructors. In opts, the FIRST option is the correct one; gen.py shuffles positions."""


def P(slug, title, diff, tags, body, q, opts, exp):
    return dict(kind="pattern", slug=slug, title=title, diff=diff, tags=tags, body=body, q=q, opts=opts, exp=exp)


def X(slug, title, url, diff, tags, body, q, opts, exp, read=None):
    return dict(kind="problem", slug=slug, title=title, url=url, diff=diff, tags=tags, body=body, q=q,
                opts=opts, exp=exp, read=read)
