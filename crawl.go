package main

import (
	"fmt"
	"sync"
	"time"
)

var blockingChannel = make(chan struct{}, 1)
var visitedSet = make(map[string]bool)
var waitGroup sync.WaitGroup

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page
	Fetch(url string) (body string, urls []string, err error)
}

type fakeResult struct {
	body string
	urls []string
}

type fakeFetcher map[string]*fakeResult

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"http://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"http://golang.org/pkg/",
			"http://golang.org/cmd/",
		},
	},
	"http://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"http://golang.org/",
			"http://golang.org/cmd/",
			"http://golang.org/pkg/fmt/",
			"http://golang.org/pkg/os/",
		},
	},
	"http://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
	"http://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	fmt.Printf("Fetching: %s\n", url)

	time.Sleep(500 * time.Millisecond)

	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}

	return "", nil, fmt.Errorf("not found: %s", url)
}

func Crawl(url string, depth int, fetcher Fetcher) {
	// TODO: Fetch URLs in parallel.
	// TODO: Don't fetch the same URL twice. (visited_set has to be a thread-safe data structure)
	if depth <= 0 {
		return
	}

	// For goroutines to put data and block (can only carry one datum; size = 1)
	// locks when there is data; unlocks when there is not any
	blockingChannel <- *new(struct{})

	if visited, _ := visitedSet[url]; !visited {
		body, urls, err := fetcher.Fetch(url)
		visitedSet[url] = true
		if err != nil {
			// 1. Print the error
			fmt.Println(err)
			// 2. Kill the goroutine (decrease counter by 1)
			waitGroup.Done()
			// 3. Unlock the channel (release its data)
			<-blockingChannel
			// 4. End the function by returning
			return
		}

		fmt.Printf("found: [depth:%d] %s %s\n", depth, url, body)
		for _, u := range urls {
			// Incraese the counter to wait for this new goroutine
			waitGroup.Add(1)
			go Crawl(u, depth-1, fetcher)
		}
	}

	waitGroup.Done()  // Decrease counter by 1
	<-blockingChannel // Unlock the channel
}

func main() {
	// Initialize the counter
	waitGroup.Add(1)

	Crawl("http://golang.org/", 4, fetcher)

	// Wait until all goroutines are done
	waitGroup.Wait()

	fmt.Println("Done!")

	for visitedURL := range visitedSet {
		fmt.Println(visitedURL)
	}
}
