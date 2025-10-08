"use client";

import { useEffect, useState } from "react";
import { loginSupabase } from "@/lib/auth/supabase";
import { getLineProfile, setupLiff } from "@/lib/liff/init";

export default function LoginPage() {
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const searchParams = new URLSearchParams(window.location.search);
    const initializeAndRedirect = async () => {
      try {
        const redirectTo = searchParams.get("redirectTo") || "/home";
        await setupLiff(redirectTo);
        const { accessToken, lineID } = await getLineProfile();
        await loginSupabase(accessToken, lineID);
        window.location.href = redirectTo;
      } catch (error) {
        console.error("[Login Page] Error:", error);
        if (error instanceof Error) {
          setError(error.message);
        } else {
          setError("予期しないエラーが発生しました");
        }
      }
    };

    initializeAndRedirect();
  }, []);

  if (error) {
    return (
      <div className="flex items-center justify-center min-h-screen p-4">
        <div className="text-center max-w-md">
          <div className="bg-red-50 border border-red-200 rounded-lg p-6">
            <h2 className="text-lg font-semibold text-red-700 mb-2">
              エラーが発生しました
            </h2>
            <p className="text-sm text-red-600 mb-4">{error}</p>
            <p className="text-xs text-gray-600">
              問題が解決しない場合は、開発者にお問い合わせください。
            </p>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="flex items-center justify-center min-h-screen">
      <div className="text-center">
        <div className="w-16 h-16 border-4 border-blue-500 border-t-transparent rounded-full animate-spin mx-auto mb-4"></div>
        <p className="text-lg font-semibold">データを読み込んでいます...</p>
        <p className="text-sm text-gray-600 mt-2">
          時間が経っても画面が変わらない場合は
          <br />
          開発者にお問い合わせください
        </p>
      </div>
    </div>
  );
}
