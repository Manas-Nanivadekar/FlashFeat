package internal

import (
	"encoding/json"
	"net"

	"github.com/Manas-Nanivadekar/flashfeat/internal/proto"
)

type EnclaveActor struct {
	conn net.Conn
}

func NewEnclaveActor() (*EnclaveActor, error) {
	c, err := DialEnclave()
	if err != nil {
		return nil, err
	}
	return &EnclaveActor{conn: c}, nil
}

func (a *EnclaveActor) SignGet(bucket, key string) (*proto.SigResponse, error) {
	req := proto.SigRequest{Method: "GET", Bucket: bucket, Key: key}
	if err := json.NewEncoder(a.conn).Encode(req); err != nil {
		return nil, err
	}
	var resp proto.SigResponse
	if err := json.NewDecoder(a.conn).Decode(&resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (a *EnclaveActor) VerifyAndUnpack(etag, sha string, body []byte) (map[string]interface{}, []byte, error) {
	req := proto.VerifyRequest{ETag: etag, SHA256: sha, Body: body}
	if err := json.NewEncoder(a.conn).Encode(req); err != nil {
		return nil, nil, err
	}
	var resp proto.VerifyResponse
	if err := json.NewDecoder(a.conn).Decode(&resp); err != nil {
		return nil, nil, err
	}
	return resp.Payload, resp.Signature, nil
}

func (a *EnclaveActor) Close() { VsockConnCloseGraceful(a.conn) }
