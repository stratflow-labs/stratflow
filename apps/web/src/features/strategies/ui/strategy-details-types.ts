import type {
  AttributeWithValues,
  Strategy,
} from "@/shared/api/gen/strategy_registry/proto/v1/strategy_types_pb";

export type StrategyDetailsEntry = {
  status: "idle" | "loading" | "loaded" | "error";
  items: AttributeWithValues[];
  total: number;
  error: string | null;
};

export type StrategyDetailsDialogProps = {
  open: boolean;
  strategy: Strategy | null;
  entry: StrategyDetailsEntry;
  expandedAttributeId: string | null;
  onClose: () => void;
  onRefresh: () => void | Promise<void>;
  onToggleAttribute: (attributeId: string) => void;
  onOpenRelationsGraph: () => void;
  onCreateAttribute: () => void;
  onCreateValue: (attributeLocalId: string) => void;
};
