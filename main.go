package main

import (
	"fmt"
	"net/http"
	"os"
)

type result struct {
	i   int
	msg string
}

func checkUrl(url string, i int, c chan result) {
	_, err := http.Get(url)
	if err != nil {
		c <- result{i, "❌ " + url + " is down!"}
		return
	}
	c <- result{i, "✅ " + url + " is up!"}
}

func main() {
	c := make(chan result)

	for i, url := range URLs {
		go checkUrl(url, i, c)
	}

	results := make([]string, len(URLs))
	for i := 0; i < len(URLs); i++ {
		r := <-c
		results[r.i] = r.msg
	}

	for _, msg := range results {
		fmt.Println(msg)
	}

	if os.Getenv("DEV") == "true" {
		fmt.Println("(ran in DEV mode)")
	}
}