import { AddRounded, RefreshRounded, SearchRounded } from "@mui/icons-material";
import {
  Button,
  CircularProgress,
  InputAdornment,
  Stack,
  TextField,
} from "@mui/material";

import { PageHeader } from "@/shared/ui/page-layout";

type StrategyScreenHeaderProps = {
  query: string;
  isRefreshing: boolean;
  isLoading: boolean;
  onQueryChange: (value: string) => void;
  onOpenCreate: () => void;
  onRefresh: () => void;
};

export const StrategyScreenHeader = ({
  query,
  isRefreshing,
  isLoading,
  onQueryChange,
  onOpenCreate,
  onRefresh,
}: StrategyScreenHeaderProps) => (
  <PageHeader
    title="Strategies"
    description="Live data from StrategyRegistryService.ListStrategies through Connect-Web."
    actions={
      <Stack
        direction={{ xs: "column", sm: "row" }}
        spacing={1.5}
        sx={{ width: "100%", justifyContent: { xl: "flex-end" } }}
      >
        <TextField
          value={query}
          onChange={(event) => onQueryChange(event.target.value)}
          placeholder="Search strategies"
          aria-label="Search strategies"
          size="small"
          fullWidth
          slotProps={{
            input: {
              startAdornment: (
                <InputAdornment position="start">
                  <SearchRounded fontSize="small" />
                </InputAdornment>
              ),
            },
          }}
          sx={{ flex: 1, minWidth: 0 }}
        />

        <Button
          variant="contained"
          startIcon={<AddRounded />}
          onClick={onOpenCreate}
          sx={{ whiteSpace: "nowrap" }}
        >
          Create strategy
        </Button>

        <Button
          variant="outlined"
          startIcon={
            isRefreshing ? (
              <CircularProgress size={16} color="inherit" />
            ) : (
              <RefreshRounded />
            )
          }
          onClick={onRefresh}
          disabled={isLoading}
          sx={{ whiteSpace: "nowrap" }}
        >
          Refresh
        </Button>
      </Stack>
    }
  />
);
