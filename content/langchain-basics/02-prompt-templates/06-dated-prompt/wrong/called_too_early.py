from datetime import date

from langchain_core.prompts import ChatPromptTemplate


def today() -> str:
    return date.today().isoformat()


def make_dated_prompt() -> ChatPromptTemplate:
    prompt = ChatPromptTemplate.from_messages([("system", "Today is {today}."), ("human", "{question}")])
    return prompt.partial(today=today())
