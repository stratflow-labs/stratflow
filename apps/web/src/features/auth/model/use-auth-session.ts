"use client";

import { useCallback, useEffect, useMemo, useRef, useState } from "react";

import {
  clearAccessToken,
  getAccessToken,
  onAccessTokenCleared,
  setAccessToken,
} from "@/shared/auth/access-token";

import { login as loginRequest, logout as logoutRequest } from "../api/session";
import { loadSessionUser } from "./session-user";
import type { AuthSession, LoginFormValues, SessionUser } from "./types";

export const useAuthSession = (): AuthSession => {
  const isMountedRef = useRef(false);
  const sessionVersionRef = useRef(0);
  const loginRequestIdRef = useRef(0);

  const [user, setUser] = useState<SessionUser | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isLoginPending, setIsLoginPending] = useState(false);
  const [isLogoutPending, setIsLogoutPending] = useState(false);

  useEffect(() => {
    isMountedRef.current = true;

    return () => {
      isMountedRef.current = false;
    };
  }, []);

  const isCurrentSession = useCallback(
    (version: number, accessToken: string) =>
      isMountedRef.current &&
      sessionVersionRef.current === version &&
      getAccessToken() === accessToken,
    [],
  );

  const clearSession = useCallback(() => {
    sessionVersionRef.current += 1;

    if (!isMountedRef.current) {
      return;
    }

    setUser(null);
    setIsLoading(false);
    setIsLoginPending(false);
    setIsLogoutPending(false);
  }, []);

  const refreshIdentity = useCallback(async (): Promise<SessionUser | null> => {
    const accessToken = getAccessToken();

    if (!accessToken) {
      if (isMountedRef.current) {
        setUser(null);
      }

      return null;
    }

    const sessionVersion = sessionVersionRef.current;

    try {
      const nextUser = await loadSessionUser(accessToken);

      if (!isCurrentSession(sessionVersion, accessToken)) {
        return null;
      }

      setUser(nextUser);
      return nextUser;
    } catch {
      if (isCurrentSession(sessionVersion, accessToken)) {
        clearAccessToken();
      }

      return null;
    }
  }, [isCurrentSession]);

  useEffect(() => {
    const unsubscribe = onAccessTokenCleared(clearSession);

    const bootstrap = async () => {
      try {
        await refreshIdentity();
      } finally {
        if (isMountedRef.current) {
          setIsLoading(false);
        }
      }
    };

    void bootstrap();

    return unsubscribe;
  }, [clearSession, refreshIdentity]);

  const login = useCallback(
    async (values: LoginFormValues) => {
      const requestId = loginRequestIdRef.current + 1;
      loginRequestIdRef.current = requestId;
      setIsLoginPending(true);

      try {
        const accessToken = await loginRequest(values);

        if (!isMountedRef.current || loginRequestIdRef.current !== requestId) {
          return;
        }

        sessionVersionRef.current += 1;
        const sessionVersion = sessionVersionRef.current;
        setAccessToken(accessToken);

        const nextUser = await loadSessionUser(accessToken);

        if (!isCurrentSession(sessionVersion, accessToken)) {
          throw new Error(
            "Session changed while signing in. Please try again.",
          );
        }

        setUser(nextUser);
        setIsLoading(false);
      } catch (error) {
        if (isMountedRef.current && loginRequestIdRef.current === requestId) {
          clearAccessToken();
        }

        throw error;
      } finally {
        if (isMountedRef.current && loginRequestIdRef.current === requestId) {
          setIsLoginPending(false);
        }
      }
    },
    [isCurrentSession],
  );

  const logout = useCallback(async () => {
    const accessToken = getAccessToken();

    setIsLogoutPending(true);
    clearAccessToken();

    try {
      await logoutRequest(accessToken);
    } catch {
      // The local session is already closed; remote cleanup is best-effort.
    }
  }, []);

  return useMemo(
    () => ({
      isAuthenticated: Boolean(user),
      isLoading,
      isLoginPending,
      isLogoutPending,
      user,
      login,
      logout,
      refreshIdentity,
    }),
    [
      isLoading,
      isLoginPending,
      isLogoutPending,
      user,
      login,
      logout,
      refreshIdentity,
    ],
  );
};
