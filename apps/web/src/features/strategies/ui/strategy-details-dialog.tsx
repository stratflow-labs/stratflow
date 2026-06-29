import { Dialog, DialogContent, DialogTitle, Stack } from "@mui/material";

import {
  StrategyAttributesTable,
  StrategyDetailsHeader,
  StrategyDetailsStatus,
  StrategyDetailsToolbar,
} from "./strategy-details-parts";
import type { StrategyDetailsDialogProps } from "./strategy-details-types";

export const StrategyDetailsDialog = ({
  open,
  strategy,
  entry,
  expandedAttributeId,
  onClose,
  onRefresh,
  onToggleAttribute,
  onOpenRelationsGraph,
  onCreateAttribute,
  onCreateValue,
}: StrategyDetailsDialogProps) => (
  <Dialog open={open} onClose={onClose} maxWidth="lg" fullWidth>
    <DialogTitle sx={{ pb: 1.5 }}>
      <StrategyDetailsHeader strategy={strategy} />
    </DialogTitle>

    <DialogContent dividers>
      <Stack spacing={2}>
        <StrategyDetailsToolbar
          strategy={strategy}
          entry={entry}
          onCreateAttribute={onCreateAttribute}
          onOpenRelationsGraph={onOpenRelationsGraph}
          onRefresh={onRefresh}
        />

        <StrategyDetailsStatus entry={entry} onRefresh={onRefresh} />

        <StrategyAttributesTable
          entry={entry}
          expandedAttributeId={expandedAttributeId}
          onToggleAttribute={onToggleAttribute}
          onCreateValue={onCreateValue}
        />
      </Stack>
    </DialogContent>
  </Dialog>
);
