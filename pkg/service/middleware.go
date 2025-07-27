package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/url"
	"sort"
	"strings"
)

// Парсим telegram init data
func parseInitData(initData, botToken string) (map[string]string, error) {
	parsed, err := url.ParseQuery(initData)
	if err != nil {
		return nil, err
	}

	// Extract hash
	hash := parsed.Get("hash")
	if hash == "" {
		return nil, errors.New("missing hash")
	}
	parsed.Del("hash")

	// Sort keys
	var keys []string
	for k := range parsed {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Build check_string
	var checkStrings []string
	for _, k := range keys {
		checkStrings = append(checkStrings, k+"="+parsed.Get(k))
	}
	checkString := strings.Join(checkStrings, "\n")

	// Compute HMAC SHA256
	secret := sha256.Sum256([]byte(botToken))
	h := hmac.New(sha256.New, secret[:])
	h.Write([]byte(checkString))
	expected := hex.EncodeToString(h.Sum(nil))

	if expected != hash {
		return nil, errors.New("invalid Telegram auth")
	}

	// Build result
	result := make(map[string]string)
	for _, k := range keys {
		result[k] = parsed.Get(k)
	}

	return result, nil
}
