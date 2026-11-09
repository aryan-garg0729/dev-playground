from cq import expect

from main import shout


def test_hi():
    'shout("hi there")'
    expect(shout("hi there"), "HI THERE!")
