# デバッグ記録：`http.Error`後に成功JSONが混在する

## 目的

空の`productId`を受け取った注文作成ハンドラーが、`400 Bad Request`とエラーメッセージだけを返す契約を確認します。修正前はHTTPステータスが`400`であるにもかかわらず、エラー本文の後ろに成功時のJSONが追記されます。ここでは、ステータスコードだけで成功・失敗を判断せず、HTTPレスポンス本文まで観測して原因を特定します。

## 最初に観測した事実

| 観測項目 | 実際の結果 | 根拠 |
| --- | --- | --- |
| HTTPステータス | `400 Bad Request` | `TestCreateOrderHandler_商品IDが空ならエラー本文だけを返す` |
| HTTP本文 | `productId は必須です\n{"id":"order-001"}\n` | 同テストの失敗出力 |
| 成功時の応答 | `201 Created`、`{"id":"order-001"}\n` | `TestCreateOrderHandler_商品IDがあれば注文IDを返す` |

```text
--- FAIL: TestCreateOrderHandler_商品IDが空ならエラー本文だけを返す (0.00s)
    handler_test.go:30: body = "productId は必須です\n{\"id\":\"order-001\"}\n", want "productId は必須です\n"
FAIL
```

## 仮説と切り分け

| 仮説 | 確認方法 | 結果 |
| --- | --- | --- |
| 入力JSONのデコードが失敗し、想定外の処理へ入った | `{"productId":""}`を明示的に渡し、テストが期待する400を確認する | 否定。デコードは成功し、必須チェックへ到達している |
| `http.Error`がハンドラー処理を終了する | `http.Error`の直後にある成功レスポンスの書き込みと、実際の本文を確認する | 否定。`http.Error`後も後続コードが実行される |
| 成功レスポンスが後続で書き込まれている | `handler.go`で`http.Error`の直後に`return`がないことと、本文末尾のJSONを照合する | 採用。成功JSONが同一レスポンスへ追記されている |

## 原因

`http.Error`はエラー本文とHTTPステータスを書き込みますが、ハンドラー関数を終了しません。公式ドキュメントも、呼び出し側が以降に`ResponseWriter`へ書き込まないようにする必要があると説明しています。[1]

修正前のハンドラーは空の`productId`で`http.Error`を呼んだ後、そのまま成功時の`201 Created`とJSONを書き込んでいました。最初に書かれたステータスは`400`のままでも、本文にはエラーと成功の内容が混在します。このため、ステータスコードだけを確認するテストでは不具合を検出できません。

## 修正

```go
if request.ProductID == "" {
    http.Error(w, "productId は必須です", http.StatusBadRequest)
    return
}
```

`http.Error`の直後に`return`を置き、エラー応答を書いた経路では成功レスポンスの書き込みに到達しないようにしました。変更は制御フローだけであり、成功ケースの仕様は変えていません。

## 再発防止テスト

`TestCreateOrderHandler_商品IDが空ならエラー本文だけを返す`は、空の`productId`に対してステータス`400`だけでなく、本文が`productId は必須です\n`と完全一致することを確認します。これにより、エラー本文へ成功JSONが追記される回帰を検出できます。

さらに`TestCreateOrderHandler_商品IDがあれば注文IDを返す`で、正常入力が`201 Created`、`application/json`、注文IDのJSONを返すことを確認します。エラー経路を直すために正常経路を壊していないことも同時に保証します。

## 再現手順

```bash
git checkout a646927
go test ./...
# TestCreateOrderHandler_商品IDが空ならエラー本文だけを返す が失敗する

git switch main
go test -v ./...
# 全テストが成功する
```

## 参考資料

[1] [net/http - Error](https://pkg.go.dev/net/http#Error)
