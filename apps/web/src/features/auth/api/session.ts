import { Code, ConnectError } from "@connectrpc/connect";

import { identityClient } from "@/shared/api/connect/clients";

type LoginInput = {
  login: string;
  password: string;
};

export const login = async (input: LoginInput): Promise<string> => {
  const response = await identityClient.login(input);
  const token = response.data?.accessToken?.trim();

  if (!token) {
    throw new Error(
      "Identity login response does not contain an access token.",
    );
  }

  return token;
};

export const verifyToken = (accessToken: string) =>
  identityClient.verifyToken({ accessToken });

export const getCurrentUser = async () => {
  const response = await identityClient.getCurrentUser({});

  if (!response.data) {
    throw new Error(
      "Identity getCurrentUser response does not contain user data.",
    );
  }

  return response.data;
};

export const logout = async (accessToken: string | null): Promise<void> => {
  try {
    await identityClient.logout(
      {},
      accessToken
        ? { headers: { Authorization: `Bearer ${accessToken}` } }
        : undefined,
    );
  } catch (error) {
    if (
      !(error instanceof ConnectError && error.code === Code.Unauthenticated)
    ) {
      throw error;
    }
  }
};
