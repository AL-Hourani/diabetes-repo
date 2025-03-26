package auth


import (
	"crypto/rand"
	"time"
	"github.com/patrickmn/go-cache"
)



var otpCache = cache.New(5*time.Minute, 10*time.Minute)

func GenerateOTP(email string) (string, error) {
	const digits = "0123456789"
	length := 6
	otp := make([]byte, length)

	_, err := rand.Read(otp)
	if err != nil {
		return "", err
	}

	for i := range otp {
		otp[i] = digits[int(otp[i])%len(digits)]
	}

	otpCode := string(otp)
	
	otpCache.Set(email, otpCode, cache.DefaultExpiration)
	return otpCode, nil
}

func VerifyOTP(email, otpCode string) bool {
	storedOTP, found := otpCache.Get(email)
	if !found {
		return false
	}
	return storedOTP == otpCode
}