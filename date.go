package main

// import (
// 	"net/http"
// 	"sync"
// 	"time"
// )

// type DownloadTask struct {
// 	URL      string
// 	FileName string
// 	Size     int64
// }

// type DownloadResult struct {
// 	Task         DownloadTask
// 	WorkerID     int
// 	Success      bool
// 	Error        error
// 	Duration     time.Duration
// 	ActualSize   int64
// 	DownloadedAt time.Time
// }

// type Downloader struct {
// 	mu          sync.Mutex
// 	results     []DownloadResult
// 	client      *http.Client
// 	workers     int
// 	downloadDir string
// }
