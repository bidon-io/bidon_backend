# New LangGraph Project

The core logic defined in `src/agent/graph.py`, showcases an single-step application that responds with a fixed string and the configuration provided.


## Getting Started

1. Install dependencies, along with the [LangGraph CLI](https://langchain-ai.github.io/langgraph/concepts/langgraph_cli/), which will be used to run the server.

```bash
uv sync --all-extras
```

2. (Optional) Customize the code and project as needed. Create a `.env` file if you need to use secrets.

```bash
cp .env.example .env
```

If you want to enable LangSmith tracing, add your LangSmith API key to the `.env` file.

```text
# .env
LANGSMITH_API_KEY=lsv2...
ANTHROPIC_API_KEY=...
LANGSMITH_API_KEY=...
LLAMA_INDEX_API_KEY=...
```

3. Start the LangGraph Server. (Or Docker-compose setup)

```shell
uv run langgraph dev
```

4. Docker-compose setup

- Build own LangGraph Docker image
```shell
uv run langgraph build -t copilot
```

```shell
docker-compose up --profile copilot
```
