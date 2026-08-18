import type { ShortenResponse } from "./types";

const API_BASE_URL = import.meta.env.VITE_API_URL ?? "http://localhost:8080";

export async function shortenUrl(originalUrl: string): Promise<ShortenResponse> {
  const res = await fetch(`${API_BASE_URL}/api/shorten`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ url: originalUrl }),
  });

  if (!res.ok) {
    const message = await res.text();
    throw new Error(message || "Failed to shorten URL");
  }

  return res.json();
}