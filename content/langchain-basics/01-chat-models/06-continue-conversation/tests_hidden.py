from cq.fakes import fake_model

from main import continue_conversation


def test_history_unchanged():
    "leaves the caller's history list unchanged"
    history = [("human", "My favourite colour is teal."), ("ai", "Nice, teal it is!")]
    before = list(history)
    continue_conversation(fake_model("teal"), history, "What is my favourite colour?")
    assert history == before, "do not append to the history you were given"


def test_empty_history():
    "works with an empty history"
    m = fake_model("Hello!")
    continue_conversation(m, [], "Hi")
    assert [type(x).__name__ for x in m.calls[0]] == ["HumanMessage"]


def test_long_history():
    "keeps every turn in order"
    history = [("human", "1"), ("ai", "2"), ("human", "3"), ("ai", "4")]
    m = fake_model("5")
    continue_conversation(m, history, "6")
    assert [x.content for x in m.calls[0]] == ["1", "2", "3", "4", "6"]
