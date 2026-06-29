"use client";

import {
  PersonAddAltRounded,
  PeopleRounded,
  RefreshRounded,
  SearchRounded,
} from "@mui/icons-material";
import {
  Alert,
  Button,
  InputAdornment,
  Stack,
  TablePagination,
  TextField,
} from "@mui/material";
import { useState } from "react";

import {
  CreateUserDialog,
  UsersTable,
  useCreateUser,
  useUsers,
} from "@/features/users";
import {
  PageHeader,
  PageLayout,
  StatBadge,
  SurfaceSection,
} from "@/shared/ui/page-layout";

const formatTotal = (total: number) => new Intl.NumberFormat().format(total);
const PAGE_SIZE_OPTIONS = [10, 25, 50, 100];

export const UsersOverview = () => {
  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);
  const {
    state,
    page,
    query,
    rowsPerPage,
    setPage,
    setQuery,
    setRowsPerPage,
    refresh,
  } = useUsers({ pageSize: 25 });
  const createUserMutation = useCreateUser();

  const isRefreshing = state.status === "loading" && state.items.length > 0;
  const hasRenderableTable = state.status !== "error" || state.items.length > 0;

  const closeCreateDialog = () => {
    if (createUserMutation.isPending) {
      return;
    }

    createUserMutation.reset();
    setIsCreateDialogOpen(false);
  };

  const handleCreateUser = async (
    input: Parameters<typeof createUserMutation.submit>[0],
  ) => {
    try {
      await createUserMutation.submit(input);
      createUserMutation.reset();
      setIsCreateDialogOpen(false);
      await refresh();
    } catch {
      // Error state is captured in the dialog hook.
    }
  };

  return (
    <>
      <PageLayout>
        <PageHeader
          title="Users"
          description="Administrative overview of all users created in the identity service."
          actions={
            <Stack direction={{ xs: "column", sm: "row" }} spacing={1.5}>
              <TextField
                size="small"
                value={query}
                onChange={(event) => setQuery(event.target.value)}
                placeholder="Search by login, name or email"
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
                variant="outlined"
                startIcon={<RefreshRounded />}
                onClick={() => void refresh()}
                disabled={state.status === "loading"}
              >
                {isRefreshing ? "Refreshing..." : "Refresh"}
              </Button>
            </Stack>
          }
        />

        <SurfaceSection
          title="Users directory"
          description="Browse accounts with server-side search and pagination."
          badge={
            <Stack direction={{ xs: "column", sm: "row" }} spacing={1}>
              <StatBadge
                icon={<PeopleRounded />}
                label={`Total: ${formatTotal(state.total)}`}
              />
              <Button
                variant="contained"
                startIcon={<PersonAddAltRounded />}
                onClick={() => setIsCreateDialogOpen(true)}
              >
                Create user
              </Button>
            </Stack>
          }
        >
          {state.status === "error" ? (
            <Alert
              severity="error"
              action={
                <Button
                  color="inherit"
                  size="small"
                  onClick={() => void refresh()}
                >
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
              <UsersTable
                items={state.items}
                isLoading={state.status === "loading"}
                query={query}
              />

              <TablePagination
                component="div"
                count={state.total}
                page={page}
                onPageChange={(_, nextPage) => setPage(nextPage)}
                rowsPerPage={rowsPerPage}
                onRowsPerPageChange={(event) => {
                  setRowsPerPage(Number(event.target.value));
                }}
                rowsPerPageOptions={PAGE_SIZE_OPTIONS}
              />
            </>
          ) : null}
        </SurfaceSection>
      </PageLayout>

      <CreateUserDialog
        key={isCreateDialogOpen ? "open" : "closed"}
        open={isCreateDialogOpen}
        isSubmitting={createUserMutation.isPending}
        error={createUserMutation.error}
        onClose={closeCreateDialog}
        onSubmit={handleCreateUser}
      />
    </>
  );
};
