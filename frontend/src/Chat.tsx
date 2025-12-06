import { useState, FormEvent, useRef, useEffect } from "react";

interface Message {
  role: "user" | "assistant";
  content: string;
}

interface Usage {
  used: number;
  remaining: number;
  limit: number;
}

export function Chat() {
  const [messages, setMessages] = useState<Message[]>([]);
  const [input, setInput] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [usage, setUsage] = useState<Usage | null>(null);
  const messagesEndRef = useRef<HTMLDivElement>(null);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    if (!input.trim() || loading) return;

    const userMessage = input.trim();
    setInput("");
    setError("");
    setLoading(true);

    // Add user message to chat
    setMessages((prev) => [...prev, { role: "user", content: userMessage }]);

    try {
      const response = await fetch("http://localhost:8080/llm", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify({ message: userMessage }),
      });

      const data = await response.json();

      if (response.ok) {
        // Add assistant response to chat
        setMessages((prev) => [
          ...prev,
          { role: "assistant", content: data.response },
        ]);
        // Update usage stats
        setUsage(data.usage);
      } else {
        setError(data.message || "Failed to get response");
      }
    } catch (err) {
      setError("Network error. Please check if you're logged in.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex flex-col p-4 max-w-4xl mx-auto">
      {/* Header */}
      <div className="bg-[#1a1a1a] rounded-lg p-4 mb-4 border border-[#333]">
        <div className="flex justify-between items-start mb-2">
          <h1 className="text-2xl font-bold">Meal Planning Assistant</h1>
          <a
            href="/preferences"
            className="text-sm text-[#646cff] hover:text-[#535bf2] transition-colors"
          >
            ⚙️ Preferences
          </a>
        </div>
        <p className="text-gray-400 text-sm">
          Tell me what ingredients you have, and I'll suggest healthy recipes!
        </p>
        {usage && (
          <div className="mt-3 text-sm">
            <span className="text-gray-400">Usage: </span>
            <span className="text-[#646cff] font-medium">
              {usage.used}/{usage.limit}
            </span>
            <span className="text-gray-500 ml-2">
              ({usage.remaining} remaining)
            </span>
          </div>
        )}
      </div>

      {/* Messages */}
      <div className="flex-1 bg-[#1a1a1a] rounded-lg border border-[#333] mb-4 overflow-hidden flex flex-col">
        <div className="flex-1 overflow-y-auto p-4 space-y-4">
          {messages.length === 0 ? (
            <div className="text-center text-gray-500 mt-8">
              <p className="mb-4">Start a conversation!</p>
              <p className="text-sm text-gray-600">
                Example: "I have chicken, rice, and broccoli"
              </p>
            </div>
          ) : (
            messages.map((msg, idx) => (
              <div
                key={idx}
                className={`flex ${
                  msg.role === "user" ? "justify-end" : "justify-start"
                }`}
              >
                <div
                  className={`max-w-[80%] rounded-lg p-4 ${
                    msg.role === "user"
                      ? "bg-[#646cff] text-white"
                      : "bg-[#2a2a2a] border border-[#444]"
                  }`}
                >
                  <div className="text-xs opacity-70 mb-1">
                    {msg.role === "user" ? "You" : "Assistant"}
                  </div>
                  <div className="whitespace-pre-wrap">{msg.content}</div>
                </div>
              </div>
            ))
          )}
          {loading && (
            <div className="flex justify-start">
              <div className="bg-[#2a2a2a] border border-[#444] rounded-lg p-4">
                <div className="text-xs opacity-70 mb-1">Assistant</div>
                <div className="flex space-x-2">
                  <div className="w-2 h-2 bg-gray-500 rounded-full animate-bounce"></div>
                  <div className="w-2 h-2 bg-gray-500 rounded-full animate-bounce delay-100"></div>
                  <div className="w-2 h-2 bg-gray-500 rounded-full animate-bounce delay-200"></div>
                </div>
              </div>
            </div>
          )}
          <div ref={messagesEndRef} />
        </div>

        {/* Error Display */}
        {error && (
          <div className="px-4 pb-2">
            <div className="bg-red-500/10 border border-red-500/50 text-red-400 px-4 py-2 rounded-lg text-sm">
              {error}
            </div>
          </div>
        )}

        {/* Input Form */}
        <form onSubmit={handleSubmit} className="p-4 border-t border-[#333]">
          <div className="flex gap-2">
            <input
              type="text"
              value={input}
              onChange={(e) => setInput(e.target.value)}
              placeholder="Type your ingredients or meal request..."
              disabled={loading}
              className="flex-1 px-4 py-3 bg-[#2a2a2a] border border-[#444] rounded-lg focus:outline-none focus:ring-2 focus:ring-[#646cff] focus:border-transparent disabled:opacity-50"
            />
            <button
              type="submit"
              disabled={loading || !input.trim()}
              className="px-6 py-3 bg-[#646cff] hover:bg-[#535bf2] disabled:bg-[#464a87] disabled:cursor-not-allowed text-white font-medium rounded-lg transition-colors"
            >
              Send
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}

export default Chat;
