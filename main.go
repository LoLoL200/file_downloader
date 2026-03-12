package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

type DownloadTask struct {
	URL      string
	FileName string
	Size     int64
}

type DownloadResult struct {
	Task         DownloadTask
	WorkerID     int
	Success      bool
	Error        error
	Duration     time.Duration
	ActualSize   int64
	DownloadedAt time.Time
}

type Downloader struct {
	mu          sync.Mutex
	results     []DownloadResult
	client      *http.Client
	workers     int
	downloadDir string
}

// Test URL-path

// Constructor
func NewDownloader(workers int) *Downloader {

	return &Downloader{
		workers:     workers,
		results:     []DownloadResult{},
		client:      &http.Client{Timeout: 30 * time.Second},
		downloadDir: "downloads",
	}
}

// Worker function
func (d *Downloader) worker(workerID int, taskChan <-chan DownloadTask, resultChan chan<- DownloadResult, wg *sync.WaitGroup) {
	defer wg.Done()

	// The main worker loop
	for task := range taskChan {
		result := d.downloadFile(task, workerID)
		resultChan <- result
	}
}

// Download Files
func (d *Downloader) downloadFile(task DownloadTask, workerID int) DownloadResult {

	// Boot start log
	fmt.Printf("⬇️ Worker %d: %s\n", workerID, task.FileName)

	// Remembering the start time
	start := time.Now()

	// The result structure is created
	result := DownloadResult{
		Task:     task,
		WorkerID: workerID,
	}
	// GET
	resp, err := d.client.Get(task.URL)
	if err != nil {
		result.Success = false
		result.Error = err
		return result
	}
	defer resp.Body.Close()

	// Create File
	file, err := os.Create(d.downloadDir + "/" + task.FileName)
	if err != nil {
		result.Success = false
		result.Error = err
		return result
	}
	defer file.Close()

	// Copy data
	size, err := io.Copy(file, resp.Body)

	// Copy error check
	if err != nil {
		result.Success = false
		result.Error = err
		return result
	}

	// Counting time
	duration := time.Since(start)

	// Successful download log
	fmt.Printf("✅ [Worker %d] %s (%.1fKB, %.1fs)\n",
		workerID,
		task.FileName,
		float64(size)/1024,
		duration.Seconds(),
	)

	// Filling in the result
	result.Success = true
	result.ActualSize = size
	result.Duration = time.Since(start)
	result.DownloadedAt = time.Now()

	return result
}

// Home function
func (d *Downloader) DownloadFiles(tasks []DownloadTask) error {
	taskChan := make(chan DownloadTask)
	resultChan := make(chan DownloadResult)

	var wg sync.WaitGroup

	// Start Worcker
	for i := 1; i <= d.workers; i++ {
		wg.Add(1)
		go d.worker(i, taskChan, resultChan, &wg)
	}

	// Sending tasks
	go func() {
		for _, task := range tasks {
			taskChan <- task
		}
		close(taskChan)
	}()

	// Getting results
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Reading with chenal "resultChan"
	for res := range resultChan {
		d.addResult(res)
	}

	return nil
}

// Safe add result
func (d *Downloader) addResult(result DownloadResult) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.results = append(d.results, result)
}
func main() {

	// List of tasks for download
	tasks := []DownloadTask{
		{"https://httpbin.org/json", "test_json.json", 1024},
		{"https://httpbin.org/xml", "test_xml.xml", 2048},
		{"https://jsonplaceholder.typicode.com/posts/1", "post_1.json", 1024},
		{"https://jsonplaceholder.typicode.com/users", "users.json", 8192},
		{"https://api.github.com/repos/golang/go", "golang_repo.json", 4096},
		{"https://httpbin.org/delay/2", "delayed_response.json", 1024},
	}

	// Number of workers
	workers := 3

	// Create downloader
	downl := NewDownloader(workers)

	// Information output
	fmt.Println("📥 ПАРАЛЕЛЬНИЙ DOWNLOADER")
	fmt.Println("=========================")
	fmt.Printf("📋 Завантажень: %d | Воркерів: %d\n\n", len(tasks), workers)

	fmt.Printf("🚀 Worker 1,%d запущені\n", workers)

	// We fix the start time
	start := time.Now()
	// Starting the downloads
	err := downl.DownloadFiles(tasks)
	if err != nil {
		fmt.Println("Download error:", err)
		return
	}

	// Getting results
	results := downl.results

	fmt.Println("\n🏁 Всі воркери завершили роботу")

	// Statistic
	var success int
	var failed int
	var totalSize int64

	// Statistics on worker
	workerStats := map[int]int{}

	// Analysis of results
	for _, r := range results {

		if r.Success {
			success++
			totalSize += r.ActualSize
			workerStats[r.WorkerID]++
		} else {
			failed++
		}
	}

	// Total time
	duration := time.Since(start).Seconds()

	// Download speed
	speed := float64(totalSize) / 1024 / duration

	// Statistics output
	fmt.Println("\n📊 СТАТИСТИКА:")
	fmt.Printf("✅ Успішно: %d | ❌ Помилки: %d\n", success, failed)
	fmt.Printf("📦 Розмір: %.1fKB | ⚡ Швидкість: %.1fKB/s\n", float64(totalSize)/1024, speed)

	fmt.Print("👷 ")

	// Statistics on worker
	for w, count := range workerStats {
		fmt.Printf("Worker %d: %d файли | ", w, count)
	}

	fmt.Println("\n📁 Збережено в: downloads/")
}

// 📥 ПАРАЛЕЛЬНИЙ DOWNLOADER
// =========================
// 📋 Завантажень: 6 | Воркерів: 3

// 🚀 Worker 1,2,3 запущені
// ⬇️ Worker 1: test_json.json
// ⬇️ Worker 2: test_xml.xml
// ✅ [Worker 1] test_json.json (0.5KB, 0.8s)
// ⬇️ Worker 1: users.json
// ✅ [Worker 2] test_xml.xml (1.2KB, 1.1s)
// 🏁 Всі воркери завершили роботу

// 📊 СТАТИСТИКА:
// ✅ Успішно: 5 | ❌ Помилки: 1
// 📦 Розмір: 12.3KB | ⚡ Швидкість: 8.2KB/s
// 👷 Worker 1: 2 файли | Worker 2: 2 файли | Worker 3: 1 файл
// 📁 Збережено в: downloads/
