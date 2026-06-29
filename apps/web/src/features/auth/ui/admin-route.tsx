"use client";

import type { ReactNode } from "react";
import { useEffect } from "react";
import { useRouter } from "next/navigation";

import { useAuth } from "../providers/auth-provider";
import { AuthStatusScreen } from "./auth-status-screen";

export const AdminRoute = ({ children }: { children: ReactNode }) => {
  const { isLoading, user } = useAuth();
  const router = useRouter();
  const isAdmin = user?.role === "admin";

  useEffect(() => {
    if (isLoading || isAdmin) {
      return;
    }

    router.replace("/");
  }, [isAdmin, isLoading, router]);

  if (isLoading) {
    return (
      <AuthStatusScreen
        title="Checking permissions"
        description="Verifying that your account can access the users workspace."
      />
    );
  }

  if (!isAdmin) {
    return (
      <AuthStatusScreen
        title="Redirecting"
        description="The users workspace is available only to administrators."
      />
    );
  }

  return <>{children}</>;
};
