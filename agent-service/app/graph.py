import os
from typing import Annotated, TypedDict, Union
from langchain_groq import ChatGroq
from langchain_core.messages import BaseMessage, HumanMessage, SystemMessage
from langgraph.graph import StateGraph, END
from langgraph.prebuilt import ToolNode
from langgraph.graph.message import add_messages
from app.tools import check_inventory, lookup_policy, create_purchase_order


# Define the State
# Memory of Agent to keep track of the conversation history.
class AgentState(TypedDict):
    messages: Annotated[list, add_messages]


# Initialize the Model (Llama 3 via Groq)
groq_api_key = os.getenv("GROQ_API_KEY")
llm = ChatGroq(
    temperature=0,
    model_name="llama-3.3-70b-versatile",
    groq_api_key=groq_api_key
)

# Bind Tools to the LLM
# Teaches Llama 3 which tools exist and how to call them.
tools = [check_inventory, lookup_policy, create_purchase_order]
llm_with_tools = llm.bind_tools(tools)


# Define the Nodes (The Logic Steps)

def agent_node(state: AgentState):
    """
    The 'Thinking' Node.
    It looks at the history (state) and decides what to do next.
    """
    messages = state["messages"]
    # Added a System Prompt to define the persona
    if not isinstance(messages[0], SystemMessage):
        system_prompt = SystemMessage(content="""
        You are an Autonomous Supply Chain Agent. 
        Your goal is to monitor inventory and purchase supplies based on company policy.
        
        RULES:
        1. Always check inventory levels first.
        2. If stock is low (<10), you must TRY to search the 'lookup_policy' tool to find supplier rules.
        3. CRITICAL: If the 'lookup_policy' tool returns "No relevant policy found" or similar, DO NOT SEARCH AGAIN. Instead, just report the low stock level to the user and warn them that policy data is missing.
        4. Never hallucinate contract terms.
        """)
        messages = [system_prompt] + messages

    # Call the LLM
    response = llm_with_tools.invoke(messages)
    return {"messages": [response]}


def should_continue(state: AgentState) -> str:
    """
    The 'Router' Node.
    Decides if we should stop (respond to user) or call a tool.
    """
    last_message = state["messages"][-1]

    # If the LLM wants to make a tool call, go to "tools" node
    if last_message.tool_calls:
        return "tools"

    # Otherwise, stop and return answer to user
    return END


# Build the Graph
workflow = StateGraph(AgentState)

# Add Nodes
workflow.add_node("agent", agent_node)
workflow.add_node("tools", ToolNode(tools))

# Set Entry Point
workflow.set_entry_point("agent")

# Add Edges (Connections)
# From Agent, we check if we should continue
workflow.add_conditional_edges(
    "agent",
    should_continue,
    {
        "tools": "tools",  # If tool call needed, go to tools
        END: END  # Else, finish
    }
)

# From Tools, we ALWAYS go back to Agent (to interpret the tool output)
workflow.add_edge("tools", "agent")

# Compile into a Runnable Application
app_graph = workflow.compile()