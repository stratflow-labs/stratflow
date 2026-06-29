const DEFAULT_REDIRECT_PATH = "/";

const LOCAL_URL_BASE = "http://local.invalid";

export const normalizeRedirectPath = (
  value: string | null | undefined,
  fallback = DEFAULT_REDIRECT_PATH,
): string => {
  const rawValue = value?.trim();

  if (!rawValue) {
    return fallback;
  }

  if (rawValue.startsWith("\\") || rawValue.startsWith("//")) {
    return fallback;
  }

  try {
    const url = new URL(rawValue, LOCAL_URL_BASE);

    if (url.origin !== LOCAL_URL_BASE) {
      return fallback;
    }

    const nextPath = `${url.pathname}${url.search}${url.hash}`;

    return nextPath.startsWith("/") ? nextPath : fallback;
  } catch {
    return fallback;
  }
};

export const getCurrentRedirectPath = (): string => {
  if (typeof window === "undefined") {
    return DEFAULT_REDIRECT_PATH;
  }

  return normalizeRedirectPath(
    `${window.location.pathname}${window.location.search}${window.location.hash}`,
  );
};

export const getRedirectFromSearchParams = (): string => {
  if (typeof window === "undefined") {
    return DEFAULT_REDIRECT_PATH;
  }

  return normalizeRedirectPath(
    new URLSearchParams(window.location.search).get("from"),
  );
};
