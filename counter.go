package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

type Counter struct {
	MaxSize int
	Query   []byte

	size     int
	input    chan string
	inputWG  sync.WaitGroup
	output   chan int
	outputWG sync.WaitGroup
	total    int
}

func NewCounter(maxSize int, query []byte) *Counter {
	return &Counter{
		MaxSize: maxSize,
		Query:   query,
		input:   make(chan string),
		output:  make(chan int),
	}
}

func (c *Counter) Start() {
	c.outputWG.Add(1)

	go func() {
		defer c.outputWG.Done()

		for n := range c.output {
			c.total += n
		}
	}()
}

func (c *Counter) Add(url string) {
	if c.size < c.MaxSize {
		c.inputWG.Add(1)
		c.size++

		go func() {
			defer c.inputWG.Done()

			for url := range c.input {
				count, err := c.Process(url)
				if err != nil {
					fmt.Printf("Error processing the '%s'!\n", url)
					return
				}

				c.output <- count
			}
		}()
	}

	c.input <- url
}

func (c *Counter) Process(url string) (int, error) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error getting the '%s'!\n", url)
		return 0, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error getting the '%s'!\n", url)
		return 0, err
	}

	return bytes.Count(data, c.Query), nil
}

func (c *Counter) Stop() int {
	close(c.input)
	c.inputWG.Wait()

	close(c.output)
	c.outputWG.Wait()

	return c.total
}
