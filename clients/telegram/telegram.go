package telegram

import (
	"encoding/json"
	errl "example/hello/lib"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

const (
	getUpdatesMethod  = "getUpdates"
	sendMessageMethod = "sendMessage"
)

type Client struct {
	host     string
	basePath string
	client   http.Client
}

func New(host, token string) *Client {
	return &Client{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
	}
}

func newBasePath(token string) string {
	return "bot" + token
}

func (c *Client) Updates(offset, limit int) ([]Updates, error) {
	query := url.Values{}
	query.Add("offset", strconv.Itoa(offset))
	query.Add("limit", strconv.Itoa(limit))

	date, err := c.doResponse(getUpdatesMethod, query)
	if err != nil {
		return nil, err
	}

	var response UpdatesResponse
	if err := json.Unmarshal(date, &response); err != nil {
		return nil, err
	}

	return response.Result, nil
}

func (c *Client) SendMessage(chatId int, text string) error {
	query := url.Values{}
	query.Add("chat_id", strconv.Itoa(chatId))
	query.Add("text", text)

	_, err := c.doResponse(sendMessageMethod, query)
	if err != nil {
		return errl.Wrap("Не смогли отправить сообщение", err)
	}

	return nil
}

func (c *Client) doResponse(method string, query url.Values) (date []byte, err error) {
	const errMsg = "Ошибка при request:"
	defer func() { err = errl.WrapIfErr(errMsg, err) }()

	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = query.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
