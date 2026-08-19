import type { ShortenResponse } from "./types";

const API_BASE_URL = import.meta.env.VITE_API_URL ?? "http://localhost:8080";

export async function shortenUrl(originalUrl: string): Promise<ShortenResponse> {
  const res = await fetch(`${API_BASE_URL}/api/shorten`, {
    method: "POST",
    credentials: "include",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ url: originalUrl }),
  });
  if (!res.ok) throw new Error(await res.text() || "Failed to shorten URL");
  return res.json();
}

export interface AuthUser {
  id: number;
  first_name: string;
  last_name: string;
  email: string;
  created_at: string;
}

export async function registerUser(
  firstName: string,
  lastName: string,
  email: string,
  password: string
): Promise<AuthUser> {
  const res = await fetch(`${API_BASE_URL}/api/register`, {
    method: "POST",
    credentials: "include",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      first_name: firstName,
      last_name: lastName,
      email,
      password,
    }),
  });
  if (!res.ok) throw new Error(await res.text() || "Registration failed");
  return res.json();
}

export async function loginUser(email: string, password: string): Promise<AuthUser> {
  const res = await fetch(`${API_BASE_URL}/api/login`, {
    method: "POST",
    credentials: "include",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password }),
  });
  if (!res.ok) throw new Error(await res.text() || "Login failed");
  return res.json();
}

export async function logoutUser(): Promise<void> {
  await fetch(`${API_BASE_URL}/api/logout`, {
    method: "POST",
    credentials: "include",
  });
}

export interface MyURL {
  short_code: string;
  short_url: string;
  original_url: string;
  created_at: string;
}

export async function getMyUrls(): Promise<MyURL[]> {
  const res = await fetch(`${API_BASE_URL}/api/my/urls`, {
    method: "GET",
    credentials: "include",
  });
  if (!res.ok) throw new Error(await res.text() || "Failed to fetch your urls");
  return res.json();
}

export interface TimelinePoint {
  date: string;
  count: number;
}

export interface Stats {
  short_code: string;
  total_clicks: number;
  timeline: TimelinePoint[];
}

export async function getStats(code: string): Promise<Stats> {
  const res = await fetch(`${API_BASE_URL}/api/stats/${code}`, {
    method: "GET",
    credentials: "include",
  });
  if (!res.ok) throw new Error(await res.text() || "Failed to fetch stats");
  return res.json();
}