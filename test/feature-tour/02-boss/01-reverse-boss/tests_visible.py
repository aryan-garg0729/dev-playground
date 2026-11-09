from cq import expect

from main import reverse_words


def test_three():
    'reverse_words("one two three")'
    expect(reverse_words("one two three"), "three two one")


def test_one():
    'reverse_words("solo")'
    expect(reverse_words("solo"), "solo")
