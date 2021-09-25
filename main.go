package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/spf13/pflag"
)

func main() {

	port := pflag.Int("port", 1337, "port to run the server on")
	storePath := pflag.String("path", "", "target path to store tokens to")
	pflag.Parse()

	if *storePath == "" {
		log.Fatal("no path set")
	}

	http.HandleFunc("/token", func(rw http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		targetUrl := params.Get("target")
		name := fmt.Sprintf("%s%s", params.Get("name"), path.Ext(targetUrl))

		tokenReq, err := http.NewRequest(http.MethodGet, targetUrl, nil)
		if err != nil {
			handleErrorResp(err, targetUrl, name, rw)
			return
		}

		tokenResp, err := http.DefaultClient.Do(tokenReq)
		if err != nil {
			handleErrorResp(err, targetUrl, name, rw)
			return
		}

		filename := fmt.Sprintf("%s/%s", *storePath, name)

		_, err = os.Stat(filename)
		if err == nil {
			log.Println(err)
			rw.WriteHeader(http.StatusOK)
			return
		}

		data, err := ioutil.ReadAll(tokenResp.Body)
		if err != nil {
			handleErrorResp(err, targetUrl, name, rw)
			return
		}

		err = os.WriteFile(filename, data, 0666)
		if err != nil {
			handleErrorResp(err, targetUrl, name, rw)
			return
		}

		rw.WriteHeader(http.StatusOK)
	})

	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	if err != nil {
		log.Fatal(err)
	}

}

func handleErrorResp(err error, targetUrl, name string, rw http.ResponseWriter) {
	rw.Header().Add("content-type", "application/json")
	rw.WriteHeader(http.StatusBadRequest)

	res := ErrorResponse{
		Error: err,
		Url:   targetUrl,
		Name:  name,
	}

	data, err := json.Marshal(res)
	if err != nil {
		raw := fmt.Sprintf("%v -> %s", res, err.Error())
		_, _ = rw.Write([]byte(raw))
		return
	}

	_, _ = rw.Write(data)
}

type ErrorResponse struct {
	Error error  `json:"error"`
	Url   string `json:"url"`
	Name  string `json:"name"`
}
