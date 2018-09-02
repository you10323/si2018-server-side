# Eureka Summer Internship 2018 API

# 技術スタック

- go 1.10.3
- dep (goの依存関係管理ツール)
- go-swagger
- goose (マイグレーションツール)
- direnv (環境変数を.envrcから読み込む)
- xorm (ORM)

# swagger

dockerで用意しています.
`docker-compose up -d` すると,
localhost:8081 で swagger-editor (エディタ),
localhost:8082 で swagger-ui (APIドキュメント) が開きます。

( http://editor.swagger.io/ でも代用可能です. )

# how to run the app

Golang `1.10.3` がインストールされている前提で、以下の手順に従ってください.

```
# 必要なライブラリの取得

go get -u bitbucket.org/liamstask/goose/cmd/goose
go get -u github.com/golang/dep/cmd/dep
go get -u github.com/go-swagger/go-swagger/cmd/swagger
go get -u github.com/direnv/direnv

# 依存関係のインストール
# (現状, initコマンドは単にdep ensureをラップしてあるだけです)
make init

# 環境変数を.envrc (direnv) で管理しています.
cp .envrc.sample .envrc
direnv allow

# ビルド
make build (swaggerのymlからgoファイルを生成

make init (生成されたgoファイルの依存関係取り込み

# サーバーを立ち上げる
make run
```

# dummy data

misc/dummy/ 下にダミーデータ生成のスクリプトを置いてます。以下makeコマンドでDBリセット & ダミー生成を行います.

```
make setup-db
```

# migration with goose

マイグレーションツールのgooseを使用しています。

```
# ./db/migrations/20180809183923_createUser.sql が作成される
goose create createHoge sql

# up
goose up

# down
goose down

# redo
goose redo
```