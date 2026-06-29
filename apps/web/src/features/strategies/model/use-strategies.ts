"use client";

import { Code, ConnectError } from "@connectrpc/connect";

import {
  usePaginatedResource,
  type PaginatedResourceState,
} from "@/shared/lib/use-paginated-resource";
import type { Strategy } from "@/shared/api/gen/strategy_registry/proto/v1/strategy_types_pb";
import { listStrategies } from "../api/strategies";

export type StrategiesState = PaginatedResourceState<Strategy>;

type UseStrategiesOptions = {
  pageSize?: number;
  sort?: string;
};

const getErrorMessage = (error: unknown): string => {
  if (error instanceof ConnectError) {
    if (error.code === Code.Unavailable || error.code === Code.Unknown) {
      return "Strategy registry is unavailable. Check that the service is running.";
    }

    if (error.code === Code.PermissionDenied) {
      return "You do not have permission to view strategies.";
    }
  }

  if (error instanceof Error && error.message) {
    return error.message;
  }

  return "Failed to load strategies.";
};

export const useStrategies = ({
  pageSize = 12,
  sort,
}: UseStrategiesOptions = {}) => {
  return usePaginatedResource<Strategy>({
    pageSize,
    sort,
    loadPage: ({ search, page, pageSize: nextPageSize, sort: nextSort }) =>
      listStrategies({
        search,
        page,
        pageSize: nextPageSize,
        sort: nextSort,
      }),
    getErrorMessage,
  });
};
