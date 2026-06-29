import { Position, type Node } from "@xyflow/react";

import type { StrategyRelationNode } from "../model/relations";

export const NODE_WIDTH = 232;
export const NODE_HEIGHT = 76;
const COLUMN_GAP = 80;
const ROW_GAP = 20;
const EDGE_GUTTER_OFFSET = 28;
const ROUTING_PADDING = 16;

export type AnchorSide = "left" | "right" | "top" | "bottom";
export type Point = { x: number; y: number };
export type Rect = { left: number; right: number; top: number; bottom: number };

const isHorizontalSegment = (from: Point, to: Point) => from.y === to.y;

export const getNodeBounds = (nodeX: number, nodeY: number): Rect => ({
  left: nodeX,
  right: nodeX + NODE_WIDTH,
  top: nodeY,
  bottom: nodeY + NODE_HEIGHT,
});

const rectContainsPoint = (rect: Rect, point: Point) =>
  point.x >= rect.left &&
  point.x <= rect.right &&
  point.y >= rect.top &&
  point.y <= rect.bottom;

const segmentIntersectsRect = (from: Point, to: Point, rect: Rect) => {
  if (isHorizontalSegment(from, to)) {
    const minX = Math.min(from.x, to.x);
    const maxX = Math.max(from.x, to.x);

    return (
      from.y > rect.top &&
      from.y < rect.bottom &&
      maxX > rect.left &&
      minX < rect.right
    );
  }

  const minY = Math.min(from.y, to.y);
  const maxY = Math.max(from.y, to.y);

  return (
    from.x > rect.left &&
    from.x < rect.right &&
    maxY > rect.top &&
    minY < rect.bottom
  );
};

const getHorizontalExitPoint = (anchor: Point, side: AnchorSide): Point => ({
  x:
    side === "left"
      ? anchor.x - EDGE_GUTTER_OFFSET
      : anchor.x + EDGE_GUTTER_OFFSET,
  y: anchor.y,
});

const getVerticalExitPoint = (anchor: Point, side: AnchorSide): Point => ({
  x: anchor.x,
  y:
    side === "top"
      ? anchor.y - EDGE_GUTTER_OFFSET
      : anchor.y + EDGE_GUTTER_OFFSET,
});

const buildPathFromPoints = (points: Point[]) =>
  points.map((point, index) => `${index === 0 ? "M" : "L"} ${point.x} ${point.y}`).join(" ");

export const buildOrthogonalPath = ({
  sourcePoint,
  sourceSide,
  targetPoint,
  targetSide,
  obstacles,
}: {
  sourcePoint: Point;
  sourceSide: AnchorSide;
  targetPoint: Point;
  targetSide: AnchorSide;
  obstacles: Rect[];
}) => {
  if (
    (sourceSide === "left" || sourceSide === "right") &&
    (targetSide === "left" || targetSide === "right")
  ) {
    const sourceExit = getHorizontalExitPoint(sourcePoint, sourceSide);
    const targetExit = getHorizontalExitPoint(targetPoint, targetSide);
    let corridorY = targetPoint.y;

    const colliding = obstacles.filter((obstacle) => {
      if (rectContainsPoint(obstacle, sourcePoint) || rectContainsPoint(obstacle, targetPoint)) {
        return false;
      }

      return (
        segmentIntersectsRect(sourceExit, { x: sourceExit.x, y: corridorY }, obstacle) ||
        segmentIntersectsRect(
          { x: sourceExit.x, y: corridorY },
          { x: targetExit.x, y: corridorY },
          obstacle,
        ) ||
        segmentIntersectsRect({ x: targetExit.x, y: corridorY }, targetExit, obstacle)
      );
    });

    if (colliding.length > 0) {
      corridorY =
        Math.max(
          targetPoint.y,
          ...colliding.map((obstacle) => obstacle.bottom + ROUTING_PADDING),
        );
    }

    return buildPathFromPoints([
      sourcePoint,
      sourceExit,
      { x: sourceExit.x, y: corridorY },
      { x: targetExit.x, y: corridorY },
      targetExit,
      targetPoint,
    ]);
  }

  const sourceExit =
    sourceSide === "top" || sourceSide === "bottom"
      ? getVerticalExitPoint(sourcePoint, sourceSide)
      : getHorizontalExitPoint(sourcePoint, sourceSide);

  const targetEntry =
    targetSide === "top" || targetSide === "bottom"
      ? getVerticalExitPoint(targetPoint, targetSide)
      : getHorizontalExitPoint(targetPoint, targetSide);

  return buildPathFromPoints([
    sourcePoint,
    sourceExit,
    { x: sourceExit.x, y: targetEntry.y },
    targetEntry,
    targetPoint,
  ]);
};

export const getNodeAnchorPoint = (
  nodeX: number,
  nodeY: number,
  side: AnchorSide,
) => {
  if (side === "left") {
    return { x: nodeX, y: nodeY + NODE_HEIGHT / 2 };
  }

  if (side === "right") {
    return { x: nodeX + NODE_WIDTH, y: nodeY + NODE_HEIGHT / 2 };
  }

  if (side === "top") {
    return { x: nodeX + NODE_WIDTH / 2, y: nodeY };
  }

  return { x: nodeX + NODE_WIDTH / 2, y: nodeY + NODE_HEIGHT };
};

export const getSideTowardPoint = (
  nodeX: number,
  nodeY: number,
  pointX: number,
  pointY: number,
): AnchorSide => {
  const centerX = nodeX + NODE_WIDTH / 2;
  const centerY = nodeY + NODE_HEIGHT / 2;
  const deltaX = pointX - centerX;
  const deltaY = pointY - centerY;

  if (Math.abs(deltaX) >= Math.abs(deltaY)) {
    return deltaX >= 0 ? "right" : "left";
  }

  return deltaY >= 0 ? "bottom" : "top";
};

export const applyColumnLayout = <T extends Record<string, unknown>>(
  nodes: Node<T, "relationNode">[],
  nodeById: Map<string, StrategyRelationNode>,
): Node<T, "relationNode">[] => {
  const attributeOrder: string[] = [];
  const attributeSeen = new Set<string>();

  nodeById.forEach((item) => {
    if (attributeSeen.has(item.localAttributeId)) {
      return;
    }

    attributeSeen.add(item.localAttributeId);
    attributeOrder.push(item.localAttributeId);
  });

  const nodesByAttribute = new Map<string, Node<T, "relationNode">[]>();
  nodes.forEach((node) => {
    const item = nodeById.get(node.id);
    if (!item) {
      return;
    }

    const current = nodesByAttribute.get(item.localAttributeId) ?? [];
    current.push(node);
    nodesByAttribute.set(item.localAttributeId, current);
  });

  return attributeOrder.flatMap((attributeId, columnIndex) =>
    (nodesByAttribute.get(attributeId) ?? []).map((node, rowIndex) => ({
      ...node,
      sourcePosition: Position.Right,
      targetPosition: Position.Left,
      position: {
        x: columnIndex * (NODE_WIDTH + COLUMN_GAP),
        y: rowIndex * (NODE_HEIGHT + ROW_GAP),
      },
      style: { width: NODE_WIDTH },
    })),
  );
};
