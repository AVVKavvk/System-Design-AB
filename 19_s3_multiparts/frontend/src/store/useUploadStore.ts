import { create } from "zustand";
import { persist } from "zustand/middleware";

interface UploadPart {
  ETag: string;
  PartNumber: number;
}

interface UploadState {
  uploadId: string | null;
  fileKey: string | null;
  completedParts: UploadPart[];
  setUploadInfo: (id: string, key: string) => void;
  addPart: (part: UploadPart) => void;
  reset: () => void;
}

export const useUploadStore = create<UploadState>()(
  persist(
    (set) => ({
      uploadId: null,
      fileKey: null,
      completedParts: [],
      setUploadInfo: (id, key) => set({ uploadId: id, fileKey: key }),
      addPart: (part) =>
        set((state) => ({
          completedParts: [...state.completedParts, part],
        })),
      reset: () => set({ uploadId: null, fileKey: null, completedParts: [] }),
    }),
    { name: "upload-storage" },
  ),
);
