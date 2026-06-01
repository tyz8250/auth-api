
# bcryptでパスワードを安全に保存する

## つまり

認証APIでは、ユーザーのパスワードをそのまま保存してはいけない。

ユーザー登録時には、平文パスワードを `bcrypt` でハッシュ化して保存する。  
ログイン時には、入力されたパスワードと保存済みハッシュを `CompareHashAndPassword` で照合する。

大事なのは、ログイン時にもう一度 `GenerateFromPassword` でハッシュを作って、文字列同士を比較してはいけないということ。

---

## なぜこれを学ぶのか

認証APIを作るなら、JWTより前に「安全なパスワード保存」を理解する必要がある。

JWTは「ログイン後に本人であることをどう証明するか」の仕組み。  
一方で bcrypt は「そもそもパスワードをどう安全に保存するか」の仕組み。

そのため、認証APIではまず bcrypt を理解することが土台になる。

---

## 使用するパッケージ

```go
import "golang.org/x/crypto/bcrypt"
```

bcrypt の中心になる関数は主に2つ。

```go
bcrypt.GenerateFromPassword()
```

```go
bcrypt.CompareHashAndPassword()
```

---

## GenerateFromPasswordとは

公式ドキュメントでは、次のように定義されている。

```go
func GenerateFromPassword(password []byte, cost int) ([]byte, error)
```

これは、

```txt
第1引数: password []byte
第2引数: cost int
戻り値: []byte, error
```

という意味。

実際のコードでは、パスワードを文字列として持っていることが多い。

```go
password := "secret123"
```

しかし `GenerateFromPassword` は `[]byte` 型を求めている。

そのため、次のように `string` から `[]byte` に変換して渡す。

```go
hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
```

ここでの疑問は、

> 公式ドキュメントでは `password []byte` と書かれているのに、コードでは `[]byte(password)` と書いている。これは違うのか？

というもの。

結論として、問題ない。

公式ドキュメントの `password []byte` は「この関数には `[]byte` 型を渡してください」という意味。  
自分のコードの `[]byte(password)` は、`string` 型の `password` を `[]byte` 型に変換しているだけ。

つまり、

```go
password := "secret123"  // string
[]byte(password)         // []byte
```

という関係。

---

## Costとは

`cost` は、bcrypt のハッシュ化処理をどれくらい重くするかを決める値。

```go
bcrypt.DefaultCost
```

を使うと、bcrypt パッケージが用意している標準的なコストでハッシュ化できる。

イメージは次のとおり。

```txt
Cost 小さい → 処理が速い → 攻撃者も試しやすい → 弱め
Cost 大きい → 処理が遅い → 攻撃者が試しにくい → 強め
```

bcrypt は、あえて計算に時間がかかるように作られている。  
理由は、もしパスワードハッシュが漏れた場合でも、攻撃者が大量のパスワード候補を高速に試しにくくするため。

---

## Costを小さくするとハッシュの数が減るのか？

最初の疑問は、

> Costを小さくすると、ハッシュの数が減るという認識でよいのか？

というもの。

これは少し違う。

Costを小さくすると、ハッシュの数が減るのではなく、ハッシュを作るための計算量が減る。

つまり、

```txt
Costを小さくする
= ハッシュ化の計算が軽くなる
= 処理が速くなる
= そのぶん総当たり攻撃に弱くなりやすい
```

という理解が近い。

bcrypt の cost は、ざっくり言うと指数的に重くなる。

```txt
Cost 10 → ある重さ
Cost 11 → だいたい約2倍重い
Cost 12 → さらに約2倍重い
```

という感覚。

最初は `bcrypt.DefaultCost` を使えばよい。  
学習段階で無理に cost を変更する必要はない。

---

## ユーザー登録時の使い方

ユーザー登録時は、受け取った平文パスワードを bcrypt でハッシュ化する。

```go
package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := "secret123"

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	fmt.Println("password:", password)
	fmt.Println("hash:", string(hash))
}
```

このとき、保存するのは平文パスワードではない。

保存してはいけないもの。

```txt
secret123
```

保存するもの。

```txt
$2a$10$...
```

このような bcrypt のハッシュ文字列をDBに保存する。

---

## CompareHashAndPasswordとは

ログイン時には、入力されたパスワードが保存済みハッシュと一致するかを確認する。

使う関数はこれ。

```go
err := bcrypt.CompareHashAndPassword(savedHash, []byte(password))
```

この関数は、

```txt
保存済みハッシュ
入力された平文パスワード
```

を受け取り、正しいパスワードかどうかを確認する。

一致した場合は、

```go
err == nil
```

になる。

一致しなかった場合は、

```go
err != nil
```

になる。

そのため、次のように書ける。

```go
err = bcrypt.CompareHashAndPassword(hash, []byte(password))
if err != nil {
	fmt.Println("password does not match")
	return
}

fmt.Println("password matches")
```

このコードの意味は、

```txt
ハッシュとパスワードを比較する
↓
一致しなかったら err が入る
↓
err != nil なので "password does not match" と表示して終了
↓
一致していたら err は nil
↓
if に入らず次へ進む
↓
"password matches" と表示する
```

ということ。

---

## ログイン時にやってはいけないこと

次のように、ログイン時にもう一度ハッシュを作って、保存済みハッシュと文字列比較してはいけない。

```go
newHash, _ := bcrypt.GenerateFromPassword([]byte(inputPassword), bcrypt.DefaultCost)

if string(newHash) == string(savedHash) {
	// OK?
}
```

一見すると、

```txt
入力されたパスワードから newHash を作る
DBに保存している savedHash と比較する
同じならログイン成功
```

のように見える。

しかし、これは bcrypt ではうまくいかない。

理由は、bcrypt は同じパスワードでも毎回違うハッシュを作るから。

たとえば、同じ `"secret123"` を2回ハッシュ化しても、次のように違う文字列になる。

```txt
1回目:
$2a$10$abc...xxxx

2回目:
$2a$10$def...yyyy
```

どちらも元のパスワードは `"secret123"` だが、ハッシュ文字列は違う。

---

## なぜ毎回違うハッシュになるのか

bcrypt は内部で `salt` というランダムな値を使っている。

salt とは、パスワードに混ぜるランダムな値のこと。

イメージとしては、

```txt
password = "secret123"
salt     = ランダムな値
```

これを組み合わせてハッシュを作る。

```txt
hash = bcrypt("secret123" + salt)
```

そのため、同じパスワードでも salt が違えば、できあがるハッシュも変わる。

```txt
"secret123" + saltA → hashA
"secret123" + saltB → hashB
```

だから、次のような文字列比較は基本的に一致しない。

```go
if string(newHash) == string(savedHash) {
	// ほぼ一致しない
}
```

---

## ではどうやって比較するのか

正しくは、`CompareHashAndPassword` を使う。

```go
err := bcrypt.CompareHashAndPassword([]byte(savedHash), []byte(inputPassword))
if err != nil {
	fmt.Println("login failed")
	return
}

fmt.Println("login ok")
```

bcrypt のハッシュ文字列には、単なるハッシュ結果だけでなく、次のような情報も含まれている。

```txt
bcryptのバージョン
cost
salt
ハッシュ結果
```

そのため、`CompareHashAndPassword` は保存済みハッシュの中から必要な情報を読み取り、入力されたパスワードが正しいかどうかを検証してくれる。

つまり、ログイン時にやるべきことは、

```txt
新しいハッシュを作って文字列比較する
```

ではなく、

```txt
保存済みハッシュに対して、入力されたパスワードが正しいか確認する
```

ということ。

---

## signup と login で役割を分ける

### signup時

ユーザー登録時は、平文パスワードをハッシュ化して保存する。

```go
password := "secret123"

hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
if err != nil {
	panic(err)
}

savedHash := string(hash)
```

この `savedHash` をDBなどに保存する。

---

### login時

ログイン時は、入力されたパスワードを新しくハッシュ化して文字列比較しない。

悪い例。

```go
newHash, _ := bcrypt.GenerateFromPassword([]byte(inputPassword), bcrypt.DefaultCost)

if string(newHash) == savedHash {
	fmt.Println("login ok")
}
```

正しい例。

```go
err := bcrypt.CompareHashAndPassword([]byte(savedHash), []byte(inputPassword))
if err != nil {
	fmt.Println("login failed")
	return
}

fmt.Println("login ok")
```

---

## たとえ

`GenerateFromPassword` は、毎回違う鍵付きの箱を作るようなもの。

同じパスワードでも、毎回違う見た目の箱になる。

```txt
同じパスワードでも
箱A、箱B、箱C が毎回違う見た目になる
```

だから、

```txt
箱A == 箱B ?
```

と見た目で比較しても一致しない。

一方で、`CompareHashAndPassword` は、

```txt
この箱は、このパスワードで開けられるか？
```

を確認してくれる関数。

そのため bcrypt では、ハッシュ文字列同士を比較するのではなく、保存済みハッシュに対して入力パスワードが正しいかを確認する。

---

## 最小理解

今の段階では、次の理解で十分。

```txt
GenerateFromPassword
= ユーザー登録時に、保存用ハッシュを作る

CompareHashAndPassword
= ログイン時に、保存済みハッシュと入力パスワードを照合する

Cost
= ハッシュ化処理の重さ。小さいと速いが弱め、大きいと遅いが強め

[]byte(password)
= string型のパスワードを、bcryptが求める[]byte型に変換している

ログイン時にやってはいけないこと
= GenerateFromPasswordで新しいハッシュを作り、保存済みハッシュと文字列比較すること
```

---

## 今日の実験コード

`auth-api` の中に practice として作る。

```txt
auth-api/
├── go.mod
├── go.sum
└── practice/
    └── bcrypt/
        └── main.go
```

`practice/bcrypt/main.go`

```go
package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := "secret123"

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	fmt.Println("password:", password)
	fmt.Println("hash:", string(hash))

	err = bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err != nil {
		fmt.Println("password does not match")
		return
	}

	fmt.Println("password matches")
}
```

実行する。

```bash
go run ./practice/bcrypt
```

---

## 次にやること

- [x] `practice/bcrypt/main.go` で bcrypt を動かす
- [x] 同じパスワードを2回ハッシュ化して、ハッシュ文字列が違うことを確認する
- [ ] わざと違うパスワードを比較して、`password does not match` になることを確認する
- [ ] signup 時に `PasswordHash` として保存する流れを考える
- [ ] login 時に `CompareHashAndPassword` で照合する流れを考える
