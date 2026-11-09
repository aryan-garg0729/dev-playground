from cq import expect

from main import fizzbuzz


def test_one():
    "fizzbuzz(1)"
    expect(fizzbuzz(1), ["1"])


def test_fizz():
    "fizzbuzz(5)[2]"
    expect(fizzbuzz(5)[2], "Fizz")


def test_buzz():
    "fizzbuzz(5)[4]"
    expect(fizzbuzz(5)[4], "Buzz")
