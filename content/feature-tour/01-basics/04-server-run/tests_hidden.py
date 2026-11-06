from main import shout


def test_empty():
    """handles an empty string

    input: shout("")
    """
    assert shout("") == "!"


def test_mixed():
    """capitalises mixed case

    input: shout("CodeQuest")
    """
    assert shout("CodeQuest") == "CODEQUEST!"
