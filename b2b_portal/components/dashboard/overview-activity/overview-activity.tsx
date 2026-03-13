"use client";

import * as React from "react";
import { LuLoader, LuBuilding2, LuUsers, LuLayers, LuClipboardList, LuUserPlus } from "react-icons/lu";

type ActivityItem = {
  id: string;
  icon: React.ReactNode;
  title: string;
  subtitle: string;
  time: string;
};

function timeAgo(dateStr: string): string {
  if (!dateStr) return "";
  try {
    const diff = Date.now() - new Date(dateStr).getTime();
    const mins = Math.floor(diff / 60000);
    if (mins < 1)  return "just now";
    if (mins < 60) return `${mins}m ago`;
    const hrs = Math.floor(mins / 60);
    if (hrs < 24)  return `${hrs}h ago`;
    const days = Math.floor(hrs / 24);
    if (days < 7)  return `${days}d ago`;
    return new Date(dateStr).toLocaleDateString();
  } catch { return ""; }
}

function iconFor(type: string): React.ReactNode {
  switch (type) {
    case "org":        return <LuBuilding2 className="size-4 text-primary" />;
    case "member":     return <LuUserPlus  className="size-4 text-violet-500" />;
    case "employee":   return <LuUsers     className="size-4 text-blue-500" />;
    case "department": return <LuLayers    className="size-4 text-amber-500" />;
    case "po":         return <LuClipboardList className="size-4 text-green-500" />;
    default:           return <LuBuilding2 className="size-4 text-muted-foreground" />;
  }
}

type ActivityApiResponse = {
  ok: boolean;
  activities?: Array<{ id: string; type: string; title: string; subtitle: string; createdAt: string }>;
  message?: string;
};

const OverviewActivity = () => {
  const [activities, setActivities] = React.useState<ActivityItem[]>([]);
  const [loading, setLoading] = React.useState(true);
  const [error, setError] = React.useState("");

  React.useEffect(() => {
    let cancelled = false;
    void fetch("/api/dashboard/activity", { cache: "no-store" })
      .then((r) => r.json())
      .then((payload: ActivityApiResponse) => {
        if (cancelled) return;
        if (payload.ok && payload.activities) {
          setActivities(payload.activities.map((a) => ({
            id:       a.id,
            icon:     iconFor(a.type),
            title:    a.title,
            subtitle: a.subtitle,
            time:     timeAgo(a.createdAt),
          })));
        } else {
          setError(payload.message ?? "Failed to load activity");
        }
      })
      .catch(() => { if (!cancelled) setError("Failed to load activity"); })
      .finally(() => { if (!cancelled) setLoading(false); });
    return () => { cancelled = true; };
  }, []);

  return (
    <div className="portal-panel p-5 space-y-3">
      <h2 className="text-lg font-semibold text-foreground">Recent Activity</h2>
      {loading ? (
        <div className="flex items-center justify-center py-8">
          <LuLoader className="size-5 animate-spin text-muted-foreground" />
        </div>
      ) : error ? (
        <p className="text-sm text-muted-foreground py-4 text-center">{error}</p>
      ) : activities.length === 0 ? (
        <p className="text-sm text-muted-foreground py-4 text-center">No recent activity.</p>
      ) : (
        <div className="space-y-1">
          {activities.map((a) => (
            <div key={a.id} className="flex items-start gap-3 rounded-lg px-2 py-2.5 hover:bg-muted/40 transition-colors">
              <div className="mt-0.5 flex size-7 shrink-0 items-center justify-center rounded-full bg-muted">
                {a.icon}
              </div>
              <div className="flex-1 min-w-0">
                <p className="text-sm font-medium text-foreground truncate">{a.title}</p>
                <p className="text-xs text-muted-foreground truncate">{a.subtitle}</p>
              </div>
              <span className="shrink-0 text-xs text-muted-foreground">{a.time}</span>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default OverviewActivity;
