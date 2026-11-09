from cq.fakes import fake_model

from main import stream_tokens


def test_strings():
    "every chunk is text, not a message object"
    chunks = stream_tokens(fake_model("red green blue"), "colours")
    assert chunks and all(isinstance(c, str) for c in chunks)


def test_order():
    "keeps the chunks in order"
    assert "".join(stream_tokens(fake_model("a b c d e"), "letters")) == "a b c d e"
