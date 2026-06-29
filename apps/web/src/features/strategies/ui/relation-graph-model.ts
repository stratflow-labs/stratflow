import type {
  GraphDraftAttribute,
  GraphDraftRelationRef,
  StrategyRelationNode,
} from "../model/relations";
import { buildRelationSourceKey } from "../model/relations";
import {
  applyColumnLayout,
  getNodeBounds,
} from "./relation-graph-routing";
import {
  defaultRelationMarker,
  defaultRelationStyle,
  type RelationFlowEdge,
  type RelationFlowNode,
} from "./relation-graph-elements";

type BuildRelationGraphModelInput = {
  nodes: StrategyRelationNode[];
  graphDraft: GraphDraftAttribute[];
  selectedSourceRelations: GraphDraftRelationRef[];
  sourceAttributeLocalId: string | null;
  sourceValueLocalId: string | null;
  armedSourceNodeId: string | null;
  contextNodeId: string | null;
  selectedEdgeId: string | null;
};

export const buildRelationGraphModel = ({
  nodes,
  graphDraft,
  selectedSourceRelations,
  sourceAttributeLocalId,
  sourceValueLocalId,
  armedSourceNodeId,
  contextNodeId,
  selectedEdgeId,
}: BuildRelationGraphModelInput) => {
  const sourceNodeId =
    sourceAttributeLocalId && sourceValueLocalId
      ? buildRelationSourceKey(sourceAttributeLocalId, sourceValueLocalId)
      : null;
  const nodeById = new Map<string, StrategyRelationNode>();

  nodes.forEach((item) => {
    nodeById.set(
      buildRelationSourceKey(item.localAttributeId, item.localValueId),
      item,
    );
  });

  const hasSelectedSource =
    sourceNodeId !== null && nodeById.has(sourceNodeId);

  const connectedNodeIds = new Set(
    hasSelectedSource
      ? selectedSourceRelations.map((relation) =>
          buildRelationSourceKey(
            relation.attributeLocalId,
            relation.valueLocalId,
          ),
        )
      : [],
  );

  const flowNodes = nodes.map((item) => {
    const nodeId = buildRelationSourceKey(
      item.localAttributeId,
      item.localValueId,
    );

    return {
      id: nodeId,
      type: "relationNode",
      position: { x: 0, y: 0 },
      data: {
        attributeSlug: item.attributeSlug,
        valueLabel: item.valueLabel,
        isSource: hasSelectedSource && nodeId === sourceNodeId,
        isConnected: connectedNodeIds.has(nodeId),
        isSelected: armedSourceNodeId === nodeId || contextNodeId === nodeId,
      },
    } satisfies RelationFlowNode;
  });

  const positionedFlowNodes = applyColumnLayout(flowNodes, nodeById);
  const positionedNodeById = new Map(
    positionedFlowNodes.map((node) => [node.id, node] as const),
  );

  const flowEdges = graphDraft.flatMap((attribute) =>
    attribute.values.flatMap((draftValue) => {
      const relationSourceKey = buildRelationSourceKey(
        attribute.localId,
        draftValue.localId,
      );
      const sourceFlowNode = positionedNodeById.get(relationSourceKey);
      if (!sourceFlowNode) {
        return [];
      }

      return draftValue.relations.flatMap((relation) => {
        const targetNodeId = buildRelationSourceKey(
          relation.attributeLocalId,
          relation.valueLocalId,
        );
        const targetFlowNode = positionedNodeById.get(targetNodeId);

        if (!targetFlowNode || targetNodeId === relationSourceKey) {
          return [];
        }

        const edgeId = `rel:${relationSourceKey}->${targetNodeId}`;

        return [
          {
            id: edgeId,
            type: "relationEdge",
            source: relationSourceKey,
            target: targetNodeId,
            data: {
              sourceSide: "right",
              targetSide: "left",
              obstacles: positionedFlowNodes
                .filter(
                  (candidate) =>
                    candidate.id !== relationSourceKey &&
                    candidate.id !== targetNodeId,
                )
                .map((candidate) =>
                  getNodeBounds(candidate.position.x, candidate.position.y),
                ),
            },
            reconnectable: true,
            selectable: true,
            selected: selectedEdgeId === edgeId,
            markerEnd: defaultRelationMarker(selectedEdgeId === edgeId),
            style: defaultRelationStyle(selectedEdgeId === edgeId),
          } satisfies RelationFlowEdge,
        ];
      });
    }),
  );

  return {
    sourceNodeId,
    nodeById,
    connectedNodeIds,
    flowNodes: positionedFlowNodes,
    flowEdges,
  };
};
