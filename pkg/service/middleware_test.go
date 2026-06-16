package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/url"
	"sort"
	"strings"
	"testing"
)

func buildInitData(t *testing.T, botToken string, values url.Values) string {
	t.Helper()

	values.Del("hash")

	var keys []string
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var parts []string
	for _, k := range keys {
		parts = append(parts, k+"="+values.Get(k))
	}
	checkString := strings.Join(parts, "\n")

	secret := sha256.Sum256([]byte(botToken))
	mac := hmac.New(sha256.New, secret[:])
	mac.Write([]byte(checkString))
	values.Set("hash", hex.EncodeToString(mac.Sum(nil)))

	return values.Encode()
}

func TestParseInitDataValid(t *testing.T) {
	botToken := "123456:ABC-DEF"
	values := url.Values{}
	values.Set("auth_date", "1710000000")
	values.Set("user", `{"id":123,"first_name":"Test","username":"tester"}`)

	initData := buildInitData(t, botToken, values)
	result, err := parseInitData(initData, botToken)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["user"] == "" {
		t.Fatal("expected user in parsed data")
	}
}

func TestParseInitDataInvalidHash(t *testing.T) {
	_, err := parseInitData("user=%7B%7D&hash=bad", "token")
	if err == nil {
		t.Fatal("expected invalid hash error")
	}
}

func TestParseInitDataMissingHash(t *testing.T) {
	_, err := parseInitData("user=%7B%7D", "token")
	if err == nil {
		t.Fatal("expected missing hash error")
	}
}
