"use client";

import { createContext, useContext, type ReactNode } from "react";

import type { AuthSession } from "../model/types";
import { useAuthSession } from "../model/use-auth-session";

const AuthContext = createContext<AuthSession | null>(null);

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const session = useAuthSession();

  return (
    <AuthContext.Provider value={session}>{children}</AuthContext.Provider>
  );
};

export const useAuth = (): AuthSession => {
  const context = useContext(AuthContext);

  if (!context) {
    throw new Error("useAuth must be used within AuthProvider");
  }

  return context;
};
