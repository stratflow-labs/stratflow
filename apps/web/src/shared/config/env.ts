const trimTrailingSlash = (value: string): string => value.replace(/\/+$/, "");

const readEnv = (...keys: string[]): string | undefined => {
  for (const key of keys) {
    const value = process.env[key];
    if (value && value.trim() !== "") {
      return value;
    }
  }

  return undefined;
};

export const IDENTITY_CONNECT_BASE_URL = trimTrailingSlash(
  readEnv("NEXT_PUBLIC_IDENTITY_CONNECT_URL") ??
    "http://localhost:8085/connect/identity",
);

export const STRATEGY_REGISTRY_CONNECT_BASE_URL = trimTrailingSlash(
  readEnv("NEXT_PUBLIC_STRATEGY_REGISTRY_CONNECT_URL") ??
    "http://localhost:8085/connect/strategy-registry",
);
