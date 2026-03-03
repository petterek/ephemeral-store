package storage

type insertRequest struct {
	Sender   string `json:"sender"`
	DataType string `json:"datatype"`
	Value    string `json:"value"`
	TTL      int    `json:"ttl"`
}

type insertResponse struct {
	Key string `json:"key"`
}

type readResponse struct {
	Value string `json:"value"`
}

type errorResponse struct {
	Error string `json:"error"`
}
