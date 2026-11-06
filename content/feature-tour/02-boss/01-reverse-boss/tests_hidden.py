from main import reverse_words


def test_empty():
    """returns an empty string for no words

    input: reverse_words("")
    """
    assert reverse_words("") == ""


def test_spaces():
    """collapses extra spaces

    input: reverse_words("  a   b ")
    """
    assert reverse_words("  a   b ") == "b a"


def test_tabs():
    """splits on tabs and newlines too

    input: reverse_words("x\ty\nz")
    """
    assert reverse_words("x\ty\nz") == "z y x"
