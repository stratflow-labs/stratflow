import type {
  AttributeValueInline,
  AttributeValueRelationInline,
  AttributeWithValues,
} from "@/shared/api/gen/strategy_registry/proto/v1/strategy_types_pb";

export type GraphDraftRelationRef = {
  attributeLocalId: string;
  valueLocalId: string;
};

export type GraphDraftValue = {
  localId: string;
  id?: string;
  attributeLocalId: string;
  slug: string;
  value: string;
  relations: GraphDraftRelationRef[];
};

export type GraphDraftAttribute = {
  localId: string;
  id?: string;
  strategyId?: string;
  slug: string;
  name: string;
  description: string;
  values: GraphDraftValue[];
};

export type StrategyRelationNode = {
  localAttributeId: string;
  localValueId: string;
  attributeId?: string;
  attributeName: string;
  attributeSlug: string;
  valueId?: string;
  valueSlug: string;
  valueLabel: string;
  relations: GraphDraftRelationRef[];
};

const buildAttributeLocalId = (attribute: AttributeWithValues): string =>
  attribute.id || `attr-slug:${attribute.slug}`;

const buildValueLocalId = (
  attributeLocalId: string,
  value: AttributeValueInline,
): string => value.id || `${attributeLocalId}::value-slug:${value.slug}`;

export const buildRelationSourceKey = (
  attributeLocalId: string,
  valueLocalId: string,
): string => `${attributeLocalId}::${valueLocalId}`;

export const relationSnapshotKey = (relations: GraphDraftRelationRef[]): string =>
  JSON.stringify(
    relations
      .map((relation) => buildRelationSourceKey(relation.attributeLocalId, relation.valueLocalId))
      .sort(),
  );

export const mapInlineRelationsToDraftRefs = (
  attributes: AttributeWithValues[],
  relations: AttributeValueRelationInline[],
): GraphDraftRelationRef[] => {
  const attributeLocalIdById = new Map<string, string>();
  const valueLocalIdById = new Map<string, { attributeLocalId: string; valueLocalId: string }>();

  attributes.forEach((attribute) => {
    const attributeLocalId = buildAttributeLocalId(attribute);
    attributeLocalIdById.set(attribute.id, attributeLocalId);

    attribute.values.forEach((value) => {
      valueLocalIdById.set(value.id, {
        attributeLocalId,
        valueLocalId: buildValueLocalId(attributeLocalId, value),
      });
    });
  });

  return relations.flatMap((relation) => {
    const target = valueLocalIdById.get(relation.toValueId);
    const attributeLocalId =
      attributeLocalIdById.get(relation.toAttributeId) ?? target?.attributeLocalId;

    if (!target || !attributeLocalId) {
      return [];
    }

    return [
      {
        attributeLocalId,
        valueLocalId: target.valueLocalId,
      },
    ];
  });
};

export const buildGraphDraft = (
  attributes: AttributeWithValues[],
): GraphDraftAttribute[] =>
  attributes.map((attribute) => {
    const attributeLocalId = buildAttributeLocalId(attribute);

    return {
      localId: attributeLocalId,
      id: attribute.id || undefined,
      strategyId: attribute.strategyId || undefined,
      slug: attribute.slug,
      name: attribute.name,
      description: attribute.description,
      values: attribute.values.map((value) => ({
        localId: buildValueLocalId(attributeLocalId, value),
        id: value.id || undefined,
        attributeLocalId,
        slug: value.slug,
        value: value.value,
        relations: mapInlineRelationsToDraftRefs(attributes, value.relations),
      })),
    };
  });

export const cloneGraphDraft = (
  draft: GraphDraftAttribute[],
): GraphDraftAttribute[] =>
  draft.map((attribute) => ({
    ...attribute,
    values: attribute.values.map((value) => ({
      ...value,
      relations: value.relations.map((relation) => ({ ...relation })),
    })),
  }));

export const buildRelationNodes = (
  draft: GraphDraftAttribute[],
): StrategyRelationNode[] =>
  draft.flatMap((attribute) =>
    attribute.values.map((value) => ({
      localAttributeId: attribute.localId,
      localValueId: value.localId,
      attributeId: attribute.id,
      attributeName: attribute.name,
      attributeSlug: attribute.slug,
      valueId: value.id,
      valueSlug: value.slug,
      valueLabel: value.value,
      relations: value.relations.map((relation) => ({ ...relation })),
    })),
  );

export const removeNodeRelations = (
  draft: GraphDraftAttribute[],
  targetValueLocalId: string,
): GraphDraftAttribute[] =>
  draft.map((attribute) => ({
    ...attribute,
    values: attribute.values.map((value) => ({
      ...value,
      relations: value.relations.filter(
        (relation) => relation.valueLocalId !== targetValueLocalId,
      ),
    })),
  }));

export const removeAttributeRelations = (
  draft: GraphDraftAttribute[],
  targetAttributeLocalId: string,
): GraphDraftAttribute[] =>
  draft.map((attribute) => ({
    ...attribute,
    values: attribute.values.map((value) => ({
      ...value,
      relations: value.relations.filter(
        (relation) => relation.attributeLocalId !== targetAttributeLocalId,
      ),
    })),
  }));

export const findAttributeDraft = (
  draft: GraphDraftAttribute[],
  attributeLocalId: string,
): GraphDraftAttribute | null =>
  draft.find((attribute) => attribute.localId === attributeLocalId) ?? null;

export const findValueDraft = (
  draft: GraphDraftAttribute[],
  valueLocalId: string,
): GraphDraftValue | null => {
  for (const attribute of draft) {
    const value = attribute.values.find((item) => item.localId === valueLocalId);
    if (value) {
      return value;
    }
  }

  return null;
};
