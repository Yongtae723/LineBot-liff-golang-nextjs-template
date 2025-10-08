# LIFF アプリケーション

LINE Front-end Framework (LIFF) を使用したチャットアプリケーションです。

## 📋 概要

このアプリは、LINEユーザーがWebブラウザ上でGemini AIと会話できるチャットインターフェースを提供します。LINE BotとLIFFで会話履歴が完全に同期されます。

## 🚀 セットアップ

### 1. 環境変数の設定

`.env.example`をコピーして`.env.local`を作成し、必要な値を設定します:

```bash
cp .env.example .env.local
```

```.env.local
NEXT_PUBLIC_LIFF_ID=your-liff-id
NEXT_PUBLIC_BACKEND_URL=http://localhost:8080
NEXT_PUBLIC_SUPABASE_URL=http://localhost:54321
NEXT_PUBLIC_SUPABASE_ANON_KEY=your-supabase-anon-key
```

### 2. 依存関係のインストール

```bash
npm install
```

### 3. 開発サーバーの起動

```bash
npm run dev
```

アプリは [http://localhost:3000](http://localhost:3000) で起動します。

## 📦 主要な技術スタック

- **Next.js 15** - App Router
- **TypeScript** - 型安全性
- **Tailwind CSS** - スタイリング
- **shadcn/ui** - UIコンポーネント
- **@line/liff** - LINE LIFF SDK
- **@supabase/ssr** - Supabase認証

## 🏗️ プロジェクト構造

```
liff/
├── app/
│   ├── layout.tsx          # ルートレイアウト
│   ├── page.tsx            # LIFF初期化ページ
│   ├── login/
│   │   └── page.tsx        # 認証処理ページ
│   ├── home/
│   │   └── page.tsx        # チャット画面（認証後）
│   └── globals.css         # グローバルスタイル
├── components/
│   ├── chat/
│   │   ├── MessageList.tsx      # メッセージリスト
│   │   ├── MessageBubble.tsx    # メッセージバブル
│   │   └── InputBar.tsx         # 入力フォーム
│   └── ui/                      # shadcn/ui コンポーネント
├── lib/
│   ├── auth/
│   │   └── supabase.ts     # Supabase認証（Server Action）
│   ├── liff/
│   │   └── init.ts         # LIFF初期化ロジック
│   ├── supabase/
│   │   ├── client.ts       # Supabaseクライアント（ブラウザ）
│   │   └── server.ts       # Supabaseクライアント（サーバー）
│   └── api/
│       └── client.ts       # Backend API クライアント
├── types/
│   └── conversation.ts     # 型定義
├── middleware.ts           # 認証middleware
└── .env.example            # 環境変数サンプル
```

## 🔧 開発コマンド

```bash
npm run dev      # 開発サーバー起動
npm run build    # 本番ビルド
npm run start    # 本番サーバー起動
npm run format   # コードフォーマット（自動修正）
npm run lint     # リント実行（チェックのみ）
npm run check    # フォーマット + 型チェック
```

### Biome (Formatter & Linter)
このプロジェクトでは[Biome](https://biomejs.dev/)を使用してコードの品質を保っています。

- **自動フォーマット**: `npm run format`でコードスタイルを統一
- **リント**: `npm run lint`でコード品質をチェック
- **設定**: `biome.json`で設定をカスタマイズ可能

## 🌐 LIFF設定

### LINE Developers Console での設定

1. [LINE Developers Console](https://developers.line.biz/console/) にログイン
2. チャネルを選択
3. "LIFF" タブから新しいLIFFアプリを追加
4. 以下の設定を行う:
   - **Endpoint URL**: `https://your-domain.com` (本番) / `https://your-tunnel-url.ngrok.io` (開発)
   - **Scope**: `profile` と `openid`
   - **Module Mode**: OFF（推奨）

5. 発行されたLIFF IDを`.env.local`の`NEXT_PUBLIC_LIFF_ID`に設定

## 🔐 認証フロー

### ページ構成
- **`/` (トップページ)**: LIFF初期化専用ページ
- **`/login`**: 認証処理ページ（Backend API + Supabase）
- **`/home`**: チャット画面（認証後のみアクセス可能）

### 認証の流れ
1. LIFFアプリを開く（`/`）
2. LIFF SDK初期化
3. LINE Login（未ログインの場合、`/home`へリダイレクトURIを設定）
4. ログイン成功後、middlewareが未認証を検知して`/login`へリダイレクト
5. `/login`ページで：
   - LINE Access Token取得
   - Backend API `/api/v1/user/register` を呼び出し
   - Backend APIから`line_id`を取得
   - Supabaseに`signInWithPassword`（email: line_id, password: line_id）
   - `/home`へリダイレクト
6. チャット画面表示

### アーキテクチャのポイント
- **Middleware認証**: Next.jsのmiddlewareで全ページの認証状態をチェック
- **Server Action**: Supabase認証は`loginSupabase`関数（Server Action）で実行
- **Backend統合**: ユーザー登録とSupabase Auth Userの作成はBackend APIが担当
- **セキュリティ**: 
  - Supabaseセッションを使用してAPI呼び出しを保護
  - LINE Access Tokenの検証はBackend APIで実施
  - パスワードはLINE IDをそのまま使用（簡略化）

## 📱 機能

- ✅ LIFF認証統合
- ✅ リアルタイムチャット
- ✅ 会話履歴の表示
- ✅ LINE Botとの会話同期
- ✅ レスポンシブデザイン

## 🐛 トラブルシューティング

### LIFF認証エラー

**症状**: "LIFF認証に失敗しました"

**解決策**:
- `NEXT_PUBLIC_LIFF_ID`が正しく設定されているか確認
- LINE Developers ConsoleでLIFF Endpoint URLが正しいか確認
- LINEアプリ内またはLIFFブラウザで開いているか確認

### Backend接続エラー

**症状**: "初期化に失敗しました"

**解決策**:
- Backend APIが起動しているか確認（`http://localhost:8080/health`）
- `NEXT_PUBLIC_BACKEND_URL`が正しく設定されているか確認
- CORS設定が適切か確認

### Supabase認証エラー

**症状**: "認証に失敗しました"

**解決策**:
- Supabaseが起動しているか確認
- `NEXT_PUBLIC_SUPABASE_URL`と`NEXT_PUBLIC_SUPABASE_ANON_KEY`が正しいか確認
- Backend APIでユーザーが正常に作成されているか確認

## 📝 ライセンス

MIT License
