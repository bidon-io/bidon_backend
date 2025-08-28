"""Minimal, canonical React agent implementation using LangGraph.

This implementation follows the latest LangGraph best practices for creating
a React (Reasoning and Acting) agent with tool integration. State and memory
management is handled automatically by the LangGraph server/platform.
"""

from __future__ import annotations

import os

from dotenv import load_dotenv
from langchain_anthropic import ChatAnthropic
from langchain_core.tools import tool
from langgraph.prebuilt import create_react_agent

from llama_index.indices.managed.llama_cloud import LlamaCloudIndex
from llama_index.core import Settings

Settings.llm = None

load_dotenv()

@tool
def search_documentation(query: str) -> str:
    """Search documentation for information.

    Args:
        query: The search query string

    Returns:
        A placeholder response (to be implemented later)
    """
    index = LlamaCloudIndex(
        name="bidon-docs",
        project_name="Default",
        organization_id="73814bd1-d8f5-4308-b041-2a55e0b0bfde",
        api_key=os.getenv("LLAMA_INDEX_API_KEY"),
    )
    response = index.as_query_engine().query(query)
    return response.response


def create_agent():
    """Create and configure the React agent with all required components.

    State and memory management is handled automatically by the LangGraph server.
    No custom checkpointer is needed as the platform provides built-in persistence.

    Returns:
        Configured LangGraph agent ready for use
    """
    model = ChatAnthropic(
        model="claude-3-7-sonnet-latest",
        api_key=os.getenv("ANTHROPIC_API_KEY"),
        temperature=0
    )

    tools = [search_documentation]
    agent = create_react_agent(
        model=model,
        tools=tools,
        prompt="You are a helpful AI assistant. Use the available tools to help answer questions."
    )

    return agent

# Create the agent instance
graph = create_agent()
