from langchain_core.messages import HumanMessage, SystemMessage

from main import build_messages


def test_other_persona():
    "works for another persona and question"
    msgs = build_messages("a chef", "What is umami?")
    assert isinstance(msgs[0], SystemMessage) and isinstance(msgs[1], HumanMessage)
    assert msgs[0].content == "You are a chef. Answer in one sentence."
    assert msgs[1].content == "What is umami?"


def test_exactly_two():
    "returns exactly two messages"
    assert len(build_messages("a poet", "Hi")) == 2
