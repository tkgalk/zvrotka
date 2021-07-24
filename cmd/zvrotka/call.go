package zvrotka

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http/httputil"
	"os"

	"github.com/alecthomas/chroma/quick"
	"github.com/galkowskit/zvrotka/http"
	"github.com/galkowskit/zvrotka/parser"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(callCmd)
}

var callCmd = &cobra.Command{
	Use:   "call",
	Short: "call sends the HTTP request",
	Run: func(cmd *cobra.Command, args []string) {
		path := os.Args[2]

		file, err := os.Open(path)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		reqs, err := parser.ReadRequests(file)
		if err != nil {
			log.Fatal(err)
		}

		var reqNames []string
		for _, r := range reqs {
			reqNames = append(reqNames, r.Name+": "+r.Url)
		}

		prompt := promptui.Select{
			Label: "Select Request",
			Items: reqNames,
		}
		idx, _, err := prompt.Run()
		if err != nil {
			log.Fatal(err)
		}

		res, err := http.SendRequest(reqs[idx])
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

		indented := string(prettyJSON.Bytes())

		fmt.Println(string(dump))
		fmt.Printf("Duration: %vms\n\n", res.Duration.Milliseconds())
		err = quick.Highlight(os.Stdout, indented, "json", "terminal16m", "dracula")
		if err != nil {
			log.Fatal(err)
		}
	},
}
