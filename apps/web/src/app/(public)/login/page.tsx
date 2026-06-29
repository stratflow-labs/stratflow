import Box from "@mui/material/Box";
import Container from "@mui/material/Container";
import Paper from "@mui/material/Paper";

import { LoginForm, PublicOnlyRoute } from "@/features/auth";

export default function LoginPage() {
  return (
    <PublicOnlyRoute>
      <Box
        component="main"
        sx={{
          minHeight: "100vh",
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
          px: 2,
        }}
      >
        <Container maxWidth="sm" sx={{ py: 2 }}>
          <Paper elevation={3} sx={{ p: { xs: 3, sm: 4 } }}>
            <LoginForm />
          </Paper>
        </Container>
      </Box>
    </PublicOnlyRoute>
  );
}
