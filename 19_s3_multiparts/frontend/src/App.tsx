import React, { useState } from "react";
import axios from "axios";
import { useUploadStore } from "./store/useUploadStore";

// const CHUNK_SIZE = 1 * 1024 * 1024; // 1MB per chunk
const CHUNK_SIZE = 10 * 1024 * 1024; // 10MB per chunk
const API_BASE = "http://localhost:8080";

export default function App() {
  const [file, setFile] = useState<File | null>(null);
  const [progress, setProgress] = useState(0);
  const { uploadId, fileKey, completedParts, setUploadInfo, addPart, reset } =
    useUploadStore();

  const uploadFile = async () => {
    if (!file) return;

    let currentUploadId = uploadId;
    let currentKey = fileKey;

    // 1. Initiate if not already started
    if (!currentUploadId) {
      const res = await axios.post(
        `${API_BASE}/initiate-upload?filename=${file.name}`,
      );
      currentUploadId = res.data.uploadId;
      currentKey = res.data.key;
      setUploadInfo(currentUploadId!, currentKey!);
    }

    const totalParts = Math.ceil(file.size / CHUNK_SIZE);

    // 2. Loop through chunks
    for (let i = 1; i <= totalParts; i++) {
      // Check if part already uploaded (Refresh Resume logic)
      if (completedParts.find((p) => p.PartNumber === i)) continue;

      const start = (i - 1) * CHUNK_SIZE;
      const end = Math.min(start + CHUNK_SIZE, file.size);
      const chunk = file.slice(start, end);

      // Get presigned URL for this part
      const { data } = await axios.post(`${API_BASE}/get-presigned-url`, {
        uploadId: currentUploadId,
        key: currentKey,
        partNumber: i,
      });

      // Upload chunk directly to S3
      const uploadRes = await axios.put(data.url, chunk, {
        headers: { "Content-Type": file.type },
      });

      const etag = uploadRes.headers.etag;
      addPart({ ETag: etag, PartNumber: i });
      setProgress(Math.round((i / totalParts) * 100));
    }

    // 3. Complete Upload
    await axios.post(`${API_BASE}/complete-upload`, {
      uploadId: currentUploadId,
      key: currentKey,
      parts: useUploadStore.getState().completedParts,
    });

    alert("Upload Complete!");
    reset();
    setProgress(0);
  };

  return (
    <div className="min-h-screen bg-gray-100 flex items-center justify-center p-5">
      <div className="bg-white p-8 rounded-lg shadow-xl w-full max-w-md">
        <h1 className="text-2xl font-bold mb-4">Resumable 100GB Upload</h1>
        <input
          type="file"
          onChange={(e) => setFile(e.target.files?.[0] || null)}
          className="mb-4 block w-full text-sm text-gray-500"
        />

        {progress > 0 && (
          <div className="w-full bg-gray-200 rounded-full h-2.5 mb-4">
            <div
              className="bg-blue-600 h-2.5 rounded-full"
              style={{ width: `${progress}%` }}
            ></div>
          </div>
        )}

        <button
          onClick={uploadFile}
          className="w-full bg-blue-600 text-white py-2 rounded hover:bg-blue-700 transition"
        >
          {uploadId ? "Resume Upload" : "Start Upload"}
        </button>

        {uploadId && (
          <p className="text-xs text-orange-500 mt-2 text-center">
            Upload in progress... will resume if refreshed.
          </p>
        )}
      </div>
    </div>
  );
}
