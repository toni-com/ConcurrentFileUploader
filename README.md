# Concurrent File Uploader (Go)

A lightweight tool written in Go that uploads files to a remote server concurrently.

This project was built to learn about Go concurrency patterns, specifically focusing on how to manage multiple network requests efficiently using **Worker Pools** and **Channels**.

## Demo

| Concurrency in Action | Final Report |
| :---: | :---: |
|<img width="331" height="322" alt="serverReceivingWorkerUploading" src="https://github.com/user-attachments/assets/7cf3765f-49c7-4798-b0cf-1a499142a395" /> | <img width="378" height="232" alt="uploadSummary" src="https://github.com/user-attachments/assets/ebc2dce5-3e18-42e9-9ce2-186ef7c307df" /> |
| The application spins up a local mock server and processes a queue of files using a fixed set of workers. | Once processing is complete, it generates a summary of the throughput and status. |

## Learnings

* **Worker Pool:** Instead of spawning a goroutine for every single file, this projects uses a fixed number of workers (e.g., 10) to process a queue of tasks.
* **Channels:** Used for safe communication between the job generator, the workers, and the result collector.
* **WaitGroups:** Ensures the application waits for all background workers to finish before exiting.
