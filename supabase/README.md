# Supabase ローカル開発環境

このディレクトリには、Supabaseのローカル開発環境に関するファイルが含まれています。

## 📋 必要なもの

- Docker Desktop
- Supabase CLI

## 🚀 セットアップ

### 1. Supabase CLIのインストール

```bash
npm install -g supabase
```

または

```bash
brew install supabase/tap/supabase
```

### 2. 初期化（初回のみ）

```bash
cd supabase
supabase init
```

### 3. ローカル環境の起動

```bash
supabase start
```

初回起動時は、Dockerイメージのダウンロードに時間がかかる場合があります（数分～10分程度）。

起動が完了すると、以下のような接続情報が表示されます：

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

### 4. Supabase Studioにアクセス

ブラウザで http://localhost:54323 を開くと、Supabase Studioが表示されます。

## 🗄️ データベーススキーマの作成

### 推奨方法: SQLファイルで一括セットアップ

このプロジェクトでは、**用意されたSQLファイルを実行してスキーマを作成し、その後マイグレーションファイルを生成する**方法を推奨しています。

#### Step 1: SQLファイルを実行

Supabase Studio (http://localhost:54323) の**SQL Editor**で以下を実行します：

**Option A: ファイルの内容をコピペ（推奨）**

```bash
# ファイルの内容を確認
cat supabase/initial_setup.sql
```

上記の内容を**SQL Editor**にコピーして実行。

**Option B: 以下のSQLを直接コピペ**

<details>
<summary>SQLスクリプト全文（クリックして展開）</summary>

```sql
-- UUID拡張を有効化
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- users テーブル作成
CREATE TABLE IF NOT EXISTS public.users (
  id UUID REFERENCES auth.users(id) ON DELETE CASCADE PRIMARY KEY,
  line_id TEXT UNIQUE NOT NULL,
  name TEXT,
  created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
  updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

-- conversations テーブル作成
CREATE TABLE IF NOT EXISTS public.conversations (
  id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
  user_id TEXT NOT NULL,
  role TEXT NOT NULL,
  content TEXT NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

-- インデックス
CREATE INDEX IF NOT EXISTS idx_users_line_id ON public.users(line_id);
CREATE INDEX IF NOT EXISTS idx_conversations_user_id ON public.conversations(user_id);
CREATE INDEX IF NOT EXISTS idx_conversations_created_at ON public.conversations(created_at DESC);

-- 制約
ALTER TABLE public.conversations 
ADD CONSTRAINT conversations_role_check 
CHECK (role IN ('user', 'assistant'));

-- RLS有効化
ALTER TABLE public.users ENABLE ROW LEVEL SECURITY;
ALTER TABLE public.conversations ENABLE ROW LEVEL SECURITY;

-- RLSポリシー (users)
CREATE POLICY "Users can view own data"
  ON public.users FOR SELECT
  USING (auth.uid() = id);

CREATE POLICY "Users can update own data"
  ON public.users FOR UPDATE
  USING (auth.uid() = id);

-- RLSポリシー (conversations)
CREATE POLICY "Users can view own conversations"
  ON public.conversations FOR SELECT
  USING (
    user_id = (SELECT line_id FROM public.users WHERE id = auth.uid())
  );

CREATE POLICY "Users can insert own conversations"
  ON public.conversations FOR INSERT
  WITH CHECK (
    user_id = (SELECT line_id FROM public.users WHERE id = auth.uid())
  );

-- トリガー
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_users_updated_at
  BEFORE UPDATE ON public.users
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at_column();
```

</details>

#### Step 2: マイグレーションファイルを生成

SQLでスキーマを作成したら、差分を検出してマイグレーションファイルを生成します：

```bash
# 現在のスキーマとの差分をマイグレーションファイルとして生成
supabase db diff -f initial_schema

# または、より詳細な名前で
supabase db diff -f create_users_and_conversations_tables
```

これで`supabase/migrations/`ディレクトリに新しいマイグレーションファイルが作成されます：

```
supabase/migrations/20250106123456_initial_schema.sql
```

#### Step 3: マイグレーションファイルを確認・編集

生成されたマイグレーションファイルを確認し、必要に応じて編集します：

```bash
# マイグレーションファイルを確認
cat supabase/migrations/20250106123456_initial_schema.sql
```

#### Step 4: マイグレーションをリセットして再適用（テスト）

```bash
# データベースをリセット（マイグレーションを最初から適用）
supabase db reset
```

### 代替方法: GUIでテーブルを作成

もしGUIで手動作成したい場合：

1. **Table Editor** → **New table** で `users` と `conversations` を作成
2. カラムを手動で追加
3. 上記の「インデックスとポリシー」のSQLを実行
4. `supabase db diff -f initial_schema` でマイグレーション生成

**注意**: GUIでの作成は時間がかかるため、SQLファイルの使用を推奨します。

## 🔑 環境変数の取得

Supabase起動時に表示される接続情報を、各サービスの`.env`ファイルに設定してください：

### backend/.env と line_bot/.env
- `SUPABASE_URL`: API URL
- `SUPABASE_KEY`: **service_role key** (⚠️ anon keyではなく)
- `SUPABASE_JWT_SECRET`: JWT secret

### liff/.env.local
- `NEXT_PUBLIC_SUPABASE_URL`: API URL
- `NEXT_PUBLIC_SUPABASE_ANON_KEY`: **anon key**

## 🛠️ よく使うコマンド

### 起動・停止

```bash
# 起動
supabase start

# 停止
supabase stop

# 再起動
supabase stop && supabase start

# ステータス確認
supabase status
```

### データベース操作

```bash
# データベースをリセット（全データ削除 + マイグレーション再適用）
supabase db reset

# 空のマイグレーションファイルを作成
supabase migration new <migration_name>

# GUIで作成したスキーマの差分からマイグレーションファイルを生成
supabase db diff -f <migration_name>

# 例: GUIでテーブルを作成した後
supabase db diff -f create_users_table

# マイグレーション履歴を確認
supabase migration list

# マイグレーションを適用（リモート環境へ）
supabase db push

# 現在のスキーマをダンプ
supabase db dump -f schema.sql

# 特定のテーブルだけダンプ
supabase db dump -f users_table.sql --data-only --schema public -t users
```

### 型定義生成

```bash
# TypeScript型定義を生成
supabase gen types typescript --local > ../liff/src/types/supabase.ts
```

### リモート環境との接続

```bash
# リモートプロジェクトにリンク
supabase link --project-ref <your-project-ref>

# リモートからマイグレーションを取得
supabase db pull

# ローカルマイグレーションをリモートに適用
supabase db push
```

## 🐛 トラブルシューティング

### Dockerが起動しない

```bash
# Docker Desktopが起動しているか確認
docker ps

# Dockerを再起動
# Docker Desktop アプリを再起動
```

### ポートが既に使用されている

```bash
# 使用中のポートを確認
lsof -i :54321  # Supabase API
lsof -i :54322  # PostgreSQL
lsof -i :54323  # Supabase Studio
lsof -i :54324  # Inbucket

# プロセスを終了
kill -9 <PID>

# または、Supabaseのポートを変更
# supabase/config.toml でポート番号を変更可能
```

### データベースがリセットされない

```bash
# 強制リセット
supabase db reset --db-url postgresql://postgres:postgres@localhost:54322/postgres

# 全て停止してから再起動
supabase stop
docker system prune -a
supabase start
```

### マイグレーションエラー

```bash
# マイグレーション履歴を確認
supabase migration list

# 特定のマイグレーションを修復
supabase migration repair <version>
```

## 📊 データベースアクセス

### psqlでの直接接続

```bash
psql postgresql://postgres:postgres@localhost:54322/postgres
```

### よく使うSQLコマンド

```sql
-- テーブル一覧
\dt

-- テーブル構造確認
\d public.users
\d public.conversations

-- データ確認
SELECT * FROM public.users;
SELECT * FROM public.conversations LIMIT 10;

-- データ削除
TRUNCATE public.conversations;
DELETE FROM public.users;
```

## 🔐 セキュリティ

ローカル開発環境のデフォルトパスワードとキーは以下の通りです：

- **Postgresパスワード**: `postgres`
- **JWT Secret**: `super-secret-jwt-token-with-at-least-32-characters-long`

**⚠️ 本番環境では絶対にこれらの値を使用しないでください！**

## 📝 開発ワークフロー（推奨）

### スキーマ変更の手順

1. **Supabase StudioのGUIでスキーマを変更**
   - テーブル作成、カラム追加、インデックス追加など

2. **差分からマイグレーションファイルを生成**
   ```bash
   supabase db diff -f <変更内容の説明>
   # 例: supabase db diff -f add_user_avatar_column
   ```

3. **生成されたマイグレーションファイルを確認**
   ```bash
   cat supabase/migrations/<timestamp>_<変更内容の説明>.sql
   ```

4. **テストのためデータベースをリセット**
   ```bash
   supabase db reset
   # 全マイグレーションが最初から適用され、正しく動作するか確認
   ```

5. **問題なければGitにコミット**
   ```bash
   git add supabase/migrations/
   git commit -m "Add migration: <変更内容の説明>"
   ```

### よくあるパターン

```bash
# パターン1: 新しいテーブルを追加
# → GUIでテーブル作成 → マイグレーション生成
supabase db diff -f create_notifications_table

# パターン2: 既存テーブルにカラム追加
# → GUIでカラム追加 → マイグレーション生成
supabase db diff -f add_user_profile_fields

# パターン3: インデックスやポリシーの追加
# → SQL Editorで実行 → マイグレーション生成
supabase db diff -f add_performance_indexes

# パターン4: 複数の変更をまとめて
# → 複数の変更を実施 → マイグレーション生成
supabase db diff -f refactor_user_schema
```

## 📝 メモ

- ローカル環境のデータは`supabase stop`後も保持されます
- `supabase db reset`でデータを完全にリセットできます
- マイグレーションファイルは`migrations/`ディレクトリに保存されます
- `seed.sql`にテストデータを記述できます
- **GUIで作成→`db diff`で生成**が最も効率的なワークフローです
- マイグレーションファイルは必ずGitで管理してください

## 🔗 参考リンク

- [Supabase CLI ドキュメント](https://supabase.com/docs/guides/cli)
- [Supabase ローカル開発](https://supabase.com/docs/guides/cli/local-development)
- [Supabase マイグレーション](https://supabase.com/docs/guides/cli/managing-environments)
- [Supabase db diff コマンド](https://supabase.com/docs/reference/cli/supabase-db-diff)

sup