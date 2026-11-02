from langchain_core.messages import BaseMessage, HumanMessage


def build_messages(persona: str, question: str) -> list[BaseMessage]:
    return [HumanMessage(f"You are {persona}. Answer in one sentence. {question}")]
