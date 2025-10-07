"use client";

import { useEffect, useRef } from "react";
import { ScrollArea } from "@/components/ui/scroll-area";
import type { Conversation } from "@/types/conversation";
import { MessageBubble } from "./MessageBubble";

interface MessageListProps {
  conversations: Conversation[];
}

export function MessageList({ conversations }: MessageListProps) {
  const scrollRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (scrollRef.current) {
      scrollRef.current.scrollTop = scrollRef.current.scrollHeight;
    }
  });

  return (
    <ScrollArea className="flex-1 p-4" ref={scrollRef}>
      <div className="space-y-4">
        {conversations.length === 0 ? (
          <div className="text-center text-gray-500 mt-8">
            <p>まだ会話がありません</p>
            <p className="text-sm mt-2">
              メッセージを送信して会話を始めましょう
            </p>
          </div>
        ) : (
          conversations.map((conv) => (
            <MessageBubble
              key={conv.id}
              role={conv.role}
              content={conv.content}
              timestamp={conv.created_at}
            />
          ))
        )}
      </div>
    </ScrollArea>
  );
}
