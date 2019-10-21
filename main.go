package main

import (
	"flag"
	"fmt"
	"strings"
)

const separator = ","

var (
	urls  = flag.String("urls", "", "comma separated list of URLs")
	k     = flag.Int("k", 5, "max degree of parallelism")
	query = flag.String("q", "", "query string")
)

func main() {
	flag.Parse()

	counter := NewCounter(*k, []byte(*query))
	counter.Start()

	for _, url := range strings.Split(*urls, separator) {
		counter.Add(url)
	}

	total := counter.Stop()
	fmt.Printf("Total: %d.\n", total)
}
