package main

import (
	// "github.com/davecgh/go-spew/spew"
	// "io/ioutil"
	// "bytes"
	"net/http"
	"net/http/httptest"
	"os"
	//"reflect"
	// "fmt"
	"bytes"
	"testing"
)

var createFileFunc = func(fileName string) (*os.File, error) {
	// make file here
	fi, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}

	return fi, nil
}

/*
func TestLoadConfig(t *testing.T) {
	var (
		l        = createLogger()
		fileName = "testFile.json"
		f, _     = createFileFunc(fileName)
	)

	defer os.Remove(fileName)

	data := "{\n	'host': 'localhost',\n	'port': '8080'\n}"
	_, err := fmt.Fprintf(f, data)
	if err != nil {
		t.Fatalf("unable to create file")
	}

	/*
		b, _ := ioutil.ReadAll(f)
		spew.Printf(b)


	result, err := loadConfig(fileName, l)
	if err != nil {
		os.Remove("testFile.json")
		t.Fatalf("Failed to load config")
	}

	var expected = &config{host: "1234", addr: "2"}
	if ok := reflect.DeepEqual(result, expected); !ok {
		t.Fatalf("config: %v, does not equal result: %v", result, expected)
	}
}
/*

/*
func TestMainHandler(t *testing.T) {
	var (
		l       = createLogger()
		client  = &http.Client{}
		handler = createHandler(l, client)
		ts      = httptest.NewServer(http.HandlerFunc(handler))
	)
}
*/

// testMainHandlerGet tests the GET componenet of the proxies client
func TestMainHandlerGet(t *testing.T) {
	var (
		l = createLogger()
	)

	var ts = httptest.NewServer(http.HandlerFunc(createHandler(l, nil)))
	defer ts.Close()

	c := ts.Client()
	_, err := c.Get(ts.URL)
	if err != nil {
		t.Fatalf("Unable to reach endpoint: %v", err.Error())
	}
}

func TestMainHandlerPost(t *testing.T) {
	var (
		l           = createLogger()
		contentType = "json"
		body        = bytes.NewBufferString("test json")
	)

	var ts = httptest.NewServer(http.HandlerFunc(createHandler(l, nil)))
	defer ts.Close()

	// var req = httptest.NewRequest("GET", ts.URL, nil)
	c := ts.Client()
	_, err := c.Post(ts.URL, contentType, body)
	if err != nil {
		t.Fatalf("Unable to reach endpoint: %v", err.Error())
	}
}

// The PUT method requests that the enclosed entity
func TestMainHandlerPut(t *testing.T) {
	var (
		l    = createLogger()
		body = bytes.NewBufferString("test json")
	)

	var ts = httptest.NewServer(http.HandlerFunc(createHandler(l, nil)))
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodPut, ts.URL, body)

	c := ts.Client()
	_, err := c.Do(req)
	if err != nil {
		t.Fatalf("Unable to reach endpoint: %v", err.Error())
	}
}
