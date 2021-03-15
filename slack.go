package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	p "path"
	"strconv"
	"time"
)

// channel is DTO of Slack channel
type channel struct {
	ID string `json:"id"`
}

// getChannels gets channel list
// >> https://api.slack.com/methods/channels.list
func getChannels() ([]*channel, error) {
	type result struct {
		OK       bool `json:"ok"`
		Channels []*channel
		Error    string `json:"error,omitempty"`
	}

	respBody, err := request(http.MethodPost, "conversations.list", nil)
	if err != nil {
		return nil, err
	}

	r := new(result)
	if err := json.Unmarshal(respBody, &r); err != nil {
		return nil, err
	}

	if !r.OK {
		return nil, fmt.Errorf("conversations.list: %s", r.Error)
	}
	return r.Channels, nil
}

// message is DTO of Slack message
type message struct {
	Text      string `json:"text"`
	TimeStamp string `json:"ts"`
}

// getMessages gets messages before "days" in specified channel
// >> https://api.slack.com/methods/conversations.history
func getMessages(channelID string, days int) ([]*message, error) {
	type result struct {
		OK       bool `json:"ok"`
		Messages []*message
		Error    string `json:"error,omitempty"`
	}

	unix := time.Now().AddDate(0, 0, -days).Unix()
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

	if !r.OK {
		return nil, fmt.Errorf("conversations.history: %s", r.Error)
	}
	return r.Messages, nil
}

// deleteMessages deletes messages
// >> https://api.slack.com/methods/chat.delete
func deleteMessages(channelID string, messages []*message) error {
	type result struct {
		OK    bool   `json:"ok"`
		Error string `json:"error,omitempty"`
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
			return fmt.Errorf("chat.delete: %s", r.Error)
		}

		time.Sleep(2 * time.Second)
	}

	return nil
}

// request is common function for sending HTTP request
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
