from main import fizzbuzz


def test_fifteen():
    """says FizzBuzz for multiples of 15

    input: fizzbuzz(15)[14]
    """
    assert fizzbuzz(15)[14] == "FizzBuzz"


def test_zero():
    """returns an empty list for 0

    input: fizzbuzz(0)
    """
    assert fizzbuzz(0) == []


def test_thirty():
    """still right at 30

    input: fizzbuzz(30)[29]
    """
    assert fizzbuzz(30)[29] == "FizzBuzz"


def test_length():
    """returns exactly n strings

    input: len(fizzbuzz(100))
    """
    assert len(fizzbuzz(100)) == 100


def test_numbers():
    """uses the number as a string otherwise

    input: fizzbuzz(7)[6]
    """
    assert fizzbuzz(7)[6] == "7"
