package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
	batch    int
	file     *os.File
}

func main() {
	var headers stringslice

	flag.Var(&headers, "H", "List of headers")

	url := flag.String("U", "", "Url")

	method := flag.String("X", "GET", "HTTP method")

	body := flag.String("d", "", "HTTP body")

	interval := flag.Int("interval", 0, "Gurl request interval")

	repeat := flag.Int("repeat", 1, "Gurl request repeat")

	batch := flag.Int("batch", 1, "Gurl request batch")

	file := flag.String("file", "", "Log file")

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

		g.batch = *batch

		if *file != "" {
			logFile, err := os.OpenFile(*file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

			if err != nil {
				panic(fmt.Sprintf("Error opening log file: %v", err))
			}

			g.file = logFile
		}

		if g.interval > 0 {
			g.requestTicker()
		} else {
			g.requestIterator()
		}
	}
}

func (g *gurl) requestTicker() {
	counter := 1
	g.requestIterator()

	ticker := time.NewTicker(time.Duration(g.interval) * time.Second)
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
			g.requestIterator()
			if counter >= g.repeat {
				close(quit)
			}
		}
	}()
}

func (g *gurl) requestIterator() {
	ch := make(chan string)
	for i := 0; i < g.batch; i++ {
		go g.request(ch)
	}

	if g.file != nil {
		log.SetOutput(g.file)
	}

	for i := 0; i < g.batch; i++ {
		log.Println(<-ch)
	}
}

func (g *gurl) request(ch chan<- string) {
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

	body, _ := ioutil.ReadAll(resp.Body)

	ch <- fmt.Sprintf("Gurl Request: \n\t Url: %v \n\t Status: %v \n\t Body: %v \n\n", g.url, resp.Status, string(body))
}
