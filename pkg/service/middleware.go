package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

const defaultTelegramAuthTTL = 24 * time.Hour

func telegramAuthTTL() time.Duration {
	hoursStr := os.Getenv("TELEGRAM_AUTH_TTL_HOURS")
	if hoursStr == "" {
		return defaultTelegramAuthTTL
	}
	hours, err := strconv.Atoi(hoursStr)
	if err != nil || hours <= 0 {
		return defaultTelegramAuthTTL
	}
	return time.Duration(hours) * time.Hour
}

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

	authDateStr := parsed.Get("auth_date")
	if authDateStr == "" {
		return nil, errors.New("missing auth_date")
	}
	authDateUnix, err := strconv.ParseInt(authDateStr, 10, 64)
	if err != nil {
		return nil, errors.New("invalid auth_date")
	}
	authDate := time.Unix(authDateUnix, 0)
	if time.Since(authDate) > telegramAuthTTL() {
		return nil, fmt.Errorf("initData expired (auth_date older than %s)", telegramAuthTTL())
	}

	// Build result
	result := make(map[string]string)
	for _, k := range keys {
		result[k] = parsed.Get(k)
	}

	return result, nil
}
