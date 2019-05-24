package cloudbuild

import (
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/xerrors"
)

// execRetry ...
type execRetry struct {
	buildID   string
	operation *Operation
}

// operationRetry ...
func operationRetry(buildID string) *execRetry {
	return &execRetry{
		buildID:   buildID,
		operation: new(Operation),
	}
}

// retry ...
func (o *Operation) retry() (*execRetry, error) {
	if _, yes := o.HasError(); yes {
		return nil, xerrors.New("found operation error")
	}
	return operationRetry(o.MetaData.Build.ID), nil
}

// Response ...
func (e *execRetry) Response() *Operation {
	return e.operation
}

func (e *execRetry) request(projectID string) (*http.Request, error) {
	path := fmt.Sprintf("https://cloudbuild.googleapis.com/v1/projects/%s/builds/%s:retry", projectID, e.buildID)
	return http.NewRequest("POST", path, nil)
}

func (e *execRetry) responseMarshler(resBody []byte) error {
	if err := json.Unmarshal(resBody, e.operation); err != nil {
		return xerrors.Errorf("build response Unmarshal error: %w", err)
	}
	return nil
}
