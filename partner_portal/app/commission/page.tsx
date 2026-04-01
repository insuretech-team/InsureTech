import DashboardLayout from "@/components/dashboard/dashboard-layout";

export default function CommissionPage() {
  return (
    <DashboardLayout>
      <div className="flex flex-col gap-6">
        <div>
          <h1 className="text-2xl font-semibold">Commission</h1>
          <p className="text-muted-foreground">
            Track commission earnings and payment history
          </p>
        </div>
        {/* Commission tracking will be implemented here */}
      </div>
    </DashboardLayout>
  );
}
