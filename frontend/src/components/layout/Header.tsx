import { Link } from "react-router-dom";
import { useAuth } from "../../context/AuthContext";

function Header() {
  const { user, logout } = useAuth();

  return (
    <header className="border-b border-ink/10 bg-paper">
      <div className="mx-auto flex max-w-4xl items-center justify-between px-6 py-4">
        <Link
          to="/"
          className="font-mono-url text-lg font-semibold tracking-tight"
        >
          sh.rt
        </Link>

        {user ? (
          <div className="flex items-center gap-4 text-sm">
            <span className="text-muted">{user.email}</span>
            <Link to="/dashboard" className="font-medium hover:text-signal">
              Dashboard
            </Link>
            <button
              onClick={logout}
              className="font-medium text-muted hover:text-ink"
            >
              Log out
            </button>
          </div>
        ) : (
          <div className="flex items-center gap-4 text-sm">
            <Link to="/login" className="font-medium hover:text-signal">
              Log in
            </Link>
            <Link
              to="/signup"
              className="rounded-lg bg-signal px-4 py-2 font-semibold text-white hover:opacity-90"
            >
              Sign up
            </Link>
          </div>
        )}
      </div>
    </header>
  );
}

export default Header;
