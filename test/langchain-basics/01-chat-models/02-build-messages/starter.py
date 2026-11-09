from langchain_core.messages import BaseMessage, HumanMessage, SystemMessage


def build_messages(persona: str, question: str) -> list[BaseMessage]:
    # Return [SystemMessage(...), HumanMessage(...)]
    pass
