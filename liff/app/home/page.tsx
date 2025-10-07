"use client";

import { useEffect, useState } from "react";
import { InputBar } from "@/components/chat/InputBar";
import { MessageList } from "@/components/chat/MessageList";
import { Card } from "@/components/ui/card";
import { apiClient } from "@/lib/api/client";
import type { Conversation } from "@/types/conversation";

export default function HomePage() {
  const [conversations, setConversations] = useState<Conversation[]>([]);
  const [loading, setLoading] = useState(true);
  const [sending, setSending] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const loadConversations = async () => {
      try {
        const data = await apiClient.getConversations();
        setConversations(data);
      } catch (err) {
        console.error("Failed to load conversations:", err);
        setError("会話履歴の読み込みに失敗しました");
      } finally {
        setLoading(false);
      }
    };

    loadConversations();
  }, []);

  const handleSendMessage = async (message: string) => {
    setSending(true);
    try {
      const userConv: Conversation = {
        id: crypto.randomUUID(),
        user_id: "",
        role: "user",
        content: message,
        created_at: new Date().toISOString(),
      };
      setConversations((prev) => [...prev, userConv]);

      const response = await apiClient.sendMessage(message);

      const assistantConv: Conversation = {
        id: crypto.randomUUID(),
        user_id: "",
        role: "assistant",
        content: response.response,
        created_at: new Date().toISOString(),
      };
      setConversations((prev) => [...prev, assistantConv]);
    } catch (error) {
      console.error("Failed to send message:", error);
      setError("メッセージの送信に失敗しました");
      setTimeout(() => setError(null), 3000);
    } finally {
      setSending(false);
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-center">
          <div className="w-16 h-16 border-4 border-blue-500 border-t-transparent rounded-full animate-spin mx-auto mb-4"></div>
          <p className="text-lg text-gray-600">読み込み中...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="flex flex-col h-screen bg-gray-50">
      <div className="bg-blue-500 text-white p-4 shadow-md">
        <h1 className="text-xl font-semibold">AI チャット</h1>
        <p className="text-sm text-blue-100">Powered by Gemini</p>
      </div>

      {error && (
        <div className="bg-red-50 border-l-4 border-red-500 p-4 mx-4 mt-4">
          <p className="text-red-700 text-sm">{error}</p>
        </div>
      )}

      <Card className="flex-1 flex flex-col m-4 shadow-lg overflow-hidden">
        <MessageList conversations={conversations} />
        <InputBar onSend={handleSendMessage} disabled={sending} />
      </Card>
    </div>
  );
}
