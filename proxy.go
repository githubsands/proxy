package proxy

import (
	"encoding/json"
	"log"
	"net/http"
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
