"use client";

import { useState, type FormEvent } from "react";
import { useRouter } from "next/navigation";
import { Alert, Button, TextField, Typography } from "@mui/material";
import Box from "@mui/material/Box";
import CircularProgress from "@mui/material/CircularProgress";
import Stack from "@mui/material/Stack";
import type { Theme } from "@mui/material/styles";

import { getLoginErrorMessage } from "../lib/error-message";
import { getRedirectFromSearchParams } from "../lib/redirect";
import { useAuth } from "../providers/auth-provider";

const loginInputSx = (theme: Theme) => {
  const autofillBackground = "var(--mui-palette-background-paper)";
  const autofillText = "var(--mui-palette-text-primary)";

  return {
    "& .MuiInputBase-root:has(input:-webkit-autofill)": {
      backgroundColor: autofillBackground,
    },
    "& .MuiOutlinedInput-root:has(input:-webkit-autofill) .MuiOutlinedInput-notchedOutline":
      {
        borderColor: theme.palette.divider,
      },
    "& input:-webkit-autofill, & input:-webkit-autofill:hover, & input:-webkit-autofill:focus, & input:-webkit-autofill:active":
      {
        backgroundColor: `${autofillBackground} !important`,
        backgroundImage: "none !important",
        boxShadow: `0 0 0 1000px ${autofillBackground} inset !important`,
        WebkitBackgroundClip: "text",
        WebkitBoxShadow: `0 0 0 1000px ${autofillBackground} inset !important`,
        WebkitTextFillColor: `${autofillText} !important`,
        caretColor: `${autofillText} !important`,
        color: `${autofillText} !important`,
        transition:
          "background-color 9999s ease-out 0s, color 9999s ease-out 0s",
      },
    "& input:autofill, & input:autofill:hover, & input:autofill:focus, & input:autofill:active":
      {
        backgroundColor: `${autofillBackground} !important`,
        backgroundImage: "none !important",
        boxShadow: `0 0 0 1000px ${autofillBackground} inset !important`,
        WebkitTextFillColor: `${autofillText} !important`,
        caretColor: `${autofillText} !important`,
        color: `${autofillText} !important`,
      },
    "& input[data-autocompleted]": {
      backgroundColor: `${autofillBackground} !important`,
    },
  };
};

export const LoginForm = () => {
  const { login, isLoginPending } = useAuth();
  const router = useRouter();

  const [loginValue, setLoginValue] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [fieldErrors, setFieldErrors] = useState<{
    login?: string;
    password?: string;
  }>({});

  const validate = () => {
    const nextErrors: typeof fieldErrors = {};
    const trimmedLogin = loginValue.trim();

    if (!trimmedLogin) {
      nextErrors.login = "Login is required";
    } else if (trimmedLogin.length < 3) {
      nextErrors.login = "Enter at least 3 characters";
    }

    if (!password) {
      nextErrors.password = "Password is required";
    }

    setFieldErrors(nextErrors);
    return nextErrors;
  };

  const onSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();

    setError(null);

    const nextErrors = validate();

    if (Object.keys(nextErrors).length > 0) {
      return;
    }

    try {
      await login({
        login: loginValue.trim(),
        password,
      });

      router.replace(getRedirectFromSearchParams());
    } catch (submitError) {
      setError(getLoginErrorMessage(submitError));
    }
  };

  return (
    <Stack spacing={2.5} component="form" onSubmit={onSubmit} noValidate>
      <Box>
        <Typography variant="h5" sx={{ fontWeight: 700 }}>
          Sign in
        </Typography>
        <Typography variant="body2" color="text.secondary" sx={{ mt: 0.5 }}>
          Use your login and password.
        </Typography>
      </Box>

      {error ? <Alert severity="error">{error}</Alert> : null}

      <TextField
        label="Login"
        value={loginValue}
        onChange={(event) => {
          setLoginValue(event.target.value);
          setFieldErrors((current) => ({ ...current, login: undefined }));
        }}
        autoComplete="username"
        fullWidth
        disabled={isLoginPending}
        error={Boolean(fieldErrors.login)}
        helperText={fieldErrors.login}
        sx={loginInputSx}
      />

      <TextField
        label="Password"
        type="password"
        value={password}
        onChange={(event) => {
          setPassword(event.target.value);
          setFieldErrors((current) => ({ ...current, password: undefined }));
        }}
        autoComplete="current-password"
        fullWidth
        disabled={isLoginPending}
        error={Boolean(fieldErrors.password)}
        helperText={fieldErrors.password}
        sx={loginInputSx}
      />

      <Button
        type="submit"
        disabled={isLoginPending}
        variant="contained"
        startIcon={
          isLoginPending ? (
            <CircularProgress size={18} color="inherit" />
          ) : undefined
        }
      >
        {isLoginPending ? "Signing in..." : "Sign in"}
      </Button>
    </Stack>
  );
};
