import { create } from "@bufbuild/protobuf";
import { act, renderHook, waitFor } from "@testing-library/react";
import { beforeEach, describe, expect, it, vi } from "vitest";

import { StrategySchema } from "@/shared/api/gen/strategy_registry/proto/v1/strategy_types_pb";

import { useStrategyDetailsDialog } from "./use-strategy-details-dialog";

const batchActionStrategyGraphMock = vi.fn();
const createAttributeMock = vi.fn();
const createAttributeValueMock = vi.fn();
const listStrategyAttributesMock = vi.fn();

vi.mock("../api/strategies", () => ({
  batchActionStrategyGraph: (...args: unknown[]) =>
    batchActionStrategyGraphMock(...args),
  createAttribute: (...args: unknown[]) => createAttributeMock(...args),
  createAttributeValue: (...args: unknown[]) =>
    createAttributeValueMock(...args),
  listStrategyAttributes: (...args: unknown[]) =>
    listStrategyAttributesMock(...args),
}));

describe("useStrategyDetailsDialog", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("saves graph relations through batch action and refreshes local graph state", async () => {
    listStrategyAttributesMock.mockResolvedValueOnce({
      items: [
        {
          id: "attribute-1",
          strategyId: "strategy-1",
          slug: "risk-level",
          name: "Risk level",
          description: "Risk bucket",
          values: [
            {
              id: "value-1",
              attributeId: "attribute-1",
              slug: "very-high",
              value: "Very High",
              relations: [],
            },
          ],
        },
      ],
      total: 1,
    });
    batchActionStrategyGraphMock.mockResolvedValue({
      data: {
        parameters: [
          {
            id: "attribute-1",
            strategyId: "strategy-1",
            slug: "risk-level",
            name: "Risk level",
            description: "Risk bucket",
            values: [
              {
                id: "value-1",
                attributeId: "attribute-1",
                slug: "very-high",
                value: "Very High",
                relations: [],
              },
            ],
          },
        ],
      },
    });

    const { result } = renderHook(() => useStrategyDetailsDialog());

    await act(async () => {
      await result.current.openStrategy(
        create(StrategySchema, {
          id: "strategy-1",
          slug: "vol-breakout",
          name: "Volatility Breakout",
          description: "Strategy",
        }),
      );
    });

    act(() => {
      result.current.openCreateAttributeDialog();
    });

    await act(async () => {
      await result.current.submitCreateDialog({
        slug: "timeframe",
        title: "Timeframe",
        description: "Execution timeframe",
      });
    });

    expect(createAttributeMock).toHaveBeenCalledWith({
      strategyRef: "strategy-1",
      slug: "timeframe",
      name: "Timeframe",
      description: "Execution timeframe",
    });
  });

  it("creates an attribute value even when the attribute draft has no backend id yet", async () => {
    listStrategyAttributesMock
      .mockResolvedValueOnce({
        items: [
          {
            id: "",
            strategyId: "strategy-1",
            slug: "risk-level",
            name: "Risk level",
            description: "Risk bucket",
            values: [],
          },
        ],
        total: 1,
      })
      .mockResolvedValueOnce({
        items: [
          {
            id: "attribute-1",
            strategyId: "strategy-1",
            slug: "risk-level",
            name: "Risk level",
            description: "Risk bucket",
            values: [
              {
                id: "value-1",
                attributeId: "attribute-1",
                slug: "very-high",
                value: "Very High",
                relations: [],
              },
            ],
          },
        ],
        total: 1,
      });

    createAttributeMock.mockResolvedValue({
      data: {
        id: "attribute-1",
      },
    });
    createAttributeValueMock.mockResolvedValue({
      data: {
        id: "value-1",
      },
    });

    const { result } = renderHook(() => useStrategyDetailsDialog());

    await act(async () => {
      await result.current.openStrategy(
        create(StrategySchema, {
        id: "strategy-1",
        slug: "vol-breakout",
        name: "Volatility Breakout",
        description: "Strategy",
        }),
      );
    });

    act(() => {
      result.current.openCreateValueDialog("risk-level");
    });

    await act(async () => {
      await result.current.submitCreateDialog({
        slug: "very-high",
        title: "Very High",
        description: "",
      });
    });

    expect(createAttributeMock).toHaveBeenCalledWith({
      strategyRef: "strategy-1",
      slug: "risk-level",
      name: "Risk level",
      description: "Risk bucket",
    });
    expect(createAttributeValueMock).toHaveBeenCalledWith({
      strategyRef: "strategy-1",
      attributeRef: "attribute-1",
      slug: "very-high",
      value: "Very High",
    });

    await waitFor(() => {
      expect(listStrategyAttributesMock).toHaveBeenCalledTimes(2);
    });
  });

  it("stores create dialog error when attribute creation fails", async () => {
    listStrategyAttributesMock.mockResolvedValue({
      items: [],
      total: 0,
    });
    createAttributeMock.mockRejectedValue(new Error("duplicate slug"));

    const { result } = renderHook(() => useStrategyDetailsDialog());

    await act(async () => {
      await result.current.openStrategy(
        create(StrategySchema, {
          id: "strategy-1",
          slug: "vol-breakout",
          name: "Volatility Breakout",
          description: "Strategy",
        }),
      );
    });

    act(() => {
      result.current.openCreateAttributeDialog();
    });

    await act(async () => {
      await result.current.submitCreateDialog({
        slug: "risk-level",
        title: "Risk level",
        description: "Risk bucket",
      });
    });

    expect(result.current.createError).toBe("duplicate slug");
    expect(result.current.isCreating).toBe(false);
  });
});
