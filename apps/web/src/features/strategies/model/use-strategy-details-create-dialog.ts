"use client";

import { useCallback, useMemo, useState } from "react";

import type { Strategy } from "@/shared/api/gen/strategy_registry/proto/v1/strategy_types_pb";

import {
  createAttribute,
  createAttributeValue,
} from "../api/strategies";
import {
  findAttributeDraft,
  type GraphDraftAttribute,
} from "./relations";
import type { CreateDialogState } from "./strategy-details-dialog-helpers";

type UseStrategyDetailsCreateDialogParams = {
  selectedStrategy: Strategy | null;
  getStrategyKey: (strategy: Pick<Strategy, "id" | "slug">) => string;
  getErrorMessage: (error: unknown) => string;
  graphDraftByStrategyKey: Record<string, GraphDraftAttribute[]>;
  reloadStrategy: (strategy: Strategy) => Promise<void>;
};

export const useStrategyDetailsCreateDialog = ({
  selectedStrategy,
  getStrategyKey,
  getErrorMessage,
  graphDraftByStrategyKey,
  reloadStrategy,
}: UseStrategyDetailsCreateDialogParams) => {
  const [createDialogState, setCreateDialogState] = useState<CreateDialogState>({
    kind: "closed",
  });
  const [createError, setCreateError] = useState<string | null>(null);
  const [isCreating, setIsCreating] = useState(false);

  const addGraphAttribute = useCallback(() => {
    if (!selectedStrategy) {
      return;
    }

    setCreateError(null);
    setCreateDialogState({ kind: "attribute" });
  }, [selectedStrategy]);

  const addGraphValue = useCallback(
    (attributeLocalId: string) => {
      if (!selectedStrategy) {
        return;
      }

      setCreateError(null);
      setCreateDialogState({ kind: "value", attributeLocalId });
    },
    [selectedStrategy],
  );

  const openCreateAttributeDialog = useCallback(() => {
    if (!selectedStrategy) {
      return;
    }

    setCreateError(null);
    setCreateDialogState({ kind: "attribute" });
  }, [selectedStrategy]);

  const openCreateValueDialog = useCallback(
    (attributeRef: string) => {
      if (!selectedStrategy) {
        return;
      }

      const strategyKey = getStrategyKey(selectedStrategy);
      const attribute = (graphDraftByStrategyKey[strategyKey] ?? []).find(
        (item) => item.localId === attributeRef || item.id === attributeRef || item.slug === attributeRef,
      );
      if (!attribute) {
        return;
      }

      setCreateError(null);
      setCreateDialogState({ kind: "value", attributeLocalId: attribute.localId });
    },
    [getStrategyKey, graphDraftByStrategyKey, selectedStrategy],
  );

  const closeCreateDialog = useCallback(() => {
    if (isCreating) {
      return;
    }

    setCreateDialogState({ kind: "closed" });
    setCreateError(null);
  }, [isCreating]);

  const submitCreateDialog = useCallback(
    async ({
      slug,
      title,
      description,
    }: {
      slug: string;
      title: string;
      description: string;
    }) => {
      if (!selectedStrategy || createDialogState.kind === "closed") {
        return;
      }

      const strategyKey = getStrategyKey(selectedStrategy);
      setIsCreating(true);
      setCreateError(null);

      try {
        if (createDialogState.kind === "attribute") {
          await createAttribute({
            strategyRef: strategyKey,
            slug,
            name: title,
            description,
          });
        } else {
          const attribute = findAttributeDraft(
            graphDraftByStrategyKey[strategyKey] ?? [],
            createDialogState.attributeLocalId,
          );

          if (!attribute) {
            throw new Error("Attribute draft is missing.");
          }

          let attributeRef = attribute.id;

          if (!attributeRef) {
            const createdAttribute = await createAttribute({
              strategyRef: strategyKey,
              slug: attribute.slug,
              name: attribute.name,
              description: attribute.description,
            });
            attributeRef = createdAttribute.data?.id;
          }

          if (!attributeRef) {
            throw new Error("Failed to create attribute before creating value.");
          }

          await createAttributeValue({
            strategyRef: strategyKey,
            attributeRef,
            slug,
            value: title,
          });
        }

        setCreateDialogState({ kind: "closed" });
        await reloadStrategy(selectedStrategy);
      } catch (error) {
        setCreateError(getErrorMessage(error));
      } finally {
        setIsCreating(false);
      }
    },
    [
      createDialogState,
      getErrorMessage,
      getStrategyKey,
      graphDraftByStrategyKey,
      reloadStrategy,
      selectedStrategy,
    ],
  );

  const createDialogMode =
    createDialogState.kind === "closed" ? null : createDialogState.kind;

  const createDialogAttribute = useMemo(
    () =>
      createDialogState.kind === "value" && selectedStrategy
        ? findAttributeDraft(
            graphDraftByStrategyKey[getStrategyKey(selectedStrategy)] ?? [],
            createDialogState.attributeLocalId,
          )
        : null,
    [createDialogState, getStrategyKey, graphDraftByStrategyKey, selectedStrategy],
  );

  return {
    addGraphAttribute,
    addGraphValue,
    openCreateAttributeDialog,
    openCreateValueDialog,
    closeCreateDialog,
    submitCreateDialog,
    createDialogMode,
    createDialogAttribute,
    createError,
    isCreating,
    resetCreateDialog() {
      setCreateDialogState({ kind: "closed" });
      setCreateError(null);
    },
  };
};
