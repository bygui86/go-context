package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	neturl "net/url"
	"time"
)

const (
	url = "http://localhost:8080/ishealthy"

	hedgedTimeoutMillisec = 1
)

func main() {
	queryWithHedgedRequestsWithContext([]string{url, url, url, url, url})

	time.Sleep(3 * time.Second)
}

/*
	Hedged Requests

	We fire a request to an external service, if the response doesn't come back within our defined timeout,
	we issue a second request.
	When a response comes back, all other requests are cancelled.
*/
func queryWithHedgedRequestsWithContext(urls []string) string {
	ch := make(chan string, len(urls))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for _, url := range urls {
		go func(u string, c chan string) {
			c <- executeQueryWithContext(u, ctx)
		}(url, ch)

		select {
		case r := <-ch:
			fmt.Println("Cancellation!")
			cancel()
			return r
		case <-time.After(hedgedTimeoutMillisec * time.Millisecond):
			fmt.Println("Timeout!")
		}
	}
	return <-ch
}

func executeQueryWithContext(url string, ctx context.Context) string {
	fmt.Printf("New request to %s\n", url)

	start := time.Now()
	parsedURL, _ := neturl.Parse(url)
	req := &http.Request{URL: parsedURL}
	req = req.WithContext(ctx)

	response, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Println(err.Error())
		return err.Error()
	}

	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Printf("Request time: %d ms from url%s\n", time.Since(start).Nanoseconds()/time.Millisecond.Nanoseconds(), url)
	return fmt.Sprintf("%s from %s", body, url)
}
