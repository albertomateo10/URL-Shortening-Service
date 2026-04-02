import type {
  CreateURLResponse,
  URLListResponse,
  ShortenedURL,
  ClicksOverTimeResponse,
  SourcesResponse,
} from "@/types/api";

const API_BASE = "/api";

async function fetchJSON<T>(url: string, init?: RequestInit): Promise<T> {
  const res = await fetch(url, init);
  if (!res.ok) {
    const body = await res.json().catch(() => ({ error: "Request failed" }));
    throw new Error(body.error || `HTTP ${res.status}`);
  }
  if (res.status === 204) return undefined as T;
  return res.json();
}

export async function createURL(url: string): Promise<CreateURLResponse> {
  return fetchJSON<CreateURLResponse>(`${API_BASE}/urls`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ url }),
  });
}

export async function listURLs(
  page = 1,
  limit = 20
): Promise<URLListResponse> {
  return fetchJSON<URLListResponse>(
    `${API_BASE}/urls?page=${page}&limit=${limit}`
  );
}

export async function getURL(id: number): Promise<ShortenedURL> {
  return fetchJSON<ShortenedURL>(`${API_BASE}/urls/${id}`);
}

export async function deleteURL(id: number): Promise<void> {
  return fetchJSON<void>(`${API_BASE}/urls/${id}`, { method: "DELETE" });
}

export async function getClicksOverTime(
  id: number,
  period = "7d"
): Promise<ClicksOverTimeResponse> {
  return fetchJSON<ClicksOverTimeResponse>(
    `${API_BASE}/urls/${id}/analytics/clicks?period=${period}`
  );
}

export async function getSources(
  id: number,
  period = "7d"
): Promise<SourcesResponse> {
  return fetchJSON<SourcesResponse>(
    `${API_BASE}/urls/${id}/analytics/sources?period=${period}`
  );
}
