package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
)

type config struct {
	host string `json: "host"`
	addr string `json: "address"`
}

func createLogger() *log.Logger {
	var (
		buf bytes.Buffer
	)

	return log.New(&buf, "", 1)
}

func loadConfig(f string, l *log.Logger) (*config, error) {
	var con *config
	cf, err := os.Open(f)
	defer cf.Close()
	if err != nil {
		l.Printf(err.Error())
		return nil, err
	}

	_ = json.NewDecoder(cf).Decode(&con)
	return con, nil
}

var (
	createHandler = func(l *log.Logger, c *http.Client) func(http.ResponseWriter, *http.Request) {
		return func(rw http.ResponseWriter, r *http.Request) {
			_, _ = httputil.DumpRequest(r, true)
			c := &http.Client{}

			switch r.Method {
			case "GET":
				_, err := c.Do(r)
				if err != nil {
					l.Printf(err.Error())
				}

				break
			case "POST":
				contentType := "json"
				fmt.Printf("\nTesting Post\n")
				_, err := c.Post(r.URL.String(), contentType, r.Body)
				if err != nil {
					l.Printf(err.Error())
				}
				defer r.Body.Close()

				break
			case "PUT":
				_, err := c.Do(r)
				fmt.Print("\nTesting PUT\n")
				if err != nil {
					l.Printf(err.Error())
				}

			default:
				l.Printf("Request %v has no specified method.\nMethod given %v", r, r.Method)
				break
			}
		}
	}
)

func main() {
	var (
		l = createLogger()
		c = &http.Client{}
	)

	config, _ := loadConfig("config.json", l)
	http.ListenAndServe(config.addr, http.HandlerFunc(createHandler(l, c)))
}
