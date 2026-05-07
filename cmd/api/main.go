package main

import (
	"auth-api/model"
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	// サーバーマルチプレクサを作成
	mux := newServerMux()

	// サーバーを起動
	log.Println("server started on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

// newServerMux は、すべてのルートを設定した新しいサーバーマルチプレクサを作成します。
func newServerMux() *http.ServeMux {
	// 新しいサーバーマルチプレクサを作成
	mux := http.NewServeMux()

	// ルートを追加
	mux.HandleFunc("POST /signup", signupHandler)

	mux.HandleFunc("POST /login", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{
			"message": "login route",
		})
	})

	mux.HandleFunc("/me", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{
			"message": "me route",
		})
	})

	return mux
}

type signupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	var req signupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid JSON",
		})
		return
	}

	user := model.User{
		ID:           1,
		Email:        req.Email,
		PasswordHash: "hashed_password",
		CreatedAt:    "2025-10-15T12:00:00Z",
		UpdatedAt:    "2025-10-15T12:00:00Z",
	}

	res := model.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	// レスポンスを返す
	writeJSON(w, http.StatusOK, res)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
