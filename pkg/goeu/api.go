package goeu

type ApiExecuteCommand struct {
	Token     string                 `json:"tk"`
	ID        string                 `json:"id"`
	Namespace string                 `json:"ns"`
	Method    string                 `json:"mt"`
	Params    map[string]interface{} `json:"pm"`
}

type ApiExecuteResult struct {
	ID      string                 `json:"id,omitempty"`
	Success bool                   `json:"ok"`
	Payload map[string]interface{} `json:"pl"`
	Error   string                 `json:"er,omitempty"`
}
