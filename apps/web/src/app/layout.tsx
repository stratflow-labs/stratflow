import type { Metadata } from "next";
import type { ReactNode } from "react";

import { AuthProvider } from "@/features/auth";
import { MuiProvider } from "@/shared/ui/providers/mui-provider";

export const metadata: Metadata = {
  title: "Stratflow",
  description: "Connect-Web frontend for Stratflow",
};

export default function RootLayout({ children }: { children: ReactNode }) {
  return (
    <html lang="en">
      <body>
        <MuiProvider>
          <AuthProvider>{children}</AuthProvider>
        </MuiProvider>
      </body>
    </html>
  );
}
