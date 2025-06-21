// feature-store/app/internal/signer.go
package internal

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
)

func Hash256(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}

func BuildSigV4Headers(creds aws.Credentials, region, bucket, key string, t time.Time) (map[string]string, error) {
	s := v4.NewSigner()
	req, _ := http.NewRequest("GET", "https://"+bucket+".s3."+region+".amazonaws.com/"+key, nil)

	if err := s.SignHTTP(context.TODO(), creds, req, "", "s3", region, t); err != nil {
		return nil, err
	}

	h := make(map[string]string, len(req.Header))
	for k, v := range req.Header {
		h[k] = v[0]
	}
	return h, nil
}

func SignAttestation(etag, sha string, ts int64) []byte {
	msg := etag + "|" + sha + "|" + strconv.FormatInt(ts, 10)
	sum := sha256.Sum256([]byte(msg))
	return sum[:]
}
