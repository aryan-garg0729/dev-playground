from datetime import date

from langchain_core.prompts import ChatPromptTemplate


def today() -> str:
    return date.today().isoformat()


def make_dated_prompt() -> ChatPromptTemplate:
    # Build the prompt, then call .partial(today=...)
    pass
