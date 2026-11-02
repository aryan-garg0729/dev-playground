from cq import expect
from cq.fakes import fake_model
from langchain_core.messages import AIMessage

from main import continue_conversation

HISTORY = [("human", "Hi, I'm Ada"), ("ai", "Hello Ada!")]


def test_reply():
    "returns the model's AIMessage"
    reply = continue_conversation(fake_model("Your name is Ada."), list(HISTORY), "What's my name?")
    assert isinstance(reply, AIMessage), f"expected an AIMessage, got {type(reply).__name__}"
    expect(reply.content, "Your name is Ada.")


def test_sends_history():
    "sends history plus the question"
    m = fake_model("Your name is Ada.")
    continue_conversation(m, list(HISTORY), "What's my name?")
    sent = m.calls[0]
    expect([type(x).__name__ for x in sent], ["HumanMessage", "AIMessage", "HumanMessage"])
    expect(sent[-1].content, "What's my name?")
