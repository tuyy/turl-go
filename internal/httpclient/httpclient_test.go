package httpclient

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestHttpClient(t *testing.T) {
	client := &http.Client{}

	buf := bytes.NewBuffer([]byte("valid body"))
	req, err := http.NewRequest("POST", "http://127.0.0.1:8080/tuyy", buf)
	if err != nil {
		panic(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
 	}
 	defer resp.Body.Close()

	b, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(b))

	//sc := bufio.NewScanner(resp.Body)
	//for sc.Scan() {
	//	text := sc.Text()
	//	if strings.Contains(text, "doctype html") {
	//		fmt.Println(text)
	//		break
	//	}
	//}
	//if err := sc.Err(); err != nil {
	//	panic(err)
	//}
}
