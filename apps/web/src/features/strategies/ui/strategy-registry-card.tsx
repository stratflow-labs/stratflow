import { RuleFolderRounded } from "@mui/icons-material";
import { Alert, Button, TablePagination } from "@mui/material";

import { StatBadge, SurfaceSection } from "@/shared/ui/page-layout";
import type { Strategy } from "@/shared/api/gen/strategy_registry/proto/v1/strategy_types_pb";
import type { StrategiesState } from "../model/use-strategies";

import { StrategyTable } from "./strategy-table";

const formatTotal = (total: number) => new Intl.NumberFormat().format(total);

type StrategyRegistryCardProps = {
  state: StrategiesState;
  page: number;
  rowsPerPage: number;
  query: string;
  hasRenderableTable: boolean;
  onRetry: () => void;
  onOpenStrategy: (strategy: Strategy) => void | Promise<void>;
  onPageChange: (page: number) => void;
  onRowsPerPageChange: (rowsPerPage: number) => void;
};

export const StrategyRegistryCard = ({
  state,
  page,
  rowsPerPage,
  query,
  hasRenderableTable,
  onRetry,
  onOpenStrategy,
  onPageChange,
  onRowsPerPageChange,
}: StrategyRegistryCardProps) => (
  <SurfaceSection
    title="Strategy registry"
    description="Browse backend strategies with server-side search and pagination."
    badge={
      <StatBadge
        icon={<RuleFolderRounded />}
        label={`Total: ${formatTotal(state.total)}`}
      />
    }
  >
    {state.status === "error" ? (
      <Alert
        severity="error"
        action={
          <Button color="inherit" size="small" onClick={onRetry}>
            Retry
          </Button>
        }
        sx={{ mb: hasRenderableTable ? 2 : 0 }}
      >
        {state.error}
      </Alert>
    ) : null}

    {hasRenderableTable ? (
      <>
        <StrategyTable
          items={state.items}
          isLoading={state.status === "loading"}
          query={query}
          onOpenStrategy={onOpenStrategy}
        />

        <TablePagination
          component="div"
          count={state.total}
          page={page}
          onPageChange={(_, nextPage) => onPageChange(nextPage)}
          rowsPerPage={rowsPerPage}
          onRowsPerPageChange={(event) => {
            onRowsPerPageChange(Number(event.target.value));
          }}
          rowsPerPageOptions={[10, 25, 50, 100]}
        />
      </>
    ) : null}
  </SurfaceSection>
);
