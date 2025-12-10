"use client";
import { useState, useRef, useEffect } from 'react';
import { chatWithAgent } from '../../services/api';

type Message = {
  role: 'user' | 'agent';
  content: string;
};

export default function AgentChat() {
  const [messages, setMessages] = useState<Message[]>([
    { role: 'agent', content: 'Hello! I am your Supply Chain Agent. I can check inventory, read contracts, and draft orders. How can I help?' }
  ]);
  const [input, setInput] = useState('');
  const [loading, setLoading] = useState(false);
  const bottomRef = useRef<HTMLDivElement>(null);

  // Auto-scroll to bottom
  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages]);

  const handleSend = async () => {
    if (!input.trim()) return;

    // 1. Add User Message
    const userMsg: Message = { role: 'user', content: input };
    setMessages(prev => [...prev, userMsg]);
    setInput('');
    setLoading(true);

    try {
      // 2. Call Python Agent
      const data = await chatWithAgent(userMsg.content);
      
      // 3. Add Agent Response
      setMessages(prev => [...prev, { role: 'agent', content: data.response }]);
    } catch (err) {
      setMessages(prev => [...prev, { role: 'agent', content: 'Error: Could not connect to the Agent Brain.' }]);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="flex flex-col h-screen bg-gray-100">
      {/* Header */}
      <header className="bg-white shadow p-4 flex justify-between items-center">
        <h1 className="text-xl font-bold text-blue-800">Agent Command Center</h1>
        <a href="/" className="text-gray-600 hover:text-blue-600">← Back to Dashboard</a>
      </header>

      {/* Chat Area */}
      <div className="flex-1 overflow-y-auto p-4 space-y-4">
        {messages.map((msg, idx) => (
          <div key={idx} className={`flex ${msg.role === 'user' ? 'justify-end' : 'justify-start'}`}>
            <div className={`max-w-[70%] p-4 rounded-lg shadow ${
              msg.role === 'user' 
                ? 'bg-blue-600 text-white rounded-br-none' 
                : 'bg-white text-gray-800 rounded-bl-none border border-gray-200'
            }`}>
              <p className="whitespace-pre-wrap">{msg.content}</p>
            </div>
          </div>
        ))}
        
        {loading && (
          <div className="flex justify-start">
            <div className="bg-gray-200 text-gray-500 p-3 rounded-lg rounded-bl-none animate-pulse">
              Agent is thinking...
            </div>
          </div>
        )}
        <div ref={bottomRef} />
      </div>

      {/* Input Area */}
      <div className="bg-white p-4 border-t">
        <div className="flex gap-2 max-w-4xl mx-auto">
          <input
            type="text"
            className="flex-1 p-3 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 text-black"
            placeholder="Ask about inventory, contracts, or orders..."
            value={input}
            onChange={(e) => setInput(e.target.value)}
            onKeyDown={(e) => e.key === 'Enter' && handleSend()}
            disabled={loading}
          />
          <button
            onClick={handleSend}
            disabled={loading}
            className="bg-blue-600 text-white px-6 py-3 rounded-lg hover:bg-blue-700 disabled:opacity-50 font-bold"
          >
            Send
          </button>
        </div>
      </div>
    </div>
  );
}