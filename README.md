# LINE Bot + LIFF + Golang + Next.js Template

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.24.2-blue.svg)](https://go.dev/)
[![Next.js Version](https://img.shields.io/badge/Next.js-15.x-black.svg)](https://nextjs.org/)

LINE BotとLIFFを使用してGemini LLMと会話できるフルスタックアプリケーションのテンプレートです。LINE Botでの会話とLIFFウェブアプリでの会話が完全に同期されます。

## ✨ 主な特徴

- 🤖 **LINE Bot統合**: LINE Messaging APIを使った自然な会話
- 🌐 **LIFF Web App**: Next.js製のモダンなチャットUI
- 🧠 **Gemini LLM**: Google Geminiを使った高度な会話機能
- 🔄 **会話同期**: LINE BotとLIFFで会話履歴を完全共有
- 🔐 **堅牢な認証**: LINE認証とSupabase認証の統合
- 🚀 **本番環境対応**: Docker、Cloud Run、Cloudflare対応
- 📦 **モノレポ構成**: Go Workspaceによる効率的な開発

## 📋 必要な環境

- **Go**: 1.24.2以上
- **Node.js**: 20.x以上
- **Docker**: 最新版
- **Supabase CLI**: 最新版
- **LINE Developers Account**: [LINE Developers](https://developers.line.biz/)
- **Google Cloud Account**: Gemini API用

## 🚀 クイックスタート

### 1. リポジトリのクローン

```bash
git clone https://github.com/YOUR_USERNAME/LineBot-liff-golang-nextjs-template.git
cd LineBot-liff-golang-nextjs-template
```

### 2. Supabase CLIのインストール（初回のみ）

```bash
npm install -g supabase
```

### 3. Supabaseローカル環境の起動

```bash
cd supabase
supabase init  # 初回のみ実行
supabase start
```

起動後、以下のような接続情報が表示されます：

```
Started supabase local development setup.

         API URL: http://localhost:54321
     GraphQL URL: http://localhost:54321/graphql/v1
          DB URL: postgresql://postgres:postgres@localhost:54322/postgres
      Studio URL: http://localhost:54323
    Inbucket URL: http://localhost:54324
      JWT secret: super-secret-jwt-token-with-at-least-32-characters-long
        anon key: eyJh...
service_role key: eyJh...
```

**📝 重要**: これらの値を控えておいてください（環境変数設定で使用します）

### 4. マイグレーションを適用

```bash
# マイグレーションを適用（テーブル作成）
supabase db reset
```

これで`users`と`conversations`テーブルが作成されます。

**確認**: http://localhost:54323 のTable Editorで`users`と`conversations`テーブルが表示されればOK！

### 5. 環境変数の設定

#### **backend/.env**
```bash
ENV=local
SUPABASE_URL=http://localhost:54321
SUPABASE_KEY=eyJh...  # service_role key
SUPABASE_JWT_SECRET=super-secret-jwt-token-with-at-least-32-characters-long
GEMINI_API_KEY=your-gemini-api-key
PORT=8080
```

#### **line_bot/.env**
```bash
ENV=local
SUPABASE_URL=http://localhost:54321
SUPABASE_KEY=eyJh...  # service_role key
GEMINI_API_KEY=your-gemini-api-key
LINE_CHANNEL_SECRET=your-line-channel-secret
LINE_CHANNEL_TOKEN=your-line-channel-token
PORT=8000
```

#### **liff/.env.local**
```bash
NEXT_PUBLIC_LIFF_ID=your-liff-id
NEXT_PUBLIC_BACKEND_URL=http://localhost:8080
NEXT_PUBLIC_SUPABASE_URL=http://localhost:54321
NEXT_PUBLIC_SUPABASE_ANON_KEY=eyJh...  # anon key (from supabase start output)
```

### 6. Go依存関係のインストール

```bash
# Workspaceの同期
go work sync

# 各モジュールの依存関係インストール
cd go_pkg && go mod download && cd ..
cd backend && go mod download && cd ..
cd line_bot && go mod download && cd ..
```

### 7. LIFF依存関係のインストール

```bash
cd liff
npm install
```

### 8. サービスの起動

**Terminal 1: Backend API**
```bash
cd backend
go run mage.go run
# または: ENV=local go run cmd/main.go
```

**Terminal 2: LINE Bot**
```bash
cd line_bot
go run mage.go run
# または: ENV=local go run cmd/main.go
```

**Terminal 3: LIFF App**
```bash
cd liff
npm run dev
```

### 9. 動作確認

各サービスが起動したら、以下のURLにアクセスできます：

- **Backend API**: http://localhost:8080/health （ヘルスチェック）
- **LINE Bot**: http://localhost:8000/health （ヘルスチェック）
- **LIFF App**: http://localhost:3000 （認証ページ）
- **LIFF Chat**: http://localhost:3000/home （チャット画面、認証後）

#### LIFF アプリの動作確認
1. LINEアプリでLIFFを開く、または開発環境で http://localhost:3000 にアクセス
2. LIFF初期化後、LINE Login画面が表示されます（未ログインの場合）
3. ログイン成功後、`/login`ページで認証処理が実行されます
4. 認証完了後、自動的に`/home`へリダイレクトされます
5. チャット画面でGeminiと会話できます

**認証フローの詳細:**
```
/ (LIFF初期化) → LINE Login → /login (Backend API + Supabase認証) → /home (チャット画面)
```
- **Supabase Studio**: http://localhost:54323 （DB管理用）

## 🎯 次のステップ

テンプレートが起動できたら、以下を設定してください：

### LINE Developers設定

1. [LINE Developers Console](https://developers.line.biz/console/) でチャネル作成
2. **Messaging API**タブ:
   - Channel Secret → `LINE_CHANNEL_SECRET`
   - Channel Access Token → `LINE_CHANNEL_TOKEN`
   - Webhook URL: `https://your-domain.com/webhook` (本番環境)
3. **LIFF**タブ:
   - LIFF URL: `https://your-liff-domain.com`
   - LIFF ID → `NEXT_PUBLIC_LIFF_ID`

### Gemini API設定

1. [Google AI Studio](https://aistudio.google.com/app/apikey) でAPI Key作成
2. API Key → `GEMINI_API_KEY`

## 📦 プロジェクト構造

```
LineBot-liff-golang-nextjs-template/
├── go_pkg/          # 共通Golangパッケージ
│   ├── llm/         # Geminiクライアント
│   ├── models/      # データモデル
│   ├── repository/  # Supabaseアクセス層
│   └── mage/        # ビルドタスク
├── backend/         # Backend APIサービス
├── line_bot/        # LINE Botサービス
├── liff/            # LIFFアプリ (Next.js)
├── supabase/        # Supabaseマイグレーション
└── docs/            # ドキュメント
```

## 🛠️ 開発コマンド

### Go (backend, line_bot, go_pkg共通)

```bash
# テスト実行
go run mage.go test

# フォーマット
go run mage.go fmt

# リント
go run mage.go lint

# モック生成
go run mage.go generate

# 依存関係更新
go run mage.go update
```

### LIFF

```bash
# 開発サーバー起動
npm run dev

# ビルド
npm run build

# フォーマット
npm run format

# リント
npm run lint

# 型チェック
npm run type-check
```

### Supabase

```bash
# ローカル環境起動
supabase start

# ローカル環境停止
supabase stop

# ローカル環境リセット（データ削除）
supabase db reset

# マイグレーション作成
supabase migration new <migration_name>

# 型定義生成
supabase gen types typescript --local > liff/src/types/supabase.ts
```

## 🔧 トラブルシューティング

### Supabaseが起動しない

```bash
# 他のプロジェクトのSupabaseを停止
cd ~/path/to/other/project
supabase stop

# このプロジェクトで起動
cd /Users/yongtae/Documents/personal/code/LineBot-liff-golang-nextjs-template
supabase start
```

### マイグレーションが適用されない

```bash
# マイグレーションファイルを確認
ls supabase/migrations/

# 強制リセット
supabase db reset
```

### Go依存関係エラー

```bash
# 全モジュールを更新
cd go_pkg && go mod tidy && cd ..
cd backend && go mod tidy && cd ..
cd line_bot && go mod tidy && cd ..
go work sync
```

## 📚 ドキュメント

- [アーキテクチャ](docs/ARCHITECTURE.md) - システムアーキテクチャの詳細
- [開発ガイド](docs/DEVELOPMENT.md) - 開発者向けガイド
- [デプロイガイド](docs/DEPLOYMENT.md) - デプロイ手順
- [API仕様書](docs/API.md) - API仕様書

## 🤝 コントリビューション

プルリクエストを歓迎します！大きな変更の場合は、まずIssueを開いて変更内容を議論してください。

## 📝 ライセンス

MIT License - 詳細は[LICENSE](LICENSE)ファイルをご覧ください。

## 🙏 謝辞

- [LINE Messaging API](https://developers.line.biz/ja/services/messaging-api/)
- [LIFF (LINE Front-end Framework)](https://developers.line.biz/ja/docs/liff/)
- [Google Gemini](https://ai.google.dev/)
- [Supabase](https://supabase.com/)
- [Next.js](https://nextjs.org/)
- [Gin](https://gin-gonic.com/)

---

Made with ❤️ for the LINE Bot community

