export type SessionUser = {
  id: string;
  login: string;
  name: string;
  lastName: string;
  email: string | null;
  role: "user" | "manager" | "admin";
  isEmailVerified: boolean;
  isVerified: boolean;
};

export type LoginFormValues = {
  login: string;
  password: string;
};

export type AuthSession = {
  isAuthenticated: boolean;
  isLoading: boolean;
  isLoginPending: boolean;
  isLogoutPending: boolean;
  user: SessionUser | null;
  login: (values: LoginFormValues) => Promise<void>;
  logout: () => Promise<void>;
  refreshIdentity: () => Promise<SessionUser | null>;
};
