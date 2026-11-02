"""Helpers for CodeQuest test files.

    from cq import expect
    from main import sum_list

    def test_three():
        "sum_list([1, 2, 3])"
        expect(sum_list([1, 2, 3]), 6)

A test's docstring is its label in the check list. On a visible test, a failed
`expect` shows "got" and "expected"; on a hidden test only the label is shown.
"""

__all__ = ["Mismatch", "expect"]

_LIMIT = 200


def _short(v) -> str:
    s = repr(v)
    return s if len(s) <= _LIMIT else s[: _LIMIT - 3] + "..."


class Mismatch(AssertionError):
    def __init__(self, got, expected, message=""):
        super().__init__(message or f"got {_short(got)}, expected {_short(expected)}")
        self.got = _short(got)
        self.expected = _short(expected)
        self.message = message


_NUMBERS = (int, float)


def _same(got, expected) -> bool:
    """Equality that a learner-defined __eq__ cannot fake: types must match exactly
    (int and float may mix), containers are compared element by element."""
    tg, te = type(got), type(expected)
    if tg is not te:
        if tg in _NUMBERS and te in _NUMBERS:
            return got == expected
        return False
    if te in (list, tuple):
        return len(got) == len(expected) and all(_same(g, e) for g, e in zip(got, expected))
    if te is dict:
        return got.keys() == expected.keys() and all(_same(got[k], expected[k]) for k in expected)
    # Same concrete type, so the comparison uses that type's own __eq__.
    return expected == got


def expect(got, expected, message: str = "") -> None:
    """Fail the test unless got equals expected (same type, same value)."""
    if not _same(got, expected):
        raise Mismatch(got, expected, message)
