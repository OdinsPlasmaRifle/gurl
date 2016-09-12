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

	interval := flag.Int("gi", 0, "Gurl request interval")

	repeat := flag.Int("gr", 0, "Gurl request repeat")

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

		g.request()

		if g.interval > 0 && g.repeat > 0 {
			g.ticker()
		}
	}
}

func (g *gurl) ticker() {
	// quit := make(chan bool)

	// go func() {
	// 	ticker := time.NewTicker(time.Second * time.Duration(g.interval))
	// 	counter := 0

	// 	for {
	// 		select {
	// 		case <-ticker.C:
	// 			g.request()
	// 			counter++
	// 		case <-quit:
	// 			ticker.Stop()
	// 			return
	// 		}
	// 	}
	// }()

	// quit <- true

	timeChan := make(chan bool)

	go func() {
		<-time.After(2 * time.Hour)
		close(timeChan)
	}()

	t := time.NewTicker(time.Duration(g.interval) * time.Second)

	func() {
		for {
			select {
			case <-t.C:
			case <-timeChan:
				t.Stop()
				return
			}
			g.request()
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

	log.Println("response Status:", resp.Status)

	body, _ := ioutil.ReadAll(resp.Body)

	log.Println("response Body:", string(body))
}

// curl -H "Content-Type: application/json" -X POST -d '{"email" : "test@test.com", "password": "123"}' -H "Origin: http://example.com" --verbose http://localhost:8080/user/register
