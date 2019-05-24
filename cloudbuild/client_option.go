package cloudbuild

import "net/http"

// Option 追加設定のインターフェイス
type Option interface {
	apply(client *Client)
}

// WithHTTPClient HTTPClientを設定する
func WithHTTPClient(client *http.Client) Option {
	return withHTTPClient{client: client}
}

type withHTTPClient struct {
	client *http.Client
}

func (w withHTTPClient) apply(client *Client) {
	client.httpClient = w.client
}

// WithAccessToken AccessTokenを設定する
func WithAccessToken(token string) Option {
	return withAccessToken{token: token}
}

type withAccessToken struct {
	token string
}

func (w withAccessToken) apply(client *Client) {
	client.acceseToken = w.token
}
