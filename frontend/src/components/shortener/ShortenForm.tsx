import { useState, type FormEvent } from "react";
import { shortenUrl } from "../../lib/api-client";
import type { ShortenResponse } from "../../lib/types";
import ShareButtons from "./ShareButtons";

function ShortenForm() {
  const [originalUrl, setOriginalUrl] = useState("");
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState<ShortenResponse | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [copied, setCopied] = useState(false);

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setLoading(true);
    setError(null);
    setResult(null);

    try {
      const response = await shortenUrl(originalUrl);
      setResult(response);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Something went wrong");
    } finally {
      setLoading(false);
    }
  }

  async function handleCopy() {
    if (!result) return;
    await navigator.clipboard.writeText(result.short_url);
    setCopied(true);
    setTimeout(() => setCopied(false), 1500);
  }

  return (
    <div className="mx-auto w-full max-w-2xl px-6">
      <form
        onSubmit={handleSubmit}
        className="flex flex-col gap-3 sm:flex-row"
      >
        <input
          type="url"
          required
          placeholder="Paste a long URL"
          value={originalUrl}
          onChange={(e) => setOriginalUrl(e.target.value)}
          className="flex-1 rounded-lg border border-ink/15 bg-white px-4 py-3 font-mono-url text-sm outline-none focus:border-signal focus:ring-2 focus:ring-signal/20"
        />
        <button
          type="submit"
          disabled={loading}
          className="rounded-lg bg-signal px-6 py-3 font-semibold text-white transition-opacity hover:opacity-90 disabled:opacity-50"
        >
          {loading ? "Shortening..." : "Shorten"}
        </button>
      </form>

      {error && (
        <p className="mt-3 text-sm text-red-600" role="alert">
          {error}
        </p>
      )}

      {result && (
        <div className="mt-4 rounded-lg border border-success/30 bg-success/5 px-4 py-3">
          <div className="flex items-center justify-between">
            <a
              href={result.short_url}
              target="_blank"
              rel="noreferrer"
              className="font-mono-url text-signal underline"
            >
              {result.short_url}
            </a>
            <button
              onClick={handleCopy}
              className="text-sm font-medium text-muted hover:text-ink"
            >
              {copied ? "Copied" : "Copy"}
            </button>
          </div>
          <ShareButtons shortUrl={result.short_url} />
        </div>
      )}
    </div>
  );
}

export default ShortenForm;
