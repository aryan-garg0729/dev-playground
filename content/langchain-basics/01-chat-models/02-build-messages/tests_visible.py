from cq import expect
from langchain_core.messages import HumanMessage, SystemMessage

from main import build_messages


def test_types():
    "message types are [SystemMessage, HumanMessage]"
    msgs = build_messages("a pirate", "Why is the sea salty?")
    expect([type(m).__name__ for m in msgs], ["SystemMessage", "HumanMessage"])


def test_system_text():
    "system message sets the persona"
    msgs = build_messages("a pirate", "Why is the sea salty?")
    expect(msgs[0].content, "You are a pirate. Answer in one sentence.")


def test_question():
    "human message is the question"
    msgs = build_messages("a pirate", "Why is the sea salty?")
    expect(msgs[1].content, "Why is the sea salty?")
