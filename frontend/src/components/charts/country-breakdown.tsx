"use client";

import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from "recharts";
import type { CountryCount } from "@/types/api";

interface CountryBreakdownChartProps {
  data: CountryCount[];
}

export default function CountryBreakdownChart({ data }: CountryBreakdownChartProps) {
  if (data.length === 0) {
    return (
      <p className="text-gray-400 text-sm text-center py-12">
        No country data yet for this period.
      </p>
    );
  }

  const display = data.slice(0, 10).map((c) => ({
    ...c,
    label: c.name || c.code || "Unknown",
  }));

  return (
    <ResponsiveContainer width="100%" height={300}>
      <BarChart data={display} margin={{ top: 5, right: 20, left: 0, bottom: 5 }}>
        <CartesianGrid strokeDasharray="3 3" stroke="#e5e7eb" />
        <XAxis dataKey="label" tick={{ fontSize: 12, fill: "#6b7280" }} />
        <YAxis tick={{ fontSize: 12, fill: "#6b7280" }} allowDecimals={false} />
        <Tooltip formatter={(value) => [String(value), "Clicks"]} />
        <Bar dataKey="count" fill="#7c3aed" radius={[4, 4, 0, 0]} />
      </BarChart>
    </ResponsiveContainer>
  );
}
