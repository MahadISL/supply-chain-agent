const JAVA_API = 'http://localhost:8080/api';
const AGENT_API = 'http://localhost:8000';

export async function getProducts() {
  const res = await fetch(`${JAVA_API}/products`, { cache: 'no-store' });
  if (!res.ok) throw new Error('Failed to fetch products');
  return res.json();
}

export async function chatWithAgent(query: string) {
  const res = await fetch(`${AGENT_API}/chat`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ query }),
  });
  if (!res.ok) throw new Error('Failed to talk to agent');
  return res.json();
}