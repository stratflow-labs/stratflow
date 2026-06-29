"use client";

import MenuRoundedIcon from "@mui/icons-material/MenuRounded";
import { Box, IconButton, Toolbar } from "@mui/material";
import { alpha } from "@mui/material/styles";

import { RouteTitle } from "./route-title";

type AppShellHeaderProps = {
  onOpenMobileNavigation: () => void;
};

export const AppShellHeader = ({
  onOpenMobileNavigation,
}: AppShellHeaderProps) => (
  <Box
    component="header"
    sx={(theme) => ({
      position: "sticky",
      top: 0,
      zIndex: theme.zIndex.appBar,
      borderBottom: `1px solid ${theme.palette.divider}`,
      backdropFilter: "blur(16px)",
      backgroundColor: alpha(theme.palette.background.paper, 0.82),
    })}
  >
    <Toolbar
      sx={{
        gap: 2,
        px: { xs: 2, md: 3, xl: 4 },
      }}
    >
      <IconButton
        edge="start"
        color="inherit"
        aria-label="Open navigation"
        onClick={onOpenMobileNavigation}
        sx={{ display: { md: "none" } }}
      >
        <MenuRoundedIcon />
      </IconButton>
      <RouteTitle />
      <Box sx={{ flexGrow: 1 }} />
    </Toolbar>
  </Box>
);
