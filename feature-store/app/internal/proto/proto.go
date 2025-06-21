package proto

type SigRequest struct {
	Method string `json:"method"`
	Bucket string `json:"bucket"`
	Key    string `json:"key"`
}

type SigResponse struct {
	Headers map[string]string `json:"headers"`
	Time    int64
}

type VerifyRequest struct {
	ETag   string `json:"etag"`
	SHA256 string `json:"SHA256"`
	Body   []byte `json:"body"`
}

type VerifyResponse struct {
	Payload   map[string]interface{} `json:"payload"`
	Signature []byte                 `json:"signature"`
}
