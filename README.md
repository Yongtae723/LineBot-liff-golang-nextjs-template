# LINE Bot + LIFF + Golang + Next.js Template

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.24.2-blue.svg)](https://go.dev/)
[![Next.js Version](https://img.shields.io/badge/Next.js-15.x-black.svg)](https://nextjs.org/)

LINE BotとLIFFを使用してGemini LLMと会話できるフルスタックアプリケーションのテンプレートです。LINE Botでの会話とLIFFウェブアプリでの会話が完全に同期されます。

TODO: Blog URLの追加
このテンプレートを公開した背景や、願いは[Qiitaの記事]()を御覧ください

> 💡 **実例**: このテンプレートは、[CookForYou](https://cookforyou.jp)（料理レシピ提案サービス）の開発で培った知見を元に作成されています。 以下のQRコード、またはボタンから友達追加して、Lineを通したアプリを体験してみてください。

[![友だち追加はこちらから！](https://img.shields.io/badge/LINE-友だち追加-00C300?style=for-the-badge&logo=line&logoColor=white)](https://line.me/R/ti/p/@807rpwtd)

![QR Code](./images/QR.png)





## ✨ 主な特徴


このレポジトリコードは、lineBotで会話した内容がLIFFにも反映されることで、BotとLIFFの連携を体験する、ミニマムアプリです。

![linebot](./images/linebot.png)
↑のBotの会話を↓liffでも継続できる
![liff](./images/liff.png)


- 特徴
  - 🤖 **LINE Bot統合**: LINE Messaging APIを使った自然な会話
  - 🌐 **LIFF Web App**: Next.js製のモダンなチャットUI
  - 🧠 **Gemini LLM**: Google Geminiを使った高度な会話機能
  - 🔄 **会話同期**: LINE BotとLIFFで会話履歴を完全共有
  - 🔐 **堅牢な認証**: LINE認証とSupabase認証の統合
  - 🐳 **Docker対応**: 任意のクラウド環境にデプロイ可能
  - 📦 **モノレポ構成**: Go Workspaceによる効率的な開発

## 🏗️ アーキテクチャ

```mermaid
graph TD
    %% ユーザー層（同じ高さに配置）
    A[LINEユーザー]
    
    %% プラットフォーム層（同じ高さに配置）
    B --> D[LIFF LIFFアプリ<br/>Next.js + TypeScript]
    A --> B[LINE Platform]

    %% Golangサービス層（共通パッケージで包含）
    B --> E[LINE Botサービス<br/>Golang]
    D <--> F[Backend API<br/>Golang]
    E <--> F

    subgraph Common["共通GOパッケージ"]
        E
        F
    end

    E <--> H[(Supabase<br/>データベース)]
    F <--> H

    E <--> I[Gemini API<br/>LLMサービス]
    F <--> I

    %% スタイリングで役割を明確に
    classDef user fill:#e8f4fd,stroke:#1976d2,stroke-width:2px,color:#000
    classDef service fill:#fff3e0,stroke:#f57c00,stroke-width:2px,color:#000
    classDef liff fill:#d3d3d3,stroke:#404040,stroke-width:2px,color:#000
    classDef platform fill:#e8f5e8,stroke:#388e3c,stroke-width:2px,color:#000
    classDef shared fill:#f3e5f5,stroke:#7b1fa2,stroke-width:2px,color:#000
    classDef external fill:#fce4ec,stroke:#c2185b,stroke-width:2px,color:#000

    class A,C user
    class B platform
    class D liff
    class E,F service
    class Common shared
    class H,I external
```

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
git clone https://github.com/Yongtae723/LineBot-liff-golang-nextjs-template.git
cd LineBot-liff-golang-nextjs-template
```

### 2. Supabase CLIのインストール（初回のみ）

```bash
npm install -g supabase
```

### 3. Supabaseローカル環境の起動

```bash
cd supabase
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

これで`user`と`conversation`テーブルが作成されます。

**確認**: http://localhost:54323 のTable Editorで`user`と`conversation`テーブルが表示されればOK！

### 5. 環境変数の設定
Line DeveloperでBotとLineログインを作成して、以下の情報を取得してください。
- Channel Secret
- Channel Access Token
- LIFF ID


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
cd common && go mod download && cd ..
cd backend && go mod download && cd ..
cd line_bot && go mod download && cd ..
```

### 7. LIFF依存関係のインストール

```bash
cd liff
npm install
```

### 8. サービスの起動

アプリには3つのサービスを起動する必要があります。

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

- **Backend API**: http://localhost:8080/health
- **LINE Bot**: http://localhost:8000/health
- **LIFF App**: http://localhost:3000
- **Supabase Studio**: http://localhost:54323

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
├── common/          # 共通Golangパッケージ
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

### Go (backend, line_bot, common共通)
各folderで以下のコマンドが使えます。

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

## 🚀 デプロイ

このテンプレートはインフラ非依存です。Dockerfileが用意されているので、お好きな環境にデプロイできます：

**Backend & LINE Bot**: Cloud Run、ECS、Railway、Fly.io等  
**LIFF App**: Cloudflare Pages、Vercel、Netlify等

詳細は各ディレクトリのREADMEを参照してください。

## 📚 ドキュメント

- [common/README.md](common/README.md) - 共通パッケージの詳細
- [backend/README.md](backend/README.md) - Backend API仕様
- [line_bot/README.md](line_bot/README.md) - LINE Bot設定
- [liff/README.md](liff/README.md) - LIFF App開発ガイド
- [supabase/README.md](supabase/README.md) - データベーススキーマ

## 🤝 コントリビューション

プルリクエストを歓迎します！
(コード抽出時に、間違いが含まれてる可能性もあり。。。)
大きな変更の場合は、まずIssueを開いて変更内容を議論してください。

## 📝 ライセンス

MIT License - 詳細は[LICENSE](LICENSE)ファイルをご覧ください。

## 💼 開発・運用サポート

このテンプレートを使った開発や、本番環境へのローンチ、運用について詳しく相談したい方へ：

**フリーランスとして以下のサポートを提供しています：**

- 🛠️ **カスタマイズ開発**: テンプレートをベースにした独自機能の実装
- 🚀 **ローンチ支援**: インフラ構築、デプロイ、本番環境セットアップ
- 🔧 **技術コンサルティング**: アーキテクチャ設計、パフォーマンス最適化
- 📊 **運用サポート**: モニタリング設定、障害対応、スケーリング戦略
- 📚 **技術研修**: チーム向けのハンズオン研修、コードレビュー

### お問い合わせ

- x: https://twitter.com/Yoooongtae
- facebook: https://www.facebook.com/yongtaih1
- email: yong723.enjoy.everything@gmail.com

お気軽にご相談ください！💪
(私が連絡に気が付きづらいため、複数媒体に連絡していただけると幸いです。)

## 🙏 謝辞

- [LINE Messaging API](https://developers.line.biz/ja/services/messaging-api/)
- [LIFF (LINE Front-end Framework)](https://developers.line.biz/ja/docs/liff/)
- [Google Gemini](https://ai.google.dev/)
- [Supabase](https://supabase.com/)
- [Next.js](https://nextjs.org/)
- [Gin](https://gin-gonic.com/)

---

Made with ❤️ by [CookForYou](https://cookforyou.jp) team