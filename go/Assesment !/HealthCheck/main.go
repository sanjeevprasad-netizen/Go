package main

import (
	"context"
	"net/http"
	"sync"
	"time"
)

type health_request struct {
	url   string
	retry int
}

type health_response struct {
	url           string
	sucess        bool
	response_time time.Duration
}

func main() {
	// In a real API (e.g. /health), this ctx comes from the HTTP request (r.Context())
	ctx := context.Background()
	healthcheck(ctx)
}

func healthcheck(ctx context.Context) {

	jobs := make(chan health_request, 10)
	result := make(chan health_response, 100)

	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go check_health(ctx, jobs, result, &wg)
	}

loop:
	for j := 0; j < 100; j++ {
		select {
		case <-ctx.Done():
			break loop

		case jobs <- health_request{
			url:   urls[j],
			retry: 0,
		}:
		}
	}

	close(jobs)

	for i := 0; i < 100; i++ {
		select {
		case <-ctx.Done():
			println("context cancelled, stopping")
			return
		case res_ult := <-result:
			println("result")
			println("url: ", res_ult.url)
			println("sucess: ", res_ult.sucess)
			println("response_time: ", res_ult.response_time.String())
		}
	}
	wg.Wait()

	close(result)
	return
}

func check_health(ctx context.Context, jobs chan health_request, result chan health_response, wg *sync.WaitGroup) {
	wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case job, ok := <-jobs:
			if !ok {
				return
			}
			handlerequest(job, ctx, jobs, result)
		}
	}
}

func handlerequest(curr_job health_request, ptrctx context.Context, jobs chan health_request, result chan health_response) {

	url := curr_job.url
	ctx, cancel := context.WithTimeout(ptrctx, 3*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return
	}

	start := time.Now()
	resp, err := http.DefaultClient.Do(req)
	duration := time.Since(start)

	if err != nil {
		// Request failed (DNS error, connection refused, timeout, etc.)
		retry(curr_job, result, jobs, ptrctx)
		return
	}
	defer resp.Body.Close() // Also prevent leaking connections

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		result <- health_response{
			url:           url,
			sucess:        true,
			response_time: duration,
		}
		return
	}

	retry(curr_job, result, jobs, ptrctx)

}

func retry(curr_job health_request, result chan health_response, jobs chan health_request, ptrctx context.Context) {

	if curr_job.retry <= 2 {
		curr_job.retry++
		handlerequest(curr_job, ptrctx, jobs, result)
	} else {
		result <- health_response{
			url:           curr_job.url,
			sucess:        false,
			response_time: 0,
		}
	}
	return
}

var urls = []string{
	// Reliable & Fast Endpoints (200 OK)
	"https://www.google.com",
	"https://www.github.com",
	"https://www.microsoft.com",
	"https://www.apple.com",
	"https://www.amazon.com",
	"https://www.wikipedia.org",
	"https://www.cloudflare.com",
	"https://www.mozilla.org",
	"https://www.stackoverflow.com",
	"https://go.dev",
	"https://pkg.go.dev",
	"https://golang.org",
	"https://www.reddit.com",
	"https://www.netflix.com",
	"https://www.spotify.com",
	"https://www.linkedin.com",
	"https://www.twitter.com",
	"https://www.youtube.com",
	"https://www.docker.com",
	"https://kubernetes.io",
	"https://www.rust-lang.org",
	"https://www.python.org",
	"https://nodejs.org",
	"https://www.typescriptlang.org",
	"https://react.dev",
	"https://vuejs.org",
	"https://angular.dev",
	"https://svelte.dev",
	"https://nextjs.org",
	"https://tailwindcss.com",
	"https://getbootstrap.com",
	"https://www.postgresql.org",
	"https://www.mysql.com",
	"https://www.mongodb.com",
	"https://redis.io",
	"https://www.elastic.co",
	"https://prometheus.io",
	"https://grafana.com",
	"https://www.nginx.com",
	"https://httpd.apache.org",
	"https://caddyserver.com",
	"https://www.kernel.org",
	"https://ubuntu.com",
	"https://www.debian.org",
	"https://archlinux.org",
	"https://www.redhat.com",
	"https://azure.microsoft.com",
	"https://aws.amazon.com",
	"https://cloud.google.com",
	"https://digitalocean.com",
	"https://www.linode.com",
	"https://www.vultr.com",
	"https://vercel.com",
	"https://netlify.com",
	"https://render.com",
	"https://fly.io",
	"https://gitlab.com",
	"https://bitbucket.org",
	"https://sourcegraph.com",
	"https://news.ycombinator.com",
	"https://lobste.rs",
	"https://dev.to",
	"https://hashnode.com",
	"https://medium.com",
	"https://archive.org",
	"https://www.eff.org",
	"https://www.w3.org",
	"https://developer.mozilla.org",
	"https://web.dev",
	"https://httpbin.org/get",
	"https://httpbin.org/status/200",
	"https://jsonplaceholder.typicode.com/posts/1",
	"https://jsonplaceholder.typicode.com/users",
	"https://dummyjson.com/products/1",
	"https://api.spacexdata.com/v4/launches/latest",
	"https://pokeapi.co/api/v2/pokemon/ditto",
	"https://api.github.com/zen",
	"https://catfact.ninja/fact",
	"https://official-joke-api.appspot.com/random_joke",
	"https://api.coindesk.com/v1/bpi/currentprice.json",

	// HTTP Status Error Endpoints (4xx / 5xx to test failure & retries)
	"https://httpbin.org/status/400",
	"https://httpbin.org/status/401",
	"https://httpbin.org/status/403",
	"https://httpbin.org/status/404",
	"https://httpbin.org/status/408",
	"https://httpbin.org/status/500",
	"https://httpbin.org/status/502",
	"https://httpbin.org/status/503",
	"https://httpbin.org/status/504",
	"https://jsonplaceholder.typicode.com/invalid-path-404",

	// Latency / Timeout Test Endpoints (to test the 3-second timeout requirement)
	"https://httpbin.org/delay/1",
	"https://httpbin.org/delay/2",
	"https://httpbin.org/delay/4",
	"https://httpbin.org/delay/5",

	// Unreachable & DNS Failure Endpoints (to test network errors & retries)
	"https://invalid-non-existent-domain-test-123.com",
	"https://dns-failure-example-xyz-987.net",
	"https://unreachable-host-testing-456.org",
	"https://fake-endpoint-health-check-fail.io",
	"http://192.0.2.1",
	"http://localhost:59999/does-not-exist",
}
