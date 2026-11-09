from langchain_core.language_models import BaseChatModel


def stream_tokens(model: BaseChatModel, question: str) -> list[str]:
    chunks = []
    for chunk in model.stream(question):
        chunks.append(chunk.content)
    return chunks
