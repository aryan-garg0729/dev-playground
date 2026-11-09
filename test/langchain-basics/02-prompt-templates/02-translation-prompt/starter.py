from langchain_core.prompts import ChatPromptTemplate


def make_translation_prompt() -> ChatPromptTemplate:
    # ChatPromptTemplate.from_messages([("system", ...), ("human", ...)])
    pass


if __name__ == "__main__":
    prompt = make_translation_prompt()
    for m in prompt.format_messages(source="English", target="French", text="Good morning"):
        print(f"{m.type}: {m.content}")
