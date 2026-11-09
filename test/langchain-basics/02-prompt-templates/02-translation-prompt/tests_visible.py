from cq import expect

from main import make_translation_prompt


def test_variables():
    "input variables are source, target and text"
    expect(sorted(make_translation_prompt().input_variables), ["source", "target", "text"])


def test_format():
    "formats a system and a human message"
    msgs = make_translation_prompt().format_messages(source="English", target="German", text="Hi")
    expect([m.type for m in msgs], ["system", "human"])
    expect(msgs[0].content, "You are a translator from English to German. Reply with the translation only.")
    expect(msgs[1].content, "Hi")
