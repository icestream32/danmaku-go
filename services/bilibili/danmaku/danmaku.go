package danmaku

import (
	"compress/flate"
	"io"
	"net/http"
)

const (
	WebDanmakuURL = "https://comment.bilibili.com/"
)

func GetDanmaku(cid string) (string, error) {

	newUrlStr := WebDanmakuURL + cid + ".xml"
	req, err := http.NewRequest("GET", newUrlStr, nil)
	if err != nil {

		return "", err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {

		return "", err
	}
	defer resp.Body.Close()

	// check status code
	if resp.StatusCode != http.StatusOK {

		return "", err
	}

	reader := flate.NewReader(resp.Body)
	defer reader.Close()

	content, err := io.ReadAll(reader)
	if err != nil {

		return "", err
	}

	return string(content), nil
}
