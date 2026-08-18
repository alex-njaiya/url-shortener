import { useEffect, useState } from "react";

const LONG_URL =
  "https://example.com/blog/2024/how-to-build-a-url-shortener-from-scratch?ref=capstone&utm_source=chingu";
const SHORT_URL = "sh.rt/aZ3k";

function CompressionDemo() {
  const [compressed, setCompressed] = useState(false);

  useEffect(() => {
    const interval = setInterval(() => {
      setCompressed((prev) => !prev);
    }, 2200);
    return () => clearInterval(interval);
  }, []);

  return (
    <div
      className="flex w-full max-w-2xl flex-col items-start gap-3 rounded-xl border border-ink/10 bg-white/60 p-6 font-mono-url text-sm sm:text-base"
      aria-hidden="true"
    >
      <span className="text-sm uppercase tracking-wide text-muted">before</span>
      <span
        className={`block overflow-hidden whitespace-nowrap transition-all duration-700 ease-in-out ${
          compressed ? "max-w-[0ch] opacity-0" : "max-w-[70ch] opacity-100"
        }`}
      >
        {LONG_URL}
      </span>

      <span
        className={`block text-signal transition-all duration-700 ease-in-out ${
          compressed ? "opacity-100" : "opacity-0"
        }`}
      >
        {SHORT_URL}
      </span>
    </div>
  );
}

function Hero() {
  return (
    <section className="mx-auto flex max-w-4xl flex-col items-center gap-8 px-6 py-20 text-center sm:py-28">
      <h1 className="font-sans-head text-4xl font-semibold tracking-tight sm:text-5xl">
        Long Links, <span className="text-signal">shortened on the spot</span>
      </h1>
      <p className="max-w-xl text-lg text-muted">
        Paste a link,get something you would actually want to share. No signup
        required to shorten -- create an aacount only if you want to track cliks
        over time.
      </p>
      <CompressionDemo />
    </section>
  );
}

export default Hero;
