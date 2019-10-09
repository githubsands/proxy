package main

import (
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
)

type config struct {
	host string `json: "host"`
	addr string `json: "address"`
}

func createLogger() (*log.Logger, error) {
	f, err := os.Create("log")
	if err != nil {
		return nil, err
	}

	return log.New(f, "", 1), nil
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
	l.Printf("Creating node at %v:%v", con.host, con.addr)
	return con, nil
}

var (
	createHandler = func(l *log.Logger) func(http.ResponseWriter, *http.Request) {
		return func(rw http.ResponseWriter, r *http.Request) {
			_, _ = httputil.DumpRequest(r, true)
			c := &http.Client{}

			switch r.Method {
			case "GET":
				l.Printf("Making HTTPRequest: %v, To: %v, Dumping headers: %v", r.Method, r.URL.String(), r.Header)

				r.RequestURI = ""
				res, err := c.Do(r)
				if err != nil || res == nil {
					l.Printf(err.Error())
				}

				l.Printf("Receiving HTTPResponse:\nHeader: %v, Response: %v", spew.Sprint(r.Header), spew.Sprint(rw))
				break
			case "POST":
				contentType := "json"
				l.Printf("Making HTTPRequest: %v, To: %v, Dumping headers: %v", r.Method, r.URL.String(), r.Header)
				res, err := c.Post(r.URL.String(), contentType, r.Body)
				if err != nil {
					l.Printf(err.Error())
				}
				defer r.Body.Close()

				l.Printf("Receiving HTTPResponse: %v, Dumping response: %v", spew.Sprint(res.Header), spew.Sprint(res.Body))
				break
			case "PUT":
				res, err := c.Do(r)
				if err != nil {
					l.Printf(err.Error())
				}

				l.Printf("Receiving HTTPResponse: %v, Dumping response: %v", spew.Sprint(res.Header), spew.Sprint(res.Body))
				break
			case "HEAD":
				l.Printf("Making HTTPRequest: %v, To: %v, Dumping headers: %v", r.Method, r.URL.String(), r.Header)
				res, err := c.Do(r)
				if err != nil {
					l.Printf("HTTPResponse HEAD received obtained")
				}

				l.Printf("Receiving HTTPResponse: %v, Dumping response: %v", spew.Sprint(res.Header), spew.Sprint(res.Body))
				break
			case "DELETE", "CONNECT", "PATH":
				l.Printf("Making HTTPRequest: %v, To: %v, Dumping headers: %v", r.Method, r.URL.String(), spew.Sprint(r.Header))
				res, err := c.Do(r)
				if err != nil {
					l.Printf("HTTPResponse HEAD received obtained")
				}

				l.Printf("Receiving HTTPResponse: %v, Dumping response: %v", spew.Sprint(res.Header), spew.Sprint(res.Body))
				break

			default:
				l.Printf("Request %v has no specified method.\nMethod given %v", r, r.Method)
				break
			}
		}
	}
)

func main() {
	var (
		l, err = createLogger()
	)
	if err != nil {
		os.Exit(0)
	}

	config, _ := loadConfig("config.json", l)
	http.ListenAndServe(config.addr, http.HandlerFunc(createHandler(l)))
}
