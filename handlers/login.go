package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/firestore"
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
	// ⛳ Parse request
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.IC == "" {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// 🔐 Check JWT_SECRET
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		http.Error(w, "server misconfigured", http.StatusInternalServerError)
		return
	}

	// 🔍 Check Firestore if IC exists
	ctx := context.Background()
	fireClient, err := firestore.NewClient(ctx, os.Getenv("GOOGLE_PROJECT_ID"))
	if err != nil {
		http.Error(w, "firestore init error", http.StatusInternalServerError)
		return
	}
	defer fireClient.Close()

	userDoc, err := fireClient.Collection("users").Doc(req.IC).Get(ctx)
	if err != nil || !userDoc.Exists() {
		http.Error(w, "user not found", http.StatusUnauthorized)
		return
	}

	// 📌 Extract role from Firestore (default "user")
	role := "user"
	if r, ok := userDoc.Data()["role"].(string); ok {
		role = r
	}

	// 🧾 Generate JWT token
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
		http.Error(w, "cannot sign token", http.StatusInternalServerError)
		return
	}

	// ✅ Return token & metadata
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{
		Token:     signedToken,
		IC:        req.IC,
		Role:      role,
		ExpiresAt: expiration,
	})
}
