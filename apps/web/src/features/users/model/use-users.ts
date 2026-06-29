"use client";

import {
  usePaginatedResource,
  type PaginatedResourceState,
} from "@/shared/lib/use-paginated-resource";
import type { User } from "@/shared/api/gen/identity/proto/v1/types_pb";
import { getUsersErrorMessage } from "../lib/error-message";
import { listUsers } from "../api/users";

export type UsersState = PaginatedResourceState<User>;

type UseUsersOptions = {
  pageSize?: number;
  sort?: string;
};

export const useUsers = ({ pageSize = 25, sort }: UseUsersOptions = {}) => {
  return usePaginatedResource<User>({
    pageSize,
    sort,
    loadPage: ({ search, page, pageSize: nextPageSize, sort: nextSort }) =>
      listUsers({
        search,
        page,
        pageSize: nextPageSize,
        sort: nextSort,
      }),
    getErrorMessage: (error) => getUsersErrorMessage(error, "load"),
  });
};
