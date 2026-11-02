from langchain_core.prompts import ChatPromptTemplate


def make_dated_prompt() -> ChatPromptTemplate:
    return ChatPromptTemplate.from_messages([("system", "Today is {today}."), ("human", "{question}")])
