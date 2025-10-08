# LIFF App (Next.js)

LINE Front-end Framework (LIFF) を使用したチャットアプリケーションです。Next.js 15 (App Router) で構築されています。

## 📦 主な機能

- ✅ LIFF認証統合
- ✅ Supabase認証連携
- ✅ リアルタイムチャットUI
- ✅ 会話履歴の表示
- ✅ Gemini LLMとの対話
- ✅ レスポンシブデザイン

## 🚀 起動方法

### 1. 依存関係インストール

```bash
npm install
```

### 2. 環境変数設定

```bash
cp .env.local.example .env.local
# .env.localファイルを編集
```

```.env.local
NEXT_PUBLIC_LIFF_ID=your-liff-id
NEXT_PUBLIC_BACKEND_URL=http://localhost:8080
NEXT_PUBLIC_SUPABASE_URL=http://localhost:54321
NEXT_PUBLIC_SUPABASE_ANON_KEY=your-supabase-anon-key
```

### 3. 開発サーバー起動

```bash
npm run dev
```

アプリは http://localhost:3000 で起動します。

## 🔧 LIFF設定

### 1. LIFF アプリ作成

1. [LINE Developers Console](https://developers.line.biz/console/) にアクセス
2. Channelを選択 → **LIFF** タブ
3. **Add** をクリック

### 2. LIFF設定

- **LIFF app name**: 任意の名前
- **Size**: Full
- **Endpoint URL**: 
  - ローカル: `http://localhost:3000`
  - 本番: `https://your-domain.com`
- **Scopes**: 
  - ✅ `profile`
  - ✅ `openid`
  - ✅ `email` (オプション)

### 3. LIFF ID取得

作成後に表示される **LIFF ID** を `.env.local` に設定：

```bash
NEXT_PUBLIC_LIFF_ID=1234567890-abcdefgh
```

## 🔐 認証フロー

```
1. ユーザーがLIFFアプリを開く (/)
   ↓
2. LIFF SDK初期化 (liff.init)
   ↓
3. LINE Access Token取得 (liff.getAccessToken)
   ↓
4. Backend API /api/v1/user/register/liff を呼び出し
   ↓
5. Backend APIがLINE User情報取得 + Supabase Auth作成 + usersテーブルにUserを新規保存
   ↓
6. LINE IDを受け取る
   ↓
7. Next.js Middlewareが未認証を検知
   ↓
8. /login へリダイレクト
   ↓
9. /login でSupabase認証 (signInWithPassword)
   ↓
10. JWT取得、セッション確立
   ↓
11. /home (チャット画面) へリダイレクト
```

**ページ遷移:**
- `/` → LIFF初期化（自動）
- `/login` → Supabase認証処理（middleware経由）
- `/home` → チャット画面（認証完了後）

## 📁 プロジェクト構造

```
liff/
├── src/
│   ├── app/
│   │   ├── page.tsx                    # LIFF初期化画面 (/)
│   │   ├── login/
│   │   │   └── page.tsx                # Supabase認証処理
│   │   ├── home/
│   │   │   └── page.tsx                # チャット画面
│   │   ├── layout.tsx
│   │   └── globals.css
│   ├── components/
│   │   ├── ui/                         # shadcn/ui コンポーネント
│   │   └── chat/
│   │       ├── ChatContainer.tsx
│   │       ├── MessageBubble.tsx
│   │       └── InputBar.tsx
│   ├── lib/
│   │   ├── supabase/
│   │   │   ├── client.ts               # Supabaseクライアント
│   │   │   └── server.ts               # サーバーサイド用
│   │   ├── api/
│   │   │   └── client.ts               # Backend API クライアント
│   │   └── liff/
│   │       └── init.ts                 # LIFF初期化
│   ├── types/
│   │   ├── conversation.ts
│   │   └── supabase.ts                 # 自動生成
│   └── middleware.ts                   # 認証チェック、未認証時は/loginへ
├── public/
├── .env.local.example
├── biome.json
├── next.config.mjs
├── package.json
├── tailwind.config.ts
└── tsconfig.json
```

## 🛠️ 開発コマンド

```bash
# 開発サーバー起動
npm run dev

# ビルド
npm run build

# プロダクションサーバー起動
npm start

# フォーマット
npm run format

# リント
npm run lint

# 型チェック
npm run type-check

# Supabase型定義生成
npm run gen:types
```

## 🎨 shadcn/ui コンポーネント

このプロジェクトは shadcn/ui を使用しています。

### コンポーネント追加

```bash
npx shadcn@latest add button
npx shadcn@latest add input
npx shadcn@latest add card
```

### 既にインストール済みのコンポーネント

- Button
- Input
- Card
- Avatar
- ScrollArea

## 🔄 Supabase統合

### クライアントサイド

```typescript
import { createClient } from '@/lib/supabase/client'

const supabase = createClient()
const { data: { session } } = await supabase.auth.getSession()
```

### サーバーサイド

```typescript
import { createClient } from '@/lib/supabase/server'

const supabase = createClient()
const { data: { user } } = await supabase.auth.getUser()
```

### 型定義の自動生成

```bash
npm run gen:types
```

これにより `src/types/supabase.ts` が生成されます。

## 🌐 Backend API連携

### API クライアントの使用

```typescript
import { apiClient } from '@/lib/api/client'

// 会話履歴取得
const conversations = await apiClient.getConversations(50)

// メッセージ送信
const response = await apiClient.sendMessage('こんにちは')
```

APIクライアントは自動的にSupabase JWTをヘッダーに付与します。

## 🚀 デプロイ

### Cloudflare Pages

```bash
# ビルド
npm run build

# Cloudflare Pagesにデプロイ
# Build command: npm run build
# Build output directory: .next
```

### Vercel

```bash
# Vercel CLIインストール
npm i -g vercel

# デプロイ
vercel
```

### Netlify

```bash
# Netlify CLIインストール
npm i -g netlify-cli

# デプロイ
netlify deploy --prod
```

## 🌍 環境変数（本番環境）

本番環境では以下の環境変数を設定してください：

| 変数名 | 説明 | 例 |
|--------|------|-----|
| `NEXT_PUBLIC_LIFF_ID` | LIFF ID | `1234567890-abcdefgh` |
| `NEXT_PUBLIC_BACKEND_URL` | Backend API URL | `https://api.your-domain.com` |
| `NEXT_PUBLIC_SUPABASE_URL` | Supabase URL | `https://xxx.supabase.co` |
| `NEXT_PUBLIC_SUPABASE_ANON_KEY` | Supabase Anon Key | `eyJh...` |

## 🧪 ローカル開発のTips

### LIFF Simulatorの使用

LIFFアプリはLINEアプリ内で動作しますが、開発時はブラウザでも確認できるようにLiff Mockを使用しています。

```typescript
export async function setupLiff(redirectTo: string): Promise<void> {
  if (process.env.NODE_ENV === "development") {
    await setupMockLiff();
  } else {
    await liff.init({ liffId: LIFF_ID });
  }
```
## 📝 ライセンス

MIT License