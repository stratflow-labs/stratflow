"use client";

import { useState } from "react";

import {
  createStrategy,
  StrategyRegistryCard,
  StrategyScreenDialogs,
  StrategyScreenHeader,
  useStrategyDetailsDialog,
  useStrategies,
} from "@/features/strategies";
import { PageLayout } from "@/shared/ui/page-layout";

export const StrategiesOverview = () => {
  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);
  const [createStrategyErrorMessage, setCreateStrategyErrorMessage] = useState<
    string | null
  >(null);
  const [isCreatingStrategy, setIsCreatingStrategy] = useState(false);
  const {
    state,
    page,
    query,
    rowsPerPage,
    setPage,
    setQuery,
    setRowsPerPage,
    refresh,
  } = useStrategies({ pageSize: 25 });
  const strategyDetailsDialog = useStrategyDetailsDialog();

  const isRefreshing = state.status === "loading" && state.items.length > 0;
  const hasRenderableTable = state.status !== "error" || state.items.length > 0;

  const handleCreateStrategy = async ({
    slug,
    title,
    description,
  }: {
    slug: string;
    title: string;
    description: string;
  }) => {
    setIsCreatingStrategy(true);
    setCreateStrategyErrorMessage(null);

    try {
      await createStrategy({
        slug,
        name: title,
        description,
      });
      setIsCreateDialogOpen(false);
      await refresh();
    } catch (error) {
      setCreateStrategyErrorMessage(
        error instanceof Error ? error.message : "Failed to create strategy.",
      );
    } finally {
      setIsCreatingStrategy(false);
    }
  };

  return (
    <PageLayout>
      <StrategyScreenHeader
        query={query}
        isRefreshing={isRefreshing}
        isLoading={state.status === "loading"}
        onQueryChange={setQuery}
        onOpenCreate={() => setIsCreateDialogOpen(true)}
        onRefresh={() => void refresh()}
      />

      <StrategyRegistryCard
        state={state}
        page={page}
        rowsPerPage={rowsPerPage}
        query={query}
        hasRenderableTable={hasRenderableTable}
        onRetry={() => void refresh()}
        onOpenStrategy={strategyDetailsDialog.openStrategy}
        onPageChange={setPage}
        onRowsPerPageChange={setRowsPerPage}
      />

      <StrategyScreenDialogs
        strategyDetailsDialog={strategyDetailsDialog}
        isCreateDialogOpen={isCreateDialogOpen}
        isCreatingStrategy={isCreatingStrategy}
        createStrategyError={createStrategyErrorMessage}
        onCloseCreateStrategy={() => {
          if (isCreatingStrategy) {
            return;
          }
          setCreateStrategyErrorMessage(null);
          setIsCreateDialogOpen(false);
        }}
        onSubmitCreateStrategy={handleCreateStrategy}
      />
    </PageLayout>
  );
};
