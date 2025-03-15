package hmac

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	veracodeRequestVersionString = "vcode_request_version_1"
	dataFormat                   = "id=%s&host=%s&url=%s&method=%s"
	headerFormat                 = "%s id=%s,ts=%s,nonce=%X,sig=%X"
	veracodeHMACSHA256           = "VERACODE-HMAC-SHA-256"
)

func currentTimestamp() int64 {
	return time.Now().UnixMilli()
}

func generateNonce(size int) ([]byte, error) {
	nonce := make([]byte, size)
	_, err := rand.Read(nonce)

	if err != nil {
		return nil, err
	}

	return nonce, nil
}

func hmac256(message, key []byte) []byte {
	sha := hmac.New(sha256.New, key)
	sha.Write(message)
	return sha.Sum(nil)
}

func removeRegion(apiCredential string) string {
	if strings.Contains(apiCredential, "-") {
		return strings.Split(apiCredential, "-")[1]
	} else {
		return apiCredential
	}
}

func calculateSignature(key, nonce, timestamp, data []byte) []byte {
	encryptedNonce := hmac256(nonce, key)
	encryptedTimestamp := hmac256(timestamp, encryptedNonce)
	signingKey := hmac256([]byte(veracodeRequestVersionString), encryptedTimestamp)
	return hmac256(data, signingKey)
}

// Returns the value for the Authorization header that must be added to requests
func CalculateAuthorizationHeader(url *url.URL, httpMethod, apiKeyID, apiKeySecret string) (string, error) {
	apiKeyID = removeRegion(apiKeyID)
	apiKeySecret = removeRegion(apiKeySecret)
	nonce, err := generateNonce(16)

	if err != nil {
		return "", err
	}

	secret, err := hex.DecodeString(apiKeySecret)

	if err != nil {
		return "", err
	}

	timestamp := strconv.FormatInt(currentTimestamp(), 10)
	data := fmt.Sprintf(dataFormat, apiKeyID, url.Hostname(), url.RequestURI(), httpMethod)
	dataSignature := calculateSignature(secret, nonce, []byte(timestamp), []byte(data))

	return fmt.Sprintf(headerFormat, veracodeHMACSHA256, apiKeyID, timestamp, nonce, dataSignature), nil
}
