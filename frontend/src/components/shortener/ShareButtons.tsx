interface ShareButtonsProps {
  shortUrl: string;
}

/**
 * Each platform (besides native share) just needs a URL with the link
 * embedded -- no SDK, no auth. We open these in a new tab so the user
 * never loses their place on the shortener page.
 */
function buildShareLinks(shortUrl: string) {
  const encoded = encodeURIComponent(shortUrl);
  return {
    whatsapp: `https://wa.me/?text=${encoded}`,
    x: `https://twitter.com/intent/tweet?url=${encoded}`,
    facebook: `https://www.facebook.com/sharer/sharer.php?u=${encoded}`,
    email: `mailto:?subject=${encodeURIComponent(
      "Check this out",
    )}&body=${encoded}`,
  };
}

function ShareButtons({ shortUrl }: ShareButtonsProps) {
  const links = buildShareLinks(shortUrl);

  // navigator.share only exists in some browsers (mostly mobile).
  // Feature-detect rather than assume -- calling it where it's
  // undefined throws, and there's no polyfill for "doesn't exist".
  const canUseNativeShare =
    typeof navigator !== "undefined" && !!navigator.share;

  async function handleNativeShare() {
    try {
      await navigator.share({ url: shortUrl });
    } catch {
      // User cancelled the share sheet -- not an error worth surfacing.
    }
  }

  return (
    <div className="mt-3 flex flex-wrap items-center gap-3 text-sm">
      <span className="text-muted">Share:</span>

      <a
        href={links.whatsapp}
        target="_blank"
        rel="noreferrer"
        className="font-medium hover:text-signal"
      >
        WhatsApp
      </a>
      <a
        href={links.x}
        target="_blank"
        rel="noreferrer"
        className="font-medium hover:text-signal"
      >
        X
      </a>
      <a
        href={links.facebook}
        target="_blank"
        rel="noreferrer"
        className="font-medium hover:text-signal"
      >
        Facebook
      </a>
      <a href={links.email} className="font-medium hover:text-signal">
        Email
      </a>

      {canUseNativeShare && (
        <button
          onClick={handleNativeShare}
          className="font-medium hover:text-signal"
        >
          More...
        </button>
      )}
    </div>
  );
}

export default ShareButtons;
