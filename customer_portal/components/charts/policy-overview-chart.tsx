"use client";

import { PieChart, Pie, Cell, ResponsiveContainer } from "recharts";

type Item = {
  name: string;
  value: number;
  color: string;
};

const data: Item[] = [
  { name: "Health", value: 284, color: "#7C3AED" }, // purple
  { name: "Life", value: 116, color: "#818CF8" }, // indigo
];

export default function PolicyOverviewChart() {
  const total = data.reduce((sum, d) => sum + d.value, 0);

  return (
    <div className="rounded-lg border bg-white">
      {/* Header */}
      <div className="border-b px-4 py-3">
        <h3 className="text-normal font-semibold text-[#242424]">
          Policy Overview
        </h3>
      </div>

      <div className="px-4 py-4">
        {/* Donut */}
        <div className="h-48 w-full">
          <ResponsiveContainer width="100%" height="100%">
            <PieChart>
              <Pie
                data={data}
                dataKey="value"
                innerRadius={58}
                outerRadius={82}
                paddingAngle={2}
                stroke="#ffffff"
                strokeWidth={2}
                startAngle={90}
                endAngle={-270}
              >
                {data.map((entry) => (
                  <Cell key={entry.name} fill={entry.color} />
                ))}
              </Pie>
            </PieChart>
          </ResponsiveContainer>
        </div>

        {/* Legend (like screenshot) */}
        <div className="mt-3 space-y-3">
          {data.map((item) => {
            const percent = Math.round((item.value / total) * 100);
            return (
              <div
                key={item.name}
                className="flex items-center justify-between"
              >
                <div className="flex items-center gap-2">
                  <span
                    className="h-3 w-3 rounded-full"
                    style={{ backgroundColor: item.color }}
                  />
                  <span className="text-sm text-[#4B5563]">{item.name}</span>
                  <span className="text-sm text-[#9CA3AF]">({percent}%)</span>
                </div>
                <span className="text-sm font-medium text-[#111827]">
                  {item.value}
                </span>
              </div>
            );
          })}
        </div>
      </div>
    </div>
  );
}
