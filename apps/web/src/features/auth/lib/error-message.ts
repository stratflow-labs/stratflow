import { Code, ConnectError } from "@connectrpc/connect";

export const getLoginErrorMessage = (error: unknown): string => {
  if (error instanceof ConnectError) {
    if (error.code === Code.Unauthenticated) {
      return "Invalid login or password.";
    }

    if (error.code === Code.Unavailable || error.code === Code.Unknown) {
      return "Authentication service is unavailable.";
    }
  }

  if (error instanceof Error && error.message) {
    return error.message;
  }

  return "Unable to sign in.";
};
