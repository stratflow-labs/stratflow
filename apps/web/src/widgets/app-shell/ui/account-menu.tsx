"use client";

import { Fragment, useId, useState, type MouseEvent } from "react";
import LogoutRoundedIcon from "@mui/icons-material/LogoutRounded";
import MoreVertRoundedIcon from "@mui/icons-material/MoreVertRounded";
import {
  Avatar,
  Box,
  ButtonBase,
  ListItemIcon,
  ListItemText,
  Menu,
  MenuItem,
  Stack,
  Typography,
} from "@mui/material";
import { alpha } from "@mui/material/styles";

import { useAuth } from "@/features/auth";

const getInitials = (name?: string | null, fallback?: string | null) => {
  const value = name?.trim() || fallback?.trim() || "User";
  return value.charAt(0).toUpperCase();
};

export const AccountMenu = () => {
  const { user, logout, isLogoutPending } = useAuth();
  const menuId = useId();
  const [anchorEl, setAnchorEl] = useState<HTMLElement | null>(null);
  const open = Boolean(anchorEl);
  const displayName =
    [user?.name, user?.lastName].filter(Boolean).join(" ") ||
    user?.login ||
    "User";

  const openMenu = (event: MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget);
  };

  const closeMenu = () => {
    setAnchorEl(null);
  };

  const handleLogout = () => {
    closeMenu();
    void logout();
  };

  return (
    <Fragment>
      <ButtonBase
        aria-label="Open account menu"
        aria-controls={open ? menuId : undefined}
        aria-haspopup="menu"
        aria-expanded={open ? "true" : undefined}
        onClick={openMenu}
        sx={(theme) => ({
          width: "100%",
          p: 1,
          borderRadius: 2,
          border: `1px solid ${theme.palette.divider}`,
          bgcolor: alpha(theme.palette.background.default, 0.7),
          justifyContent: "flex-start",
        })}
      >
        <Stack
          direction="row"
          spacing={1.5}
          sx={{ width: "100%", alignItems: "center" }}
        >
          <Avatar sx={{ width: 36, height: 36 }}>
            {getInitials(user?.name, user?.login)}
          </Avatar>
          <Box sx={{ minWidth: 0, flex: 1, textAlign: "left" }}>
            <Typography variant="body2" noWrap sx={{ fontWeight: 800 }}>
              {displayName}
            </Typography>
            <Typography variant="caption" color="text.secondary" noWrap>
              {user?.email || user?.login}
            </Typography>
          </Box>
          <MoreVertRoundedIcon fontSize="small" />
        </Stack>
      </ButtonBase>

      <Menu id={menuId} anchorEl={anchorEl} open={open} onClose={closeMenu}>
        <MenuItem onClick={handleLogout} disabled={isLogoutPending}>
          <ListItemIcon>
            <LogoutRoundedIcon fontSize="small" />
          </ListItemIcon>
          <ListItemText
            primary={isLogoutPending ? "Signing out..." : "Sign out"}
          />
        </MenuItem>
      </Menu>
    </Fragment>
  );
};
