"use client";

import { useState } from "react";
import URLForm from "@/components/url-form";
import URLList from "@/components/url-list";

export default function Home() {
  const [refreshKey, setRefreshKey] = useState(0);

  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-2xl font-bold mb-2">Shorten a URL</h1>
        <p className="text-gray-600 mb-4">
          Paste a long URL to get a short, shareable link with click analytics.
        </p>
        <URLForm onCreated={() => setRefreshKey((k) => k + 1)} />
      </div>

      <div>
        <h2 className="text-lg font-semibold mb-4">Your URLs</h2>
        <div className="bg-white rounded-lg border border-gray-200 p-4">
          <URLList refreshKey={refreshKey} />
        </div>
      </div>
    </div>
  );
}
