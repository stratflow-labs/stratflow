"use client";

import type { ReactNode } from "react";
import { useEffect } from "react";
import { useRouter } from "next/navigation";

import { getRedirectFromSearchParams } from "../lib/redirect";
import { useAuth } from "../providers/auth-provider";
import { AuthStatusScreen } from "./auth-status-screen";

export const PublicOnlyRoute = ({ children }: { children: ReactNode }) => {
  const { isAuthenticated, isLoading } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (isLoading || !isAuthenticated) {
      return;
    }

    router.replace(getRedirectFromSearchParams());
  }, [isAuthenticated, isLoading, router]);

  if (isLoading) {
    return (
      <AuthStatusScreen
        title="Checking session"
        description="Restoring your session before opening the login form."
      />
    );
  }

  if (isAuthenticated) {
    return null;
  }

  return <>{children}</>;
};
