"use client";

import { useEffect, useState } from "react";
import { listURLs, deleteURL } from "@/lib/api";
import type { ShortenedURL } from "@/types/api";

interface URLListProps {
  refreshKey: number;
}

export default function URLList({ refreshKey }: URLListProps) {
  const [urls, setUrls] = useState<ShortenedURL[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadURLs();
  }, [refreshKey]);

  async function loadURLs() {
    setLoading(true);
    try {
      const data = await listURLs(1, 50);
      setUrls(data.urls || []);
    } catch {
      setUrls([]);
    } finally {
      setLoading(false);
    }
  }

  async function handleDelete(id: number) {
    try {
      await deleteURL(id);
      setUrls((prev) => prev.filter((u) => u.id !== id));
    } catch {
      // silently fail
    }
  }

  if (loading) {
    return <p className="text-gray-500 text-sm py-8 text-center">Loading...</p>;
  }

  if (urls.length === 0) {
    return (
      <p className="text-gray-500 text-sm py-8 text-center">
        No URLs shortened yet. Create one above!
      </p>
    );
  }

  return (
    <div className="overflow-x-auto">
      <table className="w-full text-sm">
        <thead>
          <tr className="border-b border-gray-200 text-left text-gray-600">
            <th className="py-3 pr-4 font-medium">Short URL</th>
            <th className="py-3 pr-4 font-medium">Original URL</th>
            <th className="py-3 pr-4 font-medium text-right">Clicks</th>
            <th className="py-3 pr-4 font-medium">Created</th>
            <th className="py-3 font-medium"></th>
          </tr>
        </thead>
        <tbody>
          {urls.map((u) => (
            <tr key={u.id} className="border-b border-gray-100 hover:bg-gray-50">
              <td className="py-3 pr-4">
                <a
                  href={u.short_url}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-blue-600 font-mono text-xs hover:underline"
                >
                  {u.short_code}
                </a>
              </td>
              <td className="py-3 pr-4 max-w-xs truncate text-gray-700">
                {u.original_url}
              </td>
              <td className="py-3 pr-4 text-right font-mono">
                {u.click_count}
              </td>
              <td className="py-3 pr-4 text-gray-500">
                {new Date(u.created_at).toLocaleDateString()}
              </td>
              <td className="py-3 flex gap-2 justify-end">
                <a
                  href={`/dashboard/${u.id}`}
                  className="px-3 py-1 text-xs border border-gray-300 rounded hover:bg-gray-100 transition-colors"
                >
                  Analytics
                </a>
                <button
                  onClick={() => handleDelete(u.id)}
                  className="px-3 py-1 text-xs border border-red-200 text-red-600 rounded hover:bg-red-50 transition-colors"
                >
                  Delete
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
