from langchain_core.language_models import BaseChatModel
from langchain_core.messages import AIMessage, HumanMessage

ROLES = {"human": HumanMessage, "ai": AIMessage}


def continue_conversation(
    model: BaseChatModel, history: list[tuple[str, str]], question: str
) -> AIMessage:
    messages = [ROLES[role](text) for role, text in history]
    messages.append(HumanMessage(question))
    return model.invoke(messages)
