export interface ShortenedURL {
  id: number;
  short_code: string;
  original_url: string;
  short_url: string;
  click_count: number;
  created_at: string;
}

export interface CreateURLRequest {
  url: string;
}

export interface CreateURLResponse {
  id: number;
  short_code: string;
  original_url: string;
  short_url: string;
  created_at: string;
}

export interface URLListResponse {
  urls: ShortenedURL[];
  total: number;
  page: number;
  limit: number;
}

export interface DailyClickCount {
  date: string;
  count: number;
}

export interface ClicksOverTimeResponse {
  url_id: number;
  period: string;
  total_clicks: number;
  clicks_per_day: DailyClickCount[];
}

export interface BrowserCount {
  name: string;
  count: number;
}

export interface CountryCount {
  code: string;
  name: string;
  count: number;
}

export interface SourcesResponse {
  url_id: number;
  period: string;
  browsers: BrowserCount[];
  countries: CountryCount[];
}

export interface APIError {
  error: string;
}
