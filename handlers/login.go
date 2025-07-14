package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"go-voting-backend/utils"

	"github.com/golang-jwt/jwt/v5"
)

type LoginRequest struct {
	IC string `json:"ic"`
}

type LoginResponse struct {
	Token     string `json:"token"`
	IC        string `json:"ic"`
	Role      string `json:"role"`
	ExpiresAt int64  `json:"expiresAt"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("➡️  收到 /login 请求")

	// ⛳ Parse request
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.IC == "" {
		log.Printf("❌ 请求解析失败: %v", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	log.Printf("🔍 登录请求的 IC: %s", req.IC)

	// ✅ 确保 Firestore 初始化成功
	if err := utils.InitFirestore(); err != nil {
		log.Printf("❌ Firestore 初始化失败: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Firestore 未初始化"})
		return
	}

	// 🔍 查找用户
	userDoc, err := utils.FirestoreClient.Collection("users").Doc(req.IC).Get(r.Context())
	if err != nil {
		log.Printf("❌ Firestore 查询失败（IC=%s）: %v", req.IC, err)
		http.Error(w, "user not found", http.StatusUnauthorized)
		return
	}
	if !userDoc.Exists() {
		log.Printf("❌ 用户不存在（IC=%s）", req.IC)
		http.Error(w, "user not found", http.StatusUnauthorized)
		return
	}
	log.Println("✅ Firestore 查询成功")

	// 📌 提取角色
	role := "user"
	if r, ok := userDoc.Data()["role"].(string); ok {
		role = r
	}
	log.Printf("🧑 角色: %s", role)

	// 🧾 生成 JWT
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Println("❌ JWT_SECRET 未设置")
		http.Error(w, "server misconfigured", http.StatusInternalServerError)
		return
	}

	expiration := time.Now().Add(24 * time.Hour).Unix()
	claims := jwt.MapClaims{
		"ic":   req.IC,
		"role": role,
		"exp":  expiration,
		"iat":  time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Printf("❌ 签名 JWT 失败: %v", err)
		http.Error(w, "cannot sign token", http.StatusInternalServerError)
		return
	}

	// ✅ 成功返回
	log.Printf("✅ 登录成功，IC: %s", req.IC)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{
		Token:     signedToken,
		IC:        req.IC,
		Role:      role,
		ExpiresAt: expiration,
	})
}
