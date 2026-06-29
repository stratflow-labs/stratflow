import { identityClient } from "@/shared/api/connect/clients";

export type ListUsersInput = {
  search?: string;
  page?: number;
  pageSize?: number;
  sort?: string;
};

export type CreateUserInput = {
  login: string;
  name: string;
  lastName?: string;
  email: string;
  role: string;
  password: string;
  gender?: number;
};

export const listUsers = async ({
  search,
  page,
  pageSize,
  sort,
}: ListUsersInput) => {
  const response = await identityClient.listUsers({
    search: search?.trim() || undefined,
    page,
    pageSize,
    sort: sort?.trim() || undefined,
  });

  if (!response.data) {
    throw new Error("Identity listUsers response does not contain list data.");
  }

  return {
    items: response.data.items,
    total: Number(response.data.total),
  };
};

export const createUser = async ({
  login,
  name,
  lastName,
  email,
  role,
  password,
  gender,
}: CreateUserInput) => {
  const response = await identityClient.createUser({
    login: login.trim(),
    name: name.trim(),
    lastName: lastName?.trim() || "",
    email: email.trim(),
    role: role.trim(),
    password,
    gender,
  });

  if (!response.data) {
    throw new Error("Identity createUser response does not contain user data.");
  }

  return response.data;
};
