from cq import expect

from main import count_words


def test_sentence():
    'count_words("the cat and the hat")'
    expect(count_words("the cat and the hat"), {"the": 2, "cat": 1, "and": 1, "hat": 1})


def test_empty():
    'count_words("")'
    expect(count_words(""), {})
