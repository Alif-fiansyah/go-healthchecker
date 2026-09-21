package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
)

type CheckResult struct {
	URL        string
	StatusCode int
	Duration   time.Duration
	Err        error
}

func worker(id int, jobs <-chan string, results chan<- CheckResult, client *http.Client, wg *sync.WaitGroup) {
	defer wg.Done()
	for targetURL := range jobs {
		start := time.Now()
		resp, err := client.Get(targetURL)
		duration := time.Since(start)

		if err != nil {
			results <- CheckResult{
				URL:      targetURL,
				Duration: duration,
				Err:      err,
			}
			continue
		}
		resp.Body.Close()

		results <- CheckResult{
			URL:        targetURL,
			StatusCode: resp.StatusCode,
			Duration:   duration,
			Err:        nil,
		}
	}
}

func main() {
	filePath := flag.String("f", "urls.txt", "Path to file containing URLs")
	workersCount := flag.Int("c", 10, "Number of concurrent workers")
	timeoutSec := flag.Int("t", 5, "HTTP request timeout in seconds")
	flag.Parse()

	file, err := os.Open(*filePath)
	if err != nil {
		fmt.Printf("%s[-] Gagal membaca file %s: %v%s\n", colorRed, *filePath, err, colorReset)
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
		fmt.Printf("%s[-] Tidak ada URL valid di file %s.%s\n", colorYellow, *filePath, colorReset)
		return
	}

	fmt.Printf("%s[+] Memeriksa %d URL menggunakan %d concurrent worker (Timeout: %ds)...%s\n\n",
		colorCyan, len(urls), *workersCount, *timeoutSec, colorReset)

	client := &http.Client{
		Timeout: time.Duration(*timeoutSec) * time.Second,
	}

	jobs := make(chan string, len(urls))
	results := make(chan CheckResult, len(urls))

	var wg sync.WaitGroup

	// Menjalankan pool worker
	for w := 1; w <= *workersCount; w++ {
		wg.Add(1)
		go worker(w, jobs, results, client, &wg)
	}

	// Mengirim tugas ke channel jobs
	for _, targetURL := range urls {
		jobs <- targetURL
	}
	close(jobs)

	// Menutup channel results setelah seluruh worker selesai
	go func() {
		wg.Wait()
		close(results)
	}()

	startTotal := time.Now()

	for res := range results {
		latency := res.Duration.Round(time.Millisecond)

		if res.Err != nil {
			fmt.Printf("%s[DOWN]%s %-40s | Error: Host unreachable/timeout (%v)\n",
				colorRed, colorReset, res.URL, latency)
		} else if res.StatusCode >= 200 && res.StatusCode < 300 {
			fmt.Printf("%s[UP]%s   %-40s | Status: %d OK (%v)\n",
				colorGreen, colorReset, res.URL, res.StatusCode, latency)
		} else {
			fmt.Printf("%s[WARN]%s %-40s | Status: %d (%v)\n",
				colorYellow, colorReset, res.URL, res.StatusCode, latency)
		}
	}

	fmt.Printf("\n%s[+] Selesai memeriksa seluruh URL dalam %v!%s\n",
		colorCyan, time.Since(startTotal).Round(time.Millisecond), colorReset)
}