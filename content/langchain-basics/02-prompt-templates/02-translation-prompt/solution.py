from langchain_core.prompts import ChatPromptTemplate


def make_translation_prompt() -> ChatPromptTemplate:
    return ChatPromptTemplate.from_messages(
        [
            ("system", "You are a translator from {source} to {target}. Reply with the translation only."),
            ("human", "{text}"),
        ]
    )


if __name__ == "__main__":
    prompt = make_translation_prompt()
    for m in prompt.format_messages(source="English", target="French", text="Good morning"):
        print(f"{m.type}: {m.content}")
