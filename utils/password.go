package utils

import (
	"crypto/rand"
	"errors"
	"math/big"
	"regexp"
	"strings"
)

func GenerateTemporaryPassword() (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%&*+-=?"
	const length = 16

	password := make([]byte, length)

	categories := []string{
		"abcdefghijklmnopqrstuvwxyz",
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		"0123456789",
		"!@#$%&*+-=?",
	}

	for i, cat := range categories {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(cat))))
		if err != nil {
			return "", err
		}
		password[i] = cat[num.Int64()]
	}

	for i := len(categories); i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		password[i] = charset[num.Int64()]
	}

	for i := len(password) - 1; i > 0; i-- {
		j, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		if err != nil {
			return "", err
		}
		password[i], password[j.Int64()] = password[j.Int64()], password[i]
	}

	return string(password), nil
}

func ValidatePassword(password string) error {
	if len(password) < 12 {
		return errors.New("senha deve ter pelo menos 12 caracteres")
	}

	if len(password) > 128 {
		return errors.New("senha deve ter no máximo 128 caracteres")
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSymbol := regexp.MustCompile(`[!@#$%&*+\-=?]`).MatchString(password)

	if !hasUpper {
		return errors.New("senha deve conter pelo menos uma letra maiúscula")
	}

	if !hasLower {
		return errors.New("senha deve conter pelo menos uma letra minúscula")
	}

	if !hasNumber {
		return errors.New("senha deve conter pelo menos um número")
	}

	if !hasSymbol {
		return errors.New("senha deve conter pelo menos um símbolo especial (!@#$%&*+-=?)")
	}

	if regexp.MustCompile(`(.)\1{2,}`).MatchString(password) {
		return errors.New("senha não pode conter mais de 2 caracteres consecutivos iguais")
	}

	sequences := []string{
		"123456", "abcdef", "qwerty", "password", "admin", "user",
		"654321", "fedcba", "ytrewq", "drowssap", "nimda", "resu",
	}

	for _, seq := range sequences {
		if strings.Contains(strings.ToLower(password), seq) {
			return errors.New("senha não pode conter sequências comuns ou palavras óbvias")
		}
	}

	return nil
}
