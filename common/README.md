# common - 共通Golangパッケージ

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

- `version.go`: ツールのバージョン定数管理
- `tasks/tasks.go`: 共通タスク定義（Generate, Fmt, Lint, Test, Update）

## 🚀 使い方

### Go Workspaceでの利用

このパッケージは、`backend`と`line_bot`から以下のように参照されます：

```go
import (
    "cookforyou.com/linebot-liff-template/common/llm"
    "cookforyou.com/linebot-liff-template/common/models"
    "cookforyou.com/linebot-liff-template/common/repository"
)
```

### LLMクライアントの使用例

```go
package main

import (
    "context"
    "cookforyou.com/linebot-liff-template/common/llm"
    "cookforyou.com/linebot-liff-template/common/models"
)

func main() {
    ctx := context.Background()
    
    // Geminiクライアントの初期化
    geminiClient, err := llm.NewGoogleGemini(
        ctx,
        "your-api-key",
        "gemini-2.5-flash-lite",
    )
    if err != nil {
        panic(err)
    }
    defer geminiClient.Close()
    
    // 会話履歴を準備
    history := []*models.Conversation{
        {
            Role:    models.RoleUser,
            Content: "こんにちは",
        },
    }
    
    // チャット実行
    response, err := geminiClient.Chat(ctx, history)
    if err != nil {
        panic(err)
    }
    
    println(response)
}
```

### リポジトリの使用例

```go
package main

import (
    "context"
    "cookforyou.com/linebot-liff-template/common/models"
    "cookforyou.com/linebot-liff-template/common/repository"
)

func main() {
    // Supabase初期化
    err := repository.InitSupabase(
        "http://localhost:54321",
        "your-service-role-key",
    )
    if err != nil {
        panic(err)
    }
    
    ctx := context.Background()
    
    // ユーザーリポジトリ
    userRepo := repository.NewUserRepo()
    user, err := userRepo.GetByLineID(ctx, "U1234567890abcdef")
    
    // 会話リポジトリ
    convRepo := repository.NewConversationRepo()
    conversations, err := convRepo.ListByUserID(ctx, user.LineID, 50)
}
```

## 🛠️ 開発コマンド

```bash
# 利用可能なタスクを表示
go run mage.go -l

# 共通タスク
go run mage.go test       # テスト実行
go run mage.go fmt        # コードフォーマット
go run mage.go lint       # リント実行
go run mage.go update     # 依存関係更新
```

## 📝 開発ガイド

### 新しいモデルの追加

1. `models/`に新しいファイルを作成
2. 構造体とタグを定義

```go
package models

type YourModel struct {
    ID        string    `json:"id" db:"id"`
    Field     string    `json:"field" db:"field"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
}
```

### 新しいリポジトリの追加

1. `repository/`にインターフェースと実装を作成

```go
package repository

type YourRepo interface {
    GetByID(ctx context.Context, id string) (*models.YourModel, error)
    Create(ctx context.Context, model *models.YourModel) error
}

type yourRepo struct {
    *BaseRepo
}

func NewYourRepo() YourRepo {
    return &yourRepo{BaseRepo: baseRepo}
}
```


### 新しい共通タスクの追加

`mage/tasks/tasks.go`に関数を追加：

```go
// YourTask runs your custom task
func YourTask() error {
    return sh.RunV("your-command")
}
```

## 🏗️ Mage構造

`common/mage/`ディレクトリには、すべてのGoサービスで共通利用できるタスクとコマンドが定義されています：

```
common/mage/
├── version.go         # ツールバージョン定数
└── tasks/
    └── tasks.go       # 共通タスク（Generate, Fmt, Lint, Test, Update）
```

各サービス（backend, line_bot）は、これらのタスクをラッパー関数として利用できます。

## 📄 ライセンス

MIT License