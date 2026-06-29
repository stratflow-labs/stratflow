import { ConnectError, Code } from "@connectrpc/connect";
import { create } from "@bufbuild/protobuf";
import { act, renderHook } from "@testing-library/react";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { StrategySchema } from "@/shared/api/gen/strategy_registry/proto/v1/strategy_types_pb";

import { useStrategies } from "./use-strategies";

const listStrategiesMock = vi.fn();

vi.mock("../api/strategies", () => ({
  listStrategies: (...args: unknown[]) => listStrategiesMock(...args),
}));

describe("useStrategies", () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.useFakeTimers();
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  it("loads strategies and exposes loaded state", async () => {
    listStrategiesMock.mockResolvedValue({
      items: [
        create(StrategySchema, {
          id: "strategy-1",
          slug: "vol-breakout",
          name: "Volatility Breakout",
          description: "Breakout strategy",
        }),
      ],
      total: 1,
    });

    const { result } = renderHook(() => useStrategies({ pageSize: 25 }));

    await act(async () => {
      await vi.runAllTimersAsync();
      await Promise.resolve();
    });

    expect(result.current.state.status).toBe("loaded");
    expect(listStrategiesMock).toHaveBeenCalledWith({
      search: "",
      page: 1,
      pageSize: 25,
      sort: undefined,
    });
    expect(result.current.state.items).toHaveLength(1);
    expect(result.current.state.total).toBe(1);
  });

  it("maps permission denied to a user-friendly error", async () => {
    listStrategiesMock.mockRejectedValue(
      new ConnectError("forbidden", Code.PermissionDenied),
    );

    const { result } = renderHook(() => useStrategies());

    await act(async () => {
      await vi.runAllTimersAsync();
      await Promise.resolve();
    });

    expect(result.current.state.status).toBe("error");
    expect(result.current.state.error).toBe(
      "You do not have permission to view strategies.",
    );
  });

  it("debounces search updates and refreshes with the latest query", async () => {
    listStrategiesMock.mockResolvedValue({
      items: [],
      total: 0,
    });

    const { result } = renderHook(() => useStrategies());

    await act(async () => {
      await vi.runAllTimersAsync();
      await Promise.resolve();
    });

    listStrategiesMock.mockClear();

    act(() => {
      result.current.setQuery("risk");
    });

    expect(listStrategiesMock).not.toHaveBeenCalled();

    await act(async () => {
      await vi.advanceTimersByTimeAsync(250);
      await Promise.resolve();
    });

    expect(listStrategiesMock).toHaveBeenCalledWith({
      search: "risk",
      page: 1,
      pageSize: 12,
      sort: undefined,
    });
  });
});
