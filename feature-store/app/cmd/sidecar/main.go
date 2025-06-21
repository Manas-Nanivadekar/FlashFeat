package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Manas-Nanivadekar/flashfeat/internal"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	smithyhttp "github.com/aws/smithy-go/transport/http"
)

var (
	region = "ap-south-1"
	bucket = os.Getenv("FLASHFEAT_BUCKET") // set in user-data
)

func sha256Hex(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}

func main() {
	if bucket == "" {
		log.Fatal("env FLASHFEAT_BUCKET not set")
	}

	// ① AWS SDK for S3 GET
	awsCfg, _ := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	s3Client := s3.NewFromConfig(awsCfg)
	down := manager.NewDownloader(s3Client)

	// ② vsock actor (single connection reused)
	actor, err := internal.NewEnclaveActor()
	if err != nil {
		log.Fatalf("vsock: %v", err)
	}
	defer actor.Close()

	http.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("ok"))
	})

	http.HandleFunc("/features", func(w http.ResponseWriter, r *http.Request) {
		fg := r.URL.Query().Get("fg")
		eid := r.URL.Query().Get("eid")
		if fg == "" || eid == "" {
			http.Error(w, "missing fg/eid", 400)
			return
		}
		key := fmt.Sprintf("fg=%s/eid=%s/ts=%d.msgpack", fg, eid, time.Now().UnixNano())

		// A: ask enclave for SigV4 headers
		sigResp, _ := actor.SignGet(bucket, key)

		// B: GET object with those headers
		buf := manager.NewWriteAtBuffer(nil)
		_, err := down.Download(context.TODO(), buf, &s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		}, func(d *manager.Downloader) {
			for k, v := range sigResp.Headers {
				// one Option per header
				d.ClientOptions = append(d.ClientOptions,
					func(o *s3.Options) {
						o.APIOptions = append(o.APIOptions,
							smithyhttp.AddHeaderValue(k, v),
						)
					},
				)
			}
		})
		if err != nil {
			http.Error(w, err.Error(), 502)
			return
		}
		body := buf.Bytes()
		sha := sha256Hex(body)

		// get ETag via HEAD (cheap)
		head, _ := s3Client.HeadObject(context.TODO(), &s3.HeadObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})

		// C: verify in enclave
		payload, sig, err := actor.VerifyAndUnpack(*head.ETag, sha, body)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		w.Header().Set("X-Flashfeat-Signature", hex.EncodeToString(sig))
		_ = json.NewEncoder(w).Encode(payload)
	})

	log.Println("sidecar listening :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
