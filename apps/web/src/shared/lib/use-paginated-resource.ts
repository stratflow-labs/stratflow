"use client";

import { useCallback, useEffect, useRef, useState } from "react";

export type PaginatedResourceStatus = "loading" | "loaded" | "error";

export type PaginatedResourceState<TItem> = {
  status: PaginatedResourceStatus;
  items: TItem[];
  total: number;
  error: string | null;
};

type PaginatedResourceResult<TItem> = {
  items: TItem[];
  total: number;
};

type UsePaginatedResourceOptions<TItem> = {
  pageSize: number;
  sort?: string;
  loadPage: (input: {
    search: string;
    page: number;
    pageSize: number;
    sort?: string;
  }) => Promise<PaginatedResourceResult<TItem>>;
  getErrorMessage: (error: unknown) => string;
  searchDebounceMs?: number;
};

export const usePaginatedResource = <TItem>({
  pageSize,
  sort,
  loadPage,
  getErrorMessage,
  searchDebounceMs = 250,
}: UsePaginatedResourceOptions<TItem>) => {
  const [query, setQuery] = useState("");
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(pageSize);
  const [state, setState] = useState<PaginatedResourceState<TItem>>({
    status: "loading",
    items: [],
    total: 0,
    error: null,
  });
  const requestIdRef = useRef(0);

  const load = useCallback(
    async (search: string, currentPage: number, currentRowsPerPage: number) => {
      const requestId = requestIdRef.current + 1;
      requestIdRef.current = requestId;

      setState((currentState) => ({
        ...currentState,
        status: "loading",
        error: null,
      }));

      try {
        const result = await loadPage({
          search,
          page: currentPage + 1,
          pageSize: currentRowsPerPage,
          sort,
        });

        if (requestIdRef.current !== requestId) {
          return;
        }

        const lastPage = Math.max(
          0,
          Math.ceil(result.total / currentRowsPerPage) - 1,
        );

        if (currentPage > lastPage) {
          setPage(lastPage);
          return;
        }

        setState({
          status: "loaded",
          items: result.items,
          total: result.total,
          error: null,
        });
      } catch (error) {
        if (requestIdRef.current !== requestId) {
          return;
        }

        setState((currentState) => ({
          status: "error",
          items: currentState.items,
          total: currentState.total,
          error: getErrorMessage(error),
        }));
      }
    },
    [getErrorMessage, loadPage, sort],
  );

  useEffect(() => {
    const debounceMs = query.trim() ? searchDebounceMs : 0;
    const timeoutId = window.setTimeout(() => {
      void load(query, page, rowsPerPage);
    }, debounceMs);

    return () => {
      window.clearTimeout(timeoutId);
    };
  }, [load, page, query, rowsPerPage, searchDebounceMs]);

  const updateQuery = useCallback((nextQuery: string) => {
    setPage(0);
    setQuery(nextQuery);
  }, []);

  const updateRowsPerPage = useCallback((nextRowsPerPage: number) => {
    setPage(0);
    setRowsPerPage(nextRowsPerPage);
  }, []);

  const refresh = useCallback(
    () => load(query, page, rowsPerPage),
    [load, page, query, rowsPerPage],
  );

  return {
    state,
    page,
    query,
    rowsPerPage,
    setPage,
    setQuery: updateQuery,
    setRowsPerPage: updateRowsPerPage,
    refresh,
  };
};
