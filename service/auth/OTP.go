package auth


import (
	"crypto/rand"
)




func GenerateOTP() (string, error) {
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

	return string(otp), nil
}