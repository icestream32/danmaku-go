package bilibili

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type SearchType int

const (
	Video SearchType = iota
	MediaBangumi
	BiliUser
)

var searchTypeMap = []string{

	"video",
	"media_bangumi",
	"bili_user",
}

func (s SearchType) String() string {

	return searchTypeMap[s]
}

const (
	WebSearchAllURL  = "https://api.bilibili.com/x/web-interface/search/all/v2"
	WebSearchTypeURL = "https://api.bilibili.com/x/web-interface/wbi/search/type"
)

// SearchByType searches for a specific type of content on Bilibili.
func SearchByType(keyword string, t SearchType) (string, error) {

	cookie, err := GetCookie()
	if err != nil {

		fmt.Printf("Failed to get cookie: %s", err)
		return "", err
	}

	query := url.Values{}
	query.Add("keyword", keyword)
	query.Add("search_type", t.String())

	newUrlStr, err := signAndGenerateURL(WebSearchTypeURL + "?" + query.Encode())
	if err != nil {

		fmt.Printf("Failed to sign and generate URL: %s", err)
		return "", err
	}

	req, err := http.NewRequest("GET", newUrlStr, nil)
	if err != nil {

		fmt.Printf("Failed to create request: %s", err)
		return "", err
	}

	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Referer", WebURL)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", cookie)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {

		fmt.Printf("Failed to send request: %s", err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {

		fmt.Printf("Failed to read response: %s", err)
		return "", err
	}

	// check status code
	if resp.StatusCode != http.StatusOK {

		return "", errors.New("failed to get 200 OK")
	}

	return string(body), nil
}

// SearchAll searches for all types of content on Bilibili.
func SearchAll(keyword string) (string, error) {

	cookie, err := GetCookie()
	if err != nil {

		fmt.Printf("Failed to get cookie: %s", err)
		return "", err
	}

	query := url.Values{}
	query.Add("keyword", keyword)

	newUrlStr, err := signAndGenerateURL(WebSearchAllURL + "?" + query.Encode())
	if err != nil {

		fmt.Printf("Failed to sign and generate URL: %s", err)
		return "", err
	}

	req, err := http.NewRequest("GET", newUrlStr, nil)
	if err != nil {

		fmt.Printf("Failed to create request: %s", err)
		return "", err
	}

	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Referer", WebURL)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", cookie)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {

		fmt.Printf("Failed to send request: %s", err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {

		fmt.Printf("Failed to read response: %s", err)
		return "", err
	}

	// check status code
	if resp.StatusCode != http.StatusOK {

		return "", errors.New("failed to get 200 OK")
	}

	return string(body), nil
}
