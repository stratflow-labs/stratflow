"use client";

import { useCallback, useState } from "react";

import { getUsersErrorMessage } from "../lib/error-message";
import { createUser, type CreateUserInput } from "../api/users";

export const useCreateUser = () => {
  const [isPending, setPending] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const submit = useCallback(async (input: CreateUserInput) => {
    setPending(true);
    setError(null);

    try {
      return await createUser(input);
    } catch (submitError) {
      const nextError = getUsersErrorMessage(submitError, "create");
      setError(nextError);
      throw submitError;
    } finally {
      setPending(false);
    }
  }, []);

  const reset = useCallback(() => {
    setError(null);
  }, []);

  return {
    isPending,
    error,
    submit,
    reset,
  };
};
