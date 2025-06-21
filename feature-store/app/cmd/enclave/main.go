package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net"
	"time"

	"github.com/Manas-Nanivadekar/flashfeat/internal"
	"github.com/Manas-Nanivadekar/flashfeat/internal/proto"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/vmihailenco/msgpack/v5"
)

var (
	region   = "ap-south-1"
	credSeed = credentials.NewStaticCredentialsProvider("AKIA...", "SECRET", "")
)

func main() {
	lis, err := internal.ListenVsock(internal.VsockPort)
	if err != nil {
		log.Fatalf("vsock listen: %v", err)
	}
	log.Println("enclave listening on vsock:", internal.VsockPort)

	for {
		conn, _ := lis.Accept()
		go handle(conn)
	}
}

func handle(c net.Conn) {
	defer internal.VsockConnCloseGraceful(c)

	dec := json.NewDecoder(c)
	enc := json.NewEncoder(c)

	for {
		var peek map[string]interface{}
		if err := dec.Decode(&peek); err != nil {
			if err != io.EOF {
				log.Printf("decode: %v", err)
			}
			return
		}

		if _, ok := peek["method"]; ok {
			var req proto.SigRequest
			jsonBytes, _ := json.Marshal(peek)
			_ = json.Unmarshal(jsonBytes, &req)

			t := time.Now().UTC()
			cred, err := credSeed.Retrieve(context.TODO())
			if err != nil {
				log.Fatal(err)
			}
			hdrs, _ := internal.BuildSigV4Headers(cred, region, req.Bucket, req.Key, t)
			enc.Encode(proto.SigResponse{Headers: hdrs, Time: t.UnixMilli()})

		} else { // VerifyRequest
			var req proto.VerifyRequest
			jsonBytes, _ := json.Marshal(peek)
			_ = json.Unmarshal(jsonBytes, &req)

			// ① local SHA check
			if req.SHA256 != sha256Hex(req.Body) {
				enc.Encode(map[string]string{"error": "sha mismatch"})
				continue
			}
			// ② compare ETag
			if req.ETag[1:len(req.ETag)-1] != req.SHA256 { // remove quotes
				enc.Encode(map[string]string{"error": "etag mismatch"})
				continue
			}
			// ③ decode msgpack to JSON
			var out map[string]interface{}
			if err := msgpack.Unmarshal(req.Body, &out); err != nil {
				enc.Encode(map[string]string{"error": "msgpack decode"})
				continue
			}
			sig := internal.SignAttestation(req.ETag, req.SHA256, time.Now().UnixMilli())
			enc.Encode(proto.VerifyResponse{Payload: out, Signature: sig})
		}
	}
}

func sha256Hex(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}
