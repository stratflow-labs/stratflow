import { timestampDate } from "@bufbuild/protobuf/wkt";
import {
  Alert,
  Chip,
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

import type { User } from "@/shared/api/gen/identity/proto/v1/types_pb";

const formatDate = (value?: Parameters<typeof timestampDate>[0] | null) => {
  if (!value) {
    return "—";
  }

  const date = timestampDate(value);
  return Number.isNaN(date.getTime()) ? "—" : date.toLocaleString();
};

const renderRole = (role?: string) => role?.trim() || "—";

const renderVerification = (value: boolean) => (
  <Chip
    size="small"
    color={value ? "success" : "default"}
    label={value ? "Verified" : "Pending"}
    variant={value ? "filled" : "outlined"}
  />
);

type UsersTableProps = {
  items: User[];
  isLoading: boolean;
  query: string;
};

export const UsersTable = ({ items, isLoading, query }: UsersTableProps) => {
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
        <Typography color="text.secondary">Loading users...</Typography>
      </Stack>
    );
  }

  return (
    <TableContainer sx={{ width: "100%", overflowX: "auto" }}>
      <Table
        size="small"
        aria-label="Users"
        aria-busy={isLoading}
        sx={{ minWidth: "74rem", tableLayout: "fixed" }}
      >
        <TableHead>
          <TableRow>
            <TableCell sx={{ width: "18ch" }}>Login</TableCell>
            <TableCell sx={{ width: "22ch" }}>Name</TableCell>
            <TableCell sx={{ width: "28ch" }}>Email</TableCell>
            <TableCell sx={{ width: "12ch" }}>Role</TableCell>
            <TableCell sx={{ width: "14ch" }}>Email status</TableCell>
            <TableCell sx={{ width: "18ch" }}>Created</TableCell>
            <TableCell sx={{ width: "18ch" }}>Updated</TableCell>
          </TableRow>
        </TableHead>

        <TableBody>
          {items.length === 0 ? (
            <TableRow>
              <TableCell colSpan={7}>
                <Alert severity="info">
                  {query.trim()
                    ? "No users match the current search query."
                    : "No users returned by the identity service yet."}
                </Alert>
              </TableCell>
            </TableRow>
          ) : (
            items.map((user, index) => (
              <TableRow
                key={user.id || `${user.login}-${index}`}
                hover
                sx={{ "&:last-child td": { borderBottom: 0 } }}
              >
                <TableCell
                  sx={{
                    fontFamily: "monospace",
                    overflow: "hidden",
                    textOverflow: "ellipsis",
                    whiteSpace: "nowrap",
                  }}
                  title={user.login || undefined}
                >
                  {user.login || "—"}
                </TableCell>

                <TableCell
                  sx={{
                    fontWeight: 700,
                    overflow: "hidden",
                    textOverflow: "ellipsis",
                    whiteSpace: "nowrap",
                  }}
                  title={
                    [user.name, user.lastName].filter(Boolean).join(" ") ||
                    undefined
                  }
                >
                  {[user.name, user.lastName].filter(Boolean).join(" ") || "—"}
                </TableCell>

                <TableCell
                  sx={{
                    overflow: "hidden",
                    textOverflow: "ellipsis",
                    whiteSpace: "nowrap",
                  }}
                  title={user.email || undefined}
                >
                  {user.email || "—"}
                </TableCell>

                <TableCell>{renderRole(user.role)}</TableCell>
                <TableCell>
                  {renderVerification(user.isEmailVerified)}
                </TableCell>
                <TableCell>{formatDate(user.createdAt)}</TableCell>
                <TableCell>{formatDate(user.updatedAt)}</TableCell>
              </TableRow>
            ))
          )}
        </TableBody>
      </Table>
    </TableContainer>
  );
};
