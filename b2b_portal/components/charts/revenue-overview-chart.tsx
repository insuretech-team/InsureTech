"use client";

import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  Tooltip,
  ResponsiveContainer,
} from "recharts";

const data = [
  { month: "Jan", value: 28000 },
  { month: "Feb", value: 48000 },
  { month: "Mar", value: 47000 },
  { month: "Apr", value: 52000 },
  { month: "May", value: 61000 },
  { month: "Jun", value: 72000 },
  { month: "Jul", value: 73997 },
  { month: "Aug", value: 71000 },
  { month: "Sep", value: 71000 },
  { month: "Oct", value: 52000 },
  { month: "Nov", value: 48000 },
];

const RevenueOverviewChart = () => {
  return (
    <div className="rounded-lg border bg-white">
      {/* Header */}
      <div className="border-b px-4 py-3 flex items-center justify-between">
        <h3 className="text-normal font-semibold text-[#242424]">
          Revenue Report
        </h3>
        <select className="rounded-md border px-2 py-1 text-sm">
          <option>Monthly</option>
          <option>Yearly</option>
        </select>
      </div>

      {/* Chart */}
      <div className="px-4 py-4 h-[360px]">
        <ResponsiveContainer width="100%" height="100%">
          <LineChart data={data}>
            <XAxis
              dataKey="month"
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
              contentStyle={{
                borderRadius: 8,
                border: "1px solid #E5E7EB",
                fontSize: 12,
              }}
              formatter={(value: number | string | undefined) => [
                `${Number(value ?? 0).toLocaleString()}`,
                "Revenue",
              ]}
            />
            <Line
              type="monotone"
              dataKey="value"
              stroke="var(--primary)"
              strokeWidth={2}
              dot={{
                r: 4,
                stroke: "var(--primary)",
                strokeWidth: 2,
                fill: "#fff",
              }}
              activeDot={{
                r: 6,
                stroke: "var(--primary)",
                strokeWidth: 2,
                fill: "#fff",
              }}
            />
          </LineChart>
        </ResponsiveContainer>
      </div>
    </div>
  );
};

export default RevenueOverviewChart;
