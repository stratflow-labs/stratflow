"use client";

import { Code, ConnectError } from "@connectrpc/connect";
import { useCallback, useMemo, useState } from "react";

import type {
  Strategy,
} from "@/shared/api/gen/strategy_registry/proto/v1/strategy_types_pb";

import {
  listStrategyAttributes,
} from "../api/strategies";
import {
  type GraphDraftAttribute,
} from "./relations";
import {
  buildLoadedEntry,
  buildLoadedGraphState,
  INITIAL_ENTRY,
  INITIAL_GRAPH_SAVE_STATE,
  type GraphSaveState,
  type StrategyDetailsEntry,
} from "./strategy-details-dialog-helpers";
import { useStrategyDetailsCreateDialog } from "./use-strategy-details-create-dialog";
import { useStrategyDetailsGraph } from "./use-strategy-details-graph";

const getStrategyKey = (strategy: Pick<Strategy, "id" | "slug">): string =>
  strategy.id || strategy.slug;

const getErrorMessage = (error: unknown): string => {
  if (error instanceof ConnectError) {
    if (error.code === Code.Unavailable || error.code === Code.Unknown) {
      return "Strategy attributes are unavailable. Check that the service is running.";
    }

    if (error.code === Code.PermissionDenied) {
      return "You do not have permission to view strategy attributes.";
    }
  }

  if (error instanceof Error && error.message) {
    return error.message;
  }

  return "Failed to load strategy attributes.";
};

const loadStrategyGraph = async (strategyKey: string) => {
  const result = await listStrategyAttributes({
    strategyRef: strategyKey,
    page: 1,
    pageSize: 100,
    sort: "created_at_desc",
  });

  return {
    entry: buildLoadedEntry(result.items, result.total),
    graphState: buildLoadedGraphState(result.items),
  };
};

export const useStrategyDetailsDialog = () => {
  const [selectedStrategy, setSelectedStrategy] = useState<Strategy | null>(null);
  const [entryByStrategyKey, setEntryByStrategyKey] = useState<
    Record<string, StrategyDetailsEntry>
  >({});
  const [expandedAttributeIdsByStrategyKey, setExpandedAttributeIdsByStrategyKey] =
    useState<Record<string, string | null>>({});
  const [graphBaselineByStrategyKey, setGraphBaselineByStrategyKey] = useState<
    Record<string, GraphDraftAttribute[]>
  >({});
  const [graphDraftByStrategyKey, setGraphDraftByStrategyKey] = useState<
    Record<string, GraphDraftAttribute[]>
  >({});
  const [graphSourceKeyByStrategyKey, setGraphSourceKeyByStrategyKey] = useState<
    Record<string, string | null>
  >({});
  const [graphOpenByStrategyKey, setGraphOpenByStrategyKey] = useState<
    Record<string, boolean>
  >({});
  const [graphSaveStateByStrategyKey, setGraphSaveStateByStrategyKey] = useState<
    Record<string, GraphSaveState>
  >({});

  const openStrategy = useCallback(async (strategy: Strategy) => {
    const strategyKey = getStrategyKey(strategy);
    if (!strategyKey) {
      return;
    }

    setSelectedStrategy(strategy);
    setEntryByStrategyKey((current) => ({
      ...current,
      [strategyKey]: {
        ...(current[strategyKey] ?? INITIAL_ENTRY),
        status: "loading",
        error: null,
      },
    }));

    try {
      const { entry, graphState } = await loadStrategyGraph(strategyKey);

      setEntryByStrategyKey((current) => ({
        ...current,
        [strategyKey]: entry,
      }));
      setGraphBaselineByStrategyKey((current) => ({
        ...current,
        [strategyKey]: graphState.baseline,
      }));
      setGraphDraftByStrategyKey((current) => ({
        ...current,
        [strategyKey]: graphState.draft,
      }));
      setGraphSaveStateByStrategyKey((current) => ({
        ...current,
        [strategyKey]: INITIAL_GRAPH_SAVE_STATE,
      }));
    } catch (error) {
      setEntryByStrategyKey((current) => ({
        ...current,
        [strategyKey]: {
          ...(current[strategyKey] ?? INITIAL_ENTRY),
          status: "error",
          error: getErrorMessage(error),
        },
      }));
    }
  }, []);

  const createDialog = useStrategyDetailsCreateDialog({
    selectedStrategy,
    getStrategyKey,
    getErrorMessage,
    graphDraftByStrategyKey,
    reloadStrategy: openStrategy,
  });

  const graph = useStrategyDetailsGraph({
    selectedStrategy,
    getStrategyKey,
    getErrorMessage,
    graphBaselineByStrategyKey,
    setGraphBaselineByStrategyKey,
    graphDraftByStrategyKey,
    setGraphDraftByStrategyKey,
    graphSourceKeyByStrategyKey,
    setGraphSourceKeyByStrategyKey,
    graphOpenByStrategyKey,
    setGraphOpenByStrategyKey,
    graphSaveStateByStrategyKey,
    setGraphSaveStateByStrategyKey,
    setEntryByStrategyKey,
  });

  const closeStrategy = useCallback(() => {
    setSelectedStrategy(null);
    createDialog.resetCreateDialog();
  }, [createDialog]);

  const refreshSelectedStrategy = useCallback(async () => {
    if (!selectedStrategy) {
      return;
    }

    await openStrategy(selectedStrategy);
  }, [openStrategy, selectedStrategy]);

  const toggleAttribute = useCallback(
    (attributeId: string) => {
      if (!selectedStrategy) {
        return;
      }

      const strategyKey = getStrategyKey(selectedStrategy);
      setExpandedAttributeIdsByStrategyKey((current) => ({
        ...current,
        [strategyKey]:
          current[strategyKey] === attributeId ? null : attributeId,
      }));
    },
    [selectedStrategy],
  );
  const selectedStrategyKey = selectedStrategy ? getStrategyKey(selectedStrategy) : null;

  const entry = selectedStrategyKey
    ? (entryByStrategyKey[selectedStrategyKey] ?? INITIAL_ENTRY)
    : INITIAL_ENTRY;

  const expandedAttributeId = selectedStrategyKey
    ? (expandedAttributeIdsByStrategyKey[selectedStrategyKey] ?? null)
    : null;

  const isOpen = Boolean(selectedStrategy);

  return useMemo(
    () => ({
      isOpen,
      selectedStrategy,
      entry,
      expandedAttributeId,
      relationNodes: graph.relationNodes,
      graphDraft: graph.graphDraft,
      selectedGraphSourceKey: graph.selectedGraphSourceKey,
      isGraphOpen: graph.isGraphOpen,
      graphSaveState: graph.graphSaveState,
      dirtyGraphSourceCount: graph.dirtyGraphSourceCount,
      createDialogMode: createDialog.createDialogMode,
      createDialogAttribute: createDialog.createDialogAttribute,
      createError: createDialog.createError,
      isCreating: createDialog.isCreating,
      openStrategy,
      closeStrategy,
      refreshSelectedStrategy,
      toggleAttribute,
      openRelationsGraph: graph.openRelationsGraph,
      closeRelationsGraph: graph.closeRelationsGraph,
      openCreateAttributeDialog: createDialog.openCreateAttributeDialog,
      openCreateValueDialog: createDialog.openCreateValueDialog,
      closeCreateDialog: createDialog.closeCreateDialog,
      submitCreateDialog: createDialog.submitCreateDialog,
      selectGraphSource: graph.selectGraphSource,
      updateGraphRelations: graph.updateGraphRelations,
      addGraphAttribute: createDialog.addGraphAttribute,
      renameGraphAttribute: graph.renameGraphAttribute,
      deleteGraphAttribute: graph.deleteGraphAttribute,
      addGraphValue: createDialog.addGraphValue,
      renameGraphValue: graph.renameGraphValue,
      deleteGraphValue: graph.deleteGraphValue,
      saveGraphRelations: graph.saveGraphRelations,
    }),
    [
      closeStrategy,
      createDialog,
      entry,
      expandedAttributeId,
      graph,
      isOpen,
      openStrategy,
      refreshSelectedStrategy,
      selectedStrategy,
      toggleAttribute,
    ],
  );
};
