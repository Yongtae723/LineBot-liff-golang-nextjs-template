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
```
POST /api/v1/user/register
Content-Type: application/json

{
  "access_token": "LINE_ACCESS_TOKEN"
}
```

### 会話履歴取得
```
GET /api/v1/conversations?limit=50
Authorization: Bearer SUPABASE_JWT
```

### 会話送信
```
POST /api/v1/conversations
Authorization: Bearer SUPABASE_JWT
Content-Type: application/json

{
  "message": "こんにちは"
}
```

### ヘルスチェック
```
GET /health
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

### 1. ユーザー登録
1. LIFF AppがLINE Access Tokenを取得
2. `/api/v1/user/register`にAccess Tokenを送信
3. Backend APIがLINE APIでプロフィール取得
4. Supabase Authにユーザー作成
5. `users`テーブルにLINE User IDを保存

### 2. 認証済みリクエスト
1. LIFF AppがSupabase JWTを取得
2. `Authorization: Bearer JWT`ヘッダーで送信
3. `middleware.Auth()`がJWTを検証
4. `user_id`をcontextに設定
5. Handlerで`user_id`を使用

## 🌍 環境変数

| 変数名 | 説明 | デフォルト |
|--------|------|-----------|
| `ENV` | 環境（local/development/staging/production） | `local` |
| `PORT` | サーバーポート | `8080` |
| `SUPABASE_URL` | Supabase API URL | - |
| `SUPABASE_KEY` | Supabase service role key | - |
| `SUPABASE_JWT_SECRET` | JWT検証用シークレット | - |
| `GEMINI_API_KEY` | Google Gemini API Key | - |
| `COMMON_PASSWORD_PREFIX` | ユーザー登録時のパスワードPrefix | `linebot_` |

## 📝 ライセンス

MIT License

