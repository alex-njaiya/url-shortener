function Footer() {
  const year = new Date().getFullYear();

  return (
    <footer className="border-t border-ink/10">
      <div className="mx-auto flex max-w-4xl flex-col items-center gap-2 px-6 py-8 text-sm text-muted sm:flex-row sm:justify-between">
        <span className="font-mono-url">sh.rt</span>
        <span>Made with ❤️</span>
        <span>&copy; {year} sh.rt. All rights reserved.</span>
      </div>
    </footer>
  );
}

export default Footer;