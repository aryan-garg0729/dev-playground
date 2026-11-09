from langchain_core.language_models import BaseChatModel


def stream_tokens(model: BaseChatModel, question: str) -> list[str]:
    return list(model.stream(question))
