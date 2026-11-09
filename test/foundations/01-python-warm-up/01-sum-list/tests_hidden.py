from main import sum_list


def test_negative():
    """handles negative numbers

    input: sum_list([-1000, 999, -4])
    """
    assert sum_list([-1000, 999, -4]) == -5


def test_long():
    """handles a long list

    input: sum_list(list(range(10000)))
    """
    assert sum_list(list(range(10000))) == 49995000
