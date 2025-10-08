import { createClient } from "@/lib/supabase/client";
import type { Conversation, ConversationResponse } from "@/types/conversation";

export class ApiClient {
  private baseUrl: string;

  constructor() {
    this.baseUrl = process.env.NEXT_PUBLIC_BACKEND_URL || "";
  }

  private async getAuthHeader(): Promise<string> {
    const supabase = createClient();
    const {
      data: { session },
    } = await supabase.auth.getSession();

    if (!session?.access_token) {
      throw new Error("No access token");
    }

    return `Bearer ${session.access_token}`;
  }

  async registerUser(accessToken: string): Promise<{ line_id: string }> {
    const response = await fetch(`${this.baseUrl}/api/v1/user/register`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ access_token: accessToken }),
    });

    if (!response.ok) {
      throw new Error("Failed to register user");
    }

    return response.json();
  }

  async getConversations(limit = 50): Promise<Conversation[]> {
    const authHeader = await this.getAuthHeader();

    const response = await fetch(
      `${this.baseUrl}/api/v1/conversations?limit=${limit}`,
      {
        headers: {
          Authorization: authHeader,
        },
      }
    );

    if (!response.ok) {
      throw new Error("Failed to fetch conversations");
    }

    return response.json();
  }

  async sendMessage(message: string): Promise<ConversationResponse> {
    const authHeader = await this.getAuthHeader();

    const response = await fetch(`${this.baseUrl}/api/v1/conversations`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: authHeader,
      },
      body: JSON.stringify({ message }),
    });

    if (!response.ok) {
      throw new Error("Failed to send message");
    }

    return response.json();
  }
}

export const apiClient = new ApiClient();
