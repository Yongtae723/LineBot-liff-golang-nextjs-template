export type Role = "user" | "assistant";

export interface Conversation {
  id: string;
  user_id: string;
  role: Role;
  content: string;
  created_at: string;
}

export interface ConversationResponse {
  response: string;
}
