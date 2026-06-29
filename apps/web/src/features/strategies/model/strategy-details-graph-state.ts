"use client";

import {
  removeAttributeRelations,
  removeNodeRelations,
  type GraphDraftAttribute,
  type GraphDraftRelationRef,
} from "./relations";
import { updateValueRelationsInDraft } from "./strategy-details-dialog-helpers";

export const updateStrategyGraphRelations = (
  draft: GraphDraftAttribute[],
  sourceKey: string,
  relations: GraphDraftRelationRef[],
): GraphDraftAttribute[] => {
  const [, valueLocalId] = sourceKey.split("::");
  if (!valueLocalId) {
    return draft;
  }

  return updateValueRelationsInDraft(draft, valueLocalId, relations);
};

export const renameStrategyGraphAttribute = (
  draft: GraphDraftAttribute[],
  attributeLocalId: string,
  nextAttribute: {
    slug: string;
    name: string;
    description: string;
  },
): GraphDraftAttribute[] =>
  draft.map((item) =>
    item.localId === attributeLocalId
      ? {
          ...item,
          slug: nextAttribute.slug,
          name: nextAttribute.name,
          description: nextAttribute.description,
        }
      : item,
  );

export const deleteStrategyGraphAttribute = (
  draft: GraphDraftAttribute[],
  attributeLocalId: string,
): GraphDraftAttribute[] =>
  removeAttributeRelations(
    draft.filter((attribute) => attribute.localId !== attributeLocalId),
    attributeLocalId,
  );

export const renameStrategyGraphValue = (
  draft: GraphDraftAttribute[],
  valueLocalId: string,
  nextValue: {
    slug: string;
    value: string;
  },
): GraphDraftAttribute[] =>
  draft.map((attribute) => ({
    ...attribute,
    values: attribute.values.map((item) =>
      item.localId === valueLocalId
        ? {
            ...item,
            slug: nextValue.slug,
            value: nextValue.value,
          }
        : item,
    ),
  }));

export const deleteStrategyGraphValue = (
  draft: GraphDraftAttribute[],
  valueLocalId: string,
): GraphDraftAttribute[] =>
  removeNodeRelations(
    draft.map((attribute) => ({
      ...attribute,
      values: attribute.values.filter((value) => value.localId !== valueLocalId),
    })),
    valueLocalId,
  );

export const clearGraphSourceIfAttributeRemoved = (
  sourceKey: string | null | undefined,
  attributeLocalId: string,
): string | null =>
  sourceKey?.startsWith(`${attributeLocalId}::`) ? null : (sourceKey ?? null);

export const clearGraphSourceIfValueRemoved = (
  sourceKey: string | null | undefined,
  valueLocalId: string,
): string | null =>
  sourceKey?.endsWith(`::${valueLocalId}`) ? null : (sourceKey ?? null);
