package main

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode/utf8"

	"golang.org/x/crypto/argon2"
)

const (
	passwordArgon2Memory      = 64 * 1024
	passwordArgon2Iterations  = 3
	passwordArgon2Parallelism = 2
	passwordSaltLength        = 16
	passwordHashLength        = 32
	minimumPasswordLength     = 8
	maximumPasswordBytes      = 1024
)

func validateNewPassword(password string) error {
	if len(password) > maximumPasswordBytes {
		return fmt.Errorf("password exceeds %d bytes", maximumPasswordBytes)
	}
	if utf8.RuneCountInString(password) < minimumPasswordLength {
		return fmt.Errorf("password must contain at least %d characters", minimumPasswordLength)
	}
	return nil
}

func hashPassword(password string) (string, error) {
	salt := make([]byte, passwordSaltLength)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", fmt.Errorf("generate password salt: %w", err)
	}
	hash := argon2.IDKey([]byte(password), salt, passwordArgon2Iterations, passwordArgon2Memory, passwordArgon2Parallelism, passwordHashLength)
	return fmt.Sprintf(
		"$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		passwordArgon2Memory,
		passwordArgon2Iterations,
		passwordArgon2Parallelism,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	), nil
}

func verifyPassword(storedHash, password string) (verified, needsUpgrade bool, err error) {
	if strings.HasPrefix(storedHash, "$argon2id$") {
		parameters, salt, expectedHash, err := parseArgon2Hash(storedHash)
		if err != nil {
			return false, false, err
		}
		actualHash := argon2.IDKey([]byte(password), salt, parameters.iterations, parameters.memory, parameters.parallelism, uint32(len(expectedHash)))
		verified = subtle.ConstantTimeCompare(actualHash, expectedHash) == 1
		needsUpgrade = verified && (parameters.memory < passwordArgon2Memory ||
			parameters.iterations < passwordArgon2Iterations ||
			parameters.parallelism != passwordArgon2Parallelism ||
			len(expectedHash) != passwordHashLength)
		return verified, needsUpgrade, nil
	}

	expectedLegacyHash := sha256.Sum256([]byte(password))
	storedBytes, err := hex.DecodeString(storedHash)
	if err != nil || len(storedBytes) != len(expectedLegacyHash) {
		return false, false, nil
	}
	legacyVerified := subtle.ConstantTimeCompare(storedBytes, expectedLegacyHash[:]) == 1
	return legacyVerified, legacyVerified, nil
}

type argon2Parameters struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
}

func parseArgon2Hash(encoded string) (argon2Parameters, []byte, []byte, error) {
	parts := strings.Split(encoded, "$")
	if len(parts) != 6 || parts[0] != "" || parts[1] != "argon2id" || parts[2] != "v=19" {
		return argon2Parameters{}, nil, nil, errors.New("invalid Argon2id password hash format")
	}

	var parameters argon2Parameters
	for _, parameter := range strings.Split(parts[3], ",") {
		keyValue := strings.SplitN(parameter, "=", 2)
		if len(keyValue) != 2 {
			return argon2Parameters{}, nil, nil, errors.New("invalid Argon2id parameter")
		}
		value, err := strconv.ParseUint(keyValue[1], 10, 32)
		if err != nil {
			return argon2Parameters{}, nil, nil, fmt.Errorf("parse Argon2id parameter: %w", err)
		}
		switch keyValue[0] {
		case "m":
			parameters.memory = uint32(value)
		case "t":
			parameters.iterations = uint32(value)
		case "p":
			if value > 255 {
				return argon2Parameters{}, nil, nil, errors.New("Argon2id parallelism is out of range")
			}
			parameters.parallelism = uint8(value)
		default:
			return argon2Parameters{}, nil, nil, errors.New("unknown Argon2id parameter")
		}
	}
	if parameters.memory < 8*1024 || parameters.memory > 256*1024 ||
		parameters.iterations == 0 || parameters.iterations > 10 ||
		parameters.parallelism == 0 || parameters.parallelism > 8 {
		return argon2Parameters{}, nil, nil, errors.New("Argon2id parameters are outside supported limits")
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil || len(salt) < 16 || len(salt) > 64 {
		return argon2Parameters{}, nil, nil, errors.New("invalid Argon2id salt")
	}
	expectedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil || len(expectedHash) < 16 || len(expectedHash) > 64 {
		return argon2Parameters{}, nil, nil, errors.New("invalid Argon2id hash")
	}
	return parameters, salt, expectedHash, nil
}
