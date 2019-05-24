package cloudbuild

import (
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/xerrors"
)

// List ...
type List struct {
	Build         []*BuildConf `json:"builds,omitempty"`
	NextPageToken string       `json:"nextPageToken,omitempty"`
}

// execLists ...
type execLists struct {
	pageSize  int
	pageToken string
	filter    string
	list      *List
}

// operationList ...
func operationList(pageSize int, pageToken, filter string) *execLists {
	return &execLists{
		pageSize:  pageSize,
		pageToken: pageToken,
		filter:    filter,
		list:      new(List),
	}
}

// Response ...
func (e *execLists) Response() *List {
	return e.list
}

func (e *execLists) request(projectID string) (*http.Request, error) {
	path := fmt.Sprintf("https://cloudbuild.googleapis.com/v1/projects/%s/builds", projectID)
	return http.NewRequest("GET", path, nil)
}

func (e *execLists) responseMarshler(resBody []byte) error {
	if err := json.Unmarshal(resBody, e.list); err != nil {
		return xerrors.Errorf("build response Unmarshal error: %w", err)
	}
	return nil
}
