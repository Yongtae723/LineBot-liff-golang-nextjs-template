# common - Shared Go Package

このパッケージは、LINE BotとBackend APIの両方で使用される共通のGoコードを含んでいます。

## 📦 パッケージ構成

### llm/
Google Gemini LLMクライアントの実装

- `gemini.go`: Geminiクライアントの実装
- `http.go`: HTTPリトライロジック

### models/
データモデル定義

- `user.go`: ユーザーデータ
- `conversation.go`: 会話データ

### repository/
Supabaseデータアクセス層

- `base.go`: Supabase初期化
- `user.go`: ユーザーリポジトリ
- `conversation.go`: 会話リポジトリ
- `auth.go`: Supabase Auth Admin API

### mage/
Mageビルドタスク定義

- **commands.go**: 共通コマンド（`RunMockery`, `RunMockgen`）
- **version.go**: ツールのバージョン定数管理
- **generate.go**: Generator登録機構
- **generate/generate.go**: Generatorsスライス
- **tasks/tasks.go**: 共通タスク定義（Generate, Fmt, Lint, Test, Update）

## 🚀 開発コマンド

```bash
# 利用可能なタスクを表示
go run github.com/magefile/mage@latest -l

# 共通タスク
go run github.com/magefile/mage@latest generate   # モック生成
go run github.com/magefile/mage@latest test       # テスト実行
go run github.com/magefile/mage@latest fmt        # コードフォーマット
go run github.com/magefile/mage@latest lint       # リント実行
go run github.com/magefile/mage@latest update     # 依存関係更新
```

## 🏗️ Mage構造

`common/mage/`ディレクトリには、すべてのGoサービスで共通利用できるタスクとコマンドが定義されています：

```
common/mage/
├── commands.go        # 共通コマンド（RunMockery, RunMockgen）
├── version.go         # ツールバージョン定数
├── generate.go        # Generator登録機構
├── generate/
│   └── generate.go    # Generatorsスライス
└── tasks/
    └── tasks.go       # 共通タスク（Generate, Fmt, Lint, Test, Update）
```

各サービス（backend, line_bot）は、これらのタスクをラッパー関数として利用できます。cookLabと同じ構造で、再利用性の高いタスク管理を実現しています。

## 📝 開発ガイド

### 新しいモデルの追加
1. `models/`に新しいファイルを作成
2. 構造体とタグを定義
3. 必要に応じてconstを定義

### 新しいリポジトリの追加
1. `repository/`にインターフェースと実装を作成
2. `.mockery.yaml`にインターフェースを追加
3. `go run github.com/magefile/mage@latest generate`でモックを生成

### 新しい共通タスクの追加
1. `mage/tasks/tasks.go`に関数を追加
2. 各サービスの`magefiles/magefile.go`でラッパー関数を作成

## 📄 ライセンス

MIT License
