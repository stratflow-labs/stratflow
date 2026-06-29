"use client";

import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import type { OnReconnect } from "@xyflow/react";

import {
  buildRelationSourceKey,
  type GraphDraftAttribute,
  type GraphDraftRelationRef,
  type StrategyRelationNode,
} from "../model/relations";
import { buildRelationGraphModel } from "./relation-graph-model";
import type { RelationFlowEdge, RelationFlowNode } from "./relation-graph-elements";
import {
  buildOrthogonalPath,
  getNodeAnchorPoint,
  getNodeBounds,
  getSideTowardPoint,
  type AnchorSide,
} from "./relation-graph-routing";

type UseRelationGraphEditorParams = {
  sourceAttributeLocalId: string | null;
  sourceValueLocalId: string | null;
  nodes: StrategyRelationNode[];
  graphDraft: GraphDraftAttribute[];
  onSelectSource: (sourceKey: string) => void;
  onChangeSourceRelations: (
    sourceKey: string,
    relations: GraphDraftRelationRef[],
  ) => void;
};

export const useRelationGraphEditor = ({
  sourceAttributeLocalId,
  sourceValueLocalId,
  nodes,
  graphDraft,
  onSelectSource,
  onChangeSourceRelations,
}: UseRelationGraphEditorParams) => {
  const containerRef = useRef<HTMLDivElement | null>(null);
  const [armedSourceNodeId, setArmedSourceNodeId] = useState<string | null>(null);
  const [contextNodeId, setContextNodeId] = useState<string | null>(null);
  const [selectedEdgeId, setSelectedEdgeId] = useState<string | null>(null);
  const [cursorPosition, setCursorPosition] = useState<{ x: number; y: number } | null>(
    null,
  );
  const [menuPosition, setMenuPosition] = useState<{
    mouseX: number;
    mouseY: number;
  } | null>(null);

  const graphDraftByValueLocalId = useMemo(() => {
    const map = new Map<string, GraphDraftAttribute["values"][number]>();
    graphDraft.forEach((attribute) => {
      attribute.values.forEach((value) => {
        map.set(value.localId, value);
      });
    });
    return map;
  }, [graphDraft]);

  const selectedSourceRelations = useMemo(
    () =>
      sourceValueLocalId
        ? graphDraftByValueLocalId.get(sourceValueLocalId)?.relations ?? []
        : [],
    [graphDraftByValueLocalId, sourceValueLocalId],
  );

  const graphModel = useMemo(
    () =>
      buildRelationGraphModel({
        nodes,
        graphDraft,
        selectedSourceRelations,
        sourceAttributeLocalId,
        sourceValueLocalId,
        armedSourceNodeId,
        contextNodeId,
        selectedEdgeId,
      }),
    [
      nodes,
      graphDraft,
      selectedSourceRelations,
      sourceAttributeLocalId,
      sourceValueLocalId,
      armedSourceNodeId,
      contextNodeId,
      selectedEdgeId,
    ],
  );

  const effectiveArmedSourceNodeId =
    armedSourceNodeId && graphModel.nodeById.has(armedSourceNodeId)
      ? armedSourceNodeId
      : null;
  const effectiveContextNodeId =
    contextNodeId && graphModel.nodeById.has(contextNodeId)
      ? contextNodeId
      : null;

  const contextNode = effectiveContextNodeId
    ? (graphModel.nodeById.get(effectiveContextNodeId) ?? null)
    : null;
  const armedSourceFlowNode = effectiveArmedSourceNodeId
    ? (graphModel.flowNodes.find((node) => node.id === effectiveArmedSourceNodeId) ??
      null)
    : null;
  const contextNodeIsSource =
    graphModel.sourceNodeId !== null &&
    effectiveContextNodeId === graphModel.sourceNodeId;
  const contextNodeIsConnected = effectiveContextNodeId
    ? graphModel.connectedNodeIds.has(effectiveContextNodeId)
    : false;

  const updateRelations = useCallback(
    (sourceKey: string, nextRelations: GraphDraftRelationRef[]) => {
      const [currentSourceAttributeLocalId, currentSourceValueLocalId] =
        sourceKey.split("::");
      if (!currentSourceAttributeLocalId || !currentSourceValueLocalId) {
        return;
      }

      onChangeSourceRelations(
        sourceKey,
        nextRelations
          .filter(
            (relation) =>
              relation.attributeLocalId !== currentSourceAttributeLocalId ||
              relation.valueLocalId !== currentSourceValueLocalId,
          )
          .filter(
            (relation, index, current) =>
              current.findIndex(
                (item) =>
                  item.attributeLocalId === relation.attributeLocalId &&
                  item.valueLocalId === relation.valueLocalId,
              ) === index,
          ),
      );
    },
    [onChangeSourceRelations],
  );

  const closeNodeMenu = useCallback(() => {
    setMenuPosition(null);
  }, []);

  const updateCursorPosition = useCallback((clientX: number, clientY: number) => {
    const container = containerRef.current;
    if (!container) {
      return;
    }

    const bounds = container.getBoundingClientRect();
    setCursorPosition({
      x: clientX - bounds.left,
      y: clientY - bounds.top,
    });
  }, []);

  const addRelation = useCallback(
    (sourceKey: string, relation: GraphDraftRelationRef) => {
      const [, sourceValueLocalIdForRelation] = sourceKey.split("::");
      const currentRelations =
        (sourceValueLocalIdForRelation
          ? graphDraftByValueLocalId.get(sourceValueLocalIdForRelation)?.relations
          : undefined) ?? [];
      updateRelations(sourceKey, [...currentRelations, relation]);
    },
    [graphDraftByValueLocalId, updateRelations],
  );

  const removeRelation = useCallback(
    (sourceKey: string, relation: GraphDraftRelationRef) => {
      const [, sourceValueLocalIdForRelation] = sourceKey.split("::");
      const currentRelations =
        (sourceValueLocalIdForRelation
          ? graphDraftByValueLocalId.get(sourceValueLocalIdForRelation)?.relations
          : undefined) ?? [];
      updateRelations(
        sourceKey,
        currentRelations.filter(
          (item) =>
            item.attributeLocalId !== relation.attributeLocalId ||
            item.valueLocalId !== relation.valueLocalId,
        ),
      );
    },
    [graphDraftByValueLocalId, updateRelations],
  );

  const removeRelationByEdgeId = useCallback(
    (edgeId: string) => {
      const [sourcePart, targetPart] = edgeId.replace("rel:", "").split("->");
      if (!sourcePart || !targetPart) {
        return;
      }

      const [targetAttributeId, targetValueId] = targetPart.split("::");
      if (!targetAttributeId || !targetValueId) {
        return;
      }

      onSelectSource(sourcePart);
      removeRelation(sourcePart, {
        attributeLocalId: targetAttributeId,
        valueLocalId: targetValueId,
      });
    },
    [onSelectSource, removeRelation],
  );

  const handleUseAsSource = useCallback(() => {
    if (!contextNode) {
      closeNodeMenu();
      return;
    }

    const sourceKey = buildRelationSourceKey(
      contextNode.localAttributeId,
      contextNode.localValueId,
    );
    onSelectSource(sourceKey);
    setArmedSourceNodeId(sourceKey);
    closeNodeMenu();
  }, [closeNodeMenu, contextNode, onSelectSource]);

  const handleAddRelation = useCallback(() => {
    if (!contextNode) {
      closeNodeMenu();
      return;
    }

    if (!graphModel.sourceNodeId) {
      handleUseAsSource();
      return;
    }

    if (!contextNodeIsSource && !contextNodeIsConnected) {
      addRelation(graphModel.sourceNodeId, {
        attributeLocalId: contextNode.localAttributeId,
        valueLocalId: contextNode.localValueId,
      });
    }

    closeNodeMenu();
  }, [
    addRelation,
    closeNodeMenu,
    contextNode,
    contextNodeIsConnected,
    contextNodeIsSource,
    graphModel.sourceNodeId,
    handleUseAsSource,
  ]);

  const handleDeleteRelation = useCallback(() => {
    if (
      !contextNode ||
      contextNodeIsSource ||
      !contextNodeIsConnected ||
      !graphModel.sourceNodeId
    ) {
      closeNodeMenu();
      return;
    }

    removeRelation(graphModel.sourceNodeId, {
      attributeLocalId: contextNode.localAttributeId,
      valueLocalId: contextNode.localValueId,
    });
    closeNodeMenu();
  }, [
    closeNodeMenu,
    contextNode,
    contextNodeIsConnected,
    contextNodeIsSource,
    graphModel.sourceNodeId,
    removeRelation,
  ]);

  const handleDeleteAll = useCallback(() => {
    if (!graphModel.sourceNodeId) {
      closeNodeMenu();
      return;
    }

    updateRelations(graphModel.sourceNodeId, []);
    setArmedSourceNodeId(null);
    setSelectedEdgeId(null);
    closeNodeMenu();
  }, [closeNodeMenu, graphModel.sourceNodeId, updateRelations]);

  const handleReconnect: OnReconnect<RelationFlowEdge> = useCallback(
    (oldEdge, connection) => {
      const relationSourceKey = oldEdge.source;
      if (!relationSourceKey) {
        return;
      }

      const targetNodeId = connection.target;
      if (!targetNodeId || targetNodeId === relationSourceKey) {
        return;
      }

      const targetNode = graphModel.nodeById.get(targetNodeId);
      if (!targetNode) {
        return;
      }

      const [, sourceValueIdForEdge] = relationSourceKey.split("::");
      if (!sourceValueIdForEdge) {
        return;
      }

      const currentRelations =
        graphDraftByValueLocalId.get(sourceValueIdForEdge)?.relations ?? [];

      const nextRelations = currentRelations
        .filter((item) => {
          const edgeId = `rel:${relationSourceKey}->${buildRelationSourceKey(
            item.attributeLocalId,
            item.valueLocalId,
          )}`;
          return edgeId !== oldEdge.id;
        })
        .concat({
          attributeLocalId: targetNode.localAttributeId,
          valueLocalId: targetNode.localValueId,
        });

      onSelectSource(relationSourceKey);
      setArmedSourceNodeId(relationSourceKey);
      updateRelations(relationSourceKey, nextRelations);
      setSelectedEdgeId(`rel:${relationSourceKey}->${targetNodeId}`);
    },
    [graphDraftByValueLocalId, graphModel.nodeById, onSelectSource, updateRelations],
  );

  useEffect(() => {
    if (!selectedEdgeId) {
      return;
    }

    const handleKeyDown = (event: KeyboardEvent) => {
      if (event.key !== "Delete" && event.key !== "Backspace") {
        return;
      }

      event.preventDefault();
      removeRelationByEdgeId(selectedEdgeId);
      setSelectedEdgeId(null);
    };

    window.addEventListener("keydown", handleKeyDown);
    return () => {
      window.removeEventListener("keydown", handleKeyDown);
    };
  }, [removeRelationByEdgeId, selectedEdgeId]);

  const ghostArrowPath = useMemo(() => {
    if (!armedSourceFlowNode || !cursorPosition) {
      return null;
    }

    const sourceSide = getSideTowardPoint(
      armedSourceFlowNode.position.x,
      armedSourceFlowNode.position.y,
      cursorPosition.x,
      cursorPosition.y,
    );
    const sourcePoint = getNodeAnchorPoint(
      armedSourceFlowNode.position.x,
      armedSourceFlowNode.position.y,
      sourceSide,
    );
    const targetPoint = { x: cursorPosition.x, y: cursorPosition.y };
    const targetSide: AnchorSide =
      Math.abs(targetPoint.x - sourcePoint.x) >= Math.abs(targetPoint.y - sourcePoint.y)
        ? targetPoint.x >= sourcePoint.x
          ? "left"
          : "right"
        : targetPoint.y >= sourcePoint.y
          ? "top"
          : "bottom";

    const obstacles = graphModel.flowNodes
      .filter((node) => node.id !== armedSourceFlowNode.id)
      .map((node) => getNodeBounds(node.position.x, node.position.y));

    return buildOrthogonalPath({
      sourcePoint,
      sourceSide,
      targetPoint,
      targetSide,
      obstacles,
    });
  }, [armedSourceFlowNode, cursorPosition, graphModel.flowNodes]);

  const handleNodeClick = useCallback(
    (event: React.MouseEvent, node: RelationFlowNode) => {
      event.preventDefault();
      event.stopPropagation();
      closeNodeMenu();
      setSelectedEdgeId(null);

      if (!effectiveArmedSourceNodeId) {
        setArmedSourceNodeId(node.id);
        updateCursorPosition(event.clientX, event.clientY);
        onSelectSource(node.id);
        return;
      }

      if (effectiveArmedSourceNodeId === node.id) {
        setArmedSourceNodeId(null);
        return;
      }

      const item = graphModel.nodeById.get(node.id);
      const nodeIsConnected = graphModel.connectedNodeIds.has(node.id);
      if (!item) {
        return;
      }

      if (graphModel.sourceNodeId && node.id !== graphModel.sourceNodeId) {
        if (nodeIsConnected) {
          removeRelation(graphModel.sourceNodeId, {
            attributeLocalId: item.localAttributeId,
            valueLocalId: item.localValueId,
          });
        } else {
          addRelation(graphModel.sourceNodeId, {
            attributeLocalId: item.localAttributeId,
            valueLocalId: item.localValueId,
          });
        }
      }

      setArmedSourceNodeId(null);
      setCursorPosition(null);
    },
    [
      addRelation,
      closeNodeMenu,
      effectiveArmedSourceNodeId,
      graphModel.connectedNodeIds,
      graphModel.nodeById,
      graphModel.sourceNodeId,
      onSelectSource,
      removeRelation,
      updateCursorPosition,
    ],
  );

  const handleNodeContextMenu = useCallback(
    (event: React.MouseEvent, node: RelationFlowNode) => {
      event.preventDefault();
      event.stopPropagation();
      setSelectedEdgeId(null);
      setCursorPosition(null);
      setContextNodeId(node.id);
      setMenuPosition({ mouseX: event.clientX + 2, mouseY: event.clientY - 6 });
    },
    [],
  );

  const handleEdgeClick = useCallback(
    (event: React.MouseEvent, edge: RelationFlowEdge) => {
      event.preventDefault();
      event.stopPropagation();
      closeNodeMenu();
      setContextNodeId(null);
      setArmedSourceNodeId(null);
      onSelectSource(edge.source);
      setSelectedEdgeId(edge.id);
    },
    [closeNodeMenu, onSelectSource],
  );

  const handlePaneClick = useCallback(() => {
    setArmedSourceNodeId(null);
    setContextNodeId(null);
    setSelectedEdgeId(null);
    setCursorPosition(null);
    setMenuPosition(null);
  }, []);

  return {
    containerRef,
    graphModel,
    selectedSourceRelations,
    ghostArrowPath,
    menuPosition,
    contextNode,
    contextNodeIsSource,
    contextNodeIsConnected,
    effectiveArmedSourceNodeId,
    closeNodeMenu,
    handleUseAsSource,
    handleAddRelation,
    handleDeleteRelation,
    handleDeleteAll,
    handleReconnect,
    handleNodeClick,
    handleNodeContextMenu,
    handleEdgeClick,
    handlePaneClick,
    updateCursorPosition,
    setArmedSourceNodeId,
    setSelectedEdgeId,
    setCursorPosition,
    removeRelationByEdgeId,
    updateRelations,
    addRelation,
    removeRelation,
  };
};
