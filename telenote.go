package telenote

import (
	"errors"
	"fmt"
	"io/ioutil"
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
// the given bot token.
func NewNotifier(token string) *Notifier {
	return &Notifier{
		Token:  token,
		Client: nil,
	}
}

// WithToken sets the bot token and returns the Notifier structure.
func (n *Notifier) WithToken(token string) *Notifier {
	n.Token = token
	return n
}

// WithClient sets the http client and returns the Notifier structure.
func (n *Notifier) WithClient(client *http.Client) *Notifier {
	n.Client = client
	return n
}

// Notify sends the text to the given chat and returns an error on failure.
func (n *Notifier) Notify(chatId int64, text string, opts ...Option) error {
	options := NewOptions()
	for _, o := range opts {
		if err := o(options); err != nil {
			return err
		}
	}

	if len(n.Token) == 0 {
		return errors.New("bot token not set")
	} else if len(text) == 0 {
		return errors.New("empty text")
	}

	uri := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", n.Token)
	data := url.Values{}
	data.Set("chat_id", fmt.Sprintf("%d", chatId))
	data.Set("text", text)
	data.Set("parse_mode", options.ParseMode)
	data.Set("disable_web_page_preview", fmt.Sprintf("%v", options.DisableWebPreview))

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
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf("unexpected status: %s: %s", res.Status, string(body))
	}
	return nil
}
