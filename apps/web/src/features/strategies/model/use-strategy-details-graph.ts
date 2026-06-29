"use client";

import { useCallback, useMemo, type Dispatch, type SetStateAction } from "react";

import type { Strategy } from "@/shared/api/gen/strategy_registry/proto/v1/strategy_types_pb";

import { batchActionStrategyGraph } from "../api/strategies";
import {
  buildRelationNodes,
  type GraphDraftAttribute,
  type GraphDraftRelationRef,
} from "./relations";
import {
  promptRenameAttribute,
  promptRenameValue,
} from "./strategy-details-graph-prompts";
import {
  clearGraphSourceIfAttributeRemoved,
  clearGraphSourceIfValueRemoved,
  deleteStrategyGraphAttribute,
  deleteStrategyGraphValue,
  renameStrategyGraphAttribute,
  renameStrategyGraphValue,
  updateStrategyGraphRelations,
} from "./strategy-details-graph-state";
import {
  buildLoadedEntry,
  buildLoadedGraphState,
  INITIAL_GRAPH_SAVE_STATE,
  type GraphSaveState,
  type StrategyDetailsEntry,
} from "./strategy-details-dialog-helpers";
import { buildGraphActions } from "./strategy-details-graph-actions";

type UseStrategyDetailsGraphParams = {
  selectedStrategy: Strategy | null;
  getStrategyKey: (strategy: Pick<Strategy, "id" | "slug">) => string;
  getErrorMessage: (error: unknown) => string;
  graphBaselineByStrategyKey: Record<string, GraphDraftAttribute[]>;
  setGraphBaselineByStrategyKey: Dispatch<
    SetStateAction<Record<string, GraphDraftAttribute[]>>
  >;
  graphDraftByStrategyKey: Record<string, GraphDraftAttribute[]>;
  setGraphDraftByStrategyKey: Dispatch<
    SetStateAction<Record<string, GraphDraftAttribute[]>>
  >;
  graphSourceKeyByStrategyKey: Record<string, string | null>;
  setGraphSourceKeyByStrategyKey: Dispatch<
    SetStateAction<Record<string, string | null>>
  >;
  graphOpenByStrategyKey: Record<string, boolean>;
  setGraphOpenByStrategyKey: Dispatch<
    SetStateAction<Record<string, boolean>>
  >;
  graphSaveStateByStrategyKey: Record<string, GraphSaveState>;
  setGraphSaveStateByStrategyKey: Dispatch<
    SetStateAction<Record<string, GraphSaveState>>
  >;
  setEntryByStrategyKey: Dispatch<
    SetStateAction<Record<string, StrategyDetailsEntry>>
  >;
};

export const useStrategyDetailsGraph = ({
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
}: UseStrategyDetailsGraphParams) => {
  const openRelationsGraph = useCallback(() => {
    if (!selectedStrategy) {
      return;
    }

    const strategyKey = getStrategyKey(selectedStrategy);
    setGraphSourceKeyByStrategyKey((current) => ({
      ...current,
      [strategyKey]: null,
    }));
    setGraphOpenByStrategyKey((current) => ({
      ...current,
      [strategyKey]: true,
    }));
    setGraphSaveStateByStrategyKey((current) => ({
      ...current,
      [strategyKey]: INITIAL_GRAPH_SAVE_STATE,
    }));
  }, [
    getStrategyKey,
    selectedStrategy,
    setGraphOpenByStrategyKey,
    setGraphSaveStateByStrategyKey,
    setGraphSourceKeyByStrategyKey,
  ]);

  const closeRelationsGraph = useCallback(() => {
    if (!selectedStrategy) {
      return;
    }

    const strategyKey = getStrategyKey(selectedStrategy);
    setGraphOpenByStrategyKey((current) => ({
      ...current,
      [strategyKey]: false,
    }));
    setGraphSourceKeyByStrategyKey((current) => ({
      ...current,
      [strategyKey]: null,
    }));
  }, [getStrategyKey, selectedStrategy, setGraphOpenByStrategyKey, setGraphSourceKeyByStrategyKey]);

  const selectGraphSource = useCallback(
    (sourceKey: string) => {
      if (!selectedStrategy) {
        return;
      }

      const strategyKey = getStrategyKey(selectedStrategy);
      setGraphSourceKeyByStrategyKey((current) => ({
        ...current,
        [strategyKey]: sourceKey,
      }));
    },
    [getStrategyKey, selectedStrategy, setGraphSourceKeyByStrategyKey],
  );

  const updateGraphRelations = useCallback(
    (sourceKey: string, relations: GraphDraftRelationRef[]) => {
      if (!selectedStrategy) {
        return;
      }

      const strategyKey = getStrategyKey(selectedStrategy);
      const [, valueLocalId] = sourceKey.split("::");
      if (!valueLocalId) {
        return;
      }

      setGraphDraftByStrategyKey((current) => ({
        ...current,
        [strategyKey]: updateStrategyGraphRelations(
          current[strategyKey] ?? [],
          sourceKey,
          relations,
        ),
      }));
    },
    [getStrategyKey, selectedStrategy, setGraphDraftByStrategyKey],
  );

  const renameGraphAttribute = useCallback(
    (attributeLocalId: string) => {
      if (!selectedStrategy) {
        return;
      }

      const strategyKey = getStrategyKey(selectedStrategy);
      const nextAttribute = promptRenameAttribute(
        graphDraftByStrategyKey[strategyKey] ?? [],
        attributeLocalId,
      );
      if (!nextAttribute) {
        return;
      }

      setGraphDraftByStrategyKey((current) => ({
        ...current,
        [strategyKey]: renameStrategyGraphAttribute(
          current[strategyKey] ?? [],
          attributeLocalId,
          nextAttribute,
        ),
      }));
    },
    [getStrategyKey, graphDraftByStrategyKey, selectedStrategy, setGraphDraftByStrategyKey],
  );

  const deleteGraphAttribute = useCallback(
    (attributeLocalId: string) => {
      if (!selectedStrategy) {
        return;
      }

      const strategyKey = getStrategyKey(selectedStrategy);
      setGraphDraftByStrategyKey((current) => ({
        ...current,
        [strategyKey]: deleteStrategyGraphAttribute(
          current[strategyKey] ?? [],
          attributeLocalId,
        ),
      }));
      setGraphSourceKeyByStrategyKey((current) => ({
        ...current,
        [strategyKey]: clearGraphSourceIfAttributeRemoved(
          current[strategyKey],
          attributeLocalId,
        ),
      }));
    },
    [getStrategyKey, selectedStrategy, setGraphDraftByStrategyKey, setGraphSourceKeyByStrategyKey],
  );

  const renameGraphValue = useCallback(
    (valueLocalId: string) => {
      if (!selectedStrategy) {
        return;
      }

      const strategyKey = getStrategyKey(selectedStrategy);
      const nextValue = promptRenameValue(
        graphDraftByStrategyKey[strategyKey] ?? [],
        valueLocalId,
      );
      if (!nextValue) {
        return;
      }

      setGraphDraftByStrategyKey((current) => ({
        ...current,
        [strategyKey]: renameStrategyGraphValue(
          current[strategyKey] ?? [],
          valueLocalId,
          nextValue,
        ),
      }));
    },
    [getStrategyKey, graphDraftByStrategyKey, selectedStrategy, setGraphDraftByStrategyKey],
  );

  const deleteGraphValue = useCallback(
    (valueLocalId: string) => {
      if (!selectedStrategy) {
        return;
      }

      const strategyKey = getStrategyKey(selectedStrategy);
      setGraphDraftByStrategyKey((current) => ({
        ...current,
        [strategyKey]: deleteStrategyGraphValue(
          current[strategyKey] ?? [],
          valueLocalId,
        ),
      }));
      setGraphSourceKeyByStrategyKey((current) => ({
        ...current,
        [strategyKey]: clearGraphSourceIfValueRemoved(
          current[strategyKey],
          valueLocalId,
        ),
      }));
    },
    [getStrategyKey, selectedStrategy, setGraphDraftByStrategyKey, setGraphSourceKeyByStrategyKey],
  );

  const saveGraphRelations = useCallback(async () => {
    if (!selectedStrategy) {
      return;
    }

    const strategyKey = getStrategyKey(selectedStrategy);
    const baseline = graphBaselineByStrategyKey[strategyKey] ?? [];
    const draft = graphDraftByStrategyKey[strategyKey] ?? [];
    const actions = buildGraphActions(baseline, draft);

    if (actions.length === 0) {
      return;
    }

    setGraphSaveStateByStrategyKey((current) => ({
      ...current,
      [strategyKey]: {
        status: "saving",
        error: null,
      },
    }));

    try {
      const response = await batchActionStrategyGraph({
        strategyRef: strategyKey,
        actions,
      });
      const items = response.data?.parameters ?? [];
      const graphState = buildLoadedGraphState(items);

      setEntryByStrategyKey((current) => ({
        ...current,
        [strategyKey]: buildLoadedEntry(items, items.length),
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
      setGraphSaveStateByStrategyKey((current) => ({
        ...current,
        [strategyKey]: {
          status: "error",
          error: getErrorMessage(error),
        },
      }));
    }
  }, [
    getErrorMessage,
    getStrategyKey,
    graphBaselineByStrategyKey,
    graphDraftByStrategyKey,
    selectedStrategy,
    setEntryByStrategyKey,
    setGraphBaselineByStrategyKey,
    setGraphDraftByStrategyKey,
    setGraphSaveStateByStrategyKey,
  ]);

  const selectedStrategyKey = selectedStrategy ? getStrategyKey(selectedStrategy) : null;
  const graphDraft = useMemo(
    () => (selectedStrategyKey ? graphDraftByStrategyKey[selectedStrategyKey] ?? [] : []),
    [graphDraftByStrategyKey, selectedStrategyKey],
  );
  const graphBaseline = useMemo(
    () => (selectedStrategyKey ? graphBaselineByStrategyKey[selectedStrategyKey] ?? [] : []),
    [graphBaselineByStrategyKey, selectedStrategyKey],
  );
  const relationNodes = useMemo(() => buildRelationNodes(graphDraft), [graphDraft]);
  const selectedGraphSourceKey = selectedStrategyKey
    ? (graphSourceKeyByStrategyKey[selectedStrategyKey] ?? null)
    : null;
  const isGraphOpen = selectedStrategyKey
    ? Boolean(graphOpenByStrategyKey[selectedStrategyKey])
    : false;
  const graphSaveState = selectedStrategyKey
    ? (graphSaveStateByStrategyKey[selectedStrategyKey] ?? INITIAL_GRAPH_SAVE_STATE)
    : INITIAL_GRAPH_SAVE_STATE;
  const dirtyGraphSourceCount = buildGraphActions(graphBaseline, graphDraft).length;

  return {
    graphDraft,
    relationNodes,
    selectedGraphSourceKey,
    isGraphOpen,
    graphSaveState,
    dirtyGraphSourceCount,
    openRelationsGraph,
    closeRelationsGraph,
    selectGraphSource,
    updateGraphRelations,
    renameGraphAttribute,
    deleteGraphAttribute,
    renameGraphValue,
    deleteGraphValue,
    saveGraphRelations,
  };
};
