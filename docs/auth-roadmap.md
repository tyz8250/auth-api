# 認証 API ロードマップ

## Issue 1 の目的

このプロジェクトでは、Go で小さな認証 API を作りながら、HTTP request / response、レイヤー分割、DB、パスワードハッシュ、認証状態の扱いを段階的に理解する。

最初に全体像を整理しておくことで、次に実装する issue が「認証 API のどの部分なのか」を説明できるようにする。

## 最終ゴール

- `POST /signup` でユーザー登録できる
- `POST /login` でログインできる
- `GET /me` で認証済みユーザー情報を取得できる
- password は平文保存せず、DB には `password_hash` として保存する
- token または session を使って認証状態を扱う
- `handler` / `service` / `repository` / `model` の責務を分ける
- 実装した機能に対して、段階的にテストを書く
- 詰まったことや理解したことを `docs` に残す

## エンドポイントの役割

### `POST /signup`

新しいユーザーを登録する。

主な流れ:

1. request JSON から `email` と `password` を受け取る
2. 入力値を確認する
3. password をハッシュ化する
4. ユーザーを DB に保存する
5. password や `password_hash` を含めずに結果を返す

学ぶこと:

- JSON request body の読み取り
- HTTP status code
- password hash
- email 重複チェック
- handler / service / repository の分離

### `POST /login`

登録済みユーザーがログインする。

主な流れ:

1. request JSON から `email` と `password` を受け取る
2. email でユーザーを探す
3. 入力された password と保存済みの `password_hash` を比較する
4. 正しければ token または session を発行する
5. password や `password_hash` を含めずに結果を返す

学ぶこと:

- password の照合
- 認証失敗時のエラー
- token または session の発行
- login 成功と失敗のテスト

### `GET /me`

認証済みユーザーの情報を取得する。

主な流れ:

1. request から token または session 情報を取り出す
2. 認証情報が正しいか確認する
3. 認証済みユーザーを特定する
4. password や `password_hash` を含めずにユーザー情報を返す

学ぶこと:

- 認証 middleware
- token なし / token 不正 / token 正常の分岐
- 認証済みユーザー情報の扱い
- secret をレスポンスに含めない設計

## レイヤーの責務

### `handler`

HTTP request / response を扱う層。

担当すること:

- HTTP method の確認
- JSON request body の decode
- service の呼び出し
- HTTP status code の決定
- JSON response の返却

担当しないこと:

- DB への直接アクセス
- password hash の詳細
- token / session の細かい生成処理

### `service`

ユースケースや認証ロジックを扱う層。

担当すること:

- signup / login / me の処理の流れを組み立てる
- password hash や password 照合を行う
- repository interface を使ってデータを取得、保存する
- どのエラーを返すか判断する

### `repository`

DB や永続化を扱う層。

担当すること:

- ユーザーを保存する
- email でユーザーを探す
- id でユーザーを探す
- DB 固有の処理を閉じ込める

service は repository の具体実装ではなく、できるだけ interface に依存する。

### `model`

データ構造を表す層。

担当すること:

- `User` などの構造体を定義する
- DB に保存する値と、API で返す値を区別しやすくする

注意すること:

- API response に password や `password_hash` を含めない
- secret や token を不用意に model に混ぜない

## token と session の入口メモ

認証状態を扱う代表的な方法には token と session がある。

### token

ログイン成功時に token を発行し、クライアントが次回以降の request に token を付けて送る方式。

例:

- `Authorization: Bearer <token>`

特徴:

- API と相性がよい
- サーバー側で session 保存をしない構成も作れる
- token の期限、署名、失効方法を考える必要がある

### session

ログイン成功時に session id を発行し、サーバー側で session 情報を保存する方式。

例:

- cookie に session id を入れる

特徴:

- Web アプリでよく使われる
- サーバー側で session を管理しやすい
- session 保存先や cookie の安全設定を考える必要がある

このプロジェクトでは、token または session のどちらを使うかを後続 issue で決める。導入前に方針を整理してから実装する。

## 実装ロードマップ

1. 認証 API の全体像と学習ロードマップを `docs` に整理する
2. `POST /signup` の request JSON を受け取れるようにする
3. signup の入力値チェックと JSON エラー response を整理する
4. `model.User` を定義する
5. repository interface を定義し、まずはメモリ上でユーザーを保存できるようにする
6. service に signup のユースケースを移す
7. password をハッシュ化して保存する
8. signup のテストを追加する
9. `POST /login` の request JSON を受け取れるようにする
10. password 照合を行い、login 成功 / 失敗を分ける
11. token または session の方針を決める
12. login 成功時に認証情報を返す
13. 認証 middleware を追加する
14. `GET /me` で認証済みユーザー情報を返す
15. DB を導入し、repository の実装を DB に置き換える
16. endpoint ごとのテストを増やす
17. 学んだことや設計判断を `docs` に追加する

## 次に実装する issue

次は `POST /signup` の request JSON を受け取れるようにする。

目的:

- HTTP handler で JSON request body を読む流れを理解する

やること:

- signup handler の request struct を作る
- `email` と `password` を JSON から decode する
- invalid JSON のときに JSON エラーを返す

確認方法:

- `go test ./...`
- `curl` で `/signup` に JSON を送る

