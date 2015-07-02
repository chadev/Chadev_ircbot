package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func sprungeSend(t string) (string, error) {
	resp, err := http.PostForm("http://sprunge.us",
		url.Values{"sprunge": {t}})
	if err != nil {
		return "", fmt.Errorf("unable to POST to sprunge.us server: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("response from the server: %s", resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("unable to parse response from the server: %v", err)
	}

	return string(body), nil
}
