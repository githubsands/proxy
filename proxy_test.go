package main

import (
	// "bytes"
	"fmt"
	// "io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	// "reflect"
	"bytes"
	"net/url"
	"testing"
)

var createFileFunc = func(fileName string) (*os.File, error) {
	fi, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}

	return fi, nil
}

// testServerFunc returns the server the proxy handle will be proxying for test purpose
var testServerFunc = func() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var (
			b bytes.Buffer
			// writer = bufio.NewWriter(&b)
		)
		fmt.Print("Made it to third node")
		_, err := b.WriteString("test")
		if err != nil {
			fmt.Printf("Error occured: %v", err)
		}

		switch r.Method {
		case "GET":

			fmt.Printf("\n\nReceived request: %v", r)
			rw.Write(b.Bytes())
		/*

			case "POST":


			case "PUT":


			case "HEAD":

		*/
		default:
			fmt.Printf("Given http method has yet to been programmed")
		}
	}
}

func TestCreateLogger(t *testing.T) {
	_, err := createLogger()
	if err != nil {
		t.Fatalf("Unable to create logger")
	}

	_, err = os.Open("log")
	if err != nil {
		t.Fatalf("File not made under filename: log")
	}

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

	// _, _ := ioutil.ReadAll(f)

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
*/

// testMainHandlerGet tests the GET componenet of the proxies client
func TestMainHandlerGet(t *testing.T) {
	var (
		l, _ = createLogger()
	)

	// create the test server, a server to test the proxy with
	var testServerNode3 = httptest.NewServer(testServerFunc())
	defer testServerNode3.Close()

	// create the proxy server.  this is the main handler of this application
	var proxyServerNode2 = httptest.NewServer(http.HandlerFunc(createHandler(l)))
	defer proxyServerNode2.Close()

	var request = &http.Request{}
	request.RequestURI = ""
	u2, _ := url.Parse(testServerNode3.URL)
	request.URL = u2

	// add proxy server's url to our test server's proxy
	//
	// to use this outside this testing environment use your http browsers options if allowed or
	// programmatically on the os layer using linux iptables

	var testClientNode1 = new(http.Client)
	testClientNode1.Transport = &http.Transport{Proxy: func(r *http.Request) (*url.URL, error) {
		u, err := url.Parse(proxyServerNode2.URL)
		if err != nil {
			return nil, err
		}

		return u, nil
	},
	}

	_, err := testClientNode1.Do(request)
	if err != nil {
		t.Fatalf("Unable to reach endpoint: %v", err.Error())
	}
}

/*
func TestMainHandlerPost(t *testing.T) {
	var (
		l           = createLogger()
		contentType = "json"
		body        = bytes.NewBufferString("test json")
	)

	var ts = httptest.NewServer(http.HandlerFunc(createHandler(l)))
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

	var ts = httptest.NewServer(http.HandlerFunc(createHandler(l)))
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodPut, ts.URL, body)

	c := ts.Client()
	_, err := c.Do(req)
	if err != nil {
		t.Fatalf("Unable to reach endpoint: %v", err.Error())
	}
}

// TestMainHandlerHead does not return a response
func TestMainHandlerHead(t *testing.T) {
	var (
		l = createLogger()
	)

	var ts = httptest.NewServer(http.HandlerFunc(createHandler(l)))
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodHead, ts.URL, nil)

	spew.Sprintf("Dumping: %v", req)
	spew.Dump("%v", req)

	c := ts.Client()
	_, err := c.Do(req)
	if err != nil {
		t.Fatalf("Unable to reach endpoint: %v", err.Error())
	}

	spew.Sprintf("Dumping: %v", l)
}
*/

/*
func TestRecognizeContentType(t *testing.T) {

	// create test json response
	jsonRes := &http.Response{}
	jsonRes.Body =

	var testList = []*http.Response{
		&http.Response{}
	}


		= recognizeContentType(

}
*/
