import { Code, ConnectError, type Interceptor } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";

import { clearAccessToken, getAccessToken } from "@/shared/auth/access-token";

const authInterceptor: Interceptor = (next) => async (req) => {
  const token = getAccessToken();

  if (token) {
    req.header.set("Authorization", `Bearer ${token}`);
  }

  try {
    return await next(req);
  } catch (error) {
    if (error instanceof ConnectError && error.code === Code.Unauthenticated) {
      clearAccessToken();
    }

    throw error;
  }
};

export const createBrowserTransport = (baseUrl: string) =>
  createConnectTransport({
    baseUrl,
    interceptors: [authInterceptor],
  });
