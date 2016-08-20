# apimock

Mock API Server
if you prepare `actions.json` file, `apimock` use this file as response JSON data.

## Build

### clone from github in GOPATH directory
```shell
mkdir -p ${GOPATH}/src/github.com/yhanada
cd ${GOPATH}/src/github.com/yhanada
git clone git@github.com:yhanada/apimock.git
cd apimock
```

### build

```shell
make build
```

### run

```shell
make run
```

## API response data
URIごとにレスポンスデータを定義することが出来ます。

URIごとに`actions.json`という名前のファイルを置くことが出来、このファイルにレスポンス情報を定義することになります。
このファイルではHTTP Methodごとのエントリがあり、その中にはレスポンス用のStatusコードとBodyの内容が記述されています。

またレスポンスデータについては動的に値を変更することは出来ません。同じURIに対する呼び出しの結果は常に同じとなります。

この中に、HTTP Methodの一つ一つに対してレスポンス用のJSONデータなどを返すことができます。

## exec Options
`apimock-go`ツールオプションは次の3種類があります。

* port
    * 起動するPort番号
* root
    * モック用のAPIデータをJSON形式で置いているファイル
    * JSONファイル名は`actions.json`固定
* check
    * モックサーバは起動せず、`root`で指定したディレクトリ以下の`actions.json`ファイルの内容を表示する

# License
MIT
