from cq import expect

from main import sum_list


def test_one():
    "sum_list([5])"
    expect(sum_list([5]), 5)


def test_three():
    "sum_list([1, 2, 3])"
    expect(sum_list([1, 2, 3]), 6)


def test_empty():
    "sum_list([])"
    expect(sum_list([]), 0)
