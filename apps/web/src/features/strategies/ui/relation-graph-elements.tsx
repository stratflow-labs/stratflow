import {
  BaseEdge,
  Handle,
  MarkerType,
  Position,
  type Edge,
  type EdgeProps,
  type Node,
  type NodeProps,
} from "@xyflow/react";
import { memo } from "react";

import styles from "./relation-graph-editor.module.css";
import {
  buildOrthogonalPath,
  type AnchorSide,
  type Rect,
} from "./relation-graph-routing";

export type RelationNodeData = {
  attributeSlug: string;
  valueLabel: string;
  isSource: boolean;
  isConnected: boolean;
  isSelected: boolean;
} & Record<string, unknown>;

export type RelationFlowNode = Node<RelationNodeData, "relationNode">;
export type RelationEdgeData = {
  sourceSide: AnchorSide;
  targetSide: AnchorSide;
  obstacles: Rect[];
};
export type RelationFlowEdge = Edge<RelationEdgeData, "relationEdge">;

const RelationNodeView = memo(({ data }: NodeProps<RelationFlowNode>) => {
  const className = [
    styles.node,
    data.isSource ? styles.nodeSource : "",
    !data.isSource && data.isConnected ? styles.nodeConnected : "",
    data.isSelected ? styles.nodeSelected : "",
  ]
    .filter(Boolean)
    .join(" ");

  return (
    <div className={className}>
      <Handle
        type="target"
        position={Position.Left}
        style={{
          width: 12,
          height: "100%",
          top: 0,
          opacity: 0,
          borderRadius: 0,
          transform: "none",
        }}
      />
      <div className={styles.attribute}>{data.attributeSlug}</div>
      <div className={styles.value}>{data.valueLabel}</div>
      <Handle
        type="source"
        position={Position.Right}
        style={{
          width: 12,
          height: "100%",
          top: 0,
          opacity: 0,
          borderRadius: 0,
          transform: "none",
        }}
      />
    </div>
  );
});

RelationNodeView.displayName = "RelationNodeView";

const RelationEdgeView = memo((props: EdgeProps<RelationFlowEdge>) => {
  const { id, sourceX, sourceY, targetX, targetY, markerEnd, data } = props;
  const path = buildOrthogonalPath({
    sourcePoint: { x: sourceX, y: sourceY },
    sourceSide: (data?.sourceSide as AnchorSide | undefined) ?? "right",
    targetPoint: { x: targetX, y: targetY },
    targetSide: (data?.targetSide as AnchorSide | undefined) ?? "left",
    obstacles: (data?.obstacles as Rect[] | undefined) ?? [],
  });

  return (
    <BaseEdge
      id={id}
      path={path}
      markerEnd={markerEnd}
      style={props.style ?? { stroke: "#2e7d32", strokeWidth: 1.8 }}
    />
  );
});

RelationEdgeView.displayName = "RelationEdgeView";

export const relationNodeTypes = { relationNode: RelationNodeView };

export const relationEdgeTypes = { relationEdge: RelationEdgeView };

export const defaultRelationMarker = (selected: boolean) => ({
  type: MarkerType.ArrowClosed,
  color: selected ? "#1565c0" : "#2e7d32",
});

export const defaultRelationStyle = (selected: boolean) =>
  selected
    ? { stroke: "#1565c0", strokeWidth: 3 }
    : { stroke: "#2e7d32", strokeWidth: 1.8 };
