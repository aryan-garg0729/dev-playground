from cq import expect
from cq.fakes import fake_model

from main import stream_tokens


def test_several_chunks():
    "returns more than one chunk"
    chunks = stream_tokens(fake_model("one two three"), "count")
    assert isinstance(chunks, list), f"expected a list, got {type(chunks).__name__}"
    assert len(chunks) > 1, "expected several chunks: are you calling .stream()?"


def test_joined():
    "chunks join back into the reply"
    expect("".join(stream_tokens(fake_model("one two three"), "count")), "one two three")
