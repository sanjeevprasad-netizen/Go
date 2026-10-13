package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

func main() {
	fmt.Println("Hello")
	//Print1to10Go()

	//useMutex()
	//openConnectionwithtimeout()

	userWorker()
}

type result struct {
	url    string
	status int
	error  string
}

type site struct {
	url string
}

func userWorker() {
	job := make(chan site, 3)
	results := make(chan result, 3)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	for i := 1; i <= 2; i++ {
		go workercraw(i, job, results, client)
	}

	sites := []site{
		{"https://google.com"},
		{"https://example.com"},
		{"https://abc.com"},
		{"https://notfound.com"},
	}
	for i := 0; i < len(sites); i++ {
		job <- sites[i]
	}
	close(job)

	for i := 0; i < len(sites); i++ {
		result := <-results
		fmt.Println(result)
	}

}

func workercraw(workerID int, job chan site, results chan result, client *http.Client) {
	for j := range job {
		fmt.Println("worker", workerID, "is fetching", j.url)
		resp, err := client.Get(j.url)

		if err != nil {
			results <- result{
				url:   j.url,
				error: err.Error(),
			}
			continue
		}

		resp.Body.Close()
		results <- result{
			url:    j.url,
			status: resp.StatusCode,
		}

	}
}
func printGoid(i int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("go id %d \n", i)

}

func Print1to10Go() {
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go printGoid(i, &wg)
	}

	wg.Wait()
}

func openconnect(done chan bool) {
	value := rand.Intn(10) + 1

	if value < 5 {
		fmt.Println("connection successfully")
	} else {
		fmt.Println("oops connection hanging")
		time.Sleep(5 * time.Second)
		fmt.Println("after 5 s")

	}
	done <- true
}

func openConnectionwithtimeout() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	done := make(chan bool, 1)

	go openconnect(done)

	select {
	case <-ctx.Done():
		fmt.Println("connection timeout")
	case <-done:
		fmt.Println("connection successfully , done")
	}

}

func useMutex() {
	var counter int
	var wg sync.WaitGroup
	var mu sync.Mutex

	wg.Add(2)

	for i := 0; i < 2; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				mu.Lock()
				counter++
				mu.Unlock()
			}
		}()
	}

	wg.Wait()
	fmt.Println(counter)
}
