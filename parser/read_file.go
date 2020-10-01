package parser

import (
	"bufio"
	"encoding/json"
	"os"
	"strings"

	"github.com/galkowskit/zvrotka/http"
)

func ReadRequests(file *os.File) ([]http.Request, error) {
	s := bufio.NewScanner(file)

	raws := readRawReqs(s)
	reqs := prepareReqs(raws)

	return reqs, nil
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

func prepareReqs(rawReqs [][]string) []http.Request {
	reqs := make([]http.Request, 0)
	for _, r := range rawReqs {
		req := http.Request{
			Name:   strings.Trim(strings.TrimPrefix(r[0], "###"), " "),
			Method: strings.Fields(r[1])[0],
			Url:    strings.Fields(r[1])[1],
		}

		headers := make([]string, 0)
		for _, h := range r[2:] {
			if strings.Contains(h, ":") && !strings.Contains(h, "\"") {
				headers = append(headers, h)
				continue
			}

			req.Body = req.Body + h
		}

		req.Headers = headers

		if req.Body != "" {
			req.Body = stringify(req.Body)
		}
		reqs = append(reqs, req)
	}

	return reqs
}
