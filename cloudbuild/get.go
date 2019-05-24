package cloudbuild

import (
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/xerrors"
)

// execGet ...
type execGet struct {
	buildID   string
	buildConf *BuildConf
}

// operationGet ...
func operationGet(buildID string) *execGet {
	return &execGet{
		buildID:   buildID,
		buildConf: new(BuildConf),
	}
}

// getProgress ...
func (o *Operation) getProgress() (*execGet, error) {
	if _, yes := o.HasError(); yes {
		return nil, xerrors.New("found operation error")
	}
	return operationGet(o.MetaData.Build.ID), nil
}

// Response ...
func (e *execGet) Response() *BuildConf {
	return e.buildConf
}

func (e *execGet) request(projectID string) (*http.Request, error) {
	path := fmt.Sprintf("https://cloudbuild.googleapis.com/v1/projects/%s/builds/%s", projectID, e.buildID)
	return http.NewRequest("GET", path, nil)
}

func (e *execGet) responseMarshler(resBody []byte) error {
	if err := json.Unmarshal(resBody, e.buildConf); err != nil {
		return xerrors.Errorf("build response Unmarshal error: %w", err)
	}
	return nil
}
