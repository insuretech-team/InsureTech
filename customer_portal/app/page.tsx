import DashboardLayout from "@/components/dashboard/dashboard-layout";
import OverviewActivity from "@/components/dashboard/overview-activity/overview-activity";
// import PolicyOverview from "@/components/dashboard/policy-overview/policy-overview";
// import QuickAccess from "@/components/dashboard/quick-access/quick-access";
import StatsCards from "@/components/dashboard/stats-cards/stats-cards";
// import UpcomingPayments from "@/components/dashboard/upcoming-payments/upcoming-payments";

export default function Home() {
  return (
    <DashboardLayout>
      <div className="flex flex-col gap-6">
        <StatsCards />
        <OverviewActivity />
        {/* <QuickAccess /> */}
        {/* <div className="grid gap-6 lg:grid-cols-2"> */}
        {/* <PolicyOverview /> */}
        {/* <UpcomingPayments /> */}
        {/* </div> */}
      </div>
    </DashboardLayout>
  );
}
