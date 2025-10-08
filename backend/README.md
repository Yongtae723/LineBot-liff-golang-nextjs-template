# Backend API Service

LINE BotとLIFFから使用されるREST APIサーバーです。

## 📦 主な機能

- ✅ ユーザー登録（LINE認証 → Supabase Auth）
- ✅ 会話履歴取得
- ✅ 会話送信（Gemini LLMとの対話）
- ✅ Supabase JWT認証

## 🚀 起動方法

### 1. 環境変数設定

```bash
cp .env.example .env
# .envファイルを編集して必要な値を設定
```

### 2. 依存関係インストール

```bash
go mod download
```

### 3. サーバー起動

```bash
# Mageを使用
go run mage.go run

# または直接実行
ENV=local go run cmd/main.go
```

サーバーは http://localhost:8080 で起動します。

## 📡 API エンドポイント

### ユーザー登録
```http
POST /api/v1/user/register/liff
Content-Type: application/json

{
  "access_token": "LINE_ACCESS_TOKEN"
}
```

**レスポンス:**
```json
{
  "line_id": "U1234567890abcdef"
}
```

### 会話履歴取得
```http
GET /api/v1/conversations?limit=50
Authorization: Bearer SUPABASE_JWT
```

**レスポンス:**
```json
[
  {
    "id": "uuid",
    "user_id": "line_user_id",
    "role": "user",
    "content": "こんにちは",
    "created_at": "2025-01-01T00:00:00Z"
  }
]
```

### 会話送信
```http
POST /api/v1/conversations
Authorization: Bearer SUPABASE_JWT
Content-Type: application/json

{
  "message": "こんにちは"
}
```

**レスポンス:**
```json
{
  "response": "こんにちは！何かお手伝いできることはありますか？"
}
```

### ヘルスチェック
```http
GET /health
```

**レスポンス:**
```json
{
  "status": "ok"
}
```

## 🛠️ 開発コマンド

```bash
# テスト実行
go run mage.go test

# フォーマット
go run mage.go fmt

# リント
go run mage.go lint

# 依存関係更新
go run mage.go update
```

## 📁 プロジェクト構造

```
backend/
├── cmd/
│   └── main.go              # エントリーポイント
├── config/
│   └── config.go            # 環境変数ローダー
├── middleware/
│   ├── auth.go              # Supabase JWT認証
│   └── cors.go              # CORS設定
├── routes/
│   ├── router.go            # ルート定義
│   ├── user.go              # ユーザー関連エンドポイント
│   └── conversation.go      # 会話関連エンドポイント
├── logic/
│   ├── user/
│   │   └── register.go      # ユーザー登録ロジック
│   └── conversation/
│       └── handler.go       # 会話処理ロジック
└── magefiles/
    └── magefile.go          # Mageタスク定義
```

## 🔒 認証フロー

### 1. ユーザー登録（LIFF経由）

```
1. LIFF AppがLINE Access Tokenを取得
2. POST /api/v1/user/register/liff にAccess Tokenを送信
3. Backend APIがLINE APIでプロフィール取得
4. Supabase Authにユーザー作成
5. usersテーブルにLINE User IDを保存
6. LINE IDを返却
```

### 2. 認証済みリクエスト

```
1. LIFF AppがSupabase JWTを取得
2. Authorization: Bearer JWT ヘッダーで送信
3. middleware.Auth()がJWTを検証
4. user_idをcontextに設定
5. Handlerでuser_idを使用
```

## 🌍 環境変数

| 変数名 | 説明 | デフォルト |
|--------|------|-----------|
| `ENV` | 環境（local/development/staging/production） | `local` |
| `PORT` | サーバーポート | `8080` |
| `SUPABASE_URL` | Supabase API URL | - |
| `SUPABASE_KEY` | Supabase service role key | - |
| `SUPABASE_JWT_SECRET` | JWT検証用シークレット | - |
| `GEMINI_API_KEY` | Google Gemini API Key | - |
| `LINE_CHANNEL_ID` | LINE Channel ID | - |

## 🐳 Docker

### マルチステージビルド

Dockerfileはマルチステージビルドを使用しており、**プロジェクトのルートディレクトリから実行する必要があります**：

```bash
# プロジェクトルートディレクトリに移動
cd /path/to/LineBot-liff-golang-nextjs-template

# イメージビルド（backendディレクトリを指定）
docker build --platform linux/amd64 -f backend/Dockerfile -t backend-api .

# コンテナ起動
docker run -p 8080:8080 --env-file backend/.env backend-api
```

### ビルドプロセス

1. **Build stage**: golang:1.24.2-alpineでrootユーザーとしてビルド実行
2. **Run stage**: alpine:latestで軽量なランタイム環境を構築
3. **セキュリティ**: 静的リンクされたバイナリで依存関係を最小化

## 🚀 デプロイ

Dockerfileが用意されているので、お好きな環境にデプロイできます：
- Google Cloud Run
- AWS ECS/Fargate
- Railway
- Fly.io

## 📝 ライセンス

MIT License