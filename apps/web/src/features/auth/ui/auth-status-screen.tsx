"use client";

import { CircularProgress, Stack, Typography } from "@mui/material";

type AuthStatusScreenProps = {
  title: string;
  description: string;
};

export const AuthStatusScreen = ({
  title,
  description,
}: AuthStatusScreenProps) => (
  <Stack
    component="main"
    spacing={2}
    sx={{
      minHeight: "100vh",
      px: 3,
      alignItems: "center",
      justifyContent: "center",
      textAlign: "center",
    }}
  >
    <CircularProgress size={28} />
    <div>
      <Typography variant="h6" sx={{ fontWeight: 700 }}>
        {title}
      </Typography>
      <Typography color="text.secondary">{description}</Typography>
    </div>
  </Stack>
);
