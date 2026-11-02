from langchain_core.prompts import ChatPromptTemplate


def make_translation_prompt() -> ChatPromptTemplate:
    return ChatPromptTemplate.from_messages(
        [
            ("system", "You are a translator from English to French. Reply with the translation only."),
            ("human", "{text}"),
        ]
    )


if __name__ == "__main__":
    prompt = make_translation_prompt()
    for m in prompt.format_messages(text="Good morning"):
        print(f"{m.type}: {m.content}")
