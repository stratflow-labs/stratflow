import Menu from "@mui/material/Menu";
import MenuItem from "@mui/material/MenuItem";

import type { StrategyRelationNode } from "../model/relations";

type RelationGraphContextMenuProps = {
  contextNode: StrategyRelationNode | null;
  contextNodeIsSource: boolean;
  contextNodeIsConnected: boolean;
  effectiveArmedSourceNodeId: string | null;
  menuPosition: { mouseX: number; mouseY: number } | null;
  selectedSourceRelationsCount: number;
  onClose: () => void;
  onUseAsSource: () => void;
  onAddRelation: () => void;
  onDeleteRelation: () => void;
  onDeleteAll: () => void;
  onAddValue: (attributeLocalId: string) => void;
  onRenameAttribute: (attributeLocalId: string) => void;
  onRenameValue: (valueLocalId: string) => void;
  onDeleteValue: (valueLocalId: string) => void;
  onDeleteAttribute: (attributeLocalId: string) => void;
  onClearSourceSelection: () => void;
};

export const RelationGraphContextMenu = ({
  contextNode,
  contextNodeIsSource,
  contextNodeIsConnected,
  effectiveArmedSourceNodeId,
  menuPosition,
  selectedSourceRelationsCount,
  onClose,
  onUseAsSource,
  onAddRelation,
  onDeleteRelation,
  onDeleteAll,
  onAddValue,
  onRenameAttribute,
  onRenameValue,
  onDeleteValue,
  onDeleteAttribute,
  onClearSourceSelection,
}: RelationGraphContextMenuProps) => (
  <Menu
    open={menuPosition !== null}
    onClose={onClose}
    anchorReference="anchorPosition"
    anchorPosition={
      menuPosition
        ? { top: menuPosition.mouseY, left: menuPosition.mouseX }
        : undefined
    }
  >
    <MenuItem onClick={onUseAsSource} disabled={!contextNode}>
      Use as source
    </MenuItem>
    <MenuItem
      onClick={onAddRelation}
      disabled={!contextNode || contextNodeIsSource || contextNodeIsConnected}
    >
      Add relation from current source
    </MenuItem>
    <MenuItem
      onClick={onDeleteRelation}
      disabled={!contextNode || contextNodeIsSource || !contextNodeIsConnected}
    >
      Remove relation
    </MenuItem>
    <MenuItem
      onClick={onDeleteAll}
      disabled={!contextNodeIsSource || selectedSourceRelationsCount === 0}
    >
      Remove all source relations
    </MenuItem>
    <MenuItem
      onClick={() => contextNode && onAddValue(contextNode.localAttributeId)}
      disabled={!contextNode}
    >
      Add value to attribute
    </MenuItem>
    <MenuItem
      onClick={() => contextNode && onRenameAttribute(contextNode.localAttributeId)}
      disabled={!contextNode}
    >
      Rename attribute
    </MenuItem>
    <MenuItem
      onClick={() => contextNode && onRenameValue(contextNode.localValueId)}
      disabled={!contextNode}
    >
      Rename value
    </MenuItem>
    <MenuItem
      onClick={() => contextNode && onDeleteValue(contextNode.localValueId)}
      disabled={!contextNode}
    >
      Delete value
    </MenuItem>
    <MenuItem
      onClick={() => contextNode && onDeleteAttribute(contextNode.localAttributeId)}
      disabled={!contextNode}
    >
      Delete attribute
    </MenuItem>
    <MenuItem onClick={onClearSourceSelection} disabled={!effectiveArmedSourceNodeId}>
      Clear source selection
    </MenuItem>
  </Menu>
);
