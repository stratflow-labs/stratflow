import { timestampDate } from "@bufbuild/protobuf/wkt";
import KeyboardArrowDownRoundedIcon from "@mui/icons-material/KeyboardArrowDownRounded";
import KeyboardArrowUpRoundedIcon from "@mui/icons-material/KeyboardArrowUpRounded";
import ShareRoundedIcon from "@mui/icons-material/ShareRounded";
import {
  Alert,
  Box,
  Button,
  ButtonBase,
  Chip,
  CircularProgress,
  Collapse,
  Divider,
  Stack,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Typography,
} from "@mui/material";

import type { AttributeWithValues, Strategy } from "@/shared/api/gen/strategy_registry/proto/v1/strategy_types_pb";

import type { StrategyDetailsEntry } from "./strategy-details-types";

const formatDate = (value?: Parameters<typeof timestampDate>[0] | null) => {
  if (!value) {
    return "—";
  }

  const date = timestampDate(value);
  return Number.isNaN(date.getTime()) ? "—" : date.toLocaleString();
};

export const StrategyDetailsHeader = ({
  strategy,
}: {
  strategy: Strategy | null;
}) => (
  <Stack spacing={0.5}>
    <Typography variant="h6" sx={{ fontWeight: 800 }}>
      {strategy?.name || strategy?.slug || "Strategy details"}
    </Typography>
    <Typography color="text.secondary" variant="body2">
      {strategy?.description || "Attributes and values for the selected strategy."}
    </Typography>
  </Stack>
);

export const StrategyDetailsToolbar = ({
  strategy,
  entry,
  onCreateAttribute,
  onOpenRelationsGraph,
  onRefresh,
}: {
  strategy: Strategy | null;
  entry: StrategyDetailsEntry;
  onCreateAttribute: () => void;
  onOpenRelationsGraph: () => void;
  onRefresh: () => void | Promise<void>;
}) => (
  <Stack
    direction={{ xs: "column", sm: "row" }}
    spacing={1.5}
    sx={{ justifyContent: "space-between", alignItems: { sm: "center" } }}
  >
    <Stack direction="row" spacing={1} useFlexGap sx={{ flexWrap: "wrap" }}>
      <Chip size="small" variant="outlined" label={strategy?.slug || "—"} />
      <Chip
        size="small"
        color="primary"
        variant="outlined"
        label={`${entry.total} attributes`}
      />
    </Stack>

    <Stack direction="row" spacing={1}>
      <Button variant="contained" onClick={onCreateAttribute}>
        Create attribute
      </Button>
      <Button
        variant="outlined"
        startIcon={<ShareRoundedIcon fontSize="small" />}
        onClick={onOpenRelationsGraph}
        disabled={entry.items.length === 0}
      >
        Open graph
      </Button>
      <Button variant="outlined" onClick={() => void onRefresh()}>
        Refresh attributes
      </Button>
    </Stack>
  </Stack>
);

export const StrategyDetailsStatus = ({
  entry,
  onRefresh,
}: {
  entry: StrategyDetailsEntry;
  onRefresh: () => void | Promise<void>;
}) => (
  <>
    {entry.status === "loading" ? (
      <Stack direction="row" spacing={1.5} sx={{ alignItems: "center", py: 2 }}>
        <CircularProgress size={20} />
        <Typography color="text.secondary">Loading attributes...</Typography>
      </Stack>
    ) : null}

    {entry.status === "error" ? (
      <Alert
        severity="error"
        action={
          <Button color="inherit" size="small" onClick={() => void onRefresh()}>
            Retry
          </Button>
        }
      >
        {entry.error}
      </Alert>
    ) : null}

    {entry.status !== "loading" && entry.status !== "error" && entry.items.length === 0 ? (
      <Alert severity="info">This strategy has no attributes yet.</Alert>
    ) : null}
  </>
);

const StrategyAttributeValues = ({
  attribute,
}: {
  attribute: AttributeWithValues;
}) =>
  attribute.values.length > 0 ? (
    <Stack spacing={1.5}>
      <Typography variant="subtitle2" sx={{ fontWeight: 700 }}>
        Attribute values
      </Typography>
      <Table size="small" aria-label="Attribute values">
        <TableHead>
          <TableRow>
            <TableCell sx={{ width: 220 }}>Value</TableCell>
            <TableCell sx={{ width: 180 }}>Slug</TableCell>
            <TableCell sx={{ width: 140 }}>Relations</TableCell>
            <TableCell sx={{ width: 160 }}>Updated</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {attribute.values.map((value) => (
            <TableRow
              key={
                value.id ||
                `${attribute.id}-${value.slug}-${value.value}`
              }
            >
              <TableCell sx={{ fontWeight: 600 }}>
                {value.value || "—"}
              </TableCell>
              <TableCell>{value.slug || "—"}</TableCell>
              <TableCell>{value.relations.length}</TableCell>
              <TableCell>{formatDate(value.updatedAt)}</TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </Stack>
  ) : (
    <Typography color="text.secondary" variant="body2">
      No values for this attribute.
    </Typography>
  );

const StrategyAttributeRow = ({
  attribute,
  isExpanded,
  isLast,
  onToggleAttribute,
  onCreateValue,
}: {
  attribute: AttributeWithValues;
  isExpanded: boolean;
  isLast: boolean;
  onToggleAttribute: (attributeId: string) => void;
  onCreateValue: (attributeLocalId: string) => void;
}) => (
  <TableRow>
    <TableCell colSpan={5} sx={{ p: 0, borderBottom: 0 }}>
      <Table size="small">
        <TableBody>
          <TableRow hover sx={{ "& > td": { borderBottom: isExpanded ? 0 : undefined } }}>
            <TableCell padding="checkbox" sx={{ width: 56 }}>
              <ButtonBase
                aria-label={
                  isExpanded
                    ? `Collapse ${attribute.name || attribute.slug || "attribute"}`
                    : `Expand ${attribute.name || attribute.slug || "attribute"}`
                }
                onClick={() => onToggleAttribute(attribute.id)}
                sx={{
                  width: 32,
                  height: 32,
                  borderRadius: 1,
                  display: "inline-flex",
                  alignItems: "center",
                  justifyContent: "center",
                }}
              >
                {isExpanded ? (
                  <KeyboardArrowUpRoundedIcon fontSize="small" />
                ) : (
                  <KeyboardArrowDownRoundedIcon fontSize="small" />
                )}
              </ButtonBase>
            </TableCell>

            <TableCell sx={{ width: 240 }}>
              <Stack spacing={0.5}>
                <Typography sx={{ fontWeight: 700 }}>
                  {attribute.name || "Untitled attribute"}
                </Typography>
                <Chip
                  size="small"
                  variant="outlined"
                  label={attribute.slug || "—"}
                  sx={{ alignSelf: "flex-start" }}
                />
              </Stack>
            </TableCell>

            <TableCell sx={{ width: 140 }}>
              <Chip
                size="small"
                color="primary"
                variant="outlined"
                label={`${attribute.values.length} values`}
              />
            </TableCell>

            <TableCell>
              <Typography color="text.secondary" variant="body2">
                {attribute.description || "No description"}
              </Typography>
            </TableCell>

            <TableCell sx={{ width: 180 }}>
              {formatDate(attribute.updatedAt)}
            </TableCell>
          </TableRow>

          <TableRow>
            <TableCell colSpan={5} sx={{ py: 0 }}>
              <Collapse in={isExpanded} timeout="auto" unmountOnExit>
                <Box sx={{ px: 2, py: 2, bgcolor: "background.default" }}>
                  <Stack spacing={1.5}>
                    <Button
                      variant="outlined"
                      size="small"
                      onClick={() => onCreateValue(attribute.id || attribute.slug)}
                      sx={{ alignSelf: "flex-start" }}
                    >
                      Create value
                    </Button>

                    <StrategyAttributeValues attribute={attribute} />
                  </Stack>
                </Box>
              </Collapse>
            </TableCell>
          </TableRow>
        </TableBody>
      </Table>

      {!isLast ? <Divider sx={{ mx: 2 }} /> : null}
    </TableCell>
  </TableRow>
);

export const StrategyAttributesTable = ({
  entry,
  expandedAttributeId,
  onToggleAttribute,
  onCreateValue,
}: {
  entry: StrategyDetailsEntry;
  expandedAttributeId: string | null;
  onToggleAttribute: (attributeId: string) => void;
  onCreateValue: (attributeLocalId: string) => void;
}) =>
  entry.items.length > 0 ? (
    <TableContainer>
      <Table size="small" aria-label="Strategy attributes">
        <TableHead>
          <TableRow>
            <TableCell sx={{ width: 56 }} />
            <TableCell sx={{ width: 240 }}>Attribute</TableCell>
            <TableCell sx={{ width: 140 }}>Values</TableCell>
            <TableCell>Description</TableCell>
            <TableCell sx={{ width: 180 }}>Updated</TableCell>
          </TableRow>
        </TableHead>

        <TableBody>
          {entry.items.map((attribute, index) => (
            <StrategyAttributeRow
              key={attribute.id || attribute.slug || index}
              attribute={attribute}
              isExpanded={expandedAttributeId === attribute.id}
              isLast={index === entry.items.length - 1}
              onToggleAttribute={onToggleAttribute}
              onCreateValue={onCreateValue}
            />
          ))}
        </TableBody>
      </Table>
    </TableContainer>
  ) : null;
