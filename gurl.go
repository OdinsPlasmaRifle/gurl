package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type stringslice []string

func (s *stringslice) String() string {
	return fmt.Sprintf("%d", *s)
}

func (s *stringslice) Set(value string) error {
	*s = append(*s, value)
	return nil
}

type gurl struct {
	headers  stringslice
	url      string
	method   string
	body     []byte
	interval int
	repeat   int
}

func main() {
	var headers stringslice

	flag.Var(&headers, "H", "List of headers")

	url := flag.String("U", "", "Url")

	method := flag.String("X", "GET", "HTTP method")

	body := flag.String("d", "", "HTTP body")

	interval := flag.Int("interval", 0, "Gurl request interval")

	repeat := flag.Int("repeat", 0, "Gurl request repeat")

	flag.Parse()

	if flag.NFlag() == 0 {
		flag.PrintDefaults()
	} else {
		g := gurl{}

		g.headers = headers

		g.url = *url

		g.method = *method

		g.body = []byte(*body)

		g.interval = *interval

		g.repeat = *repeat

		if g.interval > 0 && g.repeat > 0 {
			g.ticker()
		} else {
			g.request()
		}
	}
}

func (g *gurl) ticker() {
	counter := 1
	g.request()

	ticker := time.NewTicker(time.Second * time.Duration(g.interval))
	quit := make(chan struct{})

	func() {
		for {
			select {
			case <-ticker.C:
			case <-quit:
				ticker.Stop()
				return
			}
			counter++
			g.request()
			if counter >= g.repeat {
				close(quit)
			}
		}
	}()
}

func (g *gurl) request() {
	log.Println("Request URL: ", g.url)

	req, err := http.NewRequest(g.method, g.url, bytes.NewBuffer(g.body))

	for i := 0; i < len(g.headers); i++ {
		split := strings.Split(g.headers[i], ": ")
		req.Header.Set(split[0], split[1])
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	log.Println("Response Status:", resp.Status)

	body, _ := ioutil.ReadAll(resp.Body)

	log.Println("Response Body:", string(body))
}

// ./gurl -U="http://requestb.in/1bkrics1" -X="GET" -d="{'hello':'hello'}" -H="Test: 123" -interval=2 -repeat=2
