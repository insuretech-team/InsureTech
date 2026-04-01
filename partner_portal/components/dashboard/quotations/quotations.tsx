import DashboardLayout from "../dashboard-layout";
import { DataTable } from "./data-table/data-table";
import QuotationCard from "./card";

import { quotations } from "@/lib/quotations";
import {
  quotationColumns,
  type Quotation,
} from "@/components/dashboard/quotations/data-table/columns";

const data: Quotation[] = [
  {
    id: "01",
    quotationID: "QUO-2025-001",
    insurerName: "ABC Insurance",
    plan: "Seba",
    insuranceCategory: "Health",
    department: "Tech",
    employeeNo: 32,
    estimatedPremium: "৳ 432,650",
    quotedAmount: "৳ 432,650",
    status: "Pending",
    submissionDate: "2025-01-12",
    validUntil: "2025-02-12",
  },
  {
    id: "02",
    quotationID: "QUO-2025-002",
    insurerName: "Delta Assurance",
    plan: "Premium Care",
    insuranceCategory: "Health",
    department: "Accounts",
    employeeNo: 18,
    estimatedPremium: "৳ 210,000",
    quotedAmount: "৳ 205,000",
    status: "Approved",
    submissionDate: "2025-01-20",
    validUntil: "2025-02-20",
  },
  {
    id: "03",
    quotationID: "QUO-2025-003",
    insurerName: "Guardian Life",
    plan: "Standard",
    insuranceCategory: "Life",
    department: "Business",
    employeeNo: 45,
    estimatedPremium: "৳ 520,000",
    quotedAmount: "৳ 515,000",
    status: "Rejected",
    submissionDate: "2025-01-28",
    validUntil: "2025-02-28",
  },
];

const Quotations = () => {
  return (
    <DashboardLayout>
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-6">
        {quotations.map((quote) => (
          <QuotationCard key={quote.id} {...quote} />
        ))}
      </div>

      <div className="space-y-4 py-4">
        <DataTable columns={quotationColumns} data={data} loading={false} />
      </div>
    </DashboardLayout>
  );
};

export default Quotations;
