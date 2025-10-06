"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";

interface InputBarProps {
  onSend: (message: string) => void;
  disabled?: boolean;
}

export function InputBar({ onSend, disabled }: InputBarProps) {
  const [message, setMessage] = useState("");

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (message.trim() && !disabled) {
      onSend(message.trim());
      setMessage("");
    }
  };

  return (
    <form onSubmit={handleSubmit} className="flex gap-2 p-4 border-t bg-white">
      <Input
        type="text"
        value={message}
        onChange={(e) => setMessage(e.target.value)}
        placeholder="メッセージを入力..."
        disabled={disabled}
        className="flex-1"
      />
      <Button type="submit" disabled={disabled || !message.trim()}>
        送信
      </Button>
    </form>
  );
}
