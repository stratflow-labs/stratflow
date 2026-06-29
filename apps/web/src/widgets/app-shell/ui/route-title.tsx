"use client";

import Box from "@mui/material/Box";
import Typography from "@mui/material/Typography";
import { usePathname } from "next/navigation";

import { findAppNavItemByPath } from "../model/navigation";

export const RouteTitle = () => {
  const pathname = usePathname() ?? "/";
  const currentItem = findAppNavItemByPath(pathname);
  const title = currentItem?.label ?? "Workspace";

  return (
    <Box>
      <Typography
        variant="caption"
        color="text.secondary"
        sx={{ letterSpacing: "0.08em", textTransform: "uppercase" }}
      >
        StratFlow / {title}
      </Typography>
      <Typography
        variant="h6"
        sx={{ fontWeight: 900, letterSpacing: "-0.03em" }}
      >
        {title}
      </Typography>
    </Box>
  );
};
