import PurchaseOrders from "@/components/dashboard/purchase-orders/purchase-orders";
import type { Metadata } from "next";
import { requireServerSessionRole } from "@lib/auth/session";

export const metadata: Metadata = {
  title: "Purchase Orders | Labaid Insuretech B2B Dashboard",
};

export default async function Page() {
  await requireServerSessionRole(
    ["SYSTEM_ADMIN", "B2B_ORG_ADMIN", "BUSINESS_ADMIN", "HR_MANAGER", "FINANCE_MANAGER"],
    "/"
  );
  return <PurchaseOrders />;
}
