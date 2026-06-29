import { create } from "@bufbuild/protobuf";

import { strategyRegistryClient } from "@/shared/api/connect/clients";
import type {
  AttributeWithValues,
  StrategyGraphAction,
  Strategy,
  StrategyGraphResponse,
  UpdateAttributeValueRelationInput,
} from "@/shared/api/gen/strategy_registry/proto/v1/strategy_types_pb";
import {
  BatchActionStrategyGraphRequestSchema,
  CreateAttributeRequestSchema,
  CreateAttributeValueRequestSchema,
  CreateStrategyRequestSchema,
} from "@/shared/api/gen/strategy_registry/proto/v1/strategy_types_pb";
import {
  buildRelationTargetRefById,
  normalizeOptionalString,
  toListResult,
} from "./strategy-api-helpers";

export type ListStrategiesParams = {
  search?: string;
  page?: number;
  pageSize?: number;
  sort?: string;
};

export type StrategiesListResult = {
  items: Strategy[];
  total: number;
};

export type ListStrategyAttributesParams = {
  strategyRef: string;
  search?: string;
  page?: number;
  pageSize?: number;
  sort?: string;
};

export type StrategyAttributesListResult = {
  items: AttributeWithValues[];
  total: number;
};

export type ValueRelationsPatchInput = {
  attributeRef: string;
  valueRef: string;
  relations: UpdateAttributeValueRelationInput[];
};

export type BatchGraphActionInput = {
  strategyRef: string;
  actions: StrategyGraphAction[];
};

export type CreateStrategyInput = {
  slug: string;
  name: string;
  description: string;
};

export type CreateAttributeInput = {
  strategyRef: string;
  slug: string;
  name: string;
  description: string;
};

export type CreateAttributeValueInput = {
  strategyRef: string;
  attributeRef: string;
  slug: string;
  value: string;
};

export const listStrategies = async ({
  search,
  page = 1,
  pageSize = 12,
  sort,
}: ListStrategiesParams = {}): Promise<StrategiesListResult> => {
  const response = await strategyRegistryClient.listStrategies({
    search: normalizeOptionalString(search),
    page,
    pageSize,
    sort: normalizeOptionalString(sort),
  });

  return toListResult<Strategy>(response);
};

export const createStrategy = async ({
  slug,
  name,
  description,
}: CreateStrategyInput) => {
  const request = create(CreateStrategyRequestSchema, {
    slug,
    name,
    description,
  });

  return strategyRegistryClient.createStrategy(request);
};

export const listStrategyAttributes = async ({
  strategyRef,
  search,
  page = 1,
  pageSize = 100,
  sort,
}: ListStrategyAttributesParams): Promise<StrategyAttributesListResult> => {
  const response = await strategyRegistryClient.listAttributes({
    strategyRef,
    search: normalizeOptionalString(search),
    page,
    pageSize,
    sort: normalizeOptionalString(sort),
  });

  return toListResult<AttributeWithValues>(response);
};

export const createAttribute = async ({
  strategyRef,
  slug,
  name,
  description,
}: CreateAttributeInput) => {
  const request = create(CreateAttributeRequestSchema, {
    strategyRef,
    slug,
    name,
    description,
  });

  return strategyRegistryClient.createAttribute(request);
};

export const createAttributeValue = async ({
  strategyRef,
  attributeRef,
  slug,
  value,
}: CreateAttributeValueInput) => {
  const request = create(CreateAttributeValueRequestSchema, {
    strategyRef,
    attributeRef,
    slug,
    value,
    relations: [],
  });

  return strategyRegistryClient.createAttributeValue(request);
};

export const batchUpdateStrategyRelations = async ({
  strategyRef,
  items,
}: {
  strategyRef: string;
  items: ValueRelationsPatchInput[];
}): Promise<StrategyGraphResponse> => {
  return strategyRegistryClient.batchActionStrategyGraph({
    strategyRef,
    actions: items.map((item) => ({
      action: {
        case: "replaceRelations",
        value: {
          attributeRef: { id: item.attributeRef },
          valueRef: { id: item.valueRef },
          relations: item.relations.map((relation) =>
            buildRelationTargetRefById(relation.toAttributeId, relation.toValueId),
          ),
        },
      },
    })),
  });
};

export const batchActionStrategyGraph = async ({
  strategyRef,
  actions,
}: BatchGraphActionInput): Promise<StrategyGraphResponse> => {
  const request = create(BatchActionStrategyGraphRequestSchema, {
    strategyRef,
    actions,
  });

  return strategyRegistryClient.batchActionStrategyGraph(request);
};
