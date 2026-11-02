# codequest/python: the docker-python grading runner.
# Run only through api/internal/runner (no network, read-only root, tmpfs /work, non-root).
FROM python:3.12.16-slim

ENV PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1 \
    PIP_NO_CACHE_DIR=1 \
    PIP_DISABLE_PIP_VERSION_CHECK=1 \
    LANGSMITH_TRACING=false \
    LANGCHAIN_TRACING_V2=false

COPY requirements.txt /tmp/requirements.txt
RUN pip install -r /tmp/requirements.txt && rm /tmp/requirements.txt

COPY harness/cq /usr/local/lib/python3.12/site-packages/cq
COPY harness/harness.py harness/child.py /opt/cq/
# The root filesystem is read-only at run time, so compile bytecode now.
RUN python -m compileall -q /usr/local/lib/python3.12 /opt/cq \
 && useradd --uid 10001 --no-create-home --shell /usr/sbin/nologin runner

USER 10001
WORKDIR /work
ENTRYPOINT ["python", "-I", "/opt/cq/harness.py"]
