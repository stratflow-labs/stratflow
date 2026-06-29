"use client";

import {
  findAttributeDraft,
  findValueDraft,
  type GraphDraftAttribute,
} from "./relations";

export type RenameAttributeInput = {
  slug: string;
  name: string;
  description: string;
};

export type RenameValueInput = {
  slug: string;
  value: string;
};

export const promptRenameAttribute = (
  draft: GraphDraftAttribute[],
  attributeLocalId: string,
): RenameAttributeInput | null => {
  const attribute = findAttributeDraft(draft, attributeLocalId);
  if (!attribute) {
    return null;
  }

  const nextSlug = window.prompt("Attribute slug", attribute.slug)?.trim();
  if (!nextSlug) {
    return null;
  }

  return {
    slug: nextSlug,
    name: window.prompt("Attribute name", attribute.name)?.trim() ?? attribute.name,
    description:
      window.prompt("Attribute description", attribute.description)?.trim() ??
      attribute.description,
  };
};

export const promptRenameValue = (
  draft: GraphDraftAttribute[],
  valueLocalId: string,
): RenameValueInput | null => {
  const value = findValueDraft(draft, valueLocalId);
  if (!value) {
    return null;
  }

  const nextSlug = window.prompt("Value slug", value.slug)?.trim();
  if (!nextSlug) {
    return null;
  }

  return {
    slug: nextSlug,
    value: window.prompt("Value label", value.value)?.trim() ?? value.value,
  };
};
