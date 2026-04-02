"use client";

import { useState } from "react";
import { createURL } from "@/lib/api";
import type { CreateURLResponse } from "@/types/api";

interface URLFormProps {
  onCreated: () => void;
}

export default function URLForm({ onCreated }: URLFormProps) {
  const [url, setUrl] = useState("");
  const [result, setResult] = useState<CreateURLResponse | null>(null);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  const [copied, setCopied] = useState(false);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError("");
    setResult(null);
    setLoading(true);

    try {
      const data = await createURL(url);
      setResult(data);
      setUrl("");
      onCreated();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to create URL");
    } finally {
      setLoading(false);
    }
  }

  async function handleCopy() {
    if (!result) return;
    await navigator.clipboard.writeText(result.short_url);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  }

  return (
    <div>
      <form onSubmit={handleSubmit} className="flex gap-3">
        <input
          type="url"
          value={url}
          onChange={(e) => setUrl(e.target.value)}
          placeholder="Paste a long URL here..."
          required
          className="flex-1 px-4 py-2.5 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
        />
        <button
          type="submit"
          disabled={loading}
          className="px-6 py-2.5 bg-blue-600 text-white font-medium rounded-lg hover:bg-blue-700 disabled:opacity-50 transition-colors"
        >
          {loading ? "Shortening..." : "Shorten"}
        </button>
      </form>

      {error && (
        <div className="mt-3 p-3 bg-red-50 text-red-700 rounded-lg text-sm">
          {error}
        </div>
      )}

      {result && (
        <div className="mt-3 p-4 bg-green-50 border border-green-200 rounded-lg flex items-center justify-between">
          <div>
            <p className="text-sm text-green-700 font-medium">
              URL shortened successfully!
            </p>
            <a
              href={result.short_url}
              target="_blank"
              rel="noopener noreferrer"
              className="text-blue-600 font-mono text-sm hover:underline"
            >
              {result.short_url}
            </a>
          </div>
          <button
            onClick={handleCopy}
            className="px-4 py-1.5 text-sm border border-green-300 rounded-md hover:bg-green-100 transition-colors"
          >
            {copied ? "Copied!" : "Copy"}
          </button>
        </div>
      )}
    </div>
  );
}
