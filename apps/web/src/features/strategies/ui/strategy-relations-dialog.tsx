import {
  Alert,
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  Stack,
  Typography,
} from "@mui/material";

import type {
  GraphDraftAttribute,
  GraphDraftRelationRef,
  StrategyRelationNode,
} from "../model/relations";
import { RelationGraphEditor } from "./relation-graph-editor";

type StrategyRelationsDialogProps = {
  open: boolean;
  strategyName: string;
  nodes: StrategyRelationNode[];
  graphDraft: GraphDraftAttribute[];
  selectedSourceKey: string | null;
  dirtyCount: number;
  saveState: {
    status: "idle" | "saving" | "error";
    error: string | null;
  };
  onClose: () => void;
  onSelectSource: (sourceKey: string) => void;
  onChangeRelations: (
    sourceKey: string,
    relations: GraphDraftRelationRef[],
  ) => void;
  onAddAttribute: () => void;
  onRenameAttribute: (attributeLocalId: string) => void;
  onDeleteAttribute: (attributeLocalId: string) => void;
  onAddValue: (attributeLocalId: string) => void;
  onRenameValue: (valueLocalId: string) => void;
  onDeleteValue: (valueLocalId: string) => void;
  onSave: () => void | Promise<void>;
};

export const StrategyRelationsDialog = ({
  open,
  strategyName,
  nodes,
  graphDraft,
  selectedSourceKey,
  dirtyCount,
  saveState,
  onClose,
  onSelectSource,
  onChangeRelations,
  onAddAttribute,
  onRenameAttribute,
  onDeleteAttribute,
  onAddValue,
  onRenameValue,
  onDeleteValue,
  onSave,
}: StrategyRelationsDialogProps) => {
  const selectedSource =
    nodes.find(
      (node) => `${node.localAttributeId}::${node.localValueId}` === selectedSourceKey,
    ) ?? null;

  return (
    <Dialog open={open} onClose={onClose} maxWidth="xl" fullWidth>
      <DialogTitle sx={{ pb: 1 }}>
        <Stack spacing={0.5}>
          <Typography variant="h6" sx={{ fontWeight: 800 }}>
            Strategy relations workspace
          </Typography>
          <Typography color="text.secondary" variant="body2">
            {strategyName}
          </Typography>
        </Stack>
      </DialogTitle>

      <DialogContent dividers>
        <Stack spacing={2}>
          {saveState.status === "error" && saveState.error ? (
            <Alert severity="error">{saveState.error}</Alert>
          ) : null}

          {!selectedSource ? (
            <Alert severity="info">
              Click any node in the graph to choose the source value for outgoing
              relations.
            </Alert>
          ) : null}

          <RelationGraphEditor
            sourceAttributeLocalId={selectedSource?.localAttributeId ?? null}
            sourceValueLocalId={selectedSource?.localValueId ?? null}
            nodes={nodes}
            graphDraft={graphDraft}
            onSelectSource={onSelectSource}
            onChangeSourceRelations={onChangeRelations}
            onAddAttribute={onAddAttribute}
            onRenameAttribute={onRenameAttribute}
            onDeleteAttribute={onDeleteAttribute}
            onAddValue={onAddValue}
            onRenameValue={onRenameValue}
            onDeleteValue={onDeleteValue}
          />
        </Stack>
      </DialogContent>

      <DialogActions sx={{ px: 3, py: 2 }}>
        <Typography
          color="text.secondary"
          variant="body2"
          sx={{ mr: "auto" }}
        >
          {dirtyCount > 0
            ? `${dirtyCount} pending graph action${dirtyCount > 1 ? "s" : ""}`
            : "No unsaved graph changes"}
        </Typography>
        <Button onClick={onClose}>Close</Button>
        <Button
          variant="contained"
          onClick={() => void onSave()}
          disabled={dirtyCount === 0 || saveState.status === "saving"}
        >
          {saveState.status === "saving" ? "Saving..." : "Save graph"}
        </Button>
      </DialogActions>
    </Dialog>
  );
};
