import { create } from "@bufbuild/protobuf";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import type { ReactNode } from "react";
import { describe, expect, it, vi } from "vitest";

import { StrategySchema } from "@/shared/api/gen/strategy_registry/proto/v1/strategy_types_pb";

import { StrategyTable } from "./strategy-table";

vi.mock("@mui/material", () => ({
  Alert: ({ children }: { children: ReactNode }) => <div>{children}</div>,
  CircularProgress: () => <div>spinner</div>,
  Stack: ({ children }: { children: ReactNode }) => <div>{children}</div>,
  Table: ({ children }: { children: ReactNode }) => <table>{children}</table>,
  TableBody: ({ children }: { children: ReactNode }) => <tbody>{children}</tbody>,
  TableCell: ({
    children,
    colSpan,
  }: {
    children: ReactNode;
    colSpan?: number;
  }) => <td colSpan={colSpan}>{children}</td>,
  TableContainer: ({ children }: { children: ReactNode }) => <div>{children}</div>,
  TableHead: ({ children }: { children: ReactNode }) => <thead>{children}</thead>,
  TableRow: ({
    children,
    onClick,
  }: {
    children: ReactNode;
    onClick?: () => void;
  }) => <tr onClick={onClick}>{children}</tr>,
  Typography: ({ children }: { children: ReactNode }) => <span>{children}</span>,
}));

describe("StrategyTable", () => {
  it("renders empty search state", () => {
    render(
      <StrategyTable
        items={[]}
        isLoading={false}
        query="risk"
        onOpenStrategy={() => {}}
      />,
    );

    expect(
      screen.getByText("No strategies match the current search query."),
    ).toBeInTheDocument();
  });

  it("opens a strategy when a row is clicked", async () => {
    const user = userEvent.setup();
    const onOpenStrategy = vi.fn();
    const strategy = create(StrategySchema, {
      id: "strategy-1",
      slug: "vol-breakout",
      name: "Volatility Breakout",
      description: "Breakout strategy",
    });

    render(
      <StrategyTable
        items={[strategy]}
        isLoading={false}
        query=""
        onOpenStrategy={onOpenStrategy}
      />,
    );

    await user.click(screen.getByText("Volatility Breakout"));

    expect(onOpenStrategy).toHaveBeenCalledWith(strategy);
  });
});
