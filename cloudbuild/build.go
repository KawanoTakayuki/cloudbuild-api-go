package cloudbuild

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/xerrors"
)

// execCreate ...
type execCreate struct {
	buildConf *BuildConf
	operation *Operation
}

// operationCreate ...
func operationCreate(conf *BuildConf) *execCreate {
	return &execCreate{
		buildConf: conf,
		operation: new(Operation),
	}
}

// Response ...
func (e *execCreate) Response() *Operation {
	return e.operation
}

func (e *execCreate) request(projectID string) (*http.Request, error) {
	path := fmt.Sprintf("https://cloudbuild.googleapis.com/v1/projects/%s/builds", projectID)
	confByte, err := json.Marshal(e.buildConf)
	if err != nil {
		return nil, xerrors.Errorf("build configration json marshalize error: %w", err)
	}
	return http.NewRequest("POST", path, bytes.NewReader(confByte))
}

func (e *execCreate) responseMarshler(resBody []byte) error {
	if err := json.Unmarshal(resBody, e.operation); err != nil {
		return xerrors.Errorf("build response Unmarshal error: %w", err)
	}
	return nil
}
