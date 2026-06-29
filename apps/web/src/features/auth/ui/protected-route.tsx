"use client";

import type { ReactNode } from "react";
import { useEffect } from "react";
import { useRouter } from "next/navigation";

import { getCurrentRedirectPath } from "../lib/redirect";
import { useAuth } from "../providers/auth-provider";
import { AuthStatusScreen } from "./auth-status-screen";

export const ProtectedRoute = ({ children }: { children: ReactNode }) => {
  const { isAuthenticated, isLoading } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (isLoading || isAuthenticated) {
      return;
    }

    router.replace(
      `/login?from=${encodeURIComponent(getCurrentRedirectPath())}`,
    );
  }, [isAuthenticated, isLoading, router]);

  if (isLoading) {
    return (
      <AuthStatusScreen
        title="Checking session"
        description="Verifying your access token and loading the workspace."
      />
    );
  }

  if (!isAuthenticated) {
    return null;
  }

  return <>{children}</>;
};
