package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("usage: model <sidecar-ip>")
	}
	ip := os.Args[1]
	resp, err := http.Get(fmt.Sprintf("http://%s:8080/features?fg=demo&eid=1", ip))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)

	var m map[string]interface{}
	_ = json.Unmarshal(b, &m)
	fmt.Printf("features: %+v\n", m)
}
