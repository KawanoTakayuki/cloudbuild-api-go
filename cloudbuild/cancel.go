package cloudbuild

import (
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/xerrors"
)

// execCancel ...
type execCancel struct {
	buildID   string
	buildConf *BuildConf
}

// operationCancel ...
func operationCancel(buildID string) *execCancel {
	return &execCancel{
		buildID:   buildID,
		buildConf: new(BuildConf),
	}
}

// cancel ...
func (o *Operation) cancel() (*execCancel, error) {
	if _, yes := o.HasError(); yes {
		return nil, xerrors.New("found operation error")
	}
	return operationCancel(o.MetaData.Build.ID), nil
}

// Response ...
func (e *execCancel) Response() *BuildConf {
	return e.buildConf
}

func (e *execCancel) request(projectID string) (*http.Request, error) {
	path := fmt.Sprintf("https://cloudbuild.googleapis.com/v1/projects/%s/builds/%s:cancel", projectID, e.buildID)
	return http.NewRequest("POST", path, nil)
}

func (e *execCancel) responseMarshler(resBody []byte) error {
	if err := json.Unmarshal(resBody, e.buildConf); err != nil {
		return xerrors.Errorf("build response Unmarshal error: %w", err)
	}
	return nil
}
