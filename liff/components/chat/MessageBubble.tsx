import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import type { Role } from "@/types/conversation";

interface MessageBubbleProps {
  role: Role;
  content: string;
  timestamp: string;
}

export function MessageBubble({
  role,
  content,
  timestamp,
}: MessageBubbleProps) {
  const isUser = role === "user";

  return (
    <div className={`flex gap-3 ${isUser ? "flex-row-reverse" : "flex-row"}`}>
      <Avatar className="h-8 w-8">
        <AvatarFallback
          className={
            isUser ? "bg-blue-500 text-white" : "bg-gray-500 text-white"
          }
        >
          {isUser ? "U" : "AI"}
        </AvatarFallback>
      </Avatar>
      <div
        className={`flex flex-col max-w-[70%] ${isUser ? "items-end" : "items-start"}`}
      >
        <div
          className={`rounded-lg px-4 py-2 ${
            isUser ? "bg-blue-500 text-white" : "bg-gray-100 text-gray-900"
          }`}
        >
          <p className="text-sm whitespace-pre-wrap break-words">{content}</p>
        </div>
        <span className="text-xs text-gray-500 mt-1">
          {new Date(timestamp).toLocaleTimeString("ja-JP", {
            hour: "2-digit",
            minute: "2-digit",
          })}
        </span>
      </div>
    </div>
  );
}
