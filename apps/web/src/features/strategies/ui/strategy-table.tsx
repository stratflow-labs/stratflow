import { timestampDate } from "@bufbuild/protobuf/wkt";
import {
  Alert,
  CircularProgress,
  Stack,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Typography,
} from "@mui/material";

import type { Strategy } from "@/shared/api/gen/strategy_registry/proto/v1/strategy_types_pb";

const formatDate = (value?: Parameters<typeof timestampDate>[0] | null) => {
  if (!value) {
    return "—";
  }

  const date = timestampDate(value);

  return Number.isNaN(date.getTime()) ? "—" : date.toLocaleString();
};

const getRowKey = (
  strategy: { id?: string; slug?: string; name?: string },
  index: number,
) => strategy.id || strategy.slug || strategy.name || `strategy-${index}`;

type StrategyTableProps = {
  items: Strategy[];
  isLoading: boolean;
  query: string;
  onOpenStrategy: (strategy: Strategy) => void | Promise<void>;
};

export const StrategyTable = ({
  items,
  isLoading,
  query,
  onOpenStrategy,
}: StrategyTableProps) => {
  if (isLoading && items.length === 0) {
    return (
      <Stack
        spacing={2}
        role="status"
        sx={{
          alignItems: "center",
          justifyContent: "center",
          py: 8,
        }}
      >
        <CircularProgress size={28} />
        <Typography color="text.secondary">Loading strategies...</Typography>
      </Stack>
    );
  }

  return (
    <TableContainer sx={{ width: "100%", overflowX: "auto" }}>
      <Table
        size="small"
        aria-label="Strategies"
        aria-busy={isLoading}
        sx={{ minWidth: "64rem", tableLayout: "fixed" }}
      >
        <TableHead>
          <TableRow>
            <TableCell sx={{ width: "16ch" }}>Slug</TableCell>
            <TableCell sx={{ width: "24ch" }}>Name</TableCell>
            <TableCell>Description</TableCell>
            <TableCell sx={{ width: "18ch" }}>Created</TableCell>
            <TableCell sx={{ width: "18ch" }}>Updated</TableCell>
          </TableRow>
        </TableHead>

        <TableBody>
          {items.length === 0 ? (
            <TableRow>
              <TableCell colSpan={5}>
                <Alert severity="info">
                  {query.trim()
                    ? "No strategies match the current search query."
                    : "No strategies returned by the registry yet."}
                </Alert>
              </TableCell>
            </TableRow>
          ) : (
            items.map((strategy, index) => (
              <TableRow
                key={getRowKey(strategy, index)}
                hover
                onClick={() => void onOpenStrategy(strategy)}
                sx={{
                  cursor: "pointer",
                  "&:last-child td": { borderBottom: 0 },
                }}
              >
                <TableCell
                  sx={{
                    fontFamily: "monospace",
                    overflow: "hidden",
                    textOverflow: "ellipsis",
                    whiteSpace: "nowrap",
                  }}
                  title={strategy.slug || undefined}
                >
                  {strategy.slug || "—"}
                </TableCell>

                <TableCell
                  sx={{
                    fontWeight: 700,
                    overflow: "hidden",
                    textOverflow: "ellipsis",
                    whiteSpace: "nowrap",
                  }}
                  title={strategy.name || undefined}
                >
                  {strategy.name || "Untitled strategy"}
                </TableCell>

                <TableCell>
                  <Typography
                    variant="body2"
                    color="text.secondary"
                    noWrap
                    title={strategy.description || undefined}
                  >
                    {strategy.description || "—"}
                  </Typography>
                </TableCell>

                <TableCell>{formatDate(strategy.createdAt)}</TableCell>
                <TableCell>{formatDate(strategy.updatedAt)}</TableCell>
              </TableRow>
            ))
          )}
        </TableBody>
      </Table>
    </TableContainer>
  );
};
