import { Role } from "@/shared/api/gen/identity/proto/v1/types_pb";

import { getCurrentUser, verifyToken } from "../api/session";
import type { SessionUser } from "./types";

const mapRole = (role: Role): SessionUser["role"] => {
  switch (role) {
    case Role.ADMIN:
      return "admin";
    case Role.MANAGER:
      return "manager";
    case Role.USER:
    case Role.UNSPECIFIED:
    default:
      return "user";
  }
};

export const loadSessionUser = async (
  accessToken: string,
): Promise<SessionUser> => {
  const [tokenPayload, currentUser] = await Promise.all([
    verifyToken(accessToken),
    getCurrentUser(),
  ]);

  return {
    id: tokenPayload.userId,
    login: currentUser.login,
    name: currentUser.name,
    lastName: currentUser.lastName,
    email: currentUser.email ?? null,
    role: mapRole(tokenPayload.role),
    isEmailVerified: currentUser.isEmailVerified,
    isVerified: currentUser.isVerified,
  };
};
