package handlers

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
	"winter-home-server/internal/models"

	"winter-home-server/internal/middleware"

	"golang.org/x/crypto/bcrypt"
)

var (
	users  = make(map[string]*models.User)
	mu     sync.RWMutex
	nextID int = 1
)

func init() {
	users["admin"] = &models.User{
		ID:        0,
		Username:  "admin",
		Password:  hashPassword("admin123"),
		CreatedAt: time.Now(),
	}
}

func hashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes)
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "只允许 POST 请求", http.StatusMethodNotAllowed)
		return
	}

	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "请求格式错误", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if req.Username == "" || req.Password == "" {
		http.Error(w, "用户名和密码不能为空", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	if _, exists := users[req.Username]; exists {
		http.Error(w, "用户名已存在", http.StatusConflict)
		return
	}

	hashedPassword := hashPassword(req.Password)
	user := &models.User{
		ID:        nextID,
		Username:  req.Username,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
	}
	users[req.Username] = user
	nextID++

	token, err := middleware.GenerateToken(user.ID, user.Username)
	if err != nil {
		http.Error(w, "生成 Token 失败", http.StatusInternalServerError)
		return
	}

	resp := models.AuthResponse{
		Token:     token,
		ExpiresIn: 86400,
		User:      *user,
	}
	respondJSON(w, http.StatusCreated, resp)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "只允许 POST 请求", http.StatusMethodNotAllowed)
		return
	}

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "请求格式错误", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	mu.RLock()
	user, exists := users[req.Username]
	mu.RUnlock()

	if !exists || !checkPasswordHash(req.Password, user.Password) {
		http.Error(w, "用户名或密码错误", http.StatusUnauthorized)
		return
	}

	token, err := middleware.GenerateToken(user.ID, user.Username)
	if err != nil {
		http.Error(w, "生成 Token 失败", http.StatusInternalServerError)
		return
	}

	resp := models.AuthResponse{
		Token:     token,
		ExpiresIn: 86400,
		User:      *user,
	}
	respondJSON(w, http.StatusOK, resp)
}
