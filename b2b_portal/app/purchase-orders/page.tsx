import PurchaseOrders from "@/components/dashboard/purchase-orders/purchase-orders";
import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "Purchase Orders | Labaid Insuretech B2B Dashboard",
};

export default function Page() {
  return <PurchaseOrders />;
}
