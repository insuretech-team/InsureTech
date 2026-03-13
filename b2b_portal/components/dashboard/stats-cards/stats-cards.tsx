"use client";

import * as React from "react";
import { LuLoader } from "react-icons/lu";
import StatsCard from "./card/stats-card";
import type { DashboardStats } from "@/app/api/dashboard/stats/route";

type StatItem = {
  title: string;
  value: string | number;
  icon: string;
  bgColor: string;
  bgIcon: string;
};

function buildSuperAdminStats(stats: DashboardStats): StatItem[] {
  return [
    {
      title:   "Total Organisations",
      value:   stats.pendingOrganisations
                 ? `${stats.totalOrganisations ?? 0} (${stats.pendingOrganisations} pending)`
                 : String(stats.totalOrganisations ?? 0),
      icon:    "./stats-cards/policies-icon.svg",
      bgColor: "#D7E0FF",
      bgIcon:  "./stats-cards/policies-lg.svg",
    },
    {
      title:   "Total Employees",
      value:   String(stats.totalEmployees),
      icon:    "./stats-cards/employee-count-icon.svg",
      bgColor: "#F4EBF9",
      bgIcon:  "./stats-cards/employee-count-lg.svg",
    },
    {
      title:   "Total Departments",
      value:   String(stats.totalDepartments),
      icon:    "./stats-cards/claim-icon.svg",
      bgColor: "#DDFDDF",
      bgIcon:  "./stats-cards/claim-lg.svg",
    },
    {
      title:   "Active Purchase Orders",
      value:   String(stats.activePurchaseOrders),
      icon:    "./stats-cards/actions-icon.svg",
      bgColor: "#FDDDDE",
      bgIcon:  "./stats-cards/actions-lg.svg",
    },
  ];
}

function buildB2BAdminStats(stats: DashboardStats): StatItem[] {
  return [
    {
      title:   "Team Members",
      value:   String(stats.totalMembers ?? 0),
      icon:    "./stats-cards/policies-icon.svg",
      bgColor: "#D7E0FF",
      bgIcon:  "./stats-cards/policies-lg.svg",
    },
    {
      title:   "Departments",
      value:   String(stats.totalDepartments),
      icon:    "./stats-cards/claim-icon.svg",
      bgColor: "#DDFDDF",
      bgIcon:  "./stats-cards/claim-lg.svg",
    },
    {
      title:   "Employees",
      value:   String(stats.totalEmployees),
      icon:    "./stats-cards/employee-count-icon.svg",
      bgColor: "#F4EBF9",
      bgIcon:  "./stats-cards/employee-count-lg.svg",
    },
    {
      title:   "Active Purchase Orders",
      value:   String(stats.activePurchaseOrders),
      icon:    "./stats-cards/actions-icon.svg",
      bgColor: "#FDDDDE",
      bgIcon:  "./stats-cards/actions-lg.svg",
    },
  ];
}

const StatsCards = () => {
  const [stats, setStats] = React.useState<DashboardStats | null>(null);
  const [role, setRole] = React.useState<string>("");
  const [loading, setLoading] = React.useState(true);
  const [error, setError] = React.useState("");

  React.useEffect(() => {
    let cancelled = false;
    void fetch("/api/dashboard/stats", { cache: "no-store" })
      .then((r) => r.json())
      .then((payload: { ok: boolean; stats?: DashboardStats; role?: string; message?: string }) => {
        if (cancelled) return;
        if (payload.ok && payload.stats) {
          setStats(payload.stats);
          setRole(payload.role ?? "");
        } else {
          setError(payload.message ?? "Failed to load stats");
        }
      })
      .catch(() => { if (!cancelled) setError("Failed to load stats"); })
      .finally(() => { if (!cancelled) setLoading(false); });
    return () => { cancelled = true; };
  }, []);

  if (loading) {
    return (
      <div className="grid grid-cols-1 gap-4 md:grid-cols-2 xl:grid-cols-4">
        {[1, 2, 3, 4].map((i) => (
          <div key={i} className="portal-panel p-5 flex items-center justify-center h-28">
            <LuLoader className="size-6 animate-spin text-muted-foreground" />
          </div>
        ))}
      </div>
    );
  }

  if (error || !stats) {
    return (
      <div className="rounded-lg border border-destructive/30 bg-destructive/5 px-4 py-3 text-sm text-destructive">
        {error || "Unable to load dashboard statistics."}
      </div>
    );
  }

  const items = role === "SYSTEM_ADMIN"
    ? buildSuperAdminStats(stats)
    : buildB2BAdminStats(stats);

  return (
    <div className="grid grid-cols-1 gap-4 md:grid-cols-2 xl:grid-cols-4">
      {items.map((stat) => (
        <StatsCard key={stat.title} {...stat} />
      ))}
    </div>
  );
};

export default StatsCards;
