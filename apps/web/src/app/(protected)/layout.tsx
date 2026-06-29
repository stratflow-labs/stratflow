import type { ReactNode } from "react";

import { ProtectedRoute } from "@/features/auth";
import { AppShell } from "@/widgets/app-shell";

export default function ProtectedLayout({ children }: { children: ReactNode }) {
  return (
    <ProtectedRoute>
      <AppShell>{children}</AppShell>
    </ProtectedRoute>
  );
}
