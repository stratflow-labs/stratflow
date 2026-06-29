"use client";

import type { PropsWithChildren, ReactElement, ReactNode } from "react";
import { Box, Card, CardContent, Chip, Stack, Typography } from "@mui/material";

export const PageLayout = ({ children }: PropsWithChildren) => (
  <Stack spacing={{ xs: 3, lg: 4 }}>{children}</Stack>
);

type PageHeaderProps = PropsWithChildren<{
  title: string;
  description?: ReactNode;
  actions?: ReactNode;
}>;

export const PageHeader = ({
  title,
  description,
  actions,
}: PageHeaderProps) => (
  <Stack
    direction={{ xs: "column", xl: "row" }}
    spacing={2}
    sx={{
      alignItems: { xs: "stretch", xl: "start" },
      justifyContent: "space-between",
      gap: 2,
    }}
  >
    <Box sx={{ minWidth: 0 }}>
      <Typography
        component="h1"
        variant="h4"
        sx={{ fontWeight: 900, letterSpacing: "-0.04em" }}
      >
        {title}
      </Typography>
      {description ? (
        <Typography color="text.secondary">{description}</Typography>
      ) : null}
    </Box>

    {actions ? (
      <Box sx={{ width: { xs: "100%", xl: "auto" } }}>{actions}</Box>
    ) : null}
  </Stack>
);

type SurfaceSectionProps = PropsWithChildren<{
  title: string;
  description?: ReactNode;
  badge?: ReactNode;
}>;

export const SurfaceSection = ({
  title,
  description,
  badge,
  children,
}: SurfaceSectionProps) => (
  <Card>
    <CardContent
      sx={{
        p: { xs: 2.5, lg: 3 },
        "&:last-child": { pb: { xs: 2.5, lg: 3 } },
      }}
    >
      <Stack spacing={3}>
        <Stack
          direction={{ xs: "column", lg: "row" }}
          spacing={2}
          sx={{
            alignItems: { xs: "stretch", lg: "start" },
            justifyContent: "space-between",
          }}
        >
          <Box sx={{ minWidth: 0 }}>
            <Typography component="h2" variant="h6" sx={{ fontWeight: 800 }}>
              {title}
            </Typography>
            {description ? (
              <Typography color="text.secondary">{description}</Typography>
            ) : null}
          </Box>

          {badge ? (
            <Box sx={{ alignSelf: { xs: "start", lg: "center" } }}>{badge}</Box>
          ) : null}
        </Stack>

        {children}
      </Stack>
    </CardContent>
  </Card>
);

type StatBadgeProps = {
  icon?: ReactElement;
  label: string;
};

export const StatBadge = ({ icon, label }: StatBadgeProps) => (
  <Chip icon={icon} label={label} sx={{ fontWeight: 700 }} />
);
