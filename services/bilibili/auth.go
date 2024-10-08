package bilibili

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/tidwall/gjson"
)

var (
	mixinKeyEncTab = []int{

		46, 47, 18, 2, 53, 8, 23, 32, 15, 50, 10, 31, 58, 3, 45, 35, 27, 43, 5, 49,
		33, 9, 42, 19, 29, 28, 14, 39, 12, 38, 41, 13, 37, 48, 7, 16, 24, 55, 40,
		61, 26, 17, 0, 1, 60, 51, 30, 4, 22, 25, 54, 21, 56, 59, 6, 63, 57, 62, 11,
		36, 20, 34, 44, 52,
	}
	cache          sync.Map
	lastUpdateTime time.Time
)

const (
	WebURL       = "https://www.bilibili.com"
	WebTicketURL = "https://api.bilibili.com/bapis/bilibili.api.ticket.v1.Ticket/GenWebTicket"
	WebNavURL    = "https://api.bilibili.com/x/web-interface/nav"
	UserAgent    = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"
	key_id       = "ec02"
	key          = "XgwSnGZ1p"
)

func GenWbi(urlStr string) (string, error) {

	newUrlStr, err := signAndGenerateURL(urlStr)
	if err != nil {

		fmt.Printf("Error: %s", err)
		return "", err
	}

	return newUrlStr, nil
}

func GetCookie() (string, error) {

	c := http.Client{}
	req, err := http.NewRequest("GET", WebURL, nil)
	if err != nil {

		fmt.Println(err)
		return "", err
	}
	req.Header.Set("User-Agent", UserAgent)

	resp, err := c.Do(req)
	if err != nil {

		fmt.Println(err)
		return "", err
	}
	defer resp.Body.Close()

	cookies := resp.Header.Get("Set-Cookie")
	return cookies, nil
}

func GenWbiKeysFromTicket() (string, string) {

	ts := time.Now().Unix()
	hexsign, err := hmacSha256(key, fmt.Sprintf("ts%d", ts))
	if err != nil {

		fmt.Println(err)
		return "", ""
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
		return "", ""
	}
	req.Header.Set("User-Agent", UserAgent)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {

		fmt.Println(err)
		return "", ""
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {

		fmt.Println(err)
		return "", ""
	}

	json := string(body)

	imgURL := gjson.Get(json, "data.nav.img").String()
	subURL := gjson.Get(json, "data.nav.sub").String()

	imgKey := strings.Split(strings.Split(imgURL, "/")[len(strings.Split(imgURL, "/"))-1], ".")[0]
	subKey := strings.Split(strings.Split(subURL, "/")[len(strings.Split(subURL, "/"))-1], ".")[0]

	return imgKey, subKey
}

func GenWbiKeysFromNav() (string, string) {

	client := &http.Client{}

	req, err := http.NewRequest("GET", WebNavURL, nil)
	if err != nil {

		fmt.Printf("Error creating request: %s", err)
		return "", ""
	}

	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Referer", "https://www.bilibili.com/")

	resp, err := client.Do(req)
	if err != nil {

		fmt.Printf("Error sending request: %s", err)
		return "", ""
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {

		fmt.Printf("Error reading response: %s", err)
		return "", ""
	}

	json := string(body)

	imgURL := gjson.Get(json, "data.wbi_img.img_url").String()
	subURL := gjson.Get(json, "data.wbi_img.sub_url").String()

	imgKey := strings.Split(strings.Split(imgURL, "/")[len(strings.Split(imgURL, "/"))-1], ".")[0]
	subKey := strings.Split(strings.Split(subURL, "/")[len(strings.Split(subURL, "/"))-1], ".")[0]

	return imgKey, subKey
}

func signAndGenerateURL(urlStr string) (string, error) {

	urlObj, err := url.Parse(urlStr)
	if err != nil {

		return "", err
	}

	imgKey, subKey := getWbiKeysCached()

	query := urlObj.Query()
	params := map[string]string{}
	for k, v := range query {

		params[k] = v[0]
	}

	newParams := encWbi(params, imgKey, subKey)
	for k, v := range newParams {

		query.Set(k, v)
	}

	urlObj.RawQuery = query.Encode()
	newUrlStr := urlObj.String()

	return newUrlStr, nil
}

func encWbi(params map[string]string, imgKey, subKey string) map[string]string {

	mixinKey := getMixinKey(imgKey + subKey)
	currTime := strconv.FormatInt(time.Now().Unix(), 10)
	params["wts"] = currTime

	// Sort keys
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Remove unwanted characters
	for k, v := range params {
		v = sanitizeString(v)
		params[k] = v
	}

	// Build URL parameters
	query := url.Values{}
	for _, k := range keys {
		query.Set(k, params[k])
	}
	queryStr := query.Encode()

	// Calculate w_rid
	hash := md5.Sum([]byte(queryStr + mixinKey))
	params["w_rid"] = hex.EncodeToString(hash[:])
	return params
}

func getMixinKey(orig string) string {

	var str strings.Builder
	for _, v := range mixinKeyEncTab {
		if v < len(orig) {
			str.WriteByte(orig[v])
		}
	}
	return str.String()[:32]
}

func sanitizeString(s string) string {

	unwantedChars := []string{"!", "'", "(", ")", "*"}
	for _, char := range unwantedChars {
		s = strings.ReplaceAll(s, char, "")
	}
	return s
}

func updateCache() {

	if time.Since(lastUpdateTime).Minutes() < 10 {

		return
	}

	imgKey, subKey := GenWbiKeysFromTicket()

	cache.Store("imgKey", imgKey)
	cache.Store("subKey", subKey)

	lastUpdateTime = time.Now()
}

func getWbiKeysCached() (string, string) {

	updateCache()

	imgKeyI, _ := cache.Load("imgKey")
	subKeyI, _ := cache.Load("subKey")

	return imgKeyI.(string), subKeyI.(string)
}

func bytesToHex(bytes []byte) string {

	return hex.EncodeToString(bytes)
}

func hmacSha256(key string, message string) (string, error) {

	mac := hmac.New(sha256.New, []byte(key))

	_, err := mac.Write([]byte(message))
	if err != nil {

		return "", err
	}

	hash := mac.Sum(nil)

	return bytesToHex(hash), nil
}
