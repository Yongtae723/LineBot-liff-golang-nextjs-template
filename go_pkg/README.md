# go_pkg - Shared Go Package

このパッケージは、LINE BotとBackend APIの両方で使用される共通のGoコードを含んでいます。

## 📦 パッケージ構成

### llm/
Google Gemini LLMクライアントの実装

- `gemini.go`: Geminiクライアントの実装
- `http.go`: HTTPリトライロジック
- `llm.go`: LLMインターフェース定義

### models/
データモデル定義

- `user.go`: ユーザーデータ
- `conversation.go`: 会話データ

### repository/
Supabaseデータアクセス層

- `base.go`: Supabase初期化
- `user.go`: ユーザーリポジトリ
- `conversation.go`: 会話リポジトリ

### mage/
Mageビルドタスク定義

- `commands.go`: 共通コマンド
- `version.go`: バージョン定義

## 🚀 使用方法

### テスト実行
```bash
go run mage.go test
```

### フォーマット
```bash
go run mage.go fmt
```

### リント
```bash
go run mage.go lint
```

### モック生成
```bash
go run mage.go generate
```

## 📝 開発ガイド

### 新しいモデルの追加
1. `models/`に新しいファイルを作成
2. 構造体とタグを定義
3. 必要に応じてconstを定義

### 新しいリポジトリの追加
1. `repository/`にインターフェースと実装を作成
2. `.mockery.yaml`にインターフェースを追加
3. `go run mage.go generate`でモックを生成

## 📄 ライセンス

MIT License

