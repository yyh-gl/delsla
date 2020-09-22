package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	p "path"
	"strconv"
	"time"
)

const baseURL = "https://slack.com/api"

func main() {
	chs, err := getChannels()
	if err != nil {
		log.Fatal(err)
	}

	for _, ch := range chs {
		mss, err := getMessages(ch.ID, 1)
		if err != nil {
			log.Fatal(err)
		}

		if err := deleteMessages(ch.ID, mss); err != nil {
			log.Fatal(err)
		}

		for _, m := range mss {
			fmt.Println(m.Text)
		}
	}
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

	// TODO: return detail error message
	if !r.OK {
		return nil, errors.New("some error happened in Slack")
	}
	return r.Channels, nil
}

type message struct {
	Text      string `json:"text"`
	TimeStamp string `json:"ts"`
}

func getMessages(channelID string, beforeDays int) ([]*message, error) {
	type result struct {
		OK       bool `json:"ok"`
		Messages []*message
	}

	unix := time.Now().AddDate(0, 0, -beforeDays).Unix()
	latest := strconv.Itoa(int(unix))
	q := map[string]string{
		"latest":  latest,
		"channel": channelID,
	}
	respBody, err := request(http.MethodPost, "conversations.history", q)
	if err != nil {
		return nil, err
	}

	r := new(result)
	if err := json.Unmarshal(respBody, &r); err != nil {
		return nil, err
	}

	// TODO: return detail error message
	if !r.OK {
		return nil, errors.New("some error happened in Slack")
	}
	return r.Messages, nil
}

func deleteMessages(channelID string, messages []*message) error {
	type result struct {
		OK bool `json:"ok"`
	}

	for _, m := range messages {
		q := map[string]string{
			"channel": channelID,
			"ts":      m.TimeStamp,
		}
		respBody, err := request(http.MethodPost, "chat.delete", q)
		if err != nil {
			return err
		}

		r := new(result)
		if err := json.Unmarshal(respBody, &r); err != nil {
			return err
		}

		// TODO: return detail error message
		if !r.OK {
			return errors.New("some error happened in Slack")
		}

		time.Sleep(2 * time.Second)
	}

	return nil
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
