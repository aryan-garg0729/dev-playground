from main import double_all


def test_negative():
    """doubles negative numbers

    input: double_all([-4, 5])
    """
    assert double_all([-4, 5]) == [-8, 10]


def test_new_list():
    """leaves the input list unchanged

    input: double_all(nums) with nums = [7]
    """
    nums = [7]
    double_all(nums)
    assert nums == [7]
