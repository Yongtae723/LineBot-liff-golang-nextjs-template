# LINE Bot + LIFF + Golang + Next.js Template - 開発計画書

## 📋 目次
1. [プロジェクト概要](#1-プロジェクト概要)
2. [技術スタック](#2-技術スタック)
3. [システムアーキテクチャ](#3-システムアーキテクチャ)
4. [フォルダ構造](#4-フォルダ構造)
5. [実装計画](#5-実装計画)
6. [認証フロー](#6-認証フロー)
7. [会話同期の仕組み](#7-会話同期の仕組み)
8. [開発環境セットアップ](#8-開発環境セットアップ)
9. [デプロイ戦略](#9-デプロイ戦略)
10. [README構成](#10-readme構成)

---

## 1. プロジェクト概要

### 1.1 目的
このテンプレートリポジトリは、LINE BotとLIFF（LINE Front-end Framework）を使用して、ユーザーがGemini LLMと会話できる統合システムを構築するための再利用可能なボイラープレートです。LINE Botでの会話とLIFFウェブアプリでの会話が完全に同期され、ユーザーはどちらのインターフェースからでもシームレスに会話を継続できます。

### 1.2 主な特徴
- **統合された会話体験**: LINE BotとLIFFで同じ会話履歴を共有
- **Gemini LLM統合**: Google Geminiを使用した高度な会話機能
- **堅牢な認証システム**: LINE LIFF認証とSupabase認証の統合
- **モダンな技術スタック**: Golang（backend）+ Next.js（frontend）
- **スケーラブルなアーキテクチャ**: マイクロサービス志向の設計
- **完全なモノレポ構成**: Go Workspaceによる効率的な開発
- **本番環境対応**: Docker、Cloud Run、Cloudflareへのデプロイ設定完備

### 1.3 対象ユーザー
- LINEを活用したチャットボットを構築したい開発者
- GolangとNext.jsを使ったフルスタック開発を学びたいエンジニア
- 会話履歴の同期が必要なアプリケーションを開発したいチーム
- Supabaseを使ったバックエンド開発に興味がある開発者

---

## 2. 技術スタック

### 2.1 Backend
| 技術 | バージョン | 用途 |
|------|-----------|------|
| **Go** | 1.24.2 | メインプログラミング言語 |
| **Gin** | latest | WebフレームワークREST API構築 |
| **Go Workspaces** | - | モノレポ管理 |
| **Supabase Go Client** | latest | データベースアクセス |
| **LINE Bot SDK** | v8 | LINE Messaging API統合 |
| **Google Gemini API** | latest | LLM統合 |
| **zerolog** | latest | 構造化ログ |
| **Mage** | latest | タスクランナー |
| **mockery** | v2 | モック生成 |
| **testify** | latest | テストフレームワーク |

### 2.2 Frontend (LIFF)
| 技術 | バージョン | 用途 |
|------|-----------|------|
| **Next.js** | 15.x | React フレームワーク（App Router） |
| **React** | 19.x | UIライブラリ |
| **TypeScript** | 5.x | 型安全性 |
| **Tailwind CSS** | 3.x | スタイリング |
| **shadcn/ui** | latest | UIコンポーネント |
| **@line/liff** | latest | LIFF SDK |
| **@supabase/ssr** | latest | Supabase SSR |
| **Biome** | latest | リンター・フォーマッター |

### 2.3 Infrastructure
| 技術 | 用途 |
|------|------|
| **Supabase** | PostgreSQL、Authentication、Row Level Security |
| **Google Cloud Run** | Golang サービスのホスティング |
| **Cloudflare Pages** | LIFFアプリのホスティング |
| **Docker** | コンテナ化 |
| **GitHub Actions** | CI/CD |

---

## 3. システムアーキテクチャ

### 3.1 全体構成図

```
┌─────────────────────────────────────────────────────────────────┐
│                         LINE Platform                            │
│                                                                   │
│  ┌─────────────┐                          ┌─────────────┐       │
│  │  LINE Bot   │                          │    LIFF     │       │
│  │  (Messaging │                          │  (Web App)  │       │
│  │     API)    │                          │             │       │
│  └──────┬──────┘                          └──────┬──────┘       │
└─────────┼────────────────────────────────────────┼──────────────┘
          │                                         │
          │ Webhook                                 │ HTTPS
          │                                         │
          ▼                                         ▼
┌─────────────────┐                    ┌────────────────────┐
│   LINE Bot      │                    │   LIFF Frontend    │
│   Service       │                    │   (Next.js)        │
│   (Golang/Gin)  │                    │                    │
│   Port: 8000    │                    │   Cloudflare       │
│   Cloud Run     │                    │   Pages            │
└────────┬────────┘                    └─────────┬──────────┘
         │                                       │
         │                                       │ REST API
         │                                       │
         │              ┌────────────────────────┘
         │              │
         │              ▼
         │    ┌─────────────────┐
         │    │   Backend API   │
         │    │   (Golang/Gin)  │
         │    │   Port: 8080    │
         │    │   Cloud Run     │
         └───▶│                 │
              │  Auth Middleware│
              └────────┬────────┘
                       │
                       │ Supabase Client
                       │
                       ▼
              ┌─────────────────┐
              │   Supabase      │
              │                 │
              │  - PostgreSQL   │
              │  - Auth         │
              │  - RPC          │
              │  - RLS          │
              └────────┬────────┘
                       │
                       │ (Conversation Table)
                       │ user_id | message | role | created_at
                       │
                       └─────────────────────────────┐
                                                     │
                                                     ▼
                                          ┌──────────────────┐
                                          │  Gemini API      │
                                          │  (Google Cloud)  │
                                          └──────────────────┘

┌────────────────────────────────────────────────────────────┐
│                    Shared Go Package                        │
│  (go_pkg)                                                   │
│                                                              │
│  - models/        : データモデル                            │
│  - repository/    : Supabase データアクセス層               │
│  - llm/          : Gemini クライアント                      │
│  - mage/         : ビルドタスク定義                         │
└────────────────────────────────────────────────────────────┘
```

### 3.2 データフロー

#### 3.2.1 LINE Botでの会話フロー
```
1. User sends message via LINE
   ↓
2. LINE Platform sends webhook to LINE Bot Service
   ↓
3. LINE Bot Service:
   - Extracts user_id (LINE User ID)
   - Fetches conversation history from Supabase
   - Calls Gemini API with context
   ↓
4. Gemini API returns response
   ↓
5. LINE Bot Service:
   - Saves user message and bot response to Supabase
   - Sends response back to LINE Platform
   ↓
6. LINE Platform delivers message to user
```

#### 3.2.2 LIFF Appでの会話フロー
```
1. User opens LIFF in LINE
   ↓
2. LIFF App:
   - Initializes LIFF SDK
   - Gets LINE Access Token
   - Calls Backend API for authentication
   ↓
3. Backend API:
   - Verifies LINE token
   - Creates/retrieves Supabase user
   - Returns Supabase session token
   ↓
4. LIFF App stores session and renders chat UI
   ↓
5. User sends message via LIFF
   ↓
6. LIFF App calls Backend API with message
   ↓
7. Backend API:
   - Validates Supabase JWT
   - Fetches conversation history
   - Calls Gemini API
   - Saves to Supabase
   - Returns response
   ↓
8. LIFF App displays response in UI
```

### 3.3 会話の同期メカニズム

**共通データモデル:**
```sql
conversations (
  id UUID PRIMARY KEY,
  user_id TEXT NOT NULL,  -- LINE User ID
  role TEXT NOT NULL,     -- 'user' or 'assistant'
  content TEXT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
)
```

**同期のポイント:**
- LINE Bot ServiceとBackend APIは同じSupabaseデータベースを参照
- 両サービスは同じ`conversations`テーブルに読み書き
- `user_id`（LINE User ID）で会話を識別
- 時系列順（created_at）で会話履歴を取得

---

## 4. フォルダ構造

### 4.1 最終的なフォルダ構造

```
LineBot-liff-golang-nextjs-template/
├── .github/
│   └── workflows/
│       ├── backend-deploy.yml        # Backend API デプロイ
│       ├── line-bot-deploy.yml       # LINE Bot デプロイ
│       ├── liff-deploy.yml           # LIFF デプロイ
│       └── test.yml                  # テスト実行
│
├── go_pkg/                           # 共通 Golang パッケージ
│   ├── llm/
│   │   ├── gemini.go                 # Gemini クライアント実装
│   │   ├── llm.go                    # LLM インターフェース定義
│   │   ├── http.go                   # HTTPリトライロジック
│   │   └── mock_llmclient.go         # モック（自動生成）
│   ├── models/
│   │   ├── conversation.go           # 会話データモデル
│   │   └── user.go                   # ユーザーデータモデル
│   ├── repository/
│   │   ├── base.go                   # Supabase 初期化
│   │   ├── conversation.go           # 会話リポジトリ
│   │   ├── user.go                   # ユーザーリポジトリ
│   │   └── mock_*.go                 # モック（自動生成）
│   ├── mage/
│   │   ├── commands.go               # 共通 Mage コマンド
│   │   └── version.go                # バージョン定義
│   ├── .mockery.yaml                 # mockery 設定
│   ├── go.mod
│   ├── go.sum
│   ├── mage.go
│   └── README.md
│
├── backend/                          # Backend API サービス
│   ├── cmd/
│   │   └── main.go                   # エントリーポイント
│   ├── config/
│   │   └── config.go                 # 環境変数ローダー
│   ├── middleware/
│   │   ├── auth.go                   # Supabase JWT 認証
│   │   └── cors.go                   # CORS 設定
│   ├── routes/
│   │   ├── router.go                 # ルート定義
│   │   ├── conversation.go           # 会話エンドポイント
│   │   └── user.go                   # ユーザー登録エンドポイント
│   ├── logic/
│   │   ├── conversation/
│   │   │   ├── handler.go            # 会話ロジック
│   │   │   └── handler_test.go
│   │   └── user/
│   │       ├── register.go           # ユーザー登録ロジック
│   │       └── register_test.go
│   ├── Dockerfile
│   ├── .env.example
│   ├── go.mod
│   ├── go.sum
│   ├── mage.go
│   └── README.md
│
├── line_bot/                         # LINE Bot サービス
│   ├── cmd/
│   │   └── main.go
│   ├── config/
│   │   └── config.go
│   ├── middleware/
│   │   └── signature.go              # LINE署名検証
│   ├── routes/
│   │   ├── router.go
│   │   └── webhook.go                # Webhook ハンドラー
│   ├── logic/
│   │   └── message/
│   │       ├── handler.go            # メッセージ処理ロジック
│   │       └── handler_test.go
│   ├── Dockerfile
│   ├── .env.example
│   ├── go.mod
│   ├── go.sum
│   ├── mage.go
│   └── README.md
│
├── liff/                             # LIFF アプリ (Next.js)
│   ├── src/
│   │   ├── app/
│   │   │   ├── (authenticated)/      # 認証必須ルート
│   │   │   │   ├── chat/
│   │   │   │   │   └── page.tsx      # チャット画面
│   │   │   │   └── layout.tsx
│   │   │   ├── login/
│   │   │   │   └── page.tsx          # ログイン画面
│   │   │   ├── layout.tsx
│   │   │   └── globals.css
│   │   ├── components/
│   │   │   ├── ui/                   # shadcn/ui コンポーネント
│   │   │   ├── chat/
│   │   │   │   ├── ChatContainer.tsx
│   │   │   │   ├── MessageBubble.tsx
│   │   │   │   └── InputBar.tsx
│   │   │   └── providers/
│   │   │       └── LiffProvider.tsx  # LIFF初期化
│   │   ├── lib/
│   │   │   ├── supabase/
│   │   │   │   ├── client.ts         # Supabase クライアント
│   │   │   │   ├── server.ts         # サーバーサイド用
│   │   │   │   └── middleware.ts     # ミドルウェア
│   │   │   ├── api/
│   │   │   │   └── client.ts         # Backend API クライアント
│   │   │   └── liff/
│   │   │       └── init.ts           # LIFF初期化ロジック
│   │   ├── types/
│   │   │   ├── conversation.ts
│   │   │   └── supabase.ts           # 自動生成
│   │   └── middleware.ts             # Next.js ミドルウェア
│   ├── public/
│   ├── .env.local.example
│   ├── biome.json
│   ├── next.config.mjs
│   ├── package.json
│   ├── tailwind.config.ts
│   ├── tsconfig.json
│   ├── wrangler.toml                 # Cloudflare 設定
│   └── README.md
│
├── supabase/
│   ├── migrations/
│   │   └── 20250101000000_initial_schema.sql
│   ├── seed.sql                      # 初期データ
│   ├── config.toml
│   └── README.md
│
├── docs/
│   ├── ARCHITECTURE.md               # アーキテクチャ詳細
│   ├── DEVELOPMENT.md                # 開発ガイド
│   ├── DEPLOYMENT.md                 # デプロイガイド
│   └── API.md                        # API仕様書
│
├── .cursor/
│   └── rules/
│       ├── general.mdc               # 一般ルール
│       ├── golang.mdc                # Go開発ルール
│       └── liff.mdc                  # LIFF開発ルール
│
├── .gitignore
├── go.work                           # Go Workspace 定義
├── go.work.sum
├── LICENSE
└── README.md                         # メインREADME
```

### 4.2 cookLabからの抽出方針

| cookLabの要素 | テンプレートでの扱い | 理由 |
|--------------|-------------------|------|
| **go_pkg/llm/** | ✅ そのままコピー | Gemini統合のコア機能 |
| **go_pkg/repository/** | ✅ 簡略化してコピー | conversation, userのみに限定 |
| **go_pkg/models/** | ✅ 簡略化してコピー | conversation, userのみに限定 |
| **go_pkg/mage/** | ✅ そのままコピー | ビルドタスクの標準化 |
| **backend-service/middleware/auth** | ✅ そのままコピー | Supabase JWT認証 |
| **backend-service/config/** | ✅ そのままコピー | 環境変数管理 |
| **line_bot_golang/** | ✅ 簡略化してコピー | 会話機能のみに限定 |
| **liff-app/auth/** | ✅ そのままコピー | LINE → Supabase認証フロー |
| **liff-app/components/ui/** | ✅ shadcn/uiで再実装 | 汎用的なUIコンポーネント |
| **supabase/migrations/** | ⚠️ 最小限のスキーマのみ | users, conversationsテーブルのみ |
| **recipe関連の機能** | ❌ 除外 | テンプレートの範囲外 |
| **stripe関連の機能** | ❌ 除外 | テンプレートの範囲外 |

---

## 5. 実装計画

### 5.1 Phase 1: プロジェクト基盤構築（Step 1-5）✅ 完了

#### Step 1: リポジトリ初期化とGo Workspace設定 ✅
**作業内容:**
- [x] `.gitignore`作成（cookLabから参考）
- [x] `go.work`作成（3つのモジュール: go_pkg, backend, line_bot）
- [x] 各ディレクトリに`go.mod`初期化
- [x] ライセンス選定（MIT推奨）

**成果物:**
```
.gitignore
go.work
LICENSE
```

**実装メモ:**
- Go Workspaceで3つのモジュールを統合管理
- 共通パッケージ（go_pkg）をlocal replaceで参照
- MIT Licenseを採用

#### Step 2: go_pkg - 共通パッケージの構築 ✅
**作業内容:**
- [x] `go_pkg/models/`
  - `conversation.go`: 会話データ構造
  - `user.go`: ユーザーデータ構造
- [x] `go_pkg/repository/`
  - `base.go`: Supabase初期化（cookLabからコピー）
  - `conversation.go`: 会話CRUD操作
  - `user.go`: ユーザーCRUD操作
  - `auth.go`: Supabase Auth Admin API統合
- [x] `go_pkg/llm/`
  - `gemini.go`: Geminiクライアント（cookLabから簡略化）
  - `http.go`: HTTPリトライ（cookLabからコピー）
- [x] `go_pkg/mage/`
  - `magefile.go`: Mageコマンド（mockery生成）
- [x] `.mockery.yaml`設定
- [x] `mage.go`エントリーポイント作成
- [x] `go_pkg/README.md`作成

**実装メモ:**
- `llm.GoogleGemini`インターフェースは`Chat`メソッドのみに簡略化
- `models.Conversation`を直接使用してLLMと連携
- `repository.AuthRepo`を追加してSupabase Auth Admin API操作を分離
- Geminiモデル: `gemini-2.5-flash-lite`を使用

**重要な実装詳細:**

**models/conversation.go:**
```go
package models

import "time"

type Conversation struct {
    ID        string    `json:"id" db:"id"`
    UserID    string    `json:"user_id" db:"user_id"`
    Role      string    `json:"role" db:"role"` // "user" or "assistant"
    Content   string    `json:"content" db:"content"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
}

const (
    RoleUser      = "user"
    RoleAssistant = "assistant"
)
```

**models/user.go:**
```go
package models

import "time"

type User struct {
    ID        string    `json:"id" db:"id"`           // Supabase Auth UUID
    LineID    string    `json:"line_id" db:"line_id"` // LINE User ID
    Name      string    `json:"name" db:"name"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
```

**repository/conversation.go:**
```go
package repository

import (
    "context"
    "fmt"
    "github.com/supabase-community/postgrest-go"
    "your-module/go_pkg/models"
)

type ConversationRepo interface {
    GetByUserID(ctx context.Context, userID string, limit int) ([]*models.Conversation, error)
    Create(ctx context.Context, conv *models.Conversation) error
}

type conversationRepo struct {
    *BaseRepo
}

func NewConversationRepo() ConversationRepo {
    return &conversationRepo{BaseRepo: baseRepo}
}

func (r *conversationRepo) GetByUserID(ctx context.Context, userID string, limit int) ([]*models.Conversation, error) {
    var conversations []*models.Conversation
    _, err := r.Client.From("conversations").
        Select("*", "", false).
        Eq("user_id", userID).
        Order("created_at", &postgrest.OrderOpts{Ascending: true}).
        Limit(limit, "").
        ExecuteTo(&conversations)
    if err != nil {
        return nil, fmt.Errorf("repository: failed to get conversations: %w", err)
    }
    return conversations, nil
}

func (r *conversationRepo) Create(ctx context.Context, conv *models.Conversation) error {
    _, err := r.Client.From("conversations").
        Insert(conv, false, "", "", "").
        Execute()
    if err != nil {
        return fmt.Errorf("repository: failed to create conversation: %w", err)
    }
    return nil
}
```

#### Step 3: Supabaseスキーマ設計 ✅
**作業内容:**
- [x] Supabaseローカル環境構築
- [x] `supabase/initial_setup.sql`作成（GUI + `supabase db diff`ワークフロー用）
  - `users`テーブル
  - `conversations`テーブル
  - RLS（Row Level Security）ポリシー
  - インデックス
  - Trigger（updated_at自動更新）
- [x] `supabase/config.toml`作成
- [x] `supabase/README.md`作成（詳細な手順書）

**実装メモ:**
- GUI-First + Migration File アプローチを採用
- `supabase db diff -f migration_name`でマイグレーション生成
- RLSポリシーでユーザーごとのデータ分離を実装
- `conversations.user_id`はLINE User IDを使用（同期の要）

**初期マイグレーションスクリプト:**
```sql
-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users table (connected to auth.users)
CREATE TABLE public.users (
  id UUID REFERENCES auth.users(id) ON DELETE CASCADE PRIMARY KEY,
  line_id TEXT UNIQUE NOT NULL,
  name TEXT,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Conversations table
CREATE TABLE public.conversations (
  id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
  user_id TEXT NOT NULL, -- LINE User ID for easy lookup
  role TEXT NOT NULL CHECK (role IN ('user', 'assistant')),
  content TEXT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_conversations_user_id ON public.conversations(user_id);
CREATE INDEX idx_conversations_created_at ON public.conversations(created_at DESC);
CREATE INDEX idx_users_line_id ON public.users(line_id);

-- RLS Policies
ALTER TABLE public.users ENABLE ROW LEVEL SECURITY;
ALTER TABLE public.conversations ENABLE ROW LEVEL SECURITY;

-- Users can read their own data
CREATE POLICY "Users can view own data"
  ON public.users FOR SELECT
  USING (auth.uid() = id);

-- Users can update their own data
CREATE POLICY "Users can update own data"
  ON public.users FOR UPDATE
  USING (auth.uid() = id);

-- Users can view their own conversations (via line_id in auth metadata)
CREATE POLICY "Users can view own conversations"
  ON public.conversations FOR SELECT
  USING (
    user_id = (
      SELECT line_id FROM public.users WHERE id = auth.uid()
    )
  );

-- Users can insert their own conversations
CREATE POLICY "Users can insert own conversations"
  ON public.conversations FOR INSERT
  WITH CHECK (
    user_id = (
      SELECT line_id FROM public.users WHERE id = auth.uid()
    )
  );

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger for users table
CREATE TRIGGER update_users_updated_at
  BEFORE UPDATE ON public.users
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at_column();
```

#### Step 4: Backend API サービス構築 ✅
**作業内容:**
- [x] `backend/cmd/main.go`作成
  - Gin初期化
  - Supabase初期化
  - ルーティング設定
- [x] `backend/config/config.go`作成（環境変数ローダー）
- [x] `backend/middleware/auth.go`作成（Supabase JWT認証）
- [x] `backend/middleware/cors.go`作成
- [x] `backend/routes/router.go`作成
- [x] `backend/routes/user.go`作成（ユーザー登録エンドポイント）
- [x] `backend/routes/conversation.go`作成（会話エンドポイント）
- [x] `backend/logic/user/register.go`作成
- [x] `backend/logic/conversation/handler.go`作成
- [x] `backend/Dockerfile`作成
- [x] `backend/.env.example`作成
- [x] `backend/README.md`作成

**実装メモ:**
- **Method-Level DI**: ルートハンドラー内で依存関係を初期化
- **LINE認証統合**: LINE Access Token検証 + Supabase Auth連携
- **簡略化された認証フロー**: テンプレート用にLineIDを直接emailとして使用
- **エンドポイント**:
  - `POST /api/v1/user/register`: LINE tokenでユーザー登録
  - `GET /api/v1/conversations`: 会話履歴取得（認証必須）
  - `POST /api/v1/conversations`: 新しいメッセージ送信（認証必須）
- **Gemini統合**: gemini-2.5-flash-liteモデルを使用

**主要エンドポイント:**
```
POST   /api/v1/user/register          # LINE tokenでユーザー登録
GET    /api/v1/conversations           # 会話履歴取得
POST   /api/v1/conversations           # 新しい会話を送信
```

**backend/routes/conversation.go:**
```go
package routes

import (
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    "your-module/backend/logic/conversation"
)

type ConversationRequest struct {
    Message string `json:"message" binding:"required"`
}

type ConversationResponse struct {
    Response string `json:"response"`
}

func GetConversations(c *gin.Context) {
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
        return
    }

    limitStr := c.DefaultQuery("limit", "50")
    limit, err := strconv.Atoi(limitStr)
    if err != nil {
        limit = 50
    }

    handler := conversation.NewHandler()
    conversations, err := handler.GetHistory(c.Request.Context(), userID.(string), limit)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, conversations)
}

func PostConversation(c *gin.Context) {
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
        return
    }

    var req ConversationRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    handler := conversation.NewHandler()
    response, err := handler.ProcessMessage(c.Request.Context(), userID.(string), req.Message)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, ConversationResponse{Response: response})
}
```

**backend/logic/conversation/handler.go:**
```go
package conversation

import (
    "context"
    "fmt"
    "github.com/google/uuid"
    "your-module/go_pkg/llm"
    "your-module/go_pkg/models"
    "your-module/go_pkg/repository"
)

type Handler struct {
    convRepo repository.ConversationRepo
    userRepo repository.UserRepo
    llm      llm.GoogleGemini
}

func NewHandler() *Handler {
    ctx := context.Background()
    geminiClient, err := llm.NewGoogleGemini(
        ctx,
        os.Getenv("GEMINI_API_KEY"),
        "gemini-1.5-flash",
    )
    if err != nil {
        panic(fmt.Sprintf("failed to create gemini client: %v", err))
    }

    return &Handler{
        convRepo: repository.NewConversationRepo(),
        userRepo: repository.NewUserRepo(),
        llm:      geminiClient,
    }
}

func (h *Handler) GetHistory(ctx context.Context, userID string, limit int) ([]*models.Conversation, error) {
    user, err := h.userRepo.GetByID(ctx, userID)
    if err != nil {
        return nil, fmt.Errorf("failed to get user: %w", err)
    }

    conversations, err := h.convRepo.GetByUserID(ctx, user.LineID, limit)
    if err != nil {
        return nil, fmt.Errorf("failed to get conversations: %w", err)
    }

    return conversations, nil
}

func (h *Handler) ProcessMessage(ctx context.Context, userID string, message string) (string, error) {
    user, err := h.userRepo.GetByID(ctx, userID)
    if err != nil {
        return "", fmt.Errorf("failed to get user: %w", err)
    }

    // Save user message
    userConv := &models.Conversation{
        ID:      uuid.New().String(),
        UserID:  user.LineID,
        Role:    models.RoleUser,
        Content: message,
    }
    if err := h.convRepo.Create(ctx, userConv); err != nil {
        return "", fmt.Errorf("failed to save user message: %w", err)
    }

    // Get conversation history
    history, err := h.convRepo.GetByUserID(ctx, user.LineID, 20)
    if err != nil {
        return "", fmt.Errorf("failed to get history: %w", err)
    }

    // Build prompt with history
    prompt := buildPrompt(history)

    // Call Gemini
    response, err := h.llm.PredictForText(ctx, prompt)
    if err != nil {
        return "", fmt.Errorf("failed to get LLM response: %w", err)
    }

    // Save assistant response
    assistantConv := &models.Conversation{
        ID:      uuid.New().String(),
        UserID:  user.LineID,
        Role:    models.RoleAssistant,
        Content: response,
    }
    if err := h.convRepo.Create(ctx, assistantConv); err != nil {
        return "", fmt.Errorf("failed to save assistant message: %w", err)
    }

    return response, nil
}

func buildPrompt(history []*models.Conversation) string {
    var prompt string
    for _, conv := range history {
        if conv.Role == models.RoleUser {
            prompt += fmt.Sprintf("User: %s\n", conv.Content)
        } else {
            prompt += fmt.Sprintf("Assistant: %s\n", conv.Content)
        }
    }
    return prompt
}
```

#### Step 5: LINE Bot サービス構築 ✅
**作業内容:**
- [x] `line_bot/cmd/main.go`作成
- [x] `line_bot/config/config.go`作成
- [x] `line_bot/middleware/signature.go`作成（LINE署名検証）
- [x] `line_bot/routes/router.go`作成
- [x] `line_bot/routes/webhook.go`作成
- [x] `line_bot/logic/message.go`作成（メッセージ処理ロジック）
- [x] `line_bot/logic/follow.go`作成（フォローイベント処理）
- [x] `line_bot/Dockerfile`作成
- [x] `line_bot/.env.example`作成
- [x] `line_bot/README.md`作成

**実装メモ:**
- **Routes-Logic分離アーキテクチャ**: backendと同じパターンを採用
- **Logic層**: 文字列メッセージを返す（送信はroutes層で処理）
- **Constructor-Level DI**: logic handlersは依存関係を引数で受け取る
- **Webhook処理**: 
  - MessageEvent: ユーザーメッセージを処理してGemini応答を返信
  - FollowEvent: 新規フォロワーにLIFF登録案内を送信
- **LINE署名検証**: middlewareで実装（セキュリティ確保）
- **Gemini統合**: gemini-2.5-flash-liteモデルを使用

**line_bot/routes/webhook.go:**
```go
package routes

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
    "your-module/line_bot/logic/message"
)

func WebhookHandler(c *gin.Context) {
    events, ok := c.Get("line_events")
    if !ok {
        c.JSON(http.StatusBadRequest, gin.H{"error": "no events found"})
        return
    }

    handler := message.NewHandler()
    for _, event := range events.([]*messaging_api.Event) {
        if event.Type == messaging_api.EventType_MESSAGE {
            if messageEvent, ok := event.Message.(*messaging_api.TextMessageContent); ok {
                err := handler.HandleTextMessage(c.Request.Context(), event.Source.UserId, messageEvent.Text)
                if err != nil {
                    // Log error but continue processing
                    log.Error().Err(err).Msg("failed to handle message")
                }
            }
        }
    }

    c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
```

**line_bot/logic/message/handler.go:**
```go
package message

import (
    "context"
    "fmt"
    "os"
    "github.com/google/uuid"
    "github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
    "your-module/go_pkg/llm"
    "your-module/go_pkg/models"
    "your-module/go_pkg/repository"
)

type Handler struct {
    convRepo  repository.ConversationRepo
    llm       llm.GoogleGemini
    lineBot   *messaging_api.MessagingApiAPI
}

func NewHandler() *Handler {
    ctx := context.Background()
    geminiClient, err := llm.NewGoogleGemini(
        ctx,
        os.Getenv("GEMINI_API_KEY"),
        "gemini-1.5-flash",
    )
    if err != nil {
        panic(fmt.Sprintf("failed to create gemini client: %v", err))
    }

    bot, err := messaging_api.NewMessagingApiAPI(os.Getenv("LINE_CHANNEL_TOKEN"))
    if err != nil {
        panic(fmt.Sprintf("failed to create line bot: %v", err))
    }

    return &Handler{
        convRepo: repository.NewConversationRepo(),
        llm:      geminiClient,
        lineBot:  bot,
    }
}

func (h *Handler) HandleTextMessage(ctx context.Context, lineUserID string, text string) error {
    // Save user message
    userConv := &models.Conversation{
        ID:      uuid.New().String(),
        UserID:  lineUserID,
        Role:    models.RoleUser,
        Content: text,
    }
    if err := h.convRepo.Create(ctx, userConv); err != nil {
        return fmt.Errorf("failed to save user message: %w", err)
    }

    // Get conversation history
    history, err := h.convRepo.GetByUserID(ctx, lineUserID, 20)
    if err != nil {
        return fmt.Errorf("failed to get history: %w", err)
    }

    // Build prompt
    prompt := buildPrompt(history)

    // Call Gemini
    response, err := h.llm.PredictForText(ctx, prompt)
    if err != nil {
        return fmt.Errorf("failed to get LLM response: %w", err)
    }

    // Save assistant response
    assistantConv := &models.Conversation{
        ID:      uuid.New().String(),
        UserID:  lineUserID,
        Role:    models.RoleAssistant,
        Content: response,
    }
    if err := h.convRepo.Create(ctx, assistantConv); err != nil {
        return fmt.Errorf("failed to save assistant message: %w", err)
    }

    // Send response via LINE
    _, err = h.lineBot.ReplyMessage(&messaging_api.ReplyMessageRequest{
        ReplyToken: replyToken,
        Messages: []messaging_api.MessageInterface{
            &messaging_api.TextMessage{
                Text: response,
            },
        },
    })
    if err != nil {
        return fmt.Errorf("failed to send LINE message: %w", err)
    }

    return nil
}

func buildPrompt(history []*models.Conversation) string {
    var prompt string
    for _, conv := range history {
        if conv.Role == models.RoleUser {
            prompt += fmt.Sprintf("User: %s\n", conv.Content)
        } else {
            prompt += fmt.Sprintf("Assistant: %s\n", conv.Content)
        }
    }
    return prompt
}
```

### 5.2 Phase 2: LIFF アプリ構築（Step 6-8）🚧 次のフェーズ

#### Step 6: LIFF プロジェクト初期化 📝 未着手
**作業内容:**
- [ ] Next.js 15 プロジェクト作成（App Router）
- [ ] 必要なパッケージインストール
  - `@line/liff`
  - `@supabase/ssr`
  - `@supabase/supabase-js`
  - `shadcn/ui`
  - `tailwindcss`
  - `@radix-ui/react-*`（shadcn/ui dependencies）
- [ ] `biome.json`設定（cookLabから参考）
- [ ] `tsconfig.json`設定
- [ ] `next.config.mjs`設定
- [ ] `.env.local.example`作成
- [ ] `tailwind.config.ts`設定

**推奨手順:**
```bash
cd /Users/yongtae/Documents/personal/code/LineBot-liff-golang-nextjs-template
npx create-next-app@latest liff --typescript --tailwind --app --no-src-dir
cd liff
npm install @line/liff @supabase/ssr @supabase/supabase-js
npx shadcn@latest init
```

**package.jsonスクリプト:**
```json
{
  "scripts": {
    "dev": "next dev",
    "build": "next build",
    "start": "next start",
    "format": "biome check --write .",
    "lint": "biome check .",
    "type-check": "tsc --noEmit"
  }
}
```

#### Step 7: LIFF 認証フロー実装 📝 未着手
**作業内容:**
- [ ] `app/lib/liff/init.ts`作成（LIFF初期化ロジック）
- [ ] `app/lib/supabase/client.ts`作成（Supabaseクライアント）
- [ ] `app/lib/supabase/server.ts`作成（サーバーサイド用）
- [ ] `app/lib/api/client.ts`作成（Backend API呼び出し）
- [ ] `middleware.ts`作成（認証チェック）
- [ ] `app/login/page.tsx`作成（ログイン画面）
- [ ] `app/api/auth/signin/route.ts`作成（Supabase認証API）
- [ ] `components/providers/LiffProvider.tsx`作成（LIFF Context）

**認証フロー:**
1. LIFF初期化 → LINE Access Token取得
2. Backend API `/api/v1/user/register`呼び出し
3. Supabase `signInWithPassword`でセッション確立
4. `/chat`ページへリダイレクト

**src/lib/liff/init.ts:**
```typescript
import liff from "@line/liff";

export const initializeLiff = async (): Promise<string> => {
  const liffId = process.env.NEXT_PUBLIC_LIFF_ID;
  if (!liffId) {
    throw new Error("LIFF ID is not defined");
  }

  try {
    await liff.init({ liffId });

    if (!liff.isLoggedIn()) {
      liff.login();
      throw new Error("Not logged in, redirecting...");
    }

    const accessToken = liff.getAccessToken();
    if (!accessToken) {
      throw new Error("Failed to get access token");
    }

    return accessToken;
  } catch (error) {
    console.error("LIFF initialization failed:", error);
    throw error;
  }
};

export const getLiffProfile = async () => {
  if (!liff.isLoggedIn()) {
    throw new Error("Not logged in");
  }
  return await liff.getProfile();
};
```

**src/app/login/page.tsx:**
```typescript
"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { initializeLiff } from "@/lib/liff/init";

export default function LoginPage() {
  const router = useRouter();
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const login = async () => {
      try {
        const accessToken = await initializeLiff();

        // Call backend to register user
        const response = await fetch(
          `${process.env.NEXT_PUBLIC_BACKEND_URL}/api/v1/user/register`,
          {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
            },
            body: JSON.stringify({ access_token: accessToken }),
          }
        );

        if (!response.ok) {
          throw new Error("Failed to register user");
        }

        const { line_id } = await response.json();

        // Sign in to Supabase
        const supabaseResponse = await fetch("/api/auth/signin", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({ line_id }),
        });

        if (!supabaseResponse.ok) {
          throw new Error("Failed to sign in to Supabase");
        }

        router.push("/chat");
      } catch (err) {
        setError(err instanceof Error ? err.message : "Unknown error");
      }
    };

    login();
  }, [router]);

  if (error) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-red-500">Error: {error}</div>
      </div>
    );
  }

  return (
    <div className="flex items-center justify-center min-h-screen">
      <div>Logging in...</div>
    </div>
  );
}
```

#### Step 8: チャットUI実装 📝 未着手
**作業内容:**
- [ ] `types/conversation.ts`作成（型定義）
- [ ] `components/chat/ChatContainer.tsx`作成（チャットレイアウト）
- [ ] `components/chat/MessageBubble.tsx`作成（メッセージ表示）
- [ ] `components/chat/MessageList.tsx`作成（メッセージリスト）
- [ ] `components/chat/InputBar.tsx`作成（入力フォーム）
- [ ] `app/(authenticated)/chat/page.tsx`作成（チャット画面）
- [ ] `app/(authenticated)/layout.tsx`作成（認証済みレイアウト）
- [ ] shadcn/uiコンポーネント導入
  - `Button`
  - `Input`
  - `ScrollArea`
  - `Card`
  - `Avatar`

**実装ポイント:**
- リアルタイム更新（Supabase Realtime使用）
- オプティミスティックUI更新
- メッセージ送信中のローディング表示
- 自動スクロール（最新メッセージへ）

**src/types/conversation.ts:**
```typescript
export type Role = "user" | "assistant";

export interface Conversation {
  id: string;
  user_id: string;
  role: Role;
  content: string;
  created_at: string;
}
```

**src/lib/api/client.ts:**
```typescript
import { createClient } from "@/lib/supabase/client";
import type { Conversation } from "@/types/conversation";

export class ApiClient {
  private baseUrl: string;

  constructor() {
    this.baseUrl = process.env.NEXT_PUBLIC_BACKEND_URL || "";
  }

  private async getAuthHeader(): Promise<string> {
    const supabase = createClient();
    const {
      data: { session },
    } = await supabase.auth.getSession();

    if (!session?.access_token) {
      throw new Error("No access token");
    }

    return `Bearer ${session.access_token}`;
  }

  async getConversations(limit = 50): Promise<Conversation[]> {
    const authHeader = await this.getAuthHeader();

    const response = await fetch(
      `${this.baseUrl}/api/v1/conversations?limit=${limit}`,
      {
        headers: {
          Authorization: authHeader,
        },
      }
    );

    if (!response.ok) {
      throw new Error("Failed to fetch conversations");
    }

    return response.json();
  }

  async sendMessage(message: string): Promise<{ response: string }> {
    const authHeader = await this.getAuthHeader();

    const response = await fetch(`${this.baseUrl}/api/v1/conversations`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: authHeader,
      },
      body: JSON.stringify({ message }),
    });

    if (!response.ok) {
      throw new Error("Failed to send message");
    }

    return response.json();
  }
}

export const apiClient = new ApiClient();
```

**src/app/(authenticated)/chat/page.tsx:**
```typescript
"use client";

import { useEffect, useState } from "react";
import { ChatContainer } from "@/components/chat/ChatContainer";
import { MessageBubble } from "@/components/chat/MessageBubble";
import { InputBar } from "@/components/chat/InputBar";
import { apiClient } from "@/lib/api/client";
import type { Conversation } from "@/types/conversation";

export default function ChatPage() {
  const [conversations, setConversations] = useState<Conversation[]>([]);
  const [loading, setLoading] = useState(true);
  const [sending, setSending] = useState(false);

  useEffect(() => {
    const loadConversations = async () => {
      try {
        const data = await apiClient.getConversations();
        setConversations(data);
      } catch (error) {
        console.error("Failed to load conversations:", error);
      } finally {
        setLoading(false);
      }
    };

    loadConversations();
  }, []);

  const handleSendMessage = async (message: string) => {
    setSending(true);
    try {
      const response = await apiClient.sendMessage(message);

      // Add user message
      setConversations((prev) => [
        ...prev,
        {
          id: crypto.randomUUID(),
          user_id: "",
          role: "user",
          content: message,
          created_at: new Date().toISOString(),
        },
      ]);

      // Add assistant response
      setConversations((prev) => [
        ...prev,
        {
          id: crypto.randomUUID(),
          user_id: "",
          role: "assistant",
          content: response.response,
          created_at: new Date().toISOString(),
        },
      ]);
    } catch (error) {
      console.error("Failed to send message:", error);
    } finally {
      setSending(false);
    }
  };

  if (loading) {
    return <div>Loading...</div>;
  }

  return (
    <ChatContainer>
      <div className="flex flex-col h-full">
        <div className="flex-1 overflow-y-auto p-4 space-y-4">
          {conversations.map((conv) => (
            <MessageBubble
              key={conv.id}
              role={conv.role}
              content={conv.content}
              timestamp={conv.created_at}
            />
          ))}
        </div>
        <InputBar onSend={handleSendMessage} disabled={sending} />
      </div>
    </ChatContainer>
  );
}
```

### 5.3 Phase 3: インフラとCI/CD（Step 9-11）⏳ 後続フェーズ

#### Step 9: Docker化 📝 未着手
**作業内容:**
- [x] `backend/Dockerfile`作成済み（マルチステージビルド）
- [x] `line_bot/Dockerfile`作成済み（マルチステージビルド）
- [ ] `.dockerignore`作成（ルート）
- [ ] `docker-compose.yml`作成（ローカル開発用）
- [ ] Cloud Run用の最適化設定

**推奨Docker構成:**
- マルチステージビルド（小さいイメージサイズ）
- 非rootユーザーでの実行
- ヘルスチェックエンドポイント活用

#### Step 10: CI/CD設定 📝 未着手
**作業内容:**
- [ ] `.github/workflows/test.yml`作成（テスト実行）
- [ ] `.github/workflows/backend-deploy.yml`作成（Cloud Runデプロイ）
- [ ] `.github/workflows/line-bot-deploy.yml`作成（Cloud Runデプロイ）
- [ ] `.github/workflows/liff-deploy.yml`作成（Cloudflare Pagesデプロイ）
- [ ] GitHub Secretsの設定手順をドキュメント化

**対象環境:**
- Development（developブランチ）
- Production（mainブランチ、manual approval）

**テストワークフロー例:**
```yaml
name: Test

on:
  pull_request:
    branches: [main, develop]
  push:
    branches: [main, develop]

jobs:
  test-go:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.24.2'

      - name: Install dependencies
        run: |
          cd go_pkg && go mod download
          cd ../backend && go mod download
          cd ../line_bot && go mod download

      - name: Run tests
        run: |
          cd go_pkg && go test ./...
          cd ../backend && go test ./...
          cd ../line_bot && go test ./...

  test-liff:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Set up Node.js
        uses: actions/setup-node@v4
        with:
          node-version: '20'

      - name: Install dependencies
        working-directory: ./liff
        run: npm ci

      - name: Lint
        working-directory: ./liff
        run: npm run lint

      - name: Type check
        working-directory: ./liff
        run: npm run type-check
```

#### Step 11: ドキュメント作成 🚧 一部完了
**作業内容:**
- [x] `README.md`（メイン）作成済み - Quick Start重視
- [x] `supabase/README.md`作成済み - 詳細なセットアップ手順
- [x] `backend/README.md`作成済み
- [x] `line_bot/README.md`作成済み
- [x] `go_pkg/README.md`作成済み
- [ ] `docs/ARCHITECTURE.md`作成（システムアーキテクチャ詳細）
- [ ] `docs/DEVELOPMENT.md`作成（開発者ガイド）
- [ ] `docs/DEPLOYMENT.md`作成（デプロイ手順）
- [ ] `docs/API.md`作成（API仕様書）
- [ ] `liff/README.md`作成（LIFFアプリドキュメント）
- [ ] スクリーンショット・デモビデオの作成

**優先度:**
1. 動作するLIFFアプリの完成
2. デプロイ手順のドキュメント化
3. アーキテクチャドキュメント
4. デモ素材の作成

---

## 6. 認証フロー

### 6.1 初回ログインフロー（LIFF経由）

```
┌──────────┐
│   User   │
└────┬─────┘
     │
     │ 1. Opens LIFF in LINE
     ▼
┌─────────────────┐
│   LIFF App      │
│  (Next.js)      │
└────┬────────────┘
     │
     │ 2. liff.init()
     │ 3. liff.getAccessToken()
     │
     ▼
┌─────────────────────────────────────┐
│  POST /api/v1/user/register         │
│  Body: { access_token: "..." }      │
│                                     │
│  Backend API Service                │
└────┬────────────────────────────────┘
     │
     │ 4. Verify LINE token with LINE API
     │ 5. Get LINE user profile
     │
     ▼
┌─────────────────────────────────────┐
│  Supabase                           │
│                                     │
│  - Check if user exists             │
│    (by line_id)                     │
│  - If not, create auth user         │
│  - Create/update users table record │
│                                     │
└────┬────────────────────────────────┘
     │
     │ 6. Return line_id to LIFF
     │
     ▼
┌─────────────────┐
│   LIFF App      │
│                 │
│ 7. Sign in with │
│    Supabase     │
│    email/pass   │
│    (generated)  │
└────┬────────────┘
     │
     │ 8. Store session cookie
     │ 9. Redirect to /chat
     │
     ▼
┌─────────────────┐
│  Chat Page      │
└─────────────────┘
```

### 6.2 認証済みAPIリクエストフロー

```
┌──────────────┐
│  LIFF App    │
└──────┬───────┘
       │
       │ 1. Get Supabase session token
       │
       ▼
┌─────────────────────────────────────┐
│  GET /api/v1/conversations          │
│  Header: Authorization: Bearer XXX  │
│                                     │
│  Backend API Service                │
└────┬────────────────────────────────┘
     │
     │ 2. Extract JWT from header
     │
     ▼
┌─────────────────────────────────────┐
│  Auth Middleware                    │
│                                     │
│  - Validate JWT signature           │
│    (using SUPABASE_JWT_SECRET)      │
│  - Extract user_id from claims      │
│  - Set user_id in context           │
│                                     │
└────┬────────────────────────────────┘
     │
     │ 3. user_id available in handler
     │
     ▼
┌─────────────────────────────────────┐
│  Route Handler                      │
│                                     │
│  - Get user_id from context         │
│  - Query Supabase                   │
│  - Return data                      │
│                                     │
└─────────────────────────────────────┘
```

### 6.3 LINE Botからの会話

```
┌──────────┐
│   User   │
└────┬─────┘
     │
     │ 1. Sends message via LINE
     │
     ▼
┌─────────────────┐
│ LINE Platform   │
└────┬────────────┘
     │
     │ 2. Webhook POST
     │    /webhook
     │
     ▼
┌─────────────────────────────────────┐
│  LINE Bot Service                   │
│                                     │
│  - Verify LINE signature            │
│  - Extract LINE User ID             │
│  - Process message                  │
│                                     │
└────┬────────────────────────────────┘
     │
     │ 3. Query/Save to Supabase
     │    (using LINE User ID directly)
     │
     ▼
┌─────────────────────────────────────┐
│  Supabase                           │
│                                     │
│  conversations table                │
│  user_id = LINE User ID             │
│                                     │
└─────────────────────────────────────┘
```

**重要な設計決定:**
- LINE User IDを`conversations.user_id`として使用（認証不要で識別可能）
- LIFF AppはSupabase認証を通じて、`users`テーブルで自分の`line_id`を取得
- Backend APIは`user_id` (Supabase UUID) → `line_id` (LINE User ID) の変換を行う

---

## 7. 会話同期の仕組み

### 7.1 データモデルの統一

**Key Point: LINE User IDを共通識別子として使用**

```sql
-- conversations テーブル
CREATE TABLE conversations (
  id UUID PRIMARY KEY,
  user_id TEXT NOT NULL,  -- LINE User ID (共通識別子)
  role TEXT NOT NULL,
  content TEXT NOT NULL,
  created_at TIMESTAMPTZ
);

-- users テーブル
CREATE TABLE users (
  id UUID PRIMARY KEY,                    -- Supabase Auth UUID
  line_id TEXT UNIQUE NOT NULL,           -- LINE User ID
  name TEXT,
  created_at TIMESTAMPTZ,
  updated_at TIMESTAMPTZ
);
```

### 7.2 同期のロジック

#### LINE Bot Service
```go
// LINE WebhookからLINE User IDを直接取得
lineUserID := event.Source.UserId

// conversationsテーブルに直接保存
conversation := &models.Conversation{
    UserID: lineUserID,  // LINE User ID
    Role: "user",
    Content: message,
}
repo.Create(ctx, conversation)
```

#### Backend API Service
```go
// Supabase JWTからSupabase UUIDを取得
supabaseUserID := c.Get("user_id")  // from JWT claims

// usersテーブルでLINE User IDを取得
user, err := userRepo.GetByID(ctx, supabaseUserID)
lineUserID := user.LineID

// conversationsテーブルをLINE User IDでクエリ
conversations, err := convRepo.GetByUserID(ctx, lineUserID, limit)
```

### 7.3 同期保証の仕組み

**トランザクションは不要** (別々のインターフェースから書き込まれるため)

**整合性の保証:**
1. **一意性**: LINE User IDは一意であることがLINEプラットフォームで保証
2. **タイムスタンプ**: `created_at`で時系列順に取得可能
3. **冪等性**: 同じメッセージを複数回保存しても問題ない（UUIDで識別）

**リアルタイム同期（オプショナル）:**
- Supabase Realtimeを使用してLIFF Appで会話を自動更新
- LINE Bot経由のメッセージがLIFF Appにリアルタイムで反映

```typescript
// LIFF App
useEffect(() => {
  const supabase = createClient();
  
  const subscription = supabase
    .channel('conversations')
    .on(
      'postgres_changes',
      {
        event: 'INSERT',
        schema: 'public',
        table: 'conversations',
        filter: `user_id=eq.${lineUserID}`,
      },
      (payload) => {
        setConversations((prev) => [...prev, payload.new]);
      }
    )
    .subscribe();

  return () => {
    subscription.unsubscribe();
  };
}, [lineUserID]);
```

---

## 8. 開発環境セットアップ

### 8.1 必要なツール

| ツール | バージョン | インストール方法 |
|--------|-----------|---------------|
| **Go** | 1.24.2 | https://go.dev/dl/ |
| **Node.js** | 20.x | https://nodejs.org/ |
| **npm/yarn** | latest | Node.jsに同梱 |
| **Docker** | latest | https://www.docker.com/ |
| **Supabase CLI** | latest | `npm install -g supabase` |
| **Mage** | latest | `go install github.com/magefile/mage@latest` |
| **gcloud CLI** | latest | https://cloud.google.com/sdk/docs/install |

### 8.2 ローカル開発環境構築手順

#### Step 1: リポジトリクローン
```bash
git clone https://github.com/YOUR_ORG/LineBot-liff-golang-nextjs-template.git
cd LineBot-liff-golang-nextjs-template
```

#### Step 2: Supabaseローカル環境起動
```bash
cd supabase
supabase start
```

Supabaseの接続情報を記録:
```
API URL: http://localhost:54321
anon key: eyJh...
service_role key: eyJh...
JWT secret: super-secret-jwt-token-with-at-least-32-characters-long
```

#### Step 3: 環境変数設定

**backend/.env:**
```bash
ENV=local
SUPABASE_URL=http://localhost:54321
SUPABASE_KEY=eyJh... (service_role key)
SUPABASE_JWT_SECRET=super-secret-jwt-token-with-at-least-32-characters-long
GEMINI_API_KEY=your-gemini-api-key
PORT=8080
```

**line_bot/.env:**
```bash
ENV=local
SUPABASE_URL=http://localhost:54321
SUPABASE_KEY=eyJh... (service_role key)
GEMINI_API_KEY=your-gemini-api-key
LINE_CHANNEL_SECRET=your-line-channel-secret
LINE_CHANNEL_TOKEN=your-line-channel-token
PORT=8000
```

**liff/.env.local:**
```bash
NEXT_PUBLIC_LIFF_ID=your-liff-id
NEXT_PUBLIC_BACKEND_URL=http://localhost:8080
NEXT_PUBLIC_SUPABASE_URL=http://localhost:54321
NEXT_PUBLIC_SUPABASE_ANON_KEY=eyJh... (anon key)
COMMON_PASSWORD_PREFIX=your-common-password-prefix
```

#### Step 4: Go依存関係インストール
```bash
# Workspace sync
go work sync

# 各モジュール
cd go_pkg && go mod download && cd ..
cd backend && go mod download && cd ..
cd line_bot && go mod download && cd ..
```

#### Step 5: LIFF依存関係インストール
```bash
cd liff
npm install
```

#### Step 6: サービス起動

**Terminal 1: Backend API**
```bash
cd backend
go run mage.go run
# or: ENV=local go run cmd/main.go
```

**Terminal 2: LINE Bot**
```bash
cd line_bot
go run mage.go run
# or: ENV=local go run cmd/main.go
```

**Terminal 3: LIFF App**
```bash
cd liff
npm run dev
```

### 8.3 開発コマンド一覧

#### Go (backend, line_bot, go_pkg共通)
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

#### LIFF
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

# Supabase型定義生成
npm run gen:types
```

#### Supabase
```bash
# ローカル環境起動
supabase start

# ローカル環境停止
supabase stop

# マイグレーション作成
supabase migration new <migration_name>

# マイグレーション実行
supabase db reset

# 型定義生成
supabase gen types typescript --local > liff/src/types/supabase.ts
```

---

## 9. デプロイ戦略

### 9.1 環境構成

| 環境 | 説明 | ブランチ | 自動デプロイ |
|------|------|---------|------------|
| **Local** | 開発環境 | - | - |
| **Development** | 開発統合環境 | `develop` | ✅ |
| **Staging** | ステージング環境 | `staging` | ✅ |
| **Production** | 本番環境 | `main` | ✅ (manual approval) |

### 9.2 各サービスのデプロイ先

#### Backend API & LINE Bot
- **プラットフォーム**: Google Cloud Run
- **理由**: 
  - サーバーレス（使用量課金）
  - 自動スケーリング
  - Dockerイメージ対応
  - LINE Webhookに適している

**デプロイ手順（GitHub Actions）:**
```yaml
- name: Deploy to Cloud Run
  run: |
    gcloud run deploy ${{ env.SERVICE_NAME }} \
      --image gcr.io/${{ env.GCP_PROJECT }}/${{ env.SERVICE_NAME }}:${{ github.sha }} \
      --platform managed \
      --region asia-northeast1 \
      --allow-unauthenticated \
      --set-env-vars ENV=${{ env.ENV }} \
      --set-secrets SUPABASE_URL=supabase-url:latest,SUPABASE_KEY=supabase-key:latest
```

#### LIFF App
- **プラットフォーム**: Cloudflare Pages
- **理由**:
  - 無料枠が大きい
  - グローバルCDN
  - Next.js App Router対応（OpenNext）
  - 高速

**デプロイ手順（GitHub Actions）:**
```yaml
- name: Deploy to Cloudflare Pages
  uses: cloudflare/pages-action@v1
  with:
    apiToken: ${{ secrets.CLOUDFLARE_API_TOKEN }}
    accountId: ${{ secrets.CLOUDFLARE_ACCOUNT_ID }}
    projectName: linebot-liff-template
    directory: .vercel/output/static
```

#### Supabase
- **プラットフォーム**: Supabase Cloud
- **マイグレーション管理**: Supabase CLI + GitHub Actions

**マイグレーション自動実行:**
```yaml
- name: Run migrations
  run: |
    supabase link --project-ref ${{ secrets.SUPABASE_PROJECT_REF }}
    supabase db push
```

### 9.3 環境変数管理

#### GitHub Secrets（本番）
```
GCP_PROJECT_ID
GCP_SA_KEY
CLOUDFLARE_API_TOKEN
CLOUDFLARE_ACCOUNT_ID
SUPABASE_PROJECT_REF
SUPABASE_ACCESS_TOKEN
```

#### Cloud Secrets Manager（ランタイム）
Backend & LINE Bot で使用:
- `SUPABASE_URL`
- `SUPABASE_KEY`
- `SUPABASE_JWT_SECRET`
- `GEMINI_API_KEY`
- `LINE_CHANNEL_SECRET`
- `LINE_CHANNEL_TOKEN`

#### Cloudflare Environment Variables（LIFF）
- `NEXT_PUBLIC_LIFF_ID`
- `NEXT_PUBLIC_BACKEND_URL`
- `NEXT_PUBLIC_SUPABASE_URL`
- `NEXT_PUBLIC_SUPABASE_ANON_KEY`
- `COMMON_PASSWORD_PREFIX`

---

## 10. README構成

### 10.1 メインREADME.md構成

```markdown
# LINE Bot + LIFF + Golang + Next.js Template

[魅力的なヘッダー画像]

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.24.2-blue.svg)](https://go.dev/)
[![Next.js Version](https://img.shields.io/badge/Next.js-15.x-black.svg)](https://nextjs.org/)

## 🌟 Overview

このテンプレートは、LINE BotとLIFFを使用してGemini LLMと会話できる
フルスタックアプリケーションを構築するためのボイラープレートです。
LINE Botでの会話とLIFFウェブアプリでの会話が完全に同期されます。

### ✨ Features

- 🤖 **LINE Bot統合**: LINE Messaging APIを使った自然な会話
- 🌐 **LIFF Web App**: Next.js製のモダンなチャットUI
- 🧠 **Gemini LLM**: Google Geminiを使った高度な会話機能
- 🔄 **会話同期**: LINE BotとLIFFで会話履歴を完全共有
- 🔐 **堅牢な認証**: LINE認証とSupabase認証の統合
- 🚀 **本番環境対応**: Docker、Cloud Run、Cloudflare対応
- 📦 **モノレポ構成**: Go Workspaceによる効率的な開発
- 🧪 **テスト完備**: ユニットテストとモックの完全サポート

### 🏗️ Architecture

[アーキテクチャ図]

### 📸 Screenshots

[LINE Bot]    [LIFF App]    [会話同期のデモ]

## 🚀 Quick Start

### Prerequisites

- Go 1.24.2+
- Node.js 20+
- Docker
- Supabase CLI
- LINE Developers Account
- Google Cloud Account (for Gemini API)

### 1. Clone Repository

\```bash
git clone https://github.com/YOUR_ORG/LineBot-liff-golang-nextjs-template.git
cd LineBot-liff-golang-nextjs-template
\```

### 2. Setup Supabase

\```bash
cd supabase
supabase start
\```

### 3. Configure Environment Variables

Copy example files and fill in your credentials:

\```bash
cp backend/.env.example backend/.env
cp line_bot/.env.example line_bot/.env
cp liff/.env.local.example liff/.env.local
\```

### 4. Install Dependencies

\```bash
# Go
go work sync
cd go_pkg && go mod download && cd ..
cd backend && go mod download && cd ..
cd line_bot && go mod download && cd ..

# LIFF
cd liff && npm install
\```

### 5. Run Services

\```bash
# Terminal 1: Backend API
cd backend && go run mage.go run

# Terminal 2: LINE Bot
cd line_bot && go run mage.go run

# Terminal 3: LIFF App
cd liff && npm run dev
\```

## 📚 Documentation

- [Architecture](docs/ARCHITECTURE.md) - システムアーキテクチャの詳細
- [Development Guide](docs/DEVELOPMENT.md) - 開発者向けガイド
- [Deployment Guide](docs/DEPLOYMENT.md) - デプロイ手順
- [API Reference](docs/API.md) - API仕様書

## 🛠️ Development

### Common Commands

\```bash
# Test
go run mage.go test     # Go tests
npm run test            # LIFF tests

# Format
go run mage.go fmt      # Go format
npm run format          # LIFF format

# Lint
go run mage.go lint     # Go lint
npm run lint            # LIFF lint
\```

### Project Structure

\```
├── go_pkg/        # Shared Go packages
├── backend/       # Backend API service
├── line_bot/      # LINE Bot service
├── liff/          # LIFF App (Next.js)
├── supabase/      # Database migrations
└── docs/          # Documentation
\```

## 🌐 Deployment

### Backend & LINE Bot

Deployed to **Google Cloud Run**

\```bash
# Deploy to development
git push origin develop

# Deploy to production
git push origin main
\```

### LIFF App

Deployed to **Cloudflare Pages**

\```bash
cd liff
npm run deploy
\```

See [Deployment Guide](docs/DEPLOYMENT.md) for details.

## 🤝 Contributing

Contributions are welcome! Please read our [Contributing Guide](CONTRIBUTING.md).

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file.

## 🙏 Acknowledgments

- [LINE Messaging API](https://developers.line.biz/ja/services/messaging-api/)
- [LIFF (LINE Front-end Framework)](https://developers.line.biz/ja/docs/liff/)
- [Google Gemini](https://ai.google.dev/)
- [Supabase](https://supabase.com/)
- [Next.js](https://nextjs.org/)
- [Gin](https://gin-gonic.com/)

## 📧 Contact

- Issues: [GitHub Issues](https://github.com/YOUR_ORG/LineBot-liff-golang-nextjs-template/issues)
- Discussions: [GitHub Discussions](https://github.com/YOUR_ORG/LineBot-liff-golang-nextjs-template/discussions)

---

Made with ❤️ for the LINE Bot community
```

### 10.2 各サブディレクトリのREADME

#### go_pkg/README.md
- パッケージの概要
- 各ディレクトリの説明（models, repository, llm, mage）
- 使用方法とサンプルコード
- テストの書き方
- モック生成方法

#### backend/README.md
- サービスの概要
- エンドポイント一覧
- 環境変数説明
- ローカル実行方法
- デプロイ方法

#### line_bot/README.md
- サービスの概要
- Webhook設定方法
- 環境変数説明
- ローカル開発（ngrok使用）
- デプロイ方法

#### liff/README.md
- アプリ概要
- LIFF設定方法
- 環境変数説明
- 開発サーバー起動方法
- ビルドとデプロイ

#### supabase/README.md
- スキーマ概要
- マイグレーション実行方法
- ローカル環境構築
- RLS設定説明

---

## 11. 実装の優先順位とマイルストーン

### Milestone 1: 基盤構築 ✅ 完了
- ✅ Go Workspace設定
- ✅ go_pkg実装（models, repository, llm, mage）
- ✅ Supabaseスキーマ設計（GUI + Migration）
- ✅ Backend API基本実装
- 🎯 ゴール: Backend APIでGeminiと会話できること → **達成**

**実装内容:**
- Method-Level DI アーキテクチャ
- LINE Access Token検証
- Gemini 2.5 Flash Lite統合
- Supabase Auth Admin API統合

### Milestone 2: LINE Bot統合 ✅ 完了
- ✅ LINE Bot Service実装
- ✅ Webhook処理（Message, Follow Events）
- ✅ LINE署名検証ミドルウェア
- ✅ 会話のSupabase保存
- ✅ Routes-Logic分離アーキテクチャ
- 🎯 ゴール: LINE BotでGeminiと会話できること → **達成**

**実装内容:**
- Constructor-Level DI for logic handlers
- LIFF登録案内の自動送信
- Gemini 2.5 Flash Lite統合

### Milestone 3: LIFF App実装 🚧 次のフェーズ
- [ ] Next.js プロジェクト初期化
- [ ] LIFF認証フロー
- [ ] Supabase認証統合
- [ ] チャットUI実装
- 🎯 ゴール: LIFF AppでGeminiと会話できること

**推定作業時間:** 2-3日

### Milestone 4: 会話同期 ⏳ 保留
- ✅ 会話データモデル統一（完了）
- [ ] LINE Bot ↔ LIFF App 同期確認
- [ ] リアルタイム同期（オプション）
- 🎯 ゴール: LINE BotとLIFFで会話が完全同期すること

**前提条件:** Milestone 3完了

### Milestone 5: インフラとCI/CD ⏳ 保留
- ✅ Docker化（backend, line_bot）
- [ ] docker-compose.yml作成
- [ ] GitHub Actions設定
- [ ] Cloud Runデプロイ設定
- [ ] Cloudflare Pagesデプロイ設定
- 🎯 ゴール: 本番環境で動作すること

**推定作業時間:** 1-2日

### Milestone 6: ドキュメントと仕上げ 🚧 一部完了
- ✅ README作成（Quick Start重視）
- ✅ サブディレクトリREADME作成
- [ ] 詳細ドキュメント作成（docs/）
- [ ] スクリーンショット作成
- [ ] デモビデオ作成
- 🎯 ゴール: 誰でも使えるテンプレートとして公開

**推定作業時間:** 2-3日

---

## 現在の状態と次のアクション

### ✅ 完了した作業
1. **Go Backend（backend, line_bot, go_pkg）** - 完全実装済み
2. **Supabaseスキーマ** - 初期セットアップスクリプト完成
3. **認証フロー設計** - LINE → Backend → Supabase
4. **基本ドキュメント** - README群完成

### 🎯 次の優先タスク（Phase 2）
1. **LIFF プロジェクト初期化**（Step 6）
   - Next.js 15セットアップ
   - パッケージインストール
   - 基本設定ファイル作成

2. **LIFF 認証フロー実装**（Step 7）
   - LIFF SDK統合
   - Backend API連携
   - Supabase認証

3. **チャットUI実装**（Step 8）
   - shadcn/ui導入
   - チャットコンポーネント実装
   - リアルタイム更新

### 📊 進捗状況
- **Phase 1（基盤構築）**: 100% ✅
- **Phase 2（LIFF App）**: 0% 🚧
- **Phase 3（インフラ）**: 20%（Dockerfiles完成）⏳
- **全体進捗**: 約40%

---

## 12. 重要な設計決定と理由

### 12.1 なぜGo Workspaceを使うのか？
**理由:**
- 共通コード（go_pkg）を複数サービスで共有
- 依存関係の一元管理
- ローカル開発時のシームレスな統合

### 12.2 なぜLINE User IDを共通識別子にするのか？
**理由:**
- LINE BotはSupabase認証を経由しない
- LINE User IDは一意でLINEが保証
- 認証不要で会話履歴を保存できる
- LIFF側でLINE User IDとSupabase UUIDをマッピング

### 12.3 なぜCloudflare PagesでLIFFをホストするのか？
**理由:**
- Vercelより無料枠が大きい
- グローバルCDNで高速
- Next.js App Router対応（OpenNext経由）
- コスト効率が高い

### 12.4 なぜMiddlewareでSupabase JWT検証するのか？
**理由:**
- 認証ロジックの一元化
- セキュリティの向上
- エンドポイントごとに認証コードを書く必要がない

### 12.5 なぜMageを使うのか？
**理由:**
- Makefileより型安全
- Go開発者にとって自然
- クロスプラットフォーム対応
- 複雑なビルドタスクも記述可能

---

## 13. よくある質問（FAQ）

### Q1: LINE BotとLIFFの会話は本当に同期しますか？
**A:** はい、完全に同期します。両方とも同じSupabaseの`conversations`テーブルを参照しているため、どちらのインターフェースからでも同じ会話履歴が見えます。

### Q2: Gemini以外のLLMも使えますか？
**A:** はい、`go_pkg/llm`のインターフェースを実装すれば、OpenAI、Claude等も使用できます。

### Q3: 複数ユーザーで使えますか？
**A:** はい、LINE User IDで識別されるため、複数ユーザーが同時に利用できます。

### Q4: 会話履歴はどのくらい保存されますか？
**A:** デフォルトでは無制限ですが、古い会話を削除するバッチ処理を追加することを推奨します。

### Q5: 本番環境のコストはどのくらいですか？
**A:**
- **Supabase**: 無料プラン（500MB DB、2GB転送/月）
- **Cloud Run**: 従量課金（月200万リクエストまで無料）
- **Cloudflare Pages**: 無料プラン（無制限リクエスト）
- **Gemini API**: 従量課金（1日15リクエストまで無料）

小規模利用なら月額ほぼ無料で運用可能です。

---

## 14. トラブルシューティング

### Issue 1: LIFF認証が失敗する
**解決策:**
- LIFF IDが正しいか確認
- Backend URLが正しいか確認
- Supabase URLとキーが正しいか確認
- ブラウザのCookieを削除

### Issue 2: LINE Botが応答しない
**解決策:**
- Webhook URLが正しく設定されているか確認
- LINE署名検証が通っているか確認
- Cloud Runのログを確認
- Supabase接続が正常か確認

### Issue 3: 会話が同期しない
**解決策:**
- `conversations`テーブルの`user_id`が正しいか確認
- LINE User IDが正しく取得できているか確認
- Supabaseのログを確認

### Issue 4: Gemini APIがエラーを返す
**解決策:**
- API Keyが有効か確認
- クォータを超えていないか確認
- プロンプトが適切か確認
- リトライ機能が動作しているか確認

---

## 15. 次のステップとカスタマイズ

このテンプレートをベースに、以下のような機能を追加できます：

### 追加機能の例
- **画像認識**: Gemini Visionを使った画像分析
- **音声対応**: LINE音声メッセージの文字起こし
- **リッチメニュー**: LINE Botのリッチメニュー設定
- **通知機能**: 定期的なメッセージ配信
- **分析ダッシュボード**: 会話統計の可視化
- **多言語対応**: i18nによる国際化
- **チーム機能**: 複数ユーザーでの共同利用

### カスタマイズポイント
- **UI/UX**: LIFF AppのデザインをカスタマイズA
- **会話フロー**: システムプロンプトの調整
- **データモデル**: 追加テーブルの作成
- **認証**: OAuth2.0の追加サポート

---

## 16. コントリビューション歓迎

このテンプレートをより良くするために、以下の貢献を歓迎します：

- 🐛 **バグ報告**: Issue作成
- 💡 **機能提案**: Discussion開始
- 📝 **ドキュメント改善**: PR作成
- 🔧 **コード改善**: PR作成
- 🌐 **翻訳**: 多言語README

---

## 17. ライセンスとサポート

**ライセンス:** MIT License

このテンプレートは自由に使用、修正、配布できます。
商用利用も可能です。

**サポート:**
- GitHub Issues: バグ報告
- GitHub Discussions: 質問・提案
- Pull Requests: コントリビューション

---

**以上で詳細な開発計画書は完成です。**

このプランに従って実装を進めることで、堅牢で再利用可能な
LINE Bot + LIFF + Golang + Next.js テンプレートが完成します。

実装開始の準備ができたら、Phase 1のStep 1から順番に進めてください！

