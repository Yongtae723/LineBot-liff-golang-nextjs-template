import liff from "@line/liff";
import { LiffMockPlugin } from "@line/liff-mock";
import type { MockData } from "@line/liff-mock/dist/store/MockDataStore";
import type { LiffMock } from "@line/liff-mock/dist/type";

const LIFF_ID = process.env.NEXT_PUBLIC_LIFF_ID as string;
const NEXT_PUBLIC_MOCK_USER_LINE_ID = process.env
  .NEXT_PUBLIC_MOCK_USER_LINE_ID as string;
console.log(
  "===================NEXT_PUBLIC_MOCK_USER_LINE_ID: ",
  NEXT_PUBLIC_MOCK_USER_LINE_ID
);

if (!LIFF_ID) {
  throw new Error("LIFF_ID is not set");
}

const setupMockLiff = async (): Promise<void> => {
  liff.use(new LiffMockPlugin());
  await (liff as LiffMock).init({
    liffId: LIFF_ID,
    mock: true,
  });

  (liff as LiffMock).$mock.set((data: Partial<MockData>) => ({
    ...data,
    isInClient: true,
    getAccessToken: `local_access_token`,
    getProfile: {
      userId: process.env.NEXT_PUBLIC_MOCK_USER_LINE_ID as string,
      displayName: process.env.NEXT_PUBLIC_MOCK_USER_NAME as string,
    },
  }));
};

export async function setupLiff(redirectTo: string): Promise<void> {
  if (process.env.NODE_ENV === "development") {
    await setupMockLiff();
  } else {
    // If you use liff, .init() automatically redirect the path the user would like to access.
    await liff.init({ liffId: LIFF_ID });
  }

  // following code is for mock liff. user from line does not use this code.
  if (!liff.isLoggedIn()) {
    const redirectUri = new URL(redirectTo, window.location.origin).href;
    liff.login({ redirectUri });
    return;
  }
}

export async function getLineProfile(): Promise<{
  accessToken: string;
  lineID: string;
}> {
  const accessToken = liff.getAccessToken();
  if (!accessToken) {
    throw new Error("Line access token not available after LIFF setup");
  }
  const lineID = (await liff.getProfile()).userId;
  return { accessToken, lineID };
}
