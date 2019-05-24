package cloudbuild

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/xerrors"
)

// Client ...
type Client struct {
	porjectID   string
	httpClient  *http.Client
	acceseToken string
}

// Request ...
type Request interface {
	request(projectID string) (*http.Request, error)
	responseMarshler(resBody []byte) error
}

// NewClient 新しいクライアントを作成する
func NewClient(projectID string) *Client {
	return &Client{
		porjectID:  projectID,
		httpClient: http.DefaultClient,
	}
}

// Option オプションを設定する
func (c *Client) Option(options ...Option) *Client {
	for _, option := range options {
		option.apply(c)
	}
	return c
}

// Create ビルド構成の内容を実行する
func (c *Client) Create(conf *BuildConf) (*Operation, error) {
	create := operationCreate(conf)
	if err := c.do(create); err != nil {
		return nil, xerrors.Errorf("operation create error: %w", err)
	}
	return create.Response(), nil
}

// Get ビルドの状態を取得する
func (c *Client) Get(buildID string) (*BuildConf, error) {
	get := operationGet(buildID)
	if err := c.do(get); err != nil {
		return nil, xerrors.Errorf("operation get error: %w", err)
	}
	return get.Response(), nil
}

// Cancel ビルドをキャンセルする
func (c *Client) Cancel(buildID string) (*BuildConf, error) {
	cancel := operationCancel(buildID)
	if err := c.do(cancel); err != nil {
		return nil, xerrors.Errorf("operation cancel error: %w", err)
	}
	return cancel.Response(), nil
}

// Retry ビルドを再実行する
func (c *Client) Retry(buildID string) (*Operation, error) {
	retry := operationRetry(buildID)
	if err := c.do(retry); err != nil {
		return nil, xerrors.Errorf("operation retry error: %w", err)
	}
	return retry.Response(), nil
}

// List プロジェクトのビルド履歴を取得する
func (c *Client) List(pageSize int, pageToken, filter string) (*List, error) {
	list := operationList(pageSize, pageToken, filter)
	if err := c.do(list); err != nil {
		return nil, xerrors.Errorf("operation list error: %w", err)
	}
	return list.Response(), nil
}

// Do ...
func (c *Client) do(r Request) error {
	req, err := r.request(c.porjectID)
	if err != nil {
		return xerrors.Errorf("build request error: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.acceseToken))
	res, err := c.httpClient.Do(req)
	if err != nil {
		return xerrors.Errorf("build client do request error: %w", err)
	}
	bodyByte, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return xerrors.Errorf("build operation read error: %w", err)
	}
	if err := r.responseMarshler(bodyByte); err != nil {
		return xerrors.Errorf("build response Unmarshal error: %w", err)
	}
	return nil
}
