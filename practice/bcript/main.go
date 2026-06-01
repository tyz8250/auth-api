package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// bcryptはパスワードをハッシュ化するためのライブラリ
func main() {
	password := "secret123"

	// 平文パスワードをハッシュ化
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	// 同じパスワードを2回ハッシュ化して、ハッシュ文字列が違うことを確認する
	hash2, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	fmt.Println("password:", password)
	fmt.Println("hash:", string(hash))
	fmt.Println("hash2:", string(hash2))

	// ハッシュ化されたパスワードと平文パスワードを比較
	err = bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err != nil {
		fmt.Println("password does not match")
		return
	}

	fmt.Println("password matches")

	// わざと違うパスワードを比較して、password does not match になることを確認する
	err = bcrypt.CompareHashAndPassword(hash, []byte("wrongpassword"))
	if err != nil {
		fmt.Println("password does not match")
		return
	}

	fmt.Println("password matches")
}
