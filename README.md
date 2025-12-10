
# Enterprise Autonomous Supply Chain Orchestrator

![System Status](https://img.shields.io/badge/System-Operational-green)
![Tech Stack](https://img.shields.io/badge/Stack-Java%20%7C%20Go%20%7C%20Python%20%7C%20Next.js-blue)
![AI](https://img.shields.io/badge/AI-Agentic%20LangGraph-purple)

## Project Overview

A distributed **Agentic AI System** designed to autonomously manage enterprise supply chain operations. Unlike traditional chatbots, this system utilizes a **Reasoning Engine (LangGraph)** to perform multi-step tasks: monitoring real-time inventory, retrieving unstructured policy data (RAG) via Vector Search, and executing purchase orders through transactional APIs.

This project demonstrates a **Microservices Architecture**, utilizing the best tool for each specific task: **Java** for transactional integrity, **Go** for high-throughput ETL, and **Python** for cognitive reasoning.

---

## System Architecture

```mermaid
graph TD
    User[User / Frontend] -->|HTTP| NextJS[Next.js Dashboard]
    NextJS -->|REST| Java[Core Service (Java/Spring Boot)]
    NextJS -->|REST| Agent[Agent Service (Python/LangGraph)]
    
    subgraph "Data Persistence"
        Java -->|Read/Write| Postgres[(PostgreSQL)]
        Java -->|Cache| Redis[(Redis)]
    end
    
    subgraph "Cognitive Layer"
        Agent -->|Tool Call| Java
        Agent -->|RAG Search| Pinecone[(Pinecone Vector DB)]
        Agent -->|Inference| Groq[Llama 3 LLM]
    end
    
    subgraph "ETL Pipeline"
        Upload[PDF Contracts] -->|Stream| Go[Ingestion Service (Golang)]
        Go -->|Upsert Vectors| Pinecone
    end
```

## Key Features

### Agentic Reasoning Engine (Python + LangGraph)
*   **Stateful Decision Making:** Uses a graph-based state machine to loop through reasoning steps (Think → Act → Observe).
*   **Self-Correction:** Detects tool failures (e.g., missing policy data) and degrades gracefully without crashing.
*   **Tool Use:** Autonomously calls internal APIs to check stock levels (`check_inventory`) and draft orders (`create_purchase_order`).

### High-Performance ETL Pipeline (Go)
*   **Stream Processing:** Dedicated Golang service for parsing PDF contracts.
*   **Smart Chunking:** Implements sliding-window text chunking to preserve context for RAG.
*   **Vectorization:** Generates dense embeddings (`all-MiniLM-L6-v2`) and upserts to Pinecone.

### Core Transactional Layer (Java Spring Boot)
*   **System of Record:** Manages relational data (Products, Suppliers, Orders) with strict ACID compliance.
*   **RESTful API:** Exposes endpoints for the Agent and Frontend to interact with the ERP system.

---

## Technology Stack

| Component | Technology | Role |
| :--- | :--- | :--- |
| **Agentic AI** | Python, LangChain, LangGraph | Orchestration & Reasoning |
| **LLM Inference** | Llama 3 (via Groq), HuggingFace | Intelligence & Embeddings |
| **Core Backend** | Java 17, Spring Boot 3 | Transactional Logic & Security |
| **ETL Pipeline** | Go (Golang) | High-performance Document Ingestion |
| **Frontend** | Next.js 14, Tailwind CSS | Mission Control Dashboard |
| **Databases** | PostgreSQL, Redis, Pinecone | Relational, Cache, & Vector Storage |
| **Infra** | Docker Compose | Container Orchestration |

---

## Getting Started

### Prerequisites
*   Docker & Docker Compose
*   API Keys: Groq, Pinecone, HuggingFace

### 1. Clone the Repository
```bash
git clone https://github.com/MahadISL/supply-chain-agent.git
cd supply-chain-agent
```

### 2. Environment Setup
Create a `.env` file in the root directory:
```ini
# Databases
DB_USERNAME=admin
DB_PASSWORD=password123
DB_NAME=supply_chain_db

# AI Services
GROQ_API_KEY=gsk_...
PINECONE_API_KEY=pcsk_...
PINECONE_INDEX_HOST=https://...
HUGGINGFACE_TOKEN=hf_...

# Service Ports
INGESTION_PORT=8081
AGENT_PORT=8000
```

### 3. Launch System
```bash
docker compose up -d --build
```

### 4. Access the Dashboard
Visit **http://localhost:3000** to view the Real-Time Inventory.  
Visit **http://localhost:3000/agent** to chat with the AI.

---

## Engineering Decisions

*   **Why Go for Ingestion?** Parsing large PDF streams is CPU-intensive. Go's concurrency model and low memory footprint make it superior to Python for this specific "heavy lifting" task.
*   **Why Java for Core?** Supply chain data requires strict type safety and transactional consistency that the Spring ecosystem guarantees.
*   **Why LangGraph?** Traditional linear chains (LangChain) fail at complex reasoning. A graph architecture allows the agent to loop and retry, essential for autonomous behavior.
