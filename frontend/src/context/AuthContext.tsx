import {
  createContext,
  useContext,
  useState,
  type ReactNode,
} from "react";
import {
  registerUser,
  loginUser,
  logoutUser,
  type AuthUser,
} from "../lib/api-client";

interface AuthContextValue {
  user: AuthUser | null;
  loading: boolean;
  login: (email: string, password: string) => Promise<void>;
  signup: (
    firstName: string,
    lastName: string,
    email: string,
    password: string
  ) => Promise<void>;
  logout: () => Promise<void>;
}

const AuthContext = createContext<AuthContextValue | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<AuthUser | null>(null);

  const [loading] = useState(false);

  async function login(email: string, password: string) {
    const loggedInUser = await loginUser(email, password);
    setUser(loggedInUser);
  }

  async function signup(
    firstName: string,
    lastName: string,
    email: string,
    password: string
  ) {
    const newUser = await registerUser(firstName, lastName, email, password);
    setUser(newUser);
  }

  async function logout() {
    await logoutUser();
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
