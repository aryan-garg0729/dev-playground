from langchain_core.prompts import ChatPromptTemplate

from main import make_translation_prompt


def test_type():
    "returns a ChatPromptTemplate"
    assert isinstance(make_translation_prompt(), ChatPromptTemplate)


def test_other_languages():
    "fills every variable"
    msgs = make_translation_prompt().format_messages(source="Hindi", target="Japanese", text="Namaste")
    assert msgs[0].content == "You are a translator from Hindi to Japanese. Reply with the translation only."
    assert msgs[1].content == "Namaste"
