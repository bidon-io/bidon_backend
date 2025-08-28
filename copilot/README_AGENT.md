# LangGraph React Agent Implementation

This directory contains a minimal, canonical React agent implementation using the latest LangGraph framework, optimized for LangGraph server deployment.

## Overview

The implementation follows LangGraph best practices and includes:

- **React Agent Pattern**: Uses `create_react_agent` for the canonical ReAct (Reasoning and Acting) pattern
- **Tool Integration**: Includes a `search_documentation` tool (stub implementation)
- **Server-Managed State**: No custom checkpointer - state and memory managed by LangGraph server
- **Claude Sonnet 3.7**: Configured with Anthropic's latest model
- **Environment Configuration**: Loads API keys from `.env` file

## Architecture

### Core Components

1. **Model**: Claude Sonnet 3.7 (`claude-3-7-sonnet-latest`)
2. **Tools**: `search_documentation` tool with query parameter
3. **Checkpointer**: `InMemorySaver` for local development (auto-disabled in platform)
4. **Agent**: Created using `create_react_agent` function

### File Structure

```
src/agent/
├── graph.py          # Main agent implementation
└── __init__.py       # Module initialization

example_usage.py      # Demonstration script
README_AGENT.md       # This documentation
```

## Usage

### Basic Usage

```python
from src.agent.graph import graph

# Configure conversation thread
config = {"configurable": {"thread_id": "my-session"}}

# Invoke the agent
response = graph.invoke({
    "messages": [{"role": "user", "content": "Your question here"}]
}, config)

print(response["messages"][-1].content)
```

### Running the Example

```bash
# Run the demonstration script
uv run python example_usage.py
```

## Features Demonstrated

### ✅ Tool Integration
- The agent can call the `search_documentation` tool
- Tool calls are automatically handled by the React pattern
- Results are integrated into the conversation flow

### ✅ State Persistence
- Conversations are maintained across multiple invocations (local environment)
- Each thread ID maintains separate conversation history
- Memory persists within the same session
- Automatically disabled in LangGraph platform environments

### ✅ Multi-turn Conversations
- Agent remembers previous messages in the same thread
- Context is maintained for follow-up questions
- Natural conversation flow

### ✅ Thread Isolation
- Different thread IDs maintain separate conversations
- No cross-contamination between sessions
- Clean session management

## Platform Compatibility

The implementation automatically adapts to different environments:

### Local Development
- Uses `InMemorySaver` checkpointer for state persistence
- Full conversation memory within sessions
- Manual thread management required

### LangGraph Platform/API
- Checkpointer automatically disabled
- Platform handles persistence automatically
- Built-in thread and state management

The agent detects platform environment using:
- `LANGGRAPH_API` environment variable
- `LANGSMITH_API_URL` environment variable

## Configuration

### Environment Variables

The agent requires the following environment variable:

```bash
ANTHROPIC_API_KEY=your_anthropic_api_key_here
```

### Model Configuration

The agent is configured with:
- **Model**: `claude-3-7-sonnet-latest`
- **Temperature**: 0 (deterministic responses)
- **API Key**: Loaded from environment

## Implementation Details

### React Pattern

The implementation uses LangGraph's `create_react_agent` which provides:
- Automatic tool calling loop
- Message state management
- Built-in error handling
- Streaming support

### Tool Definition

Tools are defined using the `@tool` decorator:

```python
@tool
def search_documentation(query: str) -> str:
    """Search documentation for information."""
    # Implementation here
    return result
```

### State Management

The agent uses LangGraph's built-in state management with:
- Message history accumulation
- Automatic state updates
- Conditional checkpointing (local development only)
- Platform-managed persistence when deployed

## Next Steps

To enhance this implementation:

1. **Implement Real Search**: Replace the stub `search_documentation` with actual search functionality
2. **Add More Tools**: Extend with additional tools as needed
3. **Custom Prompts**: Implement dynamic prompting based on context
4. **Structured Output**: Add response formatting with Pydantic models
5. **Error Handling**: Implement custom error handling for tools

## Dependencies

- `langgraph>=0.6.0`
- `langchain-anthropic>=0.3.18`
- `python-dotenv>=1.1.1`

Install with:
```bash
uv add langgraph langchain-anthropic python-dotenv
```
