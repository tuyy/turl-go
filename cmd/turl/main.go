package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type Args struct {
	Url     string
	Method  string
	Body    string
	Output  string
	Input   string
	Headers map[string]string
}

var args Args

func init() {
	if len(os.Args) < 2 {
		panic("USAGE) ./turl {url} or ./turl --help")
	}

	method := flag.String("x", "GET", "http req method")
	headers := flag.String("h", "{}", "headers json str")
	flag.StringVar(&args.Body, "d", "", "http req body")
	flag.StringVar(&args.Output, "o", "", "output")
	flag.StringVar(&args.Input, "r", "", "input req file")
	flag.Parse()

	args.Method = strings.ToUpper(*method)

	err := json.Unmarshal([]byte(*headers), &args.Headers)
	if err != nil {
		panic(err)
	}

	if len(args.Input) > 0 {
		inputData, err := ioutil.ReadFile(args.Input)
		if err != nil {
			panic(err)
		}
		parseReq(string(inputData))
	} else {
		args.Url = strings.TrimSpace(os.Args[1])
		if args.Url[:4] != "http" {
			args.Url = "http://" + args.Url
		}
	}

	fmt.Printf("================== REQUEST ==================\n%#v\n\n", args)
}

func main() {
	requestHttpWithArgs()
}

// InputData Req format
// {method} {url}
// {header(key:value..)}
// {header(key:value..)}
// {header(key:value..)}
// ..empty line..
// {body}
func parseReq(inputData string) {
	tokens := strings.Split(inputData, "\n")
	methodAndUrl := strings.Split(tokens[0], " ")

	args.Method = strings.ToUpper(strings.TrimSpace(methodAndUrl[0]))
	args.Url = strings.TrimSpace(methodAndUrl[1])
	if args.Url[:4] != "http" {
		args.Url = "http://" + args.Url
	}

	isBody := false
	for _, line := range tokens[1:] {
		if line == "" {
			isBody = true
			continue
		}

		if isBody {
			args.Body += line + "\n"
		} else {
			keyValue := strings.Split(line, ":")
			key := strings.TrimSpace(keyValue[0])
			value := strings.TrimSpace(keyValue[1])
			args.Headers[key] = value
		}
	}
}

func requestHttpWithArgs() {
	client := &http.Client{}
	buf := bytes.NewBuffer([]byte(args.Body))
	req, err := http.NewRequest(args.Method, args.Url, buf)
	if err != nil {
		panic(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	fmt.Println(makeResultWithResp(resp))
}

func makeResultWithResp(resp *http.Response) string {
	defer resp.Body.Close()

	headers, err := json.MarshalIndent(resp.Header, "", "  ")
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("================== RESPONSE ==================\n%s\n%s\n\n%s\n",
		resp.Status,
		headers,
		body)
}
