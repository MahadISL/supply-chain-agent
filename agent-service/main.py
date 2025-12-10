import os
from fastapi import FastAPI, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from dotenv import load_dotenv
from pydantic import BaseModel
from langchain_core.messages import HumanMessage
from app.graph import app_graph

load_dotenv()

app = FastAPI(title="Supply Chain Agent")

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)


class ChatRequest(BaseModel):
    query: str


@app.get("/health")
def health_check():
    return {"status": "online", "service": "Agentic AI (Python)"}


@app.post("/chat")
async def chat_endpoint(request: ChatRequest):
    """
    Main interaction point.
    User sends: "Check the standing desks"
    Agent runs the LangGraph loop and returns the final answer.
    """
    try:
        # Start the Graph with the User's input
        inputs = {"messages": [HumanMessage(content=request.query)]}

        # Run the Graph (invoke)
        # The graph will loop (Agent -> Tool -> Agent) until finished
        result = app_graph.invoke(inputs)

        # Extract the final response from the AI
        last_message = result["messages"][-1]
        return {"response": last_message.content}

    except Exception as e:
        print(f"Error processing request: {e}")
        raise HTTPException(status_code=500, detail=str(e))


if __name__ == "__main__":
    import uvicorn

    port = int(os.getenv("AGENT_PORT", 8000))
    uvicorn.run(app, host="0.0.0.0", port=port)