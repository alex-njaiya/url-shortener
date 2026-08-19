import { useEffect, useState } from "react";
import { getMyUrls, getStats, type MyURL, type Stats } from "../lib/api-client";

function MiniTimeline({ timeline }: { timeline: Stats["timeline"] }) {
  if (!timeline || timeline.length === 0) {
    return <p className="text-xs text-muted">No clicks yet</p>;
  }

  const max = Math.max(...timeline.map((t) => t.count));

  return (
    <div className="flex items-end gap-1" title="Clicks per day">
      {timeline.map((point) => (
        <div
          key={point.date}
          className="w-2 rounded-sm bg-signal/70"
          style={{ height: `${Math.max(4, (point.count / max) * 32)}px` }}
          title={`${point.date}: ${point.count} click${point.count === 1 ? "" : "s"}`}
        />
      ))}
    </div>
  );
}

function DashboardPage() {
  const [urls, setUrls] = useState<MyURL[] | null>(null);
  const [statsByCode, setStatsByCode] = useState<Record<string, Stats>>({});
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    async function load() {
      try {
        const myUrls = await getMyUrls();
        setUrls(myUrls);

        // Fetch stats for every url in parallel -- fine at this
        // scale (a handful of links per user), would need a batched
        // backend endpoint if this ever needed to scale further.
        const statsEntries = await Promise.all(
          myUrls.map(async (u) => {
            try {
              const stats = await getStats(u.short_code);
              return [u.short_code, stats] as const;
            } catch {
              return null;
            }
          })
        );

        const map: Record<string, Stats> = {};
        for (const entry of statsEntries) {
          if (entry) map[entry[0]] = entry[1];
        }
        setStatsByCode(map);
      } catch (err) {
        setError(err instanceof Error ? err.message : "Failed to load your links");
      } finally {
        setLoading(false);
      }
    }

    load();
  }, []);

  if (loading) {
    return <main className="p-8 text-center text-muted">Loading your links...</main>;
  }

  if (error) {
    return (
      <main className="p-8 text-center text-red-600" role="alert">
        {error}
      </main>
    );
  }

  return (
    <main className="mx-auto max-w-4xl px-6 py-12">
      <h1 className="mb-8 text-2xl font-semibold tracking-tight">Your links</h1>

      {urls && urls.length === 0 && (
        <p className="text-muted">
          You haven't shortened any links yet. Head back to the homepage to
          create one.
        </p>
      )}

      <div className="flex flex-col gap-4">
        {urls?.map((u) => {
          const stats = statsByCode[u.short_code];
          return (
            <div
              key={u.short_code}
              className="flex flex-col gap-3 rounded-lg border border-ink/10 bg-white/60 p-4 sm:flex-row sm:items-center sm:justify-between"
            >
              <div className="min-w-0 flex-1">
                <a
                  href={u.short_url}
                  target="_blank"
                  rel="noreferrer"
                  className="font-mono-url text-signal underline"
                >
                  {u.short_url}
                </a>
                <p className="mt-1 truncate text-sm text-muted">
                  {u.original_url}
                </p>
                <p className="mt-1 text-xs text-muted">
                  {new Date(u.created_at).toLocaleDateString()}
                </p>
              </div>

              <div className="flex items-center gap-4">
                <div className="text-right">
                  <p className="text-lg font-semibold">
                    {stats?.total_clicks ?? 0}
                  </p>
                  <p className="text-xs text-muted">clicks</p>
                </div>
                <MiniTimeline timeline={stats?.timeline ?? []} />
              </div>
            </div>
          );
        })}
      </div>
    </main>
  );
}

export default DashboardPage;
