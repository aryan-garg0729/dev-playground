from langchain_core.messages import BaseMessage, HumanMessage, SystemMessage


def build_messages(persona: str, question: str) -> list[BaseMessage]:
    return [
        SystemMessage(f"You are {persona}. Answer in one sentence."),
        HumanMessage(question),
    ]
