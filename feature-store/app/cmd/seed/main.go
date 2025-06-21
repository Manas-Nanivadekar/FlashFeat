package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/vmihailenco/msgpack/v5"
)

func main() {
	if len(os.Args) < 2 {
		panic("usage: seed sample.csv")
	}
	f, _ := os.Open(os.Args[1])
	r := csv.NewReader(f)

	cfg, _ := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion("ap-south-1"),
	)
	u := manager.NewUploader(s3.NewFromConfig(cfg))

	row, _ := r.Read() // header ignored
	for {
		row, err := r.Read()
		if err != nil {
			break
		}

		record := map[string]interface{}{
			"fixed_acidity": row[0],
			// …add the rest…
		}
		buf, _ := msgpack.Marshal(record)
		key := fmt.Sprintf("fg=demo/eid=1/ts=%d.msgpack", time.Now().UnixNano())
		u.Upload(context.TODO(), &s3.PutObjectInput{
			Bucket: aws.String("flashfeat-hot-<acct>"),
			Key:    aws.String(key),
			Body:   bytes.NewReader(buf),
		})
	}
}
