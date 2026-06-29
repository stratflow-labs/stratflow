"use client";

import {
  Box,
  Divider,
  List,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  Stack,
  Typography,
} from "@mui/material";
import { alpha } from "@mui/material/styles";
import Link from "next/link";
import { usePathname } from "next/navigation";

import { useAuth } from "@/features/auth";
import {
  getVisibleAppNavItems,
  isAppNavItemPathMatch,
} from "../model/navigation";

import { AccountMenu } from "./account-menu";
import { BrandMark } from "./brand-mark";

export const SidebarContent = ({ onNavigate }: { onNavigate?: () => void }) => {
  const pathname = usePathname() ?? "/";
  const { user } = useAuth();
  const navItems = getVisibleAppNavItems(user?.role);

  return (
    <Stack sx={{ height: "100%" }}>
      <Stack
        direction="row"
        spacing={1.5}
        sx={{ px: 3, py: 2.5, alignItems: "center" }}
      >
        <BrandMark />
        <Box>
          <Typography variant="h6" sx={{ fontWeight: 900, lineHeight: 1.1 }}>
            StratFlow
          </Typography>
          <Typography variant="caption" color="text.secondary">
            Control plane
          </Typography>
        </Box>
      </Stack>

      <Divider />

      <List aria-label="Primary navigation" sx={{ p: 1.5 }}>
        {navItems.map((item) => {
          const selected = isAppNavItemPathMatch(pathname, item.path);

          return (
            <ListItemButton
              key={item.key}
              component={Link}
              href={item.path}
              selected={selected}
              onClick={onNavigate}
              sx={{
                mb: 0.5,
                borderRadius: 2,
                px: 1.5,
                py: 1,
                "&.Mui-selected": {
                  color: "primary.main",
                  bgcolor: (theme) => alpha(theme.palette.primary.main, 0.1),
                },
              }}
            >
              <ListItemIcon sx={{ color: "inherit", minWidth: 36 }}>
                {item.icon}
              </ListItemIcon>
              <ListItemText
                primary={
                  <Typography
                    variant="body2"
                    sx={{ fontWeight: selected ? 800 : 600 }}
                  >
                    {item.label}
                  </Typography>
                }
              />
            </ListItemButton>
          );
        })}
      </List>

      <Box sx={{ flexGrow: 1 }} />
      <Box sx={{ p: 1.5 }}>
        <AccountMenu />
      </Box>
    </Stack>
  );
};
