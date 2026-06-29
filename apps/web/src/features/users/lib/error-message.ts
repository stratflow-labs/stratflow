"use client";

import { Code, ConnectError } from "@connectrpc/connect";

type UsersAction = "load" | "create";

export const getUsersErrorMessage = (
  error: unknown,
  action: UsersAction,
): string => {
  if (error instanceof ConnectError) {
    if (error.code === Code.PermissionDenied) {
      return action === "create"
        ? "You do not have permission to create users."
        : "You do not have permission to view users.";
    }

    if (error.code === Code.Unauthenticated) {
      return "Your session expired. Sign in again to continue.";
    }

    if (error.code === Code.AlreadyExists) {
      return "A user with the same login or email already exists.";
    }

    if (error.code === Code.Unavailable || error.code === Code.Unknown) {
      return action === "create"
        ? "Identity service is unavailable. Unable to create the user."
        : "Identity service is unavailable. Check that the service is running.";
    }
  }

  if (error instanceof Error && error.message) {
    return error.message;
  }

  return action === "create" ? "Failed to create user." : "Failed to load users.";
};
