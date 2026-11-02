"""Deterministic stand-ins for real models, for graded LangChain tests.

A fake chat model replays a scripted list of responses, so a test can check
that learner code wires everything together, with no model or network.
"""

from typing import Any, Iterable

from langchain_core.embeddings import DeterministicFakeEmbedding
from langchain_core.language_models.fake_chat_models import GenericFakeChatModel
from langchain_core.messages import AIMessage, BaseMessage


class FakeChatModel(GenericFakeChatModel):
    """GenericFakeChatModel that also supports bind_tools / with_structured_output.

    It records every prompt it received in `.calls` so tests can inspect what
    your code actually sent to the model.
    """

    calls: list[list[BaseMessage]] = []
    bound_tools: list[Any] = []

    def bind_tools(self, tools, **kwargs):
        self.bound_tools = list(tools)
        return self

    def _generate(self, messages, stop=None, run_manager=None, **kwargs):
        self.calls.append(list(messages))
        return super()._generate(messages, stop=stop, run_manager=run_manager, **kwargs)


def ai(content: str = "", tool_calls: list[dict] | None = None) -> AIMessage:
    """Build an AIMessage. tool_calls: [{"name": ..., "args": {...}}] (ids are filled in)."""
    calls = [{"id": f"call_{i}", **tc} for i, tc in enumerate(tool_calls or [])]
    return AIMessage(content=content, tool_calls=calls)


def fake_model(*responses: str | AIMessage, repeat: bool = False) -> FakeChatModel:
    """A fake chat model that returns `responses` in order.

    fake_model("hello")                                   -> always says "hello" once
    fake_model(ai(tool_calls=[{"name": "add", "args": {"a": 1, "b": 2}}]), "3")
    repeat=True cycles through the responses forever (handy for batch / parallel tests).
    """
    msgs = [r if isinstance(r, AIMessage) else AIMessage(content=r) for r in responses]
    it: Iterable[AIMessage] = _cycle(msgs) if repeat else iter(msgs)
    return FakeChatModel(messages=it, calls=[], bound_tools=[])


def _cycle(items):
    while True:
        yield from items


def fake_embeddings(size: int = 64) -> DeterministicFakeEmbedding:
    """Same text -> same vector. Different text -> unrelated vector (no semantics!)."""
    return DeterministicFakeEmbedding(size=size)
