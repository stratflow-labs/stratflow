import Box from "@mui/material/Box";
import { alpha } from "@mui/material/styles";

export const BrandMark = () => (
  <Box
    aria-hidden="true"
    sx={(theme) => ({
      width: 40,
      height: 40,
      borderRadius: 2,
      display: "grid",
      placeItems: "center",
      color: theme.palette.primary.contrastText,
      fontWeight: 900,
      background: `linear-gradient(135deg, ${theme.palette.primary.main}, ${theme.palette.secondary.main})`,
      boxShadow: `0 10px 30px ${alpha(theme.palette.primary.main, 0.25)}`,
    })}
  >
    S
  </Box>
);
