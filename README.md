# 📥 Parallel File Downloader (Go)

A simple **parallel file downloader written in Go**.
The program downloads multiple files simultaneously using **goroutines, channels, and a worker pool pattern**.

This project demonstrates basic **Go concurrency**, HTTP requests, and file handling.

---

# ⚙️ Features

* 📥 Download files from multiple URLs
* ⚡ Parallel downloading with worker goroutines
* 📊 Download statistics (speed, size, success/fail)
* 👷 Worker pool architecture
* 🧵 Uses Go **channels** and **goroutines**
* 📁 Automatically saves files to the `downloads/` directory

---

# 🧩 Main Components

### `Downloader`

Main structure that manages:

* worker goroutines
* HTTP client
* download results

### `worker()`

A goroutine that:

1. receives tasks from a channel
2. downloads files
3. sends results to the result channel

### `downloadFile()`

Handles downloading:

* sends an HTTP GET request
* saves the file locally
* measures download time and size

### `DownloadFiles()`

The main method that:

* starts workers
* distributes tasks
* collects results

---

# 🚀 Installation & Run

Clone the repository:

```bash
git clone https://github.com/LoLoL200/file_downloader.git
```

Go to the project folder:

```bash
cd file_downloader
```

Run the program:

```bash
go run main.go
```

---

# 📊 Example Output

```
📥 PARALLEL DOWNLOADER
======================

🚀 Workers started

⬇️ Worker 1: test_json.json
⬇️ Worker 2: test_xml.xml
⬇️ Worker 3: post_1.json

✅ [Worker 1] test_json.json (0.5KB, 0.8s)

🏁 All workers finished

📊 STATISTICS
✅ Success: 5 | ❌ Failed: 1
📦 Size: 12.3KB | ⚡ Speed: 8.2KB/s
```

---

# 🛠 Technologies

* Go (Golang)
* Goroutines
* Channels
* Worker Pool pattern
* HTTP Client
* File I/O

---

# 📁 Output

Downloaded files are saved in:

```
downloads/
```

---

# 📚 Purpose

This project is a **learning example** for understanding:

* Go concurrency
* worker pool patterns
* channels and goroutines
* parallel file downloading
