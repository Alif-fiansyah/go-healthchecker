package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type CheckResult struct {
	URL        string
	StatusCode int
	Duration   time.Duration
	Err        error
}

func checkHealth(url string, client *http.Client, wg *sync.WaitGroup, results chan<- CheckResult) {
	defer wg.Done()

	start := time.Now()
	resp, err := client.Get(url)
	duration := time.Since(start)

	if err != nil {
		results <- CheckResult{
			URL:      url,
			Duration: duration,
			Err:      err,
		}
		return
	}
	defer resp.Body.Close()

	results <- CheckResult{
		URL:        url,
		StatusCode: resp.StatusCode,
		Duration:   duration,
		Err:        nil,
	}
}

func main() {
	filePath := "urls.txt"
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("[-] Gagal membaca file %s: %v\n", filePath, err)
		return
	}
	defer file.Close()

	var urls []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			urls = append(urls, line)
		}
	}

	if len(urls) == 0 {
		fmt.Println("[-] Tidak ada URL yang ditemukan di file urls.txt.")
		return
	}

	fmt.Printf("[+] Memeriksa %d URL secara bersamaan...\n\n", len(urls))

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	var wg sync.WaitGroup
	results := make(chan CheckResult, len(urls))

	startTotal := time.Now()

	for _, targetURL := range urls {
		wg.Add(1)
		go checkHealth(targetURL, client, &wg, results)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for res := range results {
		latency := res.Duration.Round(time.Millisecond)

		if res.Err != nil {
			fmt.Printf("[DOWN] %-40s | Error: Host unreachable/timeout (%v)\n", res.URL, latency)
		} else if res.StatusCode >= 200 && res.StatusCode < 300 {
			fmt.Printf("[UP]   %-40s | Status: %d OK (%v)\n", res.URL, res.StatusCode, latency)
		} else {
			fmt.Printf("[WARN] %-40s | Status: %d (%v)\n", res.URL, res.StatusCode, latency)
		}
	}

	fmt.Printf("\n[+] Selesai memeriksa seluruh URL dalam %v!\n", time.Since(startTotal).Round(time.Millisecond))
}
