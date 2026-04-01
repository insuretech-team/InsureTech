import { redirect } from "next/navigation";

import DashboardLayout from "@/components/dashboard/dashboard-layout";
import EmployeeSelfService from "@/components/employee/employee-self-service";
import { requireServerSession } from "@lib/auth/session";

export default async function EmployeePage() {
  const session = await requireServerSession();
  if (session.principal.role !== "B2B_BENEFICIARY") {
    redirect("/");
  }

  return (
    <DashboardLayout>
      <EmployeeSelfService />
    </DashboardLayout>
  );
}
