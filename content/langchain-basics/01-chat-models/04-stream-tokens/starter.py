from langchain_core.language_models import BaseChatModel


def stream_tokens(model: BaseChatModel, question: str) -> list[str]:
    chunks = []
    # Iterate over model.stream(question) and collect each chunk's text
    return chunks
