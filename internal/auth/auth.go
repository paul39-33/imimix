package auth

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashed_pw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("err generating hashed pw: %v", err)
		return "", err
	}

	str_hashed_pw := string(hashed_pw)
	return str_hashed_pw, nil
}

func CheckPasswordHash(password, hashedPassword string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		log.Printf("password does not match: %v", err)
		return false
	}
	return true
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	//create custom claims
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "imimix_app",
		Subject:   userID.String(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})

	//sign the token with the secret key
	tokenString, err := claims.SignedString([]byte(tokenSecret))
	if err != nil {
		log.Printf("err signing token: %v", err)
		return "", err
	}

	return tokenString, nil
}

func ValidateJWT(tokenString string, tokenSecret string) (uuid.UUID, error) {
	//parse the token
	token, err := jwt.ParseWithClaims(
		tokenString,
		&jwt.RegisteredClaims{},
		func(t *jwt.Token) (any, error) {
			//ensure HMAC (HS256 family)
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				log.Printf("Error ensuring HMAC signing method")
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(tokenSecret), nil
		})
	//check if theres any error or token invalid
	if err != nil || !token.Valid {
		log.Printf("err parsing token: %v", err)
		return uuid.Nil, fmt.Errorf("invalid token")
	}
	//check the claims
	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return uuid.Nil, fmt.Errorf("invalid claims type")
	}

	//get user ID from subject
	id, err := uuid.Parse(claims.Subject)
	if err != nil {
		log.Printf("err parsing user id from subject: %v", err)
		return uuid.Nil, fmt.Errorf("invalid user id in token")
	}

	return id, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	//get auth header
	authHeader := headers.Get("Authorization")
	//if no auth header
	if authHeader == "" {
		return "", fmt.Errorf("authorization header missing")
	}
	//take the TOKEN_STRING part
	auth_headers := strings.Fields(authHeader)
	token_string := auth_headers[1]
	return token_string, nil
}

func MakeRefreshToken() (string, error) {
	//generate random key and token
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		log.Printf("err generating random key for refresh token: %v", err)
		return "", err
	}

	key_string := hex.EncodeToString(key)
	return key_string, nil
}

func GetAPIKey (headers http.Header) (string, error) {
	//get auth header
	auth_header := headers.Get("Authorization")
	//if no auth header
	if auth_header == "" {
		return "", fmt.Errorf("authorization header missing")
	}
	//only take the API_KEY part
	auth_headers := strings.Fields(auth_header)
	api_key := auth_headers[1]
	return api_key, nil
}	
