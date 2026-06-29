import type { useStrategyDetailsDialog } from "../model/use-strategy-details-dialog";

import { SimpleEntityDialog } from "./simple-entity-dialog";
import { StrategyDetailsDialog } from "./strategy-details-dialog";
import { StrategyRelationsDialog } from "./strategy-relations-dialog";

type StrategyDetailsDialogController = ReturnType<typeof useStrategyDetailsDialog>;

type StrategyScreenDialogsProps = {
  strategyDetailsDialog: StrategyDetailsDialogController;
  isCreateDialogOpen: boolean;
  isCreatingStrategy: boolean;
  createStrategyError: string | null;
  onCloseCreateStrategy: () => void;
  onSubmitCreateStrategy: (input: {
    slug: string;
    title: string;
    description: string;
  }) => void | Promise<void>;
};

export const StrategyScreenDialogs = ({
  strategyDetailsDialog,
  isCreateDialogOpen,
  isCreatingStrategy,
  createStrategyError,
  onCloseCreateStrategy,
  onSubmitCreateStrategy,
}: StrategyScreenDialogsProps) => (
  <>
    <StrategyDetailsDialog
      open={strategyDetailsDialog.isOpen}
      strategy={strategyDetailsDialog.selectedStrategy}
      entry={strategyDetailsDialog.entry}
      expandedAttributeId={strategyDetailsDialog.expandedAttributeId}
      onClose={strategyDetailsDialog.closeStrategy}
      onRefresh={strategyDetailsDialog.refreshSelectedStrategy}
      onToggleAttribute={strategyDetailsDialog.toggleAttribute}
      onOpenRelationsGraph={strategyDetailsDialog.openRelationsGraph}
      onCreateAttribute={strategyDetailsDialog.openCreateAttributeDialog}
      onCreateValue={strategyDetailsDialog.openCreateValueDialog}
    />

    <StrategyRelationsDialog
      open={strategyDetailsDialog.isGraphOpen}
      strategyName={
        strategyDetailsDialog.selectedStrategy?.name ||
        strategyDetailsDialog.selectedStrategy?.slug ||
        "Strategy"
      }
      nodes={strategyDetailsDialog.relationNodes}
      graphDraft={strategyDetailsDialog.graphDraft}
      selectedSourceKey={strategyDetailsDialog.selectedGraphSourceKey}
      dirtyCount={strategyDetailsDialog.dirtyGraphSourceCount}
      saveState={strategyDetailsDialog.graphSaveState}
      onClose={strategyDetailsDialog.closeRelationsGraph}
      onSelectSource={strategyDetailsDialog.selectGraphSource}
      onChangeRelations={strategyDetailsDialog.updateGraphRelations}
      onAddAttribute={strategyDetailsDialog.addGraphAttribute}
      onRenameAttribute={strategyDetailsDialog.renameGraphAttribute}
      onDeleteAttribute={strategyDetailsDialog.deleteGraphAttribute}
      onAddValue={strategyDetailsDialog.addGraphValue}
      onRenameValue={strategyDetailsDialog.renameGraphValue}
      onDeleteValue={strategyDetailsDialog.deleteGraphValue}
      onSave={strategyDetailsDialog.saveGraphRelations}
    />

    <SimpleEntityDialog
      key={`strategy-create-${String(isCreateDialogOpen)}`}
      open={isCreateDialogOpen}
      mode="strategy"
      title="Create strategy"
      subtitle="Enter slug, title and description."
      submitLabel="Create strategy"
      isSubmitting={isCreatingStrategy}
      error={createStrategyError}
      onClose={onCloseCreateStrategy}
      onSubmit={onSubmitCreateStrategy}
    />

    <SimpleEntityDialog
      key={`details-create-${strategyDetailsDialog.createDialogMode ?? "closed"}-${strategyDetailsDialog.createDialogAttribute?.localId ?? "none"}`}
      open={strategyDetailsDialog.createDialogMode !== null}
      mode={strategyDetailsDialog.createDialogMode ?? "attribute"}
      title={
        strategyDetailsDialog.createDialogMode === "value"
          ? "Create attribute value"
          : "Create attribute"
      }
      subtitle={
        strategyDetailsDialog.createDialogMode === "value"
          ? `Create a value for ${strategyDetailsDialog.createDialogAttribute?.name || strategyDetailsDialog.createDialogAttribute?.slug || "the selected attribute"}.`
          : "Enter slug, title and description."
      }
      submitLabel={
        strategyDetailsDialog.createDialogMode === "value"
          ? "Create value"
          : "Create attribute"
      }
      isSubmitting={strategyDetailsDialog.isCreating}
      error={strategyDetailsDialog.createError}
      onClose={strategyDetailsDialog.closeCreateDialog}
      onSubmit={strategyDetailsDialog.submitCreateDialog}
    />
  </>
);
