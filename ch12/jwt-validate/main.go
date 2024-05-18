/*
 * Copyright 2024 Matthew A. Titmus
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secret = []byte("my-super-secret-password123")

func main() {
	http.HandleFunc("/", verifyHandler)
	http.ListenAndServe(":8000", nil)
}

func verifyHandler(w http.ResponseWriter, r *http.Request) {
	header := r.Header.Get("Authorization")
	token := strings.TrimPrefix(header, "Bearer ")

	if header == "" || token == header {
		log.Println("Missing authorization header")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	claims, err := verifyToken(token)
	if err != nil {
		log.Println("invalid authorization token:", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	fmt.Fprintf(w, "Welcome back, %s!\n", claims.Name)
}

type CustomClaims struct {
	Name string `json:"name"`
	jwt.RegisteredClaims
}

func keyFunc(token *jwt.Token) (any, error) {
	return secret, nil
}

func verifyToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, keyFunc)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, fmt.Errorf("unknown claims type")
	}

	if time.Now().After(claims.ExpiresAt.Time) {
		return nil, fmt.Errorf("token expired")
	}

	return claims, nil
}
