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

必要なライブラリの取得を go get で行います.

```
go get -u bitbucket.org/liamstask/goose/cmd/goose
go get -u github.com/golang/dep/cmd/dep
go get -u github.com/go-swagger/go-swagger/cmd/swagger
go get -u github.com/direnv/direnv
```

依存関係のインストール

```
make init
```

環境変数を.envrc (direnv) で管理しています.
.envrcの内容は各自割り振られたDB接続情報で書き換えてください.

```
# cp して .envrc の内容を書く
cp .envrc.sample .envrc

# 以下コマンドで、.envrcの置かれたディレクトリ配下で環境変数が有効になります.
direnv allow
```

make generate すると、si2018.ymlのswaggerの定義から、goのファイルが生成されます.
その後, 生成されたgoのファイルの依存関係取り込みのため make initを打ってください.

```
make generate
make init
```

サーバーを立ち上げる

```
make run
```
