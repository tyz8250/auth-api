# Issue #2 POST /signup の流れ

## 目的

`POST /signup` にリクエストが来たときに、JSON を読み取り、signup 用の handler で処理できる入口を作る。

この Issue では、まだ DB 保存や password hash までは深追いしない。
まずは HTTP handler の流れを理解する。

## 処理の流れ

1. `main` 関数で `mux := http.NewServeMux()` を作る
2. `mux.HandleFunc("POST /signup", signupHandler)` でルートを登録する
3. `http.ListenAndServe(":8080", mux)` でサーバを起動する
4. `POST /signup` にリクエストが来る
5. `mux` が `signupHandler` を呼ぶ
6. `signupHandler` が `r.Body` の JSON を読む
7. JSON を `signupRequest` 構造体に変換する
8. `req.Email` と `req.Password` を使って処理する
9. 最後に JSON レスポンスを返す

## 今回やること

- `POST /signup` のルートを追加する
- `signupRequest` 構造体を作る
- `json.NewDecoder(r.Body).Decode(&req)` で JSON を読む
- invalid JSON のときは JSON エラーを返す
- 成功時は仮の JSON レスポンスを返す

## 今回まだやらないこと

- DB に user を保存する
- password を hash 化する
- login を作る
- token / session を扱う
- repository / service を本格的に作る
