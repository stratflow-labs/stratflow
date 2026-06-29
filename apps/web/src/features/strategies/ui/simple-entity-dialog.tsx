"use client";

import {
  Alert,
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  Stack,
  TextField,
  Typography,
} from "@mui/material";
import { useState } from "react";

type BaseFields = {
  slug: string;
  title: string;
  description: string;
};

type SimpleEntityDialogProps = {
  open: boolean;
  mode: "strategy" | "attribute" | "value";
  title: string;
  subtitle: string;
  submitLabel: string;
  isSubmitting: boolean;
  error: string | null;
  onClose: () => void;
  onSubmit: (fields: BaseFields) => void | Promise<void>;
};

const getTitleLabel = (mode: SimpleEntityDialogProps["mode"]) =>
  mode === "value" ? "Value" : "Title";

export const SimpleEntityDialog = ({
  open,
  mode,
  title,
  subtitle,
  submitLabel,
  isSubmitting,
  error,
  onClose,
  onSubmit,
}: SimpleEntityDialogProps) => {
  const [slug, setSlug] = useState("");
  const [entityTitle, setEntityTitle] = useState("");
  const [description, setDescription] = useState("");

  const canSubmit = slug.trim().length > 0 && entityTitle.trim().length > 0;

  return (
    <Dialog open={open} onClose={isSubmitting ? undefined : onClose} maxWidth="sm" fullWidth>
      <DialogTitle sx={{ pb: 1 }}>
        <Stack spacing={0.5}>
          <Typography variant="h6" sx={{ fontWeight: 800 }}>
            {title}
          </Typography>
          <Typography color="text.secondary" variant="body2">
            {subtitle}
          </Typography>
        </Stack>
      </DialogTitle>

      <DialogContent dividers>
        <Stack spacing={2}>
          {error ? <Alert severity="error">{error}</Alert> : null}

          <TextField
            label="Slug"
            value={slug}
            onChange={(event) => setSlug(event.target.value)}
            autoFocus
            fullWidth
          />

          <TextField
            label={getTitleLabel(mode)}
            value={entityTitle}
            onChange={(event) => setEntityTitle(event.target.value)}
            fullWidth
          />

          <TextField
            label="Description"
            value={description}
            onChange={(event) => setDescription(event.target.value)}
            multiline
            minRows={3}
            fullWidth
          />
        </Stack>
      </DialogContent>

      <DialogActions sx={{ px: 3, py: 2 }}>
        <Button onClick={onClose} disabled={isSubmitting}>
          Cancel
        </Button>
        <Button
          variant="contained"
          disabled={!canSubmit || isSubmitting}
          onClick={() =>
            void onSubmit({
              slug: slug.trim(),
              title: entityTitle.trim(),
              description: description.trim(),
            })
          }
        >
          {isSubmitting ? "Saving..." : submitLabel}
        </Button>
      </DialogActions>
    </Dialog>
  );
};
