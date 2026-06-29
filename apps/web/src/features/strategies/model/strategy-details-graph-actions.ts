"use client";

import { create } from "@bufbuild/protobuf";

import type {
  StrategyGraphAction,
} from "@/shared/api/gen/strategy_registry/proto/v1/strategy_types_pb";
import {
  GraphActionCreateAttributeSchema,
  GraphActionCreateValueSchema,
  GraphActionDeleteAttributeSchema,
  GraphActionDeleteValueSchema,
  GraphActionReplaceRelationsSchema,
  GraphActionUpdateAttributeSchema,
  GraphActionUpdateValueSchema,
  GraphAttributeRefSchema,
  GraphRelationTargetRefSchema,
  GraphValueRefSchema,
  StrategyGraphActionSchema,
} from "@/shared/api/gen/strategy_registry/proto/v1/strategy_types_pb";

import {
  findValueDraft,
  relationSnapshotKey,
  type GraphDraftAttribute,
} from "./relations";

const buildGraphAttributeRef = (id?: string, slug?: string) =>
  id?.trim()
    ? create(GraphAttributeRefSchema, { id })
    : slug?.trim()
      ? create(GraphAttributeRefSchema, { slug })
      : undefined;

const buildGraphValueRef = (id?: string, slug?: string) =>
  id?.trim()
    ? create(GraphValueRefSchema, { id })
    : slug?.trim()
      ? create(GraphValueRefSchema, { slug })
      : undefined;

export const buildGraphActions = (
  baseline: GraphDraftAttribute[],
  draft: GraphDraftAttribute[],
): StrategyGraphAction[] => {
  const actions: StrategyGraphAction[] = [];

  const baselineAttributesByLocalId = new Map(
    baseline.map((attribute) => [attribute.localId, attribute] as const),
  );
  const draftAttributesByLocalId = new Map(
    draft.map((attribute) => [attribute.localId, attribute] as const),
  );

  draft.forEach((attribute) => {
    const baselineAttribute = baselineAttributesByLocalId.get(attribute.localId);
    if (!baselineAttribute) {
      actions.push({
        ...create(StrategyGraphActionSchema, {
          action: {
            case: "createAttribute",
            value: create(GraphActionCreateAttributeSchema, {
              slug: attribute.slug,
              name: attribute.name,
              description: attribute.description,
            }),
          },
        }),
      });
      return;
    }

    if (
      baselineAttribute.slug !== attribute.slug ||
      baselineAttribute.name !== attribute.name ||
      baselineAttribute.description !== attribute.description
    ) {
      actions.push({
        ...create(StrategyGraphActionSchema, {
          action: {
            case: "updateAttribute",
            value: create(GraphActionUpdateAttributeSchema, {
              attributeRef: buildGraphAttributeRef(attribute.id, baselineAttribute.slug),
              slug: baselineAttribute.slug !== attribute.slug ? attribute.slug : undefined,
              name: baselineAttribute.name !== attribute.name ? attribute.name : undefined,
              description:
                baselineAttribute.description !== attribute.description
                  ? attribute.description
                  : undefined,
            }),
          },
        }),
      });
    }
  });

  draft.forEach((attribute) => {
    const baselineAttribute = baselineAttributesByLocalId.get(attribute.localId);
    const baselineValuesByLocalId = new Map(
      (baselineAttribute?.values ?? []).map((value) => [value.localId, value] as const),
    );

    attribute.values.forEach((value) => {
      const baselineValue = baselineValuesByLocalId.get(value.localId);
      if (!baselineValue) {
        actions.push({
          ...create(StrategyGraphActionSchema, {
            action: {
              case: "createValue",
              value: create(GraphActionCreateValueSchema, {
                attributeRef: buildGraphAttributeRef(attribute.id, attribute.slug),
                slug: value.slug,
                value: value.value,
              }),
            },
          }),
        });
        return;
      }

      if (baselineValue.slug !== value.slug || baselineValue.value !== value.value) {
        actions.push({
          ...create(StrategyGraphActionSchema, {
            action: {
              case: "updateValue",
              value: create(GraphActionUpdateValueSchema, {
                attributeRef: buildGraphAttributeRef(attribute.id, attribute.slug),
                valueRef: buildGraphValueRef(value.id, baselineValue.slug),
                slug: baselineValue.slug !== value.slug ? value.slug : undefined,
                value: baselineValue.value !== value.value ? value.value : undefined,
              }),
            },
          }),
        });
      }
    });
  });

  draft.forEach((attribute) => {
    const baselineAttribute = baselineAttributesByLocalId.get(attribute.localId);
    const baselineValuesByLocalId = new Map(
      (baselineAttribute?.values ?? []).map((value) => [value.localId, value] as const),
    );

    attribute.values.forEach((value) => {
      const baselineValue = baselineValuesByLocalId.get(value.localId);
      const baselineRelations = baselineValue?.relations ?? [];

      if (relationSnapshotKey(baselineRelations) === relationSnapshotKey(value.relations)) {
        return;
      }

      const relationTargets = value.relations
        .map((relation) => {
          const targetAttribute = draftAttributesByLocalId.get(relation.attributeLocalId);
          const targetValue = findValueDraft(draft, relation.valueLocalId);

          if (!targetAttribute || !targetValue) {
            return null;
          }

          return create(GraphRelationTargetRefSchema, {
            attributeRef: buildGraphAttributeRef(
              targetAttribute.id,
              targetAttribute.slug,
            ),
            valueRef: buildGraphValueRef(targetValue.id, targetValue.slug),
          });
        })
        .filter((item): item is NonNullable<typeof item> => item !== null);

      actions.push({
        ...create(StrategyGraphActionSchema, {
          action: {
            case: "replaceRelations",
            value: create(GraphActionReplaceRelationsSchema, {
              attributeRef: buildGraphAttributeRef(attribute.id, attribute.slug),
              valueRef: buildGraphValueRef(value.id, value.slug),
              relations: relationTargets,
            }),
          },
        }),
      });
    });
  });

  baseline.forEach((attribute) => {
    const draftAttribute = draftAttributesByLocalId.get(attribute.localId);
    if (!draftAttribute) {
      actions.push({
        ...create(StrategyGraphActionSchema, {
          action: {
            case: "deleteAttribute",
            value: create(GraphActionDeleteAttributeSchema, {
              attributeRef: buildGraphAttributeRef(attribute.id, attribute.slug),
            }),
          },
        }),
      });
      return;
    }

    const draftValuesByLocalId = new Map(
      draftAttribute.values.map((value) => [value.localId, value] as const),
    );

    attribute.values.forEach((value) => {
      if (!draftValuesByLocalId.has(value.localId)) {
        actions.push({
          ...create(StrategyGraphActionSchema, {
            action: {
              case: "deleteValue",
              value: create(GraphActionDeleteValueSchema, {
                attributeRef: buildGraphAttributeRef(attribute.id, attribute.slug),
                valueRef: buildGraphValueRef(value.id, value.slug),
              }),
            },
          }),
        });
      }
    });
  });

  return actions;
};
