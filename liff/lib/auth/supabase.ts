"use server";

import { createClient } from "@/lib/supabase/server";

export const loginSupabase = async (
  accessToken: string,
  lineID: string
): Promise<void> => {
  const supabase = await createClient();
  // The user who already add bot can login with email and password.
  // Like backend implementation, you can use your own login logic.
  // since this codes are executed on the server, it is not exposed to the client.
  const { error: signInError } = await supabase.auth.signInWithPassword({
    email: `${lineID}@example.com`,
    password: lineID,
  });

  if (signInError) {
    // this is for the user who use liff before adding bot to LINE.
    try {
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_BACKEND_URL}/api/v1/user/register/liff`,
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

      const { line_id: backendLineID } = (await response.json()) as {
        line_id: string;
      };
      const { error: newSignInError } = await supabase.auth.signInWithPassword({
        email: `${backendLineID}@example.com`,
        password: backendLineID,
      });

      if (newSignInError) {
        throw new Error(
          `Failed to sign in with password: ${newSignInError.message}`
        );
      }
    } catch (error) {
      console.error("Failed during Supabase login process:", error);
      throw new Error(`Failed to complete login process: ${String(error)}`);
    }
  }
};
