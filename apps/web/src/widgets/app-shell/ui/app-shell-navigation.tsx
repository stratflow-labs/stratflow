"use client";

import { Box, Drawer } from "@mui/material";

import { DESKTOP_SIDEBAR_WIDTH, MOBILE_SIDEBAR_WIDTH } from "../model/layout";
import { SidebarContent } from "./sidebar-content";

type AppShellNavigationProps = {
  mobileOpen: boolean;
  onCloseMobileNavigation: () => void;
};

export const AppShellNavigation = ({
  mobileOpen,
  onCloseMobileNavigation,
}: AppShellNavigationProps) => (
  <>
    <Box
      component="aside"
      sx={(theme) => ({
        display: { xs: "none", md: "block" },
        position: "sticky",
        top: 0,
        alignSelf: "start",
        height: "100dvh",
        width: DESKTOP_SIDEBAR_WIDTH,
        borderRight: `1px solid ${theme.palette.divider}`,
        backgroundColor: theme.palette.background.paper,
      })}
    >
      <SidebarContent />
    </Box>

    <Drawer
      variant="temporary"
      open={mobileOpen}
      onClose={onCloseMobileNavigation}
      ModalProps={{ keepMounted: true }}
      sx={{
        display: { xs: "block", md: "none" },
        "& .MuiDrawer-paper": {
          width: MOBILE_SIDEBAR_WIDTH,
          boxSizing: "border-box",
          borderRight: 0,
        },
      }}
    >
      <SidebarContent onNavigate={onCloseMobileNavigation} />
    </Drawer>
  </>
);
