"use client";

import type { ReactNode } from "react";
import GroupRoundedIcon from "@mui/icons-material/GroupRounded";
import RuleFolderRoundedIcon from "@mui/icons-material/RuleFolderRounded";

export type AppNavItemKey = "strategy-registry" | "users";
export type AppNavRole = "user" | "manager" | "admin";

export type AppNavItem = {
  readonly key: AppNavItemKey;
  readonly label: string;
  readonly path: string;
  readonly icon: ReactNode;
  readonly visible?: boolean;
  readonly requiredRole?: AppNavRole;
};

const APP_NAV_ITEMS: readonly AppNavItem[] = [
  {
    key: "strategy-registry",
    label: "Strategy registry",
    path: "/",
    icon: <RuleFolderRoundedIcon fontSize="small" />,
    visible: true,
  },
  {
    key: "users",
    label: "Users",
    path: "/users",
    icon: <GroupRoundedIcon fontSize="small" />,
    requiredRole: "admin",
    visible: true,
  },
];

export const isAppNavItemPathMatch = (
  pathname: string,
  resourcePath: string,
): boolean => {
  if (resourcePath === "/") {
    return pathname === "/";
  }

  return pathname === resourcePath || pathname.startsWith(`${resourcePath}/`);
};

export const getVisibleAppNavItems = (
  role?: AppNavRole | null,
): readonly AppNavItem[] =>
  APP_NAV_ITEMS.filter((item) => {
    if (item.visible === false) {
      return false;
    }

    if (!item.requiredRole) {
      return true;
    }

    return item.requiredRole === role;
  });

export const findAppNavItemByPath = (
  pathname: string,
): AppNavItem | undefined =>
  APP_NAV_ITEMS.find((item) => isAppNavItemPathMatch(pathname, item.path));
