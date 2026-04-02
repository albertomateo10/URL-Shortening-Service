"use client";

import {
  AreaChart,
  Area,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from "recharts";
import type { DailyClickCount } from "@/types/api";

interface ClicksOverTimeChartProps {
  data: DailyClickCount[];
}

export default function ClicksOverTimeChart({ data }: ClicksOverTimeChartProps) {
  if (data.length === 0) {
    return (
      <p className="text-gray-400 text-sm text-center py-12">
        No click data yet for this period.
      </p>
    );
  }

  return (
    <ResponsiveContainer width="100%" height={300}>
      <AreaChart data={data} margin={{ top: 5, right: 20, left: 0, bottom: 5 }}>
        <CartesianGrid strokeDasharray="3 3" stroke="#e5e7eb" />
        <XAxis
          dataKey="date"
          tick={{ fontSize: 12, fill: "#6b7280" }}
          tickFormatter={(v: string) => {
            const d = new Date(v + "T00:00:00");
            return d.toLocaleDateString("en-US", { month: "short", day: "numeric" });
          }}
        />
        <YAxis
          tick={{ fontSize: 12, fill: "#6b7280" }}
          allowDecimals={false}
        />
        <Tooltip
          labelFormatter={(v) => {
            const d = new Date(String(v) + "T00:00:00");
            return d.toLocaleDateString("en-US", {
              month: "long",
              day: "numeric",
              year: "numeric",
            });
          }}
          formatter={(value) => [String(value), "Clicks"]}
        />
        <Area
          type="monotone"
          dataKey="count"
          stroke="#2563eb"
          fill="#3b82f6"
          fillOpacity={0.15}
          strokeWidth={2}
        />
      </AreaChart>
    </ResponsiveContainer>
  );
}
