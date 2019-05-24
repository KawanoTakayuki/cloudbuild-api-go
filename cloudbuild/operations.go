package cloudbuild

// Operation ...
type Operation struct {
	Name     string    `json:"name,omitempty"`
	MetaData *MetaData `json:"metadata,omitempty"`
	Done     bool      `json:"done"`
	//
	Error    *Error            `json:"error,omitempty"`
	Response map[string]string `json:"response,omitempty"`
}

// MetaData ...
type MetaData struct {
	Type  string     `json:"@type,omitempty"`
	Build *BuildConf `json:"build,omitempty"`
}

// Error ...
type Error struct {
	Code    int                 `json:"code"`
	Message string              `json:"message"`
	Details []map[string]string `json:"details"`
	Status  string              `json:"Status"`
}

// HasError ...
func (o *Operation) HasError() (*Error, bool) {
	if o.Error != nil {
		return o.Error, true
	}
	return nil, false
}
