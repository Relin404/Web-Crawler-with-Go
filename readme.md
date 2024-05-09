# Web Crawler

This project implements a simple web crawler written in Go. It crawls a starting URL and fetches its content and links recursively up to a specified depth.

## Requirements

- Go 1.21.5 or later (download Go)

## Running the project

### 1. Build the project

```Bash
go build
```

## 2. Run the crawler

```Bash
./webcrawler  <url> <depth>
```

- Replace `<url>` with the starting URL you want to crawl.
- Replace `<depth>` with the maximum depth of crawling (number of levels to follow links).

### Example

```Bash
./webcrawler <http://golang.org/> 2
```

This command will crawl the Go website (<http://golang.org/>) up to a depth of 2. The output will display the fetched content and links for each visited URL.

## Implementation Details

- The code uses a `fakeFetcher` for demonstration purposes. A real web crawler would need to implement a function that fetches actual web content.

- The `Crawl` function utilizes goroutines to fetch URLs concurrently.

- A `blockingChannel` and `WaitGroup` are used to synchronize access to shared resources *(visited URLs)* and ensure all goroutines finish before exiting.

## Further Development

- [ ] Implement a real web fetching function using packages like http.
- [ ] Add error handling for network issues.
- [ ] Improve concurrency management.
- [ ] Persist crawled data.
