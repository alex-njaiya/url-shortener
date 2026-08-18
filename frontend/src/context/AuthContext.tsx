import {
  createContext,
  useContext,
  useState,
  type ReactNode,
} from "react";

interface User {
    email: string;
}

interface AuthContextValue {
  user: User | null;
  loading: boolean;
  login: (email: string, password: string) => Promise<void>;
  signup: (first_name: string, last_name: string, email: string, password: string) => Promise<void>;
  logout: () => void;
}

const AuthContext = createContext<AuthContextValue | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [loading] = useState(false); // TODO: true until an initial "am I logged in" check runs

  // TODO: replace with a real call to POST /api/login once the auth
  // domain exists on the backend. For now this just fakes a session
  // so Header/ProtectedRoute have something real to react to.
  async function login(email: string, _password: string) {
    setUser({ email });
  }

  async function signup(email: string, _password: string) {
    setUser({ email });
  }

  function logout() {
    setUser(null);
  }

  return (
    <AuthContext.Provider value={{ user, loading, login, signup, logout }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return ctx;
}
