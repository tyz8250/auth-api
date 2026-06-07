# Issue #3 signup の Service 層と password_hash

## 目的

`POST /signup` の処理で、password をそのまま扱い続けず、Service 層で bcrypt hash に変換する流れを理解する。

## 今回やったこと

- `service.AuthService` を追加した
- `AuthService.Signup(email, password)` を追加した
- `Signup` の中で `bcrypt.GenerateFromPassword` を使って password を hash 化した
- `repository.MemoryUserRepository` を追加した
- `AuthService` が repository interface に依存する形にした
- hash 化した password を `PasswordHash` として repository に保存する流れにした
- handler は `AuthService.Signup` を呼び、response には `UserResponse` だけを返すようにした
- service のテストで、hash が平文 password ではなく、元 password と照合できることを確認した

## 処理の流れ

1. handler が request JSON から `email` と `password` を読む
2. handler が空文字チェックをする
3. handler が service の `Signup` を呼ぶ
4. service が password を bcrypt hash にする
5. service が `model.User.PasswordHash` に hash を入れる
6. service が repository に user を保存する
7. handler が `model.UserResponse` に変換して JSON で返す

## handler と service の分担

handler は HTTP request / response を扱う。

- JSON decode
- status code
- JSON response

service は signup の処理の中身を扱う。

- password hash
- signup のユースケース
- repository を使った保存

repository は保存方法を扱う。

- 今回は DB ではなくメモリ上に保存する
- ID と作成日時、更新日時をセットする
- 将来的に DB を入れるときは、この具体実装を置き換える

password hash の詳細を handler に置かないことで、HTTP の処理と認証ロジックを分けて読めるようになる。

## repository interface にした理由

`AuthService` は `repository.MemoryUserRepository` そのものではなく、`Create(user model.User)` を持つ interface に依存する。

これにより、service のテストでは本物の保存先ではなく、テスト用の repository を渡せる。

今回のテストでは、service が repository に渡した user の `PasswordHash` を確認している。

## テストで確認したこと

bcrypt の hash は毎回違う文字列になる。

そのため、hash 文字列を固定値で比較するのではなく、次を確認した。

- `PasswordHash` が空ではない
- `PasswordHash` が元の password そのものではない
- `bcrypt.CompareHashAndPassword` で元の password と照合できる
- repository に渡された user も平文 password ではなく hash を持っている
- memory repository に user が保存される

## 次に進むなら

次は email 重複チェックに進むとよい。

そのためには repository に `FindByEmail` のような処理を追加し、signup 前に同じ email の user がいないか確認する流れを作る。
