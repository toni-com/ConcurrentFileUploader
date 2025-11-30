# Concurrent File Uploader (Go)

A lightweight tool written in Go that uploads files to a remote server concurrently.

This project was built to learn about Go concurrency patterns, specifically focusing on how to manage multiple network requests efficiently using **Worker Pools** and **Channels**.

## Demo

**1. Concurrency in Action**
The application spins up a local mock server and processes a queue of files using a fixed set of workers.
<img width="331" height="322" alt="serverReceivingWorkerUploading" src="https://github.com/user-attachments/assets/7cf3765f-49c7-4798-b0cf-1a499142a395" />


**2. Final Report**
Once processing is complete, it generates a summary of the throughput and status.
<img width="424" height="232" alt="uploadSummary" src="https://github.com/user-attachments/assets/c0a74cf4-f7a8-42e8-80b4-f6c5e3a60718" />


## Learnings

* **Worker Pool:** Instead of spawning a goroutine for every single file, this projects uses a fixed number of workers (e.g., 10) to process a queue of tasks.
* **Channels:** Used for safe communication between the job generator, the workers, and the result collector.
* **WaitGroups:** Ensures the application waits for all background workers to finish before exiting.
