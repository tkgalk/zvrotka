package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type request struct {
    name string
    method string
    url string
    headers []string
    body string
}

func (r request) String() string {
    return fmt.Sprintf("{\nName: %s\nMethod: %s\nURL: %s\nHeaders: %v\nBody: %s\n}", r.name, r.method, r.url, r.headers, r.body)
}

func stringify(s string) string {
    in := []byte(s)
    var raw map[string]interface{}
    if err := json.Unmarshal(in, &raw); err != nil {
        panic(err)
    }
    out, err := json.Marshal(raw)
    if err != nil {
        panic(err)
    }
    return string(out)
}

func readRawReqs(s *bufio.Scanner) [][]string {
    res := make([][]string, 0)
    raw := make([]string, 0)

    for s.Scan() {
        line := s.Text()

        if strings.HasPrefix(line, "###") {
            res = append(res, raw)
            raw = make([]string, 0)
        }

        raw = append(raw, line)
    }

    res = append(res, raw)
    return res[1:]
}

func prepareReqs(rawReqs [][]string) []request {
    reqs := make([]request, 0)
    for _, r := range rawReqs {
        req := request{
            name: strings.Trim(strings.TrimPrefix(r[0], "###"), " "),
            method: strings.Fields(r[1])[0],
            url: strings.Fields(r[1])[1],
        }

        headers := make([]string, 0)
        for _, h := range r[2:] {
            if strings.Contains(h, ":") && !strings.Contains(h, "\"") {
                headers = append(headers, h)
                continue
            }

            req.body = req.body + h
        }

        req.headers = headers

        if req.body != "" {
            req.body = stringify(req.body)
        }
        reqs = append(reqs, req)
    }

    return reqs
}

func main() {
    file, err := os.Open("requests.http")
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    s := bufio.NewScanner(file)
    if err := s.Err(); err != nil {
        log.Fatal(err)
    }

    reqs := prepareReqs(readRawReqs(s))

    url, err := url.Parse(reqs[0].url)
    if err != nil {
        log.Fatal(err)
    }

    body := ioutil.NopCloser(strings.NewReader(reqs[0].body))

    req := &http.Request{
        Method: reqs[0].method,
        URL: url,
        Header: map[string][]string{},
        Body: body,
    }

    start := time.Now()
    res, err := http.DefaultClient.Do(req)
    if err != nil {
        log.Fatal(err)
    }

    duration := time.Since(start)

    data, err := ioutil.ReadAll(res.Body)
    if err != nil {
        log.Fatal(err)
    }
    defer res.Body.Close()

    fmt.Printf("%s %s\n", reqs[0].method, reqs[0].url)
    fmt.Printf("Status: %d\nDuration:%s\n", res.StatusCode, duration)
    fmt.Print(string(data))
}

