import Button from "@mui/material/Button";
import Stack from "@mui/material/Stack";
import Typography from "@mui/material/Typography";
import {
  Background,
  Controls,
  ReactFlow,
} from "@xyflow/react";

import type {
  GraphDraftAttribute,
  GraphDraftRelationRef,
  StrategyRelationNode,
} from "../model/relations";

import styles from "./relation-graph-editor.module.css";
import { RelationGraphContextMenu } from "./relation-graph-context-menu";
import {
  relationEdgeTypes,
  relationNodeTypes,
  type RelationFlowEdge,
  type RelationFlowNode,
} from "./relation-graph-elements";
import { useRelationGraphEditor } from "./use-relation-graph-editor";

import "@xyflow/react/dist/style.css";

type RelationGraphEditorProps = {
  sourceAttributeLocalId: string | null;
  sourceValueLocalId: string | null;
  nodes: StrategyRelationNode[];
  graphDraft: GraphDraftAttribute[];
  onSelectSource: (sourceKey: string) => void;
  onChangeSourceRelations: (
    sourceKey: string,
    relations: GraphDraftRelationRef[],
  ) => void;
  onAddAttribute: () => void;
  onRenameAttribute: (attributeLocalId: string) => void;
  onDeleteAttribute: (attributeLocalId: string) => void;
  onAddValue: (attributeLocalId: string) => void;
  onRenameValue: (valueLocalId: string) => void;
  onDeleteValue: (valueLocalId: string) => void;
};

export const RelationGraphEditor = ({
  sourceAttributeLocalId,
  sourceValueLocalId,
  nodes,
  graphDraft,
  onSelectSource,
  onChangeSourceRelations,
  onAddAttribute,
  onRenameAttribute,
  onDeleteAttribute,
  onAddValue,
  onRenameValue,
  onDeleteValue,
}: RelationGraphEditorProps) => {
  const {
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
  } = useRelationGraphEditor({
    sourceAttributeLocalId,
    sourceValueLocalId,
    nodes,
    graphDraft,
    onSelectSource,
    onChangeSourceRelations,
  });

  if (graphModel.flowNodes.length === 0) {
    return (
      <Stack spacing={1}>
        <Typography variant="subtitle2">Relations graph</Typography>
        <Typography color="text.secondary" variant="body2">
          Source value is missing from the loaded strategy attributes.
        </Typography>
      </Stack>
    );
  }

  return (
    <Stack spacing={1.5}>
      <Typography variant="subtitle2" sx={{ fontWeight: 700 }}>
        Relations graph
      </Typography>
      <Button
        onClick={onAddAttribute}
        variant="outlined"
        size="small"
        sx={{ alignSelf: "flex-start" }}
      >
        Add attribute
      </Button>
      <Typography color="text.secondary" variant="body2">
        Left click any node to choose source A. Then click another node to create
        a relation from A to B. Drag an arrow to another node to reconnect it.
        Right click opens the node menu. `Delete` removes the selected arrow.
      </Typography>
      {effectiveArmedSourceNodeId ? (
        <Typography color="primary" variant="body2">
          Source node selected. Pick a target node to create a relation.
        </Typography>
      ) : null}

      <div
        ref={containerRef}
        className={styles.container}
        style={{ height: 440, position: "relative" }}
        onMouseMove={(event) => {
          if (!effectiveArmedSourceNodeId) {
            return;
          }

          updateCursorPosition(event.clientX, event.clientY);
        }}
      >
        {ghostArrowPath ? (
          <svg
            width="100%"
            height="100%"
            style={{
              position: "absolute",
              inset: 0,
              pointerEvents: "none",
              zIndex: 2,
              overflow: "visible",
            }}
          >
            <defs>
              <marker
                id="relation-ghost-arrow"
                markerWidth="10"
                markerHeight="10"
                refX="8"
                refY="5"
                orient="auto"
              >
                <path d="M0,0 L10,5 L0,10 z" fill="#2e7d32" opacity="1" />
              </marker>
            </defs>
            <path
              d={ghostArrowPath}
              fill="none"
              stroke="#2e7d32"
              strokeWidth="1.8"
              opacity="1"
              markerEnd="url(#relation-ghost-arrow)"
            />
          </svg>
        ) : null}
        <ReactFlow<RelationFlowNode, RelationFlowEdge>
          className="relationGraphCanvas"
          style={{ cursor: "default" }}
          nodes={graphModel.flowNodes}
          edges={graphModel.flowEdges}
          nodeTypes={relationNodeTypes}
          edgeTypes={relationEdgeTypes}
          fitView
          fitViewOptions={{ padding: 0.24, maxZoom: 1.2 }}
          nodesDraggable={false}
          nodesConnectable
          elementsSelectable
          edgesReconnectable
          deleteKeyCode={null}
          proOptions={{ hideAttribution: true }}
          onNodeClick={handleNodeClick}
          onNodeContextMenu={handleNodeContextMenu}
          onEdgeClick={handleEdgeClick}
          onReconnect={handleReconnect}
          onPaneClick={handlePaneClick}
        >
          <Background />
          <Controls showInteractive={false} />
        </ReactFlow>
      </div>

      <RelationGraphContextMenu
        contextNode={contextNode}
        contextNodeIsSource={contextNodeIsSource}
        contextNodeIsConnected={contextNodeIsConnected}
        effectiveArmedSourceNodeId={effectiveArmedSourceNodeId}
        menuPosition={menuPosition}
        selectedSourceRelationsCount={selectedSourceRelations.length}
        onClose={closeNodeMenu}
        onUseAsSource={handleUseAsSource}
        onAddRelation={handleAddRelation}
        onDeleteRelation={handleDeleteRelation}
        onDeleteAll={handleDeleteAll}
        onAddValue={(attributeLocalId) => {
          onAddValue(attributeLocalId);
          closeNodeMenu();
        }}
        onRenameAttribute={(attributeLocalId) => {
          onRenameAttribute(attributeLocalId);
          closeNodeMenu();
        }}
        onRenameValue={(valueLocalId) => {
          onRenameValue(valueLocalId);
          closeNodeMenu();
        }}
        onDeleteValue={(valueLocalId) => {
          onDeleteValue(valueLocalId);
          closeNodeMenu();
        }}
        onDeleteAttribute={(attributeLocalId) => {
          onDeleteAttribute(attributeLocalId);
          closeNodeMenu();
        }}
        onClearSourceSelection={() => {
          setArmedSourceNodeId(null);
          closeNodeMenu();
        }}
      />
    </Stack>
  );
};
