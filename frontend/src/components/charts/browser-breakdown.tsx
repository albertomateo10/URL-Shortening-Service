"use client";

import { PieChart, Pie, Cell, Tooltip, ResponsiveContainer, Legend } from "recharts";
import type { BrowserCount } from "@/types/api";

const COLORS = ["#2563eb", "#7c3aed", "#db2777", "#ea580c", "#65a30d", "#0891b2"];

interface BrowserBreakdownChartProps {
  data: BrowserCount[];
}

export default function BrowserBreakdownChart({ data }: BrowserBreakdownChartProps) {
  if (data.length === 0) {
    return (
      <p className="text-gray-400 text-sm text-center py-12">
        No browser data yet for this period.
      </p>
    );
  }

  return (
    <ResponsiveContainer width="100%" height={300}>
      <PieChart>
        <Pie
          data={data}
          dataKey="count"
          nameKey="name"
          cx="50%"
          cy="50%"
          outerRadius={100}
          label={({ name, percent }) =>
            `${name ?? "Unknown"} ${((percent ?? 0) * 100).toFixed(0)}%`
          }
        >
          {data.map((_, i) => (
            <Cell key={i} fill={COLORS[i % COLORS.length]} />
          ))}
        </Pie>
        <Tooltip formatter={(value) => [String(value), "Clicks"]} />
        <Legend />
      </PieChart>
    </ResponsiveContainer>
  );
}
