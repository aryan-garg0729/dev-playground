from langchain_core.prompts import ChatPromptTemplate


def make_translation_prompt() -> ChatPromptTemplate:
    return ChatPromptTemplate.from_messages(
        [
            ("system", "You are a translator from {source} to {target}. Reply with the translation only."),
            ("human", "{text}"),
        ]
    )
