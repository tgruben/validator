package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"text/template"
	"time"

	"github.com/pilosa/go-pilosa"
)

type Args struct {
	ABA string
	DB  string
}

type Result struct {
	Queryid  int
	Duration time.Duration
	Args     Args
	Result   []pilosa.QueryResult
	Err      error
}

func runWorkload(t *template.Template, aba, db, host string) error {
	var buf bytes.Buffer
	args := Args{aba, db}
	err := t.Execute(&buf, args)
	if err != nil {
		return err
	}
	client, err := pilosa.NewClient(host, pilosa.OptClientSocketTimeout(time.Hour), pilosa.OptClientConnectTimeout(time.Hour))
	if err != nil {
		return err
	}
	schema, err := client.Schema()
	if err != nil {
		return err
	}
	trait := schema.Index("trait_store")
	var results []Result
	for i, group := range bytes.Split(buf.Bytes(), []byte("\n\n")) {
		r := Result{}
		r.Args = args
		r.Queryid = i
		pql := string(group)
		now := time.Now()
		resp, err := client.Query(trait.RawQuery(pql))
		if err != nil {
			r.Err = err
			continue
		} else {
			r.Duration = time.Since(now)
			r.Result = resp.Results()
		}
		results = append(results, r)

	}
	output, _ := json.Marshal(results)
	fmt.Println(string(output))

	return nil
}

func main() {
	if len(os.Args) != 5 {
		fmt.Println("Usage: validator querytemplatefilename aba db host\n")
		fmt.Println("Example: validator querytemplatefilename aba db host\n")
		os.Exit(1)
	}
	templateFile := os.Args[1] //"q2queries.tmpl"
	aba := os.Args[2]          //"307083665"
	db := os.Args[3]           //"q2db_5093"
	host := os.Args[4]         //"10.0.100.8:10101"
	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		panic(err)
	}
	err = runWorkload(tmpl, aba, db, host)
	if err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
