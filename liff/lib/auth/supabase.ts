"use server";

import { createClient } from "@/lib/supabase/server";

export const loginSupabase = async (accessToken: string): Promise<void> => {
  try {
    const response = await fetch(
      `${process.env.NEXT_PUBLIC_BACKEND_URL}/api/v1/user/register`,
      {
        method: "POST",
        body: JSON.stringify({ access_token: accessToken }),
        headers: {
          "Content-Type": "application/json",
        },
      }
    );

    if (!response.ok) {
      let errorMsg = `Server login request failed with status ${response.status}`;
      try {
        const errorData = (await response.json()) as { error?: string };
        errorMsg = `Server login failed: ${errorData?.error || "Unknown server error"}`;
      } catch (jsonError) {
        console.error("Failed to parse error response as JSON", jsonError);
      }
      throw new Error(errorMsg);
    }

    const { line_id: lineID } = (await response.json()) as { line_id: string };
    console.log("line_id: ", lineID);

    const supabase = await createClient();
    const { error: signInError } = await supabase.auth.signInWithPassword({
      email: lineID,
      password: lineID,
    });

    if (signInError) {
      throw new Error(
        `Failed to sign in with password: ${signInError.message}`
      );
    }
  } catch (error) {
    console.error("Failed during Supabase login process:", error);
    throw new Error(`Failed to complete login process: ${String(error)}`);
  }
};
