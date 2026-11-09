from langchain_core.language_models import BaseChatModel
from langchain_core.messages import AIMessage, HumanMessage


def continue_conversation(
    model: BaseChatModel, history: list[tuple[str, str]], question: str
) -> AIMessage:
    # Map "human" -> HumanMessage and "ai" -> AIMessage, then add the question
    pass
