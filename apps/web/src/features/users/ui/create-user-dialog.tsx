"use client";

import { useState } from "react";
import {
  Alert,
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  MenuItem,
  Stack,
  TextField,
  Typography,
} from "@mui/material";

import type { CreateUserInput } from "../api/users";

type CreateUserDialogProps = {
  open: boolean;
  isSubmitting: boolean;
  error: string | null;
  onClose: () => void;
  onSubmit: (input: CreateUserInput) => void | Promise<void>;
};

type FieldErrors = Partial<
  Record<"login" | "name" | "email" | "role" | "password", string>
>;

type GenderValue = "" | "1" | "2";

const initialForm = {
  login: "",
  name: "",
  lastName: "",
  email: "",
  role: "user",
  password: "",
  gender: "" as GenderValue,
};

const emailPattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

const roleOptions = [
  { value: "user", label: "User" },
  { value: "manager", label: "Manager" },
  { value: "admin", label: "Admin" },
];

const genderOptions = [
  { value: "", label: "Not specified" },
  { value: "1", label: "Male" },
  { value: "2", label: "Female" },
];

export const CreateUserDialog = ({
  open,
  isSubmitting,
  error,
  onClose,
  onSubmit,
}: CreateUserDialogProps) => {
  const [form, setForm] = useState(initialForm);
  const [fieldErrors, setFieldErrors] = useState<FieldErrors>({});

  const canSubmit =
    Boolean(form.login.trim()) &&
    Boolean(form.name.trim()) &&
    Boolean(form.email.trim()) &&
    Boolean(form.role.trim()) &&
    Boolean(form.password);

  const validate = (): FieldErrors => {
    const nextErrors: FieldErrors = {};

    if (!form.login.trim()) {
      nextErrors.login = "Login is required";
    } else if (form.login.trim().length < 3) {
      nextErrors.login = "Enter at least 3 characters";
    }

    if (!form.name.trim()) {
      nextErrors.name = "Name is required";
    }

    if (!form.email.trim()) {
      nextErrors.email = "Email is required";
    } else if (!emailPattern.test(form.email.trim())) {
      nextErrors.email = "Enter a valid email address";
    }

    if (!form.role.trim()) {
      nextErrors.role = "Role is required";
    }

    const trimmedPassword = form.password.trim();
    if (!trimmedPassword) {
      nextErrors.password = "Password is required";
    } else if (trimmedPassword.length < 5 || trimmedPassword.length > 32) {
      nextErrors.password = "Password must be between 5 and 32 characters";
    } else if (/\s/.test(form.password)) {
      nextErrors.password = "Password must not contain spaces";
    }

    return nextErrors;
  };

  const updateField = (field: keyof typeof initialForm, value: string) => {
    setForm((current) => ({ ...current, [field]: value }));
    if (field in fieldErrors) {
      setFieldErrors((current) => ({ ...current, [field]: undefined }));
    }
  };

  const handleSubmit = async () => {
    const nextErrors = validate();
    setFieldErrors(nextErrors);

    if (Object.keys(nextErrors).length > 0) {
      return;
    }

    await onSubmit({
      login: form.login.trim(),
      name: form.name.trim(),
      lastName: form.lastName.trim(),
      email: form.email.trim(),
      role: form.role.trim(),
      password: form.password,
      gender: form.gender ? Number(form.gender) : undefined,
    });
  };

  return (
    <Dialog
      open={open}
      onClose={isSubmitting ? undefined : onClose}
      maxWidth="sm"
      fullWidth
    >
      <DialogTitle sx={{ pb: 1 }}>
        <Stack spacing={0.5}>
          <Typography variant="h6" sx={{ fontWeight: 800 }}>
            Create user
          </Typography>
          <Typography color="text.secondary" variant="body2">
            Create a new identity account with administrative role assignment.
          </Typography>
        </Stack>
      </DialogTitle>

      <DialogContent dividers>
        <Stack spacing={2}>
          {error ? <Alert severity="error">{error}</Alert> : null}

          <Stack direction={{ xs: "column", sm: "row" }} spacing={2}>
            <TextField
              label="Login"
              value={form.login}
              onChange={(event) => updateField("login", event.target.value)}
              autoFocus
              fullWidth
              disabled={isSubmitting}
              error={Boolean(fieldErrors.login)}
              helperText={fieldErrors.login}
            />

            <TextField
              label="Role"
              value={form.role}
              onChange={(event) => updateField("role", event.target.value)}
              select
              fullWidth
              disabled={isSubmitting}
              error={Boolean(fieldErrors.role)}
              helperText={fieldErrors.role}
            >
              {roleOptions.map((option) => (
                <MenuItem key={option.value} value={option.value}>
                  {option.label}
                </MenuItem>
              ))}
            </TextField>
          </Stack>

          <Stack direction={{ xs: "column", sm: "row" }} spacing={2}>
            <TextField
              label="First name"
              value={form.name}
              onChange={(event) => updateField("name", event.target.value)}
              fullWidth
              disabled={isSubmitting}
              error={Boolean(fieldErrors.name)}
              helperText={fieldErrors.name}
            />

            <TextField
              label="Last name"
              value={form.lastName}
              onChange={(event) => updateField("lastName", event.target.value)}
              fullWidth
              disabled={isSubmitting}
            />
          </Stack>

          <TextField
            label="Email"
            value={form.email}
            onChange={(event) => updateField("email", event.target.value)}
            type="email"
            fullWidth
            disabled={isSubmitting}
            error={Boolean(fieldErrors.email)}
            helperText={fieldErrors.email}
          />

          <Stack direction={{ xs: "column", sm: "row" }} spacing={2}>
            <TextField
              label="Password"
              value={form.password}
              onChange={(event) => updateField("password", event.target.value)}
              type="password"
              fullWidth
              disabled={isSubmitting}
              error={Boolean(fieldErrors.password)}
              helperText={fieldErrors.password}
            />

            <TextField
              label="Gender"
              value={form.gender}
              onChange={(event) => updateField("gender", event.target.value)}
              select
              fullWidth
              disabled={isSubmitting}
            >
              {genderOptions.map((option) => (
                <MenuItem key={option.value || "empty"} value={option.value}>
                  {option.label}
                </MenuItem>
              ))}
            </TextField>
          </Stack>
        </Stack>
      </DialogContent>

      <DialogActions sx={{ px: 3, py: 2 }}>
        <Button onClick={onClose} disabled={isSubmitting}>
          Cancel
        </Button>
        <Button
          variant="contained"
          onClick={() => void handleSubmit()}
          disabled={!canSubmit || isSubmitting}
        >
          {isSubmitting ? "Creating..." : "Create user"}
        </Button>
      </DialogActions>
    </Dialog>
  );
};
