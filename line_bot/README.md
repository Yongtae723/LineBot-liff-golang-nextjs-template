# LINE Bot Service

LINE Messaging APIを使用したチャットボットサービスです。Gemini LLMと連携して自然な会話を提供します。

## 📦 主な機能

- ✅ LINE Messaging API Webhook処理
- ✅ メッセージ送受信
- ✅ Gemini LLMとの対話
- ✅ 会話履歴の保存（Supabase）
- ✅ 署名検証によるセキュリティ

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

サーバーは http://localhost:8000 で起動します。

## 📡 エンドポイント

### Webhook
```http
POST /webhook
X-Line-Signature: <signature>
```

LINE Messaging APIからのWebhookを受信します。

### ヘルスチェック
```http
GET /health
```

## 🔧 LINE Developers 設定

### 1. Messaging API設定

1. [LINE Developers Console](https://developers.line.biz/console/) にアクセス
2. Providerを選択 → Channelを選択
3. **Messaging API** タブ
4. Webhook設定:
   - Webhook URL: `https://your-domain.com/webhook`
   - Use webhook: **有効化**
   - Verify: テスト送信で確認

### 2. 必要な情報を取得

**Channel Secret**:
- Basic settings → Channel secret

**Channel Access Token**:
- Messaging API → Channel access token → **Issue**

### 3. 環境変数に設定

```bash
LINE_CHANNEL_SECRET=your-channel-secret
LINE_CHANNEL_TOKEN=your-channel-access-token
```

## 🔄 動作フロー

```
LINE User → メッセージ送信
    ↓
LINE Platform → Webhook送信
    ↓
LINE Bot Service
    ↓ 署名検証
    ↓ ユーザー確認
    ↓ 会話履歴取得
    ↓ Gemini LLM
    ↓ 会話保存
    ↓
LINE Platform → ユーザーに返信
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
line_bot/
├── cmd/
│   └── main.go              # エントリーポイント
├── config/
│   └── config.go            # 環境変数ローダー
├── middleware/
│   └── signature.go         # LINE署名検証
├── routes/
│   ├── router.go            # ルート定義
│   └── webhook.go           # Webhookハンドラー
├── logic/
│   ├── message/
│   │   └── handler.go       # メッセージ処理ロジック
│   └── follow/
│       └── handler.go       # フォローイベント処理
└── magefiles/
    └── magefile.go          # Mageタスク定義
```

## 🔒 セキュリティ

### 署名検証

LINE Platformからのリクエストは、`X-Line-Signature`ヘッダーで署名検証を行います：

```go
// middleware/signature.go
func ValidateSignature(channelSecret string) gin.HandlerFunc {
    return func(c *gin.Context) {
        signature := c.GetHeader("X-Line-Signature")
        body, _ := io.ReadAll(c.Request.Body)
        
        if !validateSignature(channelSecret, signature, body) {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid signature"})
            c.Abort()
            return
        }
        
        c.Set("body", body)
        c.Next()
    }
}
```

### Webhook URLの保護

- **HTTPSを使用**: 本番環境では必須
- **署名検証**: 全てのWebhookリクエストで実施
- **ログ記録**: 不正なリクエストを検出

## 🌍 環境変数

| 変数名 | 説明 | デフォルト |
|--------|------|-----------|
| `ENV` | 環境（local/development/staging/production） | `local` |
| `PORT` | サーバーポート | `8000` |
| `LINE_CHANNEL_SECRET` | LINE Channel Secret | - |
| `LINE_CHANNEL_TOKEN` | LINE Channel Access Token | - |
| `SUPABASE_URL` | Supabase API URL | - |
| `SUPABASE_KEY` | Supabase service role key | - |
| `GEMINI_API_KEY` | Google Gemini API Key | - |
| `BACKEND_URL` | Backend API URL | `http://localhost:8080` |
| `LIFF_APP_URL` | LIFF App URL | - |

## 🐳 Docker

### マルチステージビルド

Dockerfileはマルチステージビルドを使用しており、**プロジェクトのルートディレクトリから実行する必要があります**：

```bash
# プロジェクトルートディレクトリに移動
cd /path/to/LineBot-liff-golang-nextjs-template

# イメージビルド（line_botディレクトリを指定）
docker build --platform linux/amd64 -f line_bot/Dockerfile -t line-bot .

# コンテナ起動
docker run -p 8000:8000 --env-file line_bot/.env line-bot
```

### ビルドプロセス

1. **Build stage**: golang:1.24.2-alpineでrootユーザーとしてビルド実行
2. **Run stage**: alpine:latestで軽量なランタイム環境を構築
3. **セキュリティ**: 静的リンクされたバイナリで依存関係を最小化

## 🧪 ローカル開発（ngrok使用）

LINE Webhookはhttps://が必要なため、ローカル開発ではngrokを使用：

```bash
# ngrokインストール
brew install ngrok

# ngrok起動
ngrok http 8000

# 表示されたURLをLINE Developers Consoleに設定
# 例: https://xxxx-xx-xxx-xxx-xx.ngrok.io/webhook
```

## 🚀 デプロイ

Dockerfileが用意されているので、お好きな環境にデプロイできます：
- Google Cloud Run
- AWS ECS/Fargate
- Railway
- Fly.io

## 📝 ライセンス

MIT License