package main

import (
	"auth-api/model"
	"auth-api/repository"
	"auth-api/service"
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
	userRepository := repository.NewMemoryUserRepository()
	authService := service.NewAuthService(userRepository)

	// ルートを追加
	mux.HandleFunc("POST /signup", signupHandler(authService))

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

func signupHandler(authService *service.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req signupRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{
				"error": "invalid JSON",
			})
			return
		}
		// バリデーション(空文字対策)
		if req.Email == "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{
				"error": "email is required",
			})
			return
		}

		if req.Password == "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{
				"error": "password is required",
			})
			return
		}

		user, err := authService.Signup(req.Email, req.Password)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{
				"error": "failed to signup",
			})
			return
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
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
