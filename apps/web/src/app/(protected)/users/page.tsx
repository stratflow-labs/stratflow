import { AdminRoute } from "@/features/auth";
import { UsersOverview } from "@/widgets/users-overview";

export default function UsersPage() {
  return (
    <AdminRoute>
      <UsersOverview />
    </AdminRoute>
  );
}
