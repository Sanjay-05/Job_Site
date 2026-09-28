package helper

import (
	"crypto/rand"
	"errors"
	"fmt"
	"jobsite/models"
	"jobsite/utility"
	"math/big"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt/v5"
)

func ValidateJson(payload interface{}) error {
	validate := validator.New()

	if err := validate.Struct(payload); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			var errorMessages []string
			for _, e := range validationErrors {

				if e.Tag() == "required" {
					errorMessages = append(errorMessages, fmt.Sprintf("%s is required field", e.Field()))
				} else {
					errorMessages = append(errorMessages, fmt.Sprintf("%s failed on %s", e.Field(), e.Tag()))
				}
			}
			return errors.New((strings.Join(errorMessages, ", ")))
		}
		return err
	}
	return nil
}

// helper function to generate random number for otp to reset password
func GenerateOTP() (string, error) {
	max := big.NewInt(1000000)

	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%06d", n), nil
}

// helper function to cleanup the otp
func StartOTPCleanup(otpMutex *sync.Mutex, otpStore map[string]utility.OTPData) {
	go func() {
		for {
			time.Sleep(30 * time.Second)
			otpMutex.Lock()
			for k, v := range otpStore {
				if time.Now().After(v.Expires) {
					delete(otpStore, k)
				}
			}
			otpMutex.Unlock()
		}
	}()
}

// helper function to validate the password strength
func ValidatePassword(password string) error {
	var (
		hasUpperCase bool
		hasLowerCase bool
		hasDigit     bool
		hasSpecial   bool
	)

	for _, char := range password {
		if unicode.IsUpper(char) {
			hasUpperCase = true
		} else if unicode.IsLower(char) {
			hasLowerCase = true
		} else if unicode.IsDigit(char) {
			hasDigit = true
		} else {
			hasSpecial = true
		}
	}

	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	} else if !hasUpperCase {
		return errors.New("password must contain at least one uppercase letter")
	} else if !hasLowerCase {
		return errors.New("password must contain at least one lowercase letter")
	} else if !hasDigit {
		return errors.New("password must contain at least one number")
	} else if !hasSpecial {
		return errors.New("password must contain at least one special character")
	}

	return nil
}

func GenerateToken(user *models.User, tokenTime time.Duration, jwtSecret string) (string, error) {
	claims := models.JwtCustomClaims{
		UserID: user.ID.String(),
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now())}}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

func ParseToken(tokenString, jwtSecret string) (*models.JwtCustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.JwtCustomClaims{}, func(t *jwt.Token) (any, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*models.JwtCustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid credentials")
}
