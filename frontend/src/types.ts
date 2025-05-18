export interface Result {
  Url: string;
  Links: string[];
}

export interface CrawlResponse {
  success: boolean;
  results: Result[];
}