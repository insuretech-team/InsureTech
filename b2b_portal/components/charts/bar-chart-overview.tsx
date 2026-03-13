"use client";

import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  Tooltip,
  ResponsiveContainer,
} from "recharts";

const data = [
  { name: "Developer", value: 132000 },
  { name: "Sales", value: 76000 },
  { name: "Marketing", value: 97538 },
  { name: "Operation", value: 64000 },
  { name: "Accounts", value: 118000 },
];

const BarChartOverview = () => {
  return (
    <div className="rounded-lg border bg-white">
      {/* Header */}
      <div className="border-b px-4 py-3 flex items-center justify-between">
        <h3 className="text-normal font-semibold text-[#242424]">
          Premium by Department
        </h3>
        <select className="rounded-md border px-2 py-1 text-sm">
          <option>Monthly</option>
          <option>Yearly</option>
        </select>
      </div>

      {/* Chart */}
      <div className="px-4 py-4 h-[360px]">
        <ResponsiveContainer width="100%" height="100%">
          <BarChart data={data} barCategoryGap={20}>
            <XAxis
              dataKey="name"
              tick={{ fill: "#2b2b2b", fontSize: 12 }}
              axisLine={false}
              tickLine={false}
            />
            <YAxis
              tick={{ fill: "#2b2b2b", fontSize: 12 }}
              axisLine={false}
              tickLine={false}
            />
            <Tooltip
              cursor={{ fill: "rgba(0,0,0,0.04)" }}
              contentStyle={{
                borderRadius: 8,
                border: "1px solid #E5E7EB",
                fontSize: 12,
              }}
              formatter={(value: number | string | undefined) => [
                `Premium Amount: ${Number(value ?? 0).toLocaleString()}`,
              ]}
            />
            <Bar
              dataKey="value"
              radius={[6, 6, 0, 0]}
              fill="rgba(41,170,107,0.85)"
              activeBar={{
                fill: "var(--brand-jungle)",
              }}
            />
          </BarChart>
        </ResponsiveContainer>
      </div>
    </div>
  );
};

export default BarChartOverview;
