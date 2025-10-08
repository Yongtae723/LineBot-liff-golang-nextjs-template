"use client";

import { useEffect } from "react";
import { setupLiff } from "@/lib/liff/init";

// According to LIFF documentation, liff.init() should be called in `/`
// This page is only used for LIFF login.
// setupLiff redirects after liff.init() is called to the target path.
export default function LiffInitPage() {
  useEffect(() => {
    setupLiff("/home");
  }, []);

  return (
    <div className="flex items-center justify-center min-h-screen">
      <div className="text-center">
        <div className="w-16 h-16 border-4 border-blue-500 border-t-transparent rounded-full animate-spin mx-auto mb-4"></div>
        <p className="text-lg font-semibold">データを読み込んでいます...</p>
      </div>
    </div>
  );
}
