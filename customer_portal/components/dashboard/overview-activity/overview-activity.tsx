import PolicyOverviewChart from "@/components/charts/policy-overview-chart";

type AlertItem = {
  id: number;
  count: number;
  label: string;
  barColor: string;
  bgColor: string;
};

type ActivityItem = {
  id: number;
  title: string;
  subtitle: string;
  time: string;
  dotColor: string;
};

const alerts: AlertItem[] = [
  {
    id: 1,
    count: 6,
    label: "Policies expiring in 30 days",
    barColor: "#FF0000",
    bgColor: "#FDE2E2",
  },
  {
    id: 2,
    count: 4,
    label: "Employees without coverage",
    barColor: "#FF7A00",
    bgColor: "#FDEBDD",
  },
  {
    id: 3,
    count: 3,
    label: "Premium due",
    barColor: "#F5C400",
    bgColor: "#FBF9D9",
  },
];

const activities: ActivityItem[] = [
  {
    id: 1,
    title: "Quotation requested",
    subtitle: "Health Insurance",
    time: "15 minutes ago",
    dotColor: "#7C3AED",
  },
  {
    id: 2,
    title: "Plan updated",
    subtitle: "Health Insurance with premium ৳ 6,200",
    time: "2 hours ago",
    dotColor: "#7C3AED",
  },
  {
    id: 3,
    title: "Policy expiring",
    subtitle: "1 health plan expiring in 16 days",
    time: "8 hours ago",
    dotColor: "#FF0000",
  },
  {
    id: 4,
    title: "Employee added",
    subtitle: "15 new employees added via CSV upload",
    time: "2 days ago",
    dotColor: "#7C3AED",
  },
];

const OverviewActivity = () => {
  return (
    <div className="grid gap-4 lg:grid-cols-3">
      {/* Column 1: Policy Overview */}
      <div className="lg:col-span-1">
        <PolicyOverviewChart />
      </div>

      {/* Column 2: Alerts & Reminders */}
      <div className="rounded-lg border bg-white overflow-hidden">
        <div className="border-b px-4 py-3">
          <h3 className="text-normal font-semibold text-[#242424]">
            Alerts &amp; Reminders
          </h3>
        </div>

        <div className="p-4 space-y-4">
          {alerts.map((a) => (
            <div
              key={a.id}
              className="relative rounded-md p-4"
              style={{ backgroundColor: a.bgColor }}
            >
              <span
                className="absolute left-0 top-0 h-full w-1 rounded-l-md"
                style={{ backgroundColor: a.barColor }}
              />
              <div className="pl-3">
                <div className="text-lg font-semibold text-[#111827]">
                  {a.count}
                </div>
                <div className="text-sm text-[#6B7280]">{a.label}</div>
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* Column 3: Recent Activity */}
      <div className="rounded-lg border bg-white overflow-hidden">
        <div className="border-b px-4 py-3">
          <h3 className="text-normal font-semibold text-[#242424]">
            Recent Activity
          </h3>
        </div>

        <div className="max-h-[260px] overflow-y-auto">
          {activities.map((item, idx) => (
            <div key={item.id}>
              <div className="flex gap-3 px-4 py-4">
                <span
                  className="mt-1 h-3 w-3 rounded-full shrink-0"
                  style={{ backgroundColor: item.dotColor }}
                />
                <div className="min-w-0">
                  <div className="text-sm font-semibold text-[#111827]">
                    {item.title}
                  </div>
                  <div className="text-sm text-[#4B5563]">{item.subtitle}</div>
                  <div className="text-sm text-[#9CA3AF]">{item.time}</div>
                </div>
              </div>

              {idx !== activities.length - 1 && (
                <div className="mx-4 border-b" />
              )}
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};

export default OverviewActivity;
