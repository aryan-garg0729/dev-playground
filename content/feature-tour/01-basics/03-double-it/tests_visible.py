from cq import expect

from main import double_all


def test_three():
    "double_all([1, 2, 3])"
    expect(double_all([1, 2, 3]), [2, 4, 6])


def test_empty():
    "double_all([])"
    expect(double_all([]), [])
