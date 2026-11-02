from langchain_core.language_models import BaseChatModel
from langchain_core.messages import AIMessage


def continue_conversation(
    model: BaseChatModel, history: list[tuple[str, str]], question: str
) -> AIMessage:
    history.append(("human", question))
    return model.invoke(history)
