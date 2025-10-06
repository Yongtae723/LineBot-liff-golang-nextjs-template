import liff from "@line/liff";

export async function setupLiff(redirectTo: string): Promise<void> {
  const LIFF_ID = process.env.NEXT_PUBLIC_LIFF_ID;

  if (!LIFF_ID) {
    throw new Error("LIFF_ID is not set");
  }

  await liff.init({ liffId: LIFF_ID });

  if (!liff.isLoggedIn()) {
    const redirectUri = new URL(redirectTo, window.location.origin).href;
    liff.login({ redirectUri });
    return;
  }
}

export async function getLineAccessToken(): Promise<string> {
  const accessToken = liff.getAccessToken();
  if (!accessToken) {
    throw new Error("Line access token not available after LIFF setup");
  }
  return accessToken;
}
