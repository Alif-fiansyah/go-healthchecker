# Concurrent HTTP Health Checker

A high-performance CLI utility written in Go to concurrently verify the availability, HTTP status codes, and latency of external endpoints and web services.

Built using Go's native concurrency primitives (`goroutines`, `channels`, and `sync.WaitGroup`), the tool evaluates multiple targets simultaneously without relying on external dependencies.

---

## Features

- **Concurrent Verification**: Executes health checks across dozens of URLs concurrently via lightweight Goroutines.
- **Detailed Latency Profiling**: Calculates round-trip response duration per endpoint rounded to the millisecond.
- **Fail-Safe Timeout**: Enforces a strict 5-second HTTP transport timeout to prevent hanging connections.
- **Zero Third-Party Dependencies**: Implemented purely using Go's standard library (`net/http`, `sync`, `bufio`).
- **Configurable Target List**: Dynamically parses targets from a plain-text file, ignoring empty lines and comments.

---

## Prerequisites

- **Go**: Version 1.18 or newer.

Verify your installation:
```bash
go version
```

---

## Installation & Setup

1. **Clone the repository:**
   ```bash
   git clone [https://github.com/Alif-fiansyah/go-healthchecker.git](https://github.com/Alif-fiansyah/go-healthchecker.git)
   cd go-healthchecker
   ```

2. **Configure target URLs:**
   Add target endpoints inside `urls.txt`:
   ```text
   [https://google.com](https://google.com)
   [https://github.com](https://github.com)
   [https://httpbin.org/status/404](https://httpbin.org/status/404)
   [https://httpbin.org/status/500](https://httpbin.org/status/500)
   ```

3. **Run the utility:**
   ```bash
   go run main.go
   ```

4. **Optional: Build a Standalone Binary**
   ```bash
   go build -o healthchecker main.go
   ./healthchecker
   ```

---

## Sample Output

```text
[+] Memeriksa 6 URL secara bersamaan...

[UP]   [https://google.com](https://google.com)                       | Status: 200 OK (78ms)
[UP]   [https://github.com](https://github.com)                       | Status: 200 OK (142ms)
[WARN] [https://httpbin.org/status/404](https://httpbin.org/status/404)           | Status: 404 (310ms)
[WARN] [https://httpbin.org/status/500](https://httpbin.org/status/500)           | Status: 500 (298ms)
[DOWN] [https://domain-pasti-tidak-ada-12345.org](https://domain-pasti-tidak-ada-12345.org) | Error: Host unreachable/timeout (12ms)

[+] Selesai memeriksa seluruh URL dalam 315ms!
```

---

## Architecture Overview

- **`net/http.Client`**: Configured with strict timeouts to safeguard routine pools.
- **`sync.WaitGroup`**: Orchestrates worker synchronization and guarantees complete execution before output termination.
- **Buffered Channels**: Streams asynchronous diagnostic payloads safely across concurrent execution contexts.

---

## License

This project is licensed under the [MIT License](LICENSE).
