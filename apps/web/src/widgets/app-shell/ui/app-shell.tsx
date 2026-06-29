"use client";

import { useState, type ReactNode } from "react";
import { Box, Stack } from "@mui/material";

import { APP_SHELL_CONTENT_MAX_WIDTH } from "../model/layout";
import { AppShellHeader } from "./app-shell-header";
import { AppShellNavigation } from "./app-shell-navigation";

export const AppShell = ({ children }: { children: ReactNode }) => {
  const [isMobileNavigationOpen, setIsMobileNavigationOpen] = useState(false);

  const openMobileNavigation = () => setIsMobileNavigationOpen(true);
  const closeMobileNavigation = () => setIsMobileNavigationOpen(false);

  return (
    <Box
      sx={{
        display: "grid",
        gridTemplateColumns: {
          xs: "minmax(0, 1fr)",
          md: "auto minmax(0, 1fr)",
        },
        minHeight: "100dvh",
        bgcolor: "background.default",
      }}
    >
      <AppShellNavigation
        mobileOpen={isMobileNavigationOpen}
        onCloseMobileNavigation={closeMobileNavigation}
      />

      <Stack
        component="main"
        sx={{
          minWidth: 0,
          minHeight: "100dvh",
        }}
      >
        <AppShellHeader onOpenMobileNavigation={openMobileNavigation} />
        <Box
          sx={{
            flex: 1,
            width: "100%",
            maxWidth: APP_SHELL_CONTENT_MAX_WIDTH,
            mx: "auto",
            px: { xs: 2, md: 3, xl: 4 },
            py: { xs: 3, md: 4 },
          }}
        >
          {children}
        </Box>
      </Stack>
    </Box>
  );
};
