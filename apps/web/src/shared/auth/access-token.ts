const ACCESS_TOKEN_KEY = "qp.access_token";
const ACCESS_TOKEN_CLEARED_EVENT = "qp:access-token-cleared";

const isBrowser = () => typeof window !== "undefined";

export const getAccessToken = (): string | null => {
  if (!isBrowser()) {
    return null;
  }

  const value = window.localStorage.getItem(ACCESS_TOKEN_KEY);
  const token = value?.trim();

  return token ? token : null;
};

export const setAccessToken = (token: string): void => {
  if (!isBrowser()) {
    return;
  }

  const normalized = token.trim();

  if (!normalized) {
    clearAccessToken();
    return;
  }

  window.localStorage.setItem(ACCESS_TOKEN_KEY, normalized);
};

export const clearAccessToken = (): void => {
  if (!isBrowser()) {
    return;
  }

  window.localStorage.removeItem(ACCESS_TOKEN_KEY);
  window.dispatchEvent(new Event(ACCESS_TOKEN_CLEARED_EVENT));
};

export const onAccessTokenCleared = (handler: () => void): (() => void) => {
  if (!isBrowser()) {
    return () => undefined;
  }

  const handleStorage = (event: StorageEvent) => {
    if (
      event.storageArea === window.localStorage &&
      event.key === ACCESS_TOKEN_KEY &&
      event.newValue === null
    ) {
      handler();
    }
  };

  window.addEventListener(ACCESS_TOKEN_CLEARED_EVENT, handler);
  window.addEventListener("storage", handleStorage);

  return () => {
    window.removeEventListener(ACCESS_TOKEN_CLEARED_EVENT, handler);
    window.removeEventListener("storage", handleStorage);
  };
};
