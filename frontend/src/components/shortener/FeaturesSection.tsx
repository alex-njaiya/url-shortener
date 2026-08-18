const FEATURES = [
  {
    label: "no signup to shorten",
    detail: "Paste a link and get a short one back immediately -- an account is only needed if you want to see your history.",
  },
  {
    label: "your links, tracked",
    detail: "Once you're signed in, every click on your links is logged -- see counts and activity over time.",
  },
  {
    label: "built to redirect fast",
    detail: "Short links resolve straight to their destination -- no interstitial page, no delay.",
  },
];

function FeaturesSection() {
  return (
    <section className="mx-auto grid max-w-4xl gap-8 px-6 py-20 sm:grid-cols-3 sm:py-28">
      {FEATURES.map((feature) => (
        <div key={feature.label} className="flex flex-col gap-2">
          <h3 className="font-mono-url text-sm font-semibold text-signal">
            {feature.label}
          </h3>
          <p className="text-sm text-muted">{feature.detail}</p>
        </div>
      ))}
    </section>
  );
}

export default FeaturesSection;
