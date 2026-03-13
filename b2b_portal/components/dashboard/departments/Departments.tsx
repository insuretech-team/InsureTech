"use client";

import * as React from "react";
import DashboardLayout from "../dashboard-layout";
import { DataTable } from "./data-table/data-table";
import { buildDepartmentColumns } from "@/components/dashboard/departments/data-table/columns";
import AddDepartmentModal from "@/components/modals/add-department-modal";
import { departmentClient } from "@lib/sdk/department-client";
import { organisationClient } from "@lib/sdk/organisation-client";
import { authClient } from "@lib/sdk/auth-client";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import type { Department, Organisation } from "@lib/types/b2b";

const Departments = () => {
  const [data, setData] = React.useState<Department[]>([]);
  const [loading, setLoading] = React.useState(true);
  const [addOpen, setAddOpen] = React.useState(false);

  // Org context — mirrors the employees-table pattern
  const [organisations, setOrganisations] = React.useState<Organisation[]>([]);
  const [selectedOrgId, setSelectedOrgId] = React.useState("");
  const [resolvedOrg, setResolvedOrg] = React.useState<Organisation | null>(null);
  const [isB2BAdmin, setIsB2BAdmin] = React.useState<boolean | null>(null);

  // Step 1: resolve role → load org context
  React.useEffect(() => {
    let cancelled = false;
    void authClient.getSession().then((response) => {
      if (cancelled) return;
      const role = response.session?.principal.role ?? "";
      const isSuperAdmin = role === "SYSTEM_ADMIN";
      setIsB2BAdmin(!isSuperAdmin);

      if (isSuperAdmin) {
        void organisationClient.list().then((res) => {
          if (cancelled || !res.ok) return;
          const rows = res.organisations ?? [];
          setOrganisations(rows);
          setSelectedOrgId((cur) => cur || rows[0]?.id || "");
        });
      } else {
        void organisationClient.getMe().then((res) => {
          if (cancelled || !res.ok) return;
          if (res.organisation?.id) {
            setResolvedOrg(res.organisation);
            setOrganisations([res.organisation]);
            setSelectedOrgId(res.organisation.id);
          }
        });
      }
    }).catch(() => { if (!cancelled) setIsB2BAdmin(true); });
    return () => { cancelled = true; };
  }, []);

  // Step 2: load departments whenever selected org changes
  const reload = React.useCallback(async () => {
    if (!selectedOrgId) { setData([]); setLoading(false); return; }
    setLoading(true);
    try {
      const result = await departmentClient.list(50, 0, selectedOrgId);
      setData(result.ok ? (result.departments ?? []) : []);
    } finally {
      setLoading(false);
    }
  }, [selectedOrgId]);

  React.useEffect(() => { void reload(); }, [reload]);

  const columns = buildDepartmentColumns(reload);

  return (
    <DashboardLayout>
      <div className="space-y-4">
        <AddDepartmentModal
          open={addOpen}
          onOpenChange={setAddOpen}
          organisationId={selectedOrgId}
          onSaved={() => { setAddOpen(false); void reload(); }}
        />

        {/* Org scope selector — same pattern as employees page */}
        <div className="portal-panel p-4">
          <div className="flex flex-col gap-2 md:flex-row md:items-center md:justify-between">
            <div>
              <div className="text-sm font-medium text-foreground">Organisation Scope</div>
              <div className="text-xs text-muted-foreground">
                {isB2BAdmin
                  ? "You are viewing departments in your organisation."
                  : "Select an organisation to view its departments."}
              </div>
            </div>
            <div className="w-full md:w-80">
              {isB2BAdmin === true ? (
                <div className="rounded-md border bg-muted/30 px-3 py-2 text-sm font-medium text-foreground">
                  {resolvedOrg?.name ?? "Loading…"}
                </div>
              ) : isB2BAdmin === false ? (
                <Select value={selectedOrgId} onValueChange={setSelectedOrgId}>
                  <SelectTrigger>
                    <SelectValue placeholder="Select organisation" />
                  </SelectTrigger>
                  <SelectContent>
                    {organisations.map((org) => (
                      <SelectItem key={org.id} value={org.id}>{org.name}</SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              ) : null}
            </div>
          </div>
        </div>

        <DataTable
          columns={columns}
          data={data}
          loading={loading}
          organisationId={selectedOrgId}
          onAddClick={() => setAddOpen(true)}
          onRefresh={reload}
        />
      </div>
    </DashboardLayout>
  );
};

export default Departments;
