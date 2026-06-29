import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import type { ReactNode } from "react";
import { beforeEach, describe, expect, it, vi } from "vitest";

import { SimpleEntityDialog } from "./simple-entity-dialog";

vi.mock("@mui/material", () => ({
  Alert: ({ children }: { children: ReactNode }) => <div>{children}</div>,
  Button: ({
    children,
    disabled,
    onClick,
  }: {
    children: ReactNode;
    disabled?: boolean;
    onClick?: () => void;
  }) => (
    <button disabled={disabled} onClick={onClick}>
      {children}
    </button>
  ),
  Dialog: ({
    open,
    children,
  }: {
    open?: boolean;
    children: ReactNode;
  }) => (open ? <div>{children}</div> : null),
  DialogActions: ({ children }: { children: ReactNode }) => <div>{children}</div>,
  DialogContent: ({ children }: { children: ReactNode }) => <div>{children}</div>,
  DialogTitle: ({ children }: { children: ReactNode }) => <div>{children}</div>,
  Stack: ({ children }: { children: ReactNode }) => <div>{children}</div>,
  TextField: ({
    label,
    value,
    onChange,
  }: {
    label: string;
    value: string;
    onChange: (event: { target: { value: string } }) => void;
  }) => (
    <label>
      {label}
      <input
        aria-label={label}
        value={value}
        onChange={(event) => onChange({ target: { value: event.target.value } })}
      />
    </label>
  ),
  Typography: ({ children }: { children: ReactNode }) => <div>{children}</div>,
}));

describe("SimpleEntityDialog", () => {
  const onSubmit = vi.fn();

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("submits slug, title and description", async () => {
    const user = userEvent.setup();

    render(
      <SimpleEntityDialog
        open
        mode="strategy"
        title="Create strategy"
        subtitle="Enter slug, title and description."
        submitLabel="Create strategy"
        isSubmitting={false}
        error={null}
        onClose={() => {}}
        onSubmit={onSubmit}
      />,
    );

    await user.type(screen.getByLabelText("Slug"), "risk-reversal");
    await user.type(screen.getByLabelText("Title"), "Risk Reversal");
    await user.type(screen.getByLabelText("Description"), "New strategy");

    await user.click(screen.getByRole("button", { name: "Create strategy" }));

    expect(onSubmit).toHaveBeenCalledWith({
      slug: "risk-reversal",
      title: "Risk Reversal",
      description: "New strategy",
    });
  });
});
