# Resumable S3 Multipart Upload System

A robust, "fire-and-forget" solution for uploading massive files (up to 100GB) from a web browser to AWS S3. Built with **Golang (Echo)** for the orchestration layer and **React (Vite/Bun/Zustand)** for the frontend.

## 🚀 The Problem

Standard HTTP uploads fail for large files because:

1. **Network Instability:** A 1GB+ upload is likely to flicker; if it drops at 99%, you start from 0%.
2. **Browser Limits:** Most browsers/servers have a 2GB-5GB limit for single-stream POST requests.
3. **Memory Issues:** Reading a 100GB file into browser memory causes a crash.
4. **User Experience:** If the user refreshes the page, the upload progress is lost.

## 🛠 The Solution: S3 Multipart Upload

This project implements a **Multipart Upload strategy** which slices the file into small, manageable chunks.

### Key Features

- **Chunked Slicing:** Uses JavaScript `Blob.slice()` to read file parts without high RAM usage.
- **Parallel Uploads:** Chunks are sent directly to S3 via **Presigned URLs** (Backend never touches the bytes).
- **Refresh Resilience:** Uses **Zustand + LocalStorage** to persist the `UploadID` and `ETags`. If the user refreshes, the app resumes from the last successful chunk.
- **Secure by Design:** The S3 bucket remains private. Temporary access is granted per-chunk via expiring signatures.

---

## 🏗 System Architecture

1. **Initiate:** Frontend asks Backend for an `UploadID`.
2. **Presign:** Backend generates a unique URL for a specific `PartNumber`.
3. **Upload:** Frontend `PUT`s the file slice directly to S3 and receives an `ETag` (checksum).
4. **Complete:** Frontend sends all `ETags` + `PartNumbers` to Backend. Backend tells S3 to stitch the file together.

---

## 🧩 Core Concepts

### 1. PartNumber

Each chunk is assigned an index (1, 2, 3...). S3 uses these to know the exact sequence, allowing us to upload parts in parallel or out of order.

### 2. ETag (The "Receipt")

When a chunk reaches S3, S3 returns an **ETag** (Entity Tag). This is a unique MD5 hash of that chunk. We store this in `localStorage`. To "Complete" the upload, we must provide the correct ETag for every PartNumber.

---

## 💻 Tech Stack

- **Frontend:** React, Vite, Bun (Runtime), Zustand (State Persistence), Tailwind CSS, Axios.
- **Backend:** Golang, Echo Framework, AWS SDK v2.
- **Infrastructure:** AWS S3.

---

## ⚙️ Setup & Configuration

### Backend (.env)

```env
AWS_ACCESS_KEY_ID=your_key
AWS_SECRET_ACCESS_KEY=your_secret
AWS_REGION=ap-south-1
S3_BUCKET_NAME=your-bucket
PORT=8080
```

### Frontend (.env)

```env
VITE_API_BASE_URL=http://localhost:8080
```

## S3 CORS Policy

To allow the browser to talk to S3 and read ETags, the following CORS configuration is required:

```json
[
  {
    "AllowedHeaders": ["*"],
    "AllowedMethods": ["PUT", "POST", "GET"],
    "AllowedOrigins": ["http://localhost:5173"],
    "ExposeHeaders": ["ETag"]
  }
]
```

## Why This Matters

By moving the logic from a "Single Stream" to a "Multipart Managed" state, we achieve **Enterprise-grade reliability**. This architecture is exactly how platforms like Dropbox, Google Drive, and YouTube handle large video uploads safely.
