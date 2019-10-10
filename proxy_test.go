package proxy

import (
	// "bytes"
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
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
var testServerFunc = func(t *testing.T) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var (
			b    bytes.Buffer
			_, _ = b.WriteString("Made it to third node")
		)

		switch r.Method {
		case "GET":
			t.Logf("\n\nReceived request: %v", r)
			rw.Write(b.Bytes())
		case "POST":
			t.Logf("\n\nReceived request: %v", r)
			rw.Write(b.Bytes())
		case "PUT":
			t.Logf("\n\nReceived request: %v", r)
			rw.Write(b.Bytes())
		case "HEAD":
			t.Logf("\n\nReceived request: %v", r)
			rw.Write(b.Bytes())
		default:
			t.Logf("\n\nReceived request: %v", r)
			rw.Write(b.Bytes())
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

func TestLoadConfig(t *testing.T) {
	var (
		l, _     = createLogger()
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

func TestMainHandler(t *testing.T) {
	var (
		l, _ = createLogger()
	)

	// create the test server, a server to test the proxy with
	var testServerNode3 = httptest.NewServer(testServerFunc(t))
	defer testServerNode3.Close()

	// create the proxy server.  this is the main handler of this application
	var proxyServerNode2 = httptest.NewServer(http.HandlerFunc(createHandler(l)))
	defer proxyServerNode2.Close()

	u2, _ := url.Parse(testServerNode3.URL)

	// add other requests here
	var requests = []*http.Request{
		&http.Request{RequestURI: "", URL: u2, Method: "GET"},
		&http.Request{RequestURI: "", URL: u2, Method: "POST"},
		&http.Request{RequestURI: "", URL: u2, Method: "PUT"},
		&http.Request{RequestURI: "", URL: u2, Method: "HEAD"},
	}

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

	for i := 0; i < len(requests); i++ {
		_, err := testClientNode1.Do(requests[i])
		if err != nil {
			t.Fatalf("Unable to reach endpoint.  method: %v, error: %v", requests[i].Method, err.Error())
		}
	}
}
