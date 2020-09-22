package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	p "path"
)

const baseURL = "https://slack.com/api"

func main() {
	chs, err := getChannels()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("========================")
	fmt.Println(chs)
	fmt.Println("========================")
	os.Exit(0)
}

type channel struct {
	ID string `json:"id"`
}

func getChannels() ([]*channel, error) {
	type result struct {
		OK       bool `json:"ok"`
		Channels []*channel
	}

	respBody, err := request(http.MethodPost, "channels.list", nil)
	if err != nil {
		return nil, err
	}

	r := new(result)
	if err := json.Unmarshal(respBody, &r); err != nil {
		return nil, err
	}
	return r.Channels, nil
}

func request(method, path string, query map[string]string) ([]byte, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	u.Path = p.Join(u.Path, path)

	q := u.Query()
	q.Set("token", os.Getenv("SLACK_API_TOKEN"))
	for k, v := range query {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
