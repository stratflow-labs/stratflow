import { create } from "@bufbuild/protobuf";

import type {
  GraphRelationTargetRef,
} from "@/shared/api/gen/strategy_registry/proto/v1/strategy_types_pb";
import {
  GraphAttributeRefSchema,
  GraphRelationTargetRefSchema,
  GraphValueRefSchema,
} from "@/shared/api/gen/strategy_registry/proto/v1/strategy_types_pb";

export const normalizeOptionalString = (value?: string): string | undefined => {
  const normalized = value?.trim();
  return normalized ? normalized : undefined;
};

export const toSafeTotal = (value: bigint): number => {
  if (value <= 0n) {
    return 0;
  }

  if (value > BigInt(Number.MAX_SAFE_INTEGER)) {
    return Number.MAX_SAFE_INTEGER;
  }

  return Number(value);
};

export const toListResult = <T>(response: {
  data?: { items?: T[]; total?: bigint };
}) => ({
  items: response.data?.items ?? [],
  total: toSafeTotal(response.data?.total ?? 0n),
});

export const buildRelationTargetRefById = (
  attributeId: string,
  valueId: string,
): GraphRelationTargetRef => ({
  ...create(GraphRelationTargetRefSchema, {
    attributeRef: create(GraphAttributeRefSchema, { id: attributeId }),
    valueRef: create(GraphValueRefSchema, { id: valueId }),
  }),
});
