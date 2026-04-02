"use client";

import { useEffect, useState, use } from "react";
import { getURL, getClicksOverTime, getSources } from "@/lib/api";
import type {
  ShortenedURL,
  ClicksOverTimeResponse,
  SourcesResponse,
} from "@/types/api";
import ClicksOverTimeChart from "@/components/charts/clicks-over-time";
import BrowserBreakdownChart from "@/components/charts/browser-breakdown";
import CountryBreakdownChart from "@/components/charts/country-breakdown";

const PERIODS = [
  { value: "24h", label: "24 Hours" },
  { value: "7d", label: "7 Days" },
  { value: "30d", label: "30 Days" },
  { value: "90d", label: "90 Days" },
];

export default function DashboardPage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const { id: rawId } = use(params);
  const id = Number(rawId);

  const [url, setUrl] = useState<ShortenedURL | null>(null);
  const [clicks, setClicks] = useState<ClicksOverTimeResponse | null>(null);
  const [sources, setSources] = useState<SourcesResponse | null>(null);
  const [period, setPeriod] = useState("7d");
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    loadData();
  }, [id, period]);

  async function loadData() {
    setLoading(true);
    setError("");
    try {
      const [urlData, clicksData, sourcesData] = await Promise.all([
        getURL(id),
        getClicksOverTime(id, period),
        getSources(id, period),
      ]);
      setUrl(urlData);
      setClicks(clicksData);
      setSources(sourcesData);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to load analytics");
    } finally {
      setLoading(false);
    }
  }

  if (error) {
    return (
      <div className="p-6 text-center">
        <p className="text-red-600">{error}</p>
        <a href="/" className="text-blue-600 hover:underline text-sm mt-2 inline-block">
          Back to home
        </a>
      </div>
    );
  }

  if (loading && !url) {
    return <p className="text-gray-500 text-sm py-12 text-center">Loading analytics...</p>;
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-start justify-between">
        <div>
          <a href="/" className="text-sm text-blue-600 hover:underline">
            &larr; Back
          </a>
          <h1 className="text-2xl font-bold mt-1">Analytics Dashboard</h1>
          {url && (
            <p className="text-gray-500 text-sm mt-1 truncate max-w-lg">
              {url.original_url}
            </p>
          )}
        </div>

        {/* Period selector */}
        <div className="flex gap-1 bg-gray-100 rounded-lg p-1">
          {PERIODS.map((p) => (
            <button
              key={p.value}
              onClick={() => setPeriod(p.value)}
              className={`px-3 py-1.5 text-sm rounded-md transition-colors ${
                period === p.value
                  ? "bg-white font-medium shadow-sm"
                  : "text-gray-600 hover:text-gray-900"
              }`}
            >
              {p.label}
            </button>
          ))}
        </div>
      </div>

      {/* Summary cards */}
      {url && clicks && (
        <div className="grid grid-cols-3 gap-4">
          <div className="bg-white rounded-lg border border-gray-200 p-4">
            <p className="text-sm text-gray-500">Total Clicks</p>
            <p className="text-3xl font-bold mt-1">{clicks.total_clicks}</p>
          </div>
          <div className="bg-white rounded-lg border border-gray-200 p-4">
            <p className="text-sm text-gray-500">All-Time Clicks</p>
            <p className="text-3xl font-bold mt-1">{url.click_count}</p>
          </div>
          <div className="bg-white rounded-lg border border-gray-200 p-4">
            <p className="text-sm text-gray-500">Short Code</p>
            <p className="text-3xl font-bold mt-1 font-mono">{url.short_code}</p>
          </div>
        </div>
      )}

      {/* Clicks over time chart */}
      <div className="bg-white rounded-lg border border-gray-200 p-6">
        <h2 className="text-lg font-semibold mb-4">Clicks Over Time</h2>
        {clicks && <ClicksOverTimeChart data={clicks.clicks_per_day} />}
      </div>

      {/* Browser & Country breakdown */}
      <div className="grid grid-cols-2 gap-4">
        <div className="bg-white rounded-lg border border-gray-200 p-6">
          <h2 className="text-lg font-semibold mb-4">Browsers</h2>
          {sources && <BrowserBreakdownChart data={sources.browsers} />}
        </div>
        <div className="bg-white rounded-lg border border-gray-200 p-6">
          <h2 className="text-lg font-semibold mb-4">Countries</h2>
          {sources && <CountryBreakdownChart data={sources.countries} />}
        </div>
      </div>
    </div>
  );
}
