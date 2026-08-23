# go-api-template

Go と MySQL で REST API を開発するためのテンプレートです。Echo を HTTP フレームワークに、Bun をデータベースアクセスに使用し、OpenAPI 定義からサーバー向けコードを生成します。

サンプルとして、顧客・書籍・顧客のお気に入り書籍を扱う API を実装しています。

## 主な構成

- Go 1.26
- Echo
- MySQL 8
- Bun
- OpenAPI 3.0 + oapi-codegen
- Docker Compose
- mise（ツール・タスク管理）

## API

API 仕様は [docs/openapi/spec.yaml](docs/openapi/spec.yaml) にあります。主なエンドポイントは以下のとおりです。

| メソッド      | パス                                        | 内容                       |
| ------------- | ------------------------------------------- | -------------------------- |
| `GET`         | `/health`                                   | ヘルスチェック             |
| `GET`, `POST` | `/v1/customers`                             | 顧客の一覧取得・作成       |
| `GET`, `PUT`  | `/v1/customers/{customerId}`                | 顧客の取得・更新           |
| `POST`, `PUT` | `/v1/customers/{customerId}/favorite-books` | お気に入り書籍の登録・解除 |
| `GET`         | `/v1/books`                                 | 書籍の一覧取得             |

## ローカル開発

前提条件は Docker と Docker Compose です。

```sh
cp .env.example .env
docker compose up -d --build
```

コンテナを起動した後、マイグレーションを適用します。

```sh
mise run server migrate:up
```

必要に応じてサンプルデータを投入します。

```sh
mise run server seeder
```

API サーバーを起動します。

```sh
mise run server default
```

起動後は `http://localhost:8080/health` にアクセスできます。ポートなどを変更する場合は `.env` を編集してください。

## 開発コマンド

ルートの `mise` タスクは、サーバーコンテナ内で `server/mise.toml` のタスクを実行します。

```sh
# 単体テスト
mise run server test

# 結合テスト（db-test コンテナを使用）
mise run server test:integration

# フォーマット・lint・テスト
mise run server check

# OpenAPI 定義からコードを生成
mise run server gen:openapi

# データベーススキーマから Bun のモデルを生成
mise run server gen:db-schema
```

マイグレーション作成・適用にも `mise` タスクを利用できます。

```sh
mise run server migrate:create create_example_table
mise run server migrate:up
mise run server migrate:down 1
```

利用可能なタスクは `server/mise.toml` を参照してください。

## ディレクトリ構成

```text
.
├── compose.yaml                 # 開発用コンテナ（API / MySQL / テスト用 MySQL）
├── docs/openapi/spec.yaml       # OpenAPI 定義
└── server/
    ├── internal/cmd/             # API サーバー、シーダーのエントリポイント
    ├── internal/restapi/         # コントローラー、ミドルウェア、DI
    ├── internal/domain/          # ドメインモデル、ユースケース、リポジトリ
    ├── internal/orchestration/   # 複数ドメインにまたがる処理
    └── internal/database/        # DB 接続、マイグレーション、シード
```
