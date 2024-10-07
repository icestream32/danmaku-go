package bilibili

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	WebURL       = "https://www.bilibili.com"
	WebTicketURL = "https://api.bilibili.com/bapis/bilibili.api.ticket.v1.Ticket/GenWebTicket"
	key_id       = "ec02"
	key          = "XgwSnGZ1p"
)

func GetCookie() string {

	c := http.Client{}
	req, err := http.NewRequest("GET", WebURL, nil)
	if err != nil {

		fmt.Println(err)
		return ""
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	resp, err := c.Do(req)
	if err != nil {

		fmt.Println(err)
		return ""
	}
	defer resp.Body.Close()

	cookies := resp.Header.Get("Set-Cookie")
	return cookies
}

func GenWebTicket() {

	ts := time.Now().Unix()
	hexsign, error := hmacSha256(key, fmt.Sprintf("ts%d", ts))
	if error != nil {

		fmt.Println(error)
		return
	}

	query := url.Values{}

	params := []struct {
		key   string
		value string
	}{

		{"key_id", key_id},
		{"hexsign", hexsign},
		{"context[ts]", fmt.Sprintf("%d", ts)},
		{"csrf", ""},
	}

	for _, param := range params {

		query.Add(param.key, param.value)
	}

	url := fmt.Sprintf("%s?%s", WebTicketURL, query.Encode())
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {

		fmt.Println(err)
		return
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {

		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {

		fmt.Println(err)
		return
	}

	fmt.Println(string(body))
}

// Convert a byte array to a hex string.
func bytesToHex(bytes []byte) string {

	return hex.EncodeToString(bytes)
}

// Generate a HMAC-SHA256 hash of the given message string using the given key string.
func hmacSha256(key string, message string) (string, error) {

	// Create a new HMAC by specifying the SHA256 hash function and the key.
	mac := hmac.New(sha256.New, []byte(key))

	// Write the message to the HMAC.
	_, err := mac.Write([]byte(message))
	if err != nil {

		return "", err
	}

	// Get the final HMAC result.
	hash := mac.Sum(nil)

	// Return the hex representation of the hash.
	return bytesToHex(hash), nil
}
