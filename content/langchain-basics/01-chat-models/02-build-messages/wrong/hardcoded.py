from langchain_core.messages import BaseMessage, HumanMessage, SystemMessage


def build_messages(persona: str, question: str) -> list[BaseMessage]:
    return [
        SystemMessage("You are a pirate. Answer in one sentence."),
        HumanMessage("Why is the sea salty?"),
    ]
