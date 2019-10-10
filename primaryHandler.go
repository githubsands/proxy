package proxy

import (
	"log"
	"net/http"
)

var (
	createHandler = func(l *log.Logger) func(http.ResponseWriter, *http.Request) {
		return func(rw http.ResponseWriter, r *http.Request) {
			r.RequestURI = ""
			c := &http.Client{}

			switch r.Method {
			case "GET":
				l.Printf("Making HTTPRequest: %v, To: %v, Dumping headers: %v", r.Method, r.URL.String(), r.Header)
				res, err := c.Do(r)
				if err != nil || res == nil {
					l.Printf(err.Error())
					break
				}

				l.Printf("Receiving HTTPResponse from %v request: %v", r.Method, res)
				break
			case "POST":
				contentType := "json"
				l.Printf("Making HTTPRequest: %v, To: %v, Dumping headers: %v", r.Method, r.URL.String(), r.Header)
				res, err := c.Post(r.URL.String(), contentType, r.Body)
				if err != nil {
					l.Printf(err.Error())
					break
				}
				defer r.Body.Close()

				l.Printf("Receiving HTTPResponse from %v request: %v", r.Method, res)
				break
			case "PUT":
				l.Printf("Making HTTPRequest: %v, To: %v, Dumping headers: %v", r.Method, r.URL.String(), r.Header)
				res, err := c.Do(r)
				if err != nil {
					l.Printf(err.Error())
					break
				}

				l.Printf("Receiving HTTPResponse from %v request: %v", r.Method, res)
				break
			case "HEAD":
				l.Printf("Making HTTPRequest: %v, To: %v, Dumping headers: %v", r.Method, r.URL.String(), r.Header)
				res, err := c.Do(r)
				if err != nil {
					l.Printf("HTTPResponse HEAD received obtained")
					break
				}

				l.Printf("Receiving HTTPResponse from %v request: %v", r.Method, res)
				break
			/*
				case "DELETE", "CONNECT", "PATH":
					l.Printf("Making HTTPRequest: %v, To: %v, Dumping headers: %v", r.Method, r.URL.String(), spew.Sprint(r.Header))
					res, err := c.Do(r)
					if err != nil {
						l.Printf("HTTPResponse HEAD received obtained")
					}

					l.Printf("Receiving HTTPResponse: %v, Dumping response: %v", spew.Sprint(res.Header), spew.Sprint(res.Body))
					break
			*/

			default:
				l.Printf("Request %v has no specified method.\nMethod given %v", r, r.Method)
				break
			}
		}
	}
)
