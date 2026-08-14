# go-http-error-return-debug-lab

Goの`net/http`で、`http.Error`を呼んだあとに`return`を忘れると、エラー本文の後ろへ成功時のJSONが混在する不具合を再現する最小プロジェクトです。

## この教材で確認する契約

`POST /orders`に`productId`がない場合、HTTPステータスだけでなく、レスポンス本文もエラーメッセージだけで終わる必要があります。

| 入力 | 期待するHTTPステータス | 期待する本文 |
| --- | --- | --- |
| `{"productId":""}` | `400 Bad Request` | `productId は必須です\n` |
| `{"productId":"product-123"}` | `201 Created` | `{"id":"order-001"}\n` |

修正前は、空の`productId`に対してステータスが`400`であっても、本文が`productId は必須です\n{"id":"order-001"}\n`となります。HTTPステータスだけを確認するテストでは、この不整合を見逃します。

## 前提条件

Go 1.22以降が必要です。外部依存はありません。

## 実行方法

修正後のブランチでは、次のコマンドで全テストが通ります。

```bash
go test -v ./...
```

再現コミットでは、意図した失敗を確認できます。

```bash
git checkout a646927
go test ./...
```

修正後へ戻すには、既定ブランチへ切り替えてください。

```bash
git switch main
go test -v ./...
```

## 構成

| パス | 役割 |
| --- | --- |
| `handler.go` | 注文作成を模した最小HTTPハンドラー |
| `handler_test.go` | エラー本文と成功応答を検証する回帰テスト |
| `docs/debugging-record.md` | 観測結果、切り分け、原因、修正を残す記録 |

## 関連記事

Qiita記事の下書きは、Qiitaコンテンツリポジトリの`private/01_執筆中/06_テスト・デバッグ/デバッグ/`にあります。
