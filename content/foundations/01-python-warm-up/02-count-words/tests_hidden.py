from main import count_words


def test_case():
    """ignores case

    input: count_words("Go go GO")
    """
    assert count_words("Go go GO") == {"go": 3}


def test_whitespace():
    """splits on any whitespace

    input: count_words("a\\tb\\n a  b")
    """
    assert count_words("a\tb\n a  b") == {"a": 2, "b": 2}
