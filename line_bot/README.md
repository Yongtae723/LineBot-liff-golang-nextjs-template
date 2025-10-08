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

サーバーは http://localhost:8081 で起動します。

## 📡 エンドポイント

### Webhook
```
POST /webhook
X-Line-Signature: <signature>

LINE Messaging APIからのWebhookを受信
```

### ヘルスチェック
```
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

### 2. 必要な情報を取得

**Channel secret**:
- Basic settings → Channel secret

**Channel access token**:
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
├── handler/
│   ├── webhook.go           # Webhookハンドラー
│   └── message.go           # メッセージ処理
└── magefiles/
    └── magefile.go          # Mageタスク定義
```

## 🌍 環境変数

| 変数名 | 説明 | デフォルト |
|--------|------|-----------|
| `ENV` | 環境（local/development/staging/production） | `local` |
| `PORT` | サーバーポート | `8081` |
| `LINE_CHANNEL_SECRET` | LINE Channel Secret | - |
| `LINE_CHANNEL_TOKEN` | LINE Channel Access Token | - |
| `SUPABASE_URL` | Supabase API URL | - |
| `SUPABASE_KEY` | Supabase service role key | - |
| `GEMINI_API_KEY` | Google Gemini API Key | - |

## 🔒 セキュリティ

### 署名検証

LINE Platformからのリクエストは、`X-Line-Signature`ヘッダーで署名検証を行います：

```go
// middleware/signature.go
func ValidateSignature(channelSecret string) gin.HandlerFunc
```

### Webhook URLの保護

- **HTTPSを使用**: 本番環境では必須
- **署名検証**: 全てのWebhookリクエストで実施
- **ログ記録**: 不正なリクエストを検出

## 📝 ライセンス

MIT License

