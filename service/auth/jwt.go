package auth

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/AL-Hourani/care-center/config"
	"github.com/golang-jwt/jwt/v5"
)

func CreateJWT(secret []byte, id int) (string, error) {
	expiration := time.Second * time.Duration(config.Envs.JWTExpirationInSecond)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256 , jwt.MapClaims{
		"ID":strconv.Itoa(id),
		"expiredAt":time.Now().Add(expiration).Unix(),
	})

	tokenString ,err := token.SignedString(secret)
	if err != nil {
		return "",err
	}

	return tokenString , nil
}

type ContextKey string


const UserContextKey ContextKey = "user"


func WithJWTAuth(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// استخراج التوكن من الطلب
		tokenString := GetTokenFromRequest(r)
		if tokenString == "" {
			http.Error(w, "Unauthorized: No token provided", http.StatusUnauthorized)
			return
		}

		// التحقق من صحة التوكن
		token, err := ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}

		// التأكد من أن التوكن صالح ولم ينتهِ
		if !token.Valid {
			http.Error(w, "Unauthorized: Token is not valid", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, token)
		
		// استدعاء الدالة المحمية مع تمرير السياق الجديد
		handlerFunc(w, r.WithContext(ctx))
	}
}



func GetTokenFromRequest(r *http.Request) string {
	// الحصول على التوكن من هيدر Authorization
	tokenAuth := r.Header.Get("Authorization")
	if tokenAuth != "" {
		// التحقق إذا كان التوكن على شكل Bearer
		if len(tokenAuth) > 7 && tokenAuth[:7] == "Bearer " {
			return tokenAuth[7:] // حذف كلمة "Bearer " واستخراج التوكن
		}
		return tokenAuth
	}
	return ""
}


func ValidateToken(tokenString string) (*jwt.Token, error) {
	// التحقق من صحة التوكن وتفسيره
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// التحقق من خوارزمية التوقيع المستخدمة
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// إرجاع المفتاح السري المستخدم للتوقيع
		return []byte(config.Envs.JWTSecret), nil
	})
}


// get id fro token

func GetIDFromToken(token *jwt.Token) (int, error) {
	// استخراج claims من التوكن
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	// استخراج الحقل ID
	idStr, ok := claims["ID"].(string)
	if !ok {
		return 0, fmt.Errorf("ID not found in token")
	}

	// تحويل ID إلى int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("invalid ID format")
	}

	return id, nil
}
