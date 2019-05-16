package telenote

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// Notifier contains telegram credentials and state.
type Notifier struct {
	Token  string
	Client *http.Client
}

// NewNotifier returns a new Notifier instance initialized with
// the given access token.
func NewNotifier(token string) *Notifier {
	return &Notifier{
		Token:  token,
		Client: nil,
	}
}

// WithToken sets the access token and returns the Notifier structure.
func (n *Notifier) WithToken(token string) *Notifier {
	n.Token = token
	return n
}

// WithClient sets the http client and returns the Notifier structure.
func (n *Notifier) WithClient(client *http.Client) *Notifier {
	n.Client = client
	return n
}

// Notify send the notification and returns an error on failure.
func (n *Notifier) Notify(chatId string, text string) error {
	if len(n.Token) == 0 {
		return errors.New("access token not set")
	}
	if len(text) == 0 {
		return errors.New("empty text")
	}

	uri := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", n.Token)
	data := url.Values{}
	data.Set("chat_id", chatId)
	data.Set("text", text)
	data.Set("parse_mode", "Markdown")
	data.Set("disable_web_page_preview", "false")

	client := n.Client
	if client == nil {
		client = &http.Client{}
	}

	r, _ := http.NewRequest("POST", uri, strings.NewReader(data.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	res, err := client.Do(r)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("unexpected status: %s", res.Status)
	}
	return nil
}
