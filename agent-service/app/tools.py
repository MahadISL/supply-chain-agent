import os
import requests
import json
from langchain.tools import tool
from langchain_community.embeddings import HuggingFaceInferenceAPIEmbeddings
from pinecone import Pinecone
from typing import List, Optional

# --- CONFIGURATION ---
JAVA_URL = os.getenv("JAVA_SERVICE_URL", "http://core-service:8080/api")
HF_TOKEN = os.getenv("HUGGINGFACE_TOKEN")
PC_API_KEY = os.getenv("PINECONE_API_KEY")
PC_HOST = os.getenv("PINECONE_INDEX_HOST")

# Initialize Pinecone/Embeddings safely
try:
    embeddings = HuggingFaceInferenceAPIEmbeddings(
        api_key=HF_TOKEN,
        model_name="sentence-transformers/all-MiniLM-L6-v2"
    )
    pc = Pinecone(api_key=PC_API_KEY)
    index = pc.Index(host=PC_HOST)
    print("AI Tools: Pinecone & Embeddings Initialized Successfully")
except Exception as e:
    print(f"AI Tools Error: Failed to init Pinecone: {e}")
    index = None


# --- TOOL 1: RAG RETRIEVER ---
@tool
def lookup_policy(query: str) -> str:
    """
    Useful for checking supplier contracts, shipping policies, or rules.
    Input should be a specific question like "What is the minimum order for Apex Furniture?".
    """
    print(f"TOOL CALL: lookup_policy with query='{query}'")  # Debug Log

    if not index:
        return "System Error: Policy database is currently offline. Assume standard terms apply."

    try:
        query_vector = embeddings.embed_query(query)
        results = index.query(
            vector=query_vector,
            top_k=3,
            include_metadata=True
        )

        context = ""
        for match in results['matches']:
            text = match['metadata'].get('text', '')
            filename = match['metadata'].get('filename', 'Unknown')
            context += f"Source ({filename}): {text}\n---\n"

        if not context:
            # CRITICAL: Tell LLM to stop searching
            return "Result: No specific policy found in the database. DO NOT RETRY. Proceed with available information."

        return context
    except Exception as e:
        print(f"TOOL ERROR (lookup_policy): {e}")
        return f"Error searching policies: {str(e)}. DO NOT RETRY."


# --- TOOL 2: INVENTORY CHECKER ---
@tool
def check_inventory(product_name: Optional[str] = None) -> str:
    """
    Checks the current stock levels of products in the warehouse.
    """
    print(f"TOOL CALL: check_inventory for '{product_name}'")  # Debug Log

    try:
        # We fetch all products to be safe
        response = requests.get(f"{JAVA_URL}/products", timeout=5)
        if response.status_code == 200:
            products = response.json()
            report = "Current Inventory Data:\n"
            for p in products:
                report += f"- ID {p['id']}: {p['name']} (Stock: {p['stockQuantity']}, Min: {p['minStockLevel']}, Supplier: {p['supplier']['name']})\n"
            return report
        else:
            print(f"JAVA API ERROR: {response.status_code}")
            return f"Error: Failed to fetch inventory (Status {response.status_code}). DO NOT RETRY."
    except Exception as e:
        print(f"CONNECTION ERROR (check_inventory): {e}")
        return f"System Error: Cannot connect to Inventory Database at {JAVA_URL}. Is the Java service running? DO NOT RETRY."


# --- TOOL 3: PURCHASE ORDER CREATOR ---
@tool
def create_purchase_order(product_id: int, quantity: int, reasoning: str) -> str:
    """
    Drafts a purchase order for a specific product.
    Requires product_id and quantity.
    """
    print(f"TOOL CALL: create_purchase_order for ID {product_id}, Qty {quantity}")

    try:
        payload = {"productId": product_id, "quantity": quantity}
        headers = {'Content-Type': 'application/json'}
        response = requests.post(f"{JAVA_URL}/orders", json=payload, headers=headers, timeout=5)

        if response.status_code == 200:
            order_data = response.json()
            return f"SUCCESS: Order #{order_data['id']} created. Status: {order_data['status']}."
        else:
            return f"Failed to create order. Status: {response.status_code}"
    except Exception as e:
        return f"Error creating order: {str(e)}"