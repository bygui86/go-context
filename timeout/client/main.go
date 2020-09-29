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

	timeoutMillisec = 1
)

func main() {
	// You still receive the cancel function in case you want to manually trigger that.
	// It works the same way as a normal context cancellation.
	ctx, cancel := context.WithTimeout(context.Background(), timeoutMillisec*time.Millisecond)
	defer cancel()

	executeQueryWithContext(url, ctx)
}

func executeQueryWithContext(url string, ctx context.Context) string {
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
