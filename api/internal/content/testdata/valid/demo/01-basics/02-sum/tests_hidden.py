# SERVER_ONLY_SENTINEL: hidden tests
from main import sum_list


def test_negative():
    "handles negative numbers"
    assert sum_list([-1000, 999]) == -1  # SERVER_ONLY_SENTINEL


def test_empty():
    assert sum_list([]) == 0
