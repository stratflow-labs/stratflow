import type { AttributeWithValues } from "@/shared/api/gen/strategy_registry/proto/v1/strategy_types_pb";

import {
  buildGraphDraft,
  cloneGraphDraft,
  type GraphDraftAttribute,
} from "./relations";

export type AttributeLoadStatus = "idle" | "loading" | "loaded" | "error";

export type StrategyDetailsEntry = {
  status: AttributeLoadStatus;
  items: AttributeWithValues[];
  total: number;
  error: string | null;
};

export type GraphSaveState = {
  status: "idle" | "saving" | "error";
  error: string | null;
};

export type CreateDialogState =
  | { kind: "closed" }
  | { kind: "attribute" }
  | { kind: "value"; attributeLocalId: string };

export const INITIAL_ENTRY: StrategyDetailsEntry = {
  status: "idle",
  items: [],
  total: 0,
  error: null,
};

export const INITIAL_GRAPH_SAVE_STATE: GraphSaveState = {
  status: "idle",
  error: null,
};

export const buildLoadedEntry = (
  items: AttributeWithValues[],
  total: number,
): StrategyDetailsEntry => ({
  status: "loaded",
  items,
  total,
  error: null,
});

export const buildLoadedGraphState = (items: AttributeWithValues[]) => {
  const baseline = buildGraphDraft(items);
  return {
    baseline,
    draft: cloneGraphDraft(baseline),
  };
};

export const updateValueRelationsInDraft = (
  draft: GraphDraftAttribute[],
  valueLocalId: string,
  relations: { attributeLocalId: string; valueLocalId: string }[],
) =>
  draft.map((attribute) => ({
    ...attribute,
    values: attribute.values.map((value) =>
      value.localId === valueLocalId
        ? { ...value, relations: relations.map((relation) => ({ ...relation })) }
        : value,
    ),
  }));
