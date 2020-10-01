package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http/httputil"
	"os"

	"github.com/galkowskit/zvrotka/http"
	"github.com/galkowskit/zvrotka/parser"
)

func main() {
	path := os.Args[1]
	name := os.Args[2]

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reqs, err := parser.ReadRequests(file)
	if err != nil {
		log.Fatal(err)
	}

	var req http.Request
	for _, r := range reqs {
		if r.Name == name {
			req = r
			break
		}
	}

	res, err := http.SendRequest(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Res.Body.Close()

	data, err := ioutil.ReadAll(res.Res.Body)
	if err != nil {
		log.Fatal(err)
	}

	dump, err := httputil.DumpResponse(res.Res, true)
	if err != nil {
		log.Fatal(err)
	}

	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, data, "", "\t")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(dump))
	fmt.Printf("Duration: %vms\n\n", res.Duration.Milliseconds())
	fmt.Println(string(prettyJSON.Bytes()))
}
