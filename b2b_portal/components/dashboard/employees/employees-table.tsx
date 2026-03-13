"use client";

import * as React from "react";
import { DataTable } from "@/components/dashboard/employees/data-table/data-table";
import { buildEmployeeColumns } from "@/components/dashboard/employees/data-table/columns";
import DashboardLayout from "../dashboard-layout";
import { employeeClient } from "@lib/sdk/employee-client";
import { organisationClient } from "@lib/sdk/organisation-client";
import { authClient } from "@lib/sdk/auth-client";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import type { Employee, Organisation } from "@lib/types/b2b";

export default function EmployeesPage() {
  const [organisations, setOrganisations] = React.useState<Organisation[]>([]);
  const [selectedOrgId, setSelectedOrgId] = React.useState("");
  const [resolvedOrg, setResolvedOrg] = React.useState<Organisation | null>(null);
  const [data, setData] = React.useState<Employee[]>([]);
  const [loading, setLoading] = React.useState(true);
  // true = b2b_admin (locked to own org), false = super_admin (sees dropdown)
  const [isB2BAdmin, setIsB2BAdmin] = React.useState<boolean | null>(null);

  // Step 1: Determine role from session, then load org context accordingly.
  React.useEffect(() => {
    let cancelled = false;

    void authClient.getSession().then((response) => {
      if (cancelled) return;
      const session = response.session;
      const role = session?.principal.role ?? "";
      const isSuperAdmin = role === "SYSTEM_ADMIN";
      setIsB2BAdmin(!isSuperAdmin);

      if (isSuperAdmin) {
        // Super admin: load all orgs for the dropdown
        void organisationClient.list().then((listResult) => {
          if (cancelled || !listResult.ok) return;
          const rows = listResult.organisations ?? [];
          setOrganisations(rows);
          setSelectedOrgId((current) => current || rows[0]?.id || "");
        });
      } else {
        // B2B admin: resolve their own org from /organisations/me
        void organisationClient.getMe().then((result) => {
          if (cancelled || !result.ok) return;
          if (result.organisation?.id) {
            setResolvedOrg(result.organisation);
            setOrganisations([result.organisation]);
            setSelectedOrgId(result.organisation.id);
          }
        });
      }
    }).catch(() => {
      if (!cancelled) setIsB2BAdmin(true);
    });

    return () => { cancelled = true; };
  }, []);

  const reload = React.useCallback(async () => {
    if (!selectedOrgId) {
      setData([]);
      setLoading(false);
      return;
    }
    setLoading(true);
    try {
      // Pass business_id explicitly so the API route has it for both roles.
      // The route will also auto-resolve it from the session as a fallback.
      const result = await employeeClient.list({ businessId: selectedOrgId });
      setData(result.ok ? (result.employees ?? []) : []);
    } finally {
      setLoading(false);
    }
  }, [selectedOrgId]);

  React.useEffect(() => {
    void reload();
  }, [reload]);

  const columns = buildEmployeeColumns(reload);

  return (
    <DashboardLayout>
      <div className="space-y-4">
        <div className="portal-panel p-4">
          <div className="flex flex-col gap-2 md:flex-row md:items-center md:justify-between">
            <div>
              <div className="text-sm font-medium text-foreground">Organisation Scope</div>
              <div className="text-xs text-muted-foreground">
                {isB2BAdmin
                  ? "You are viewing employees in your organisation."
                  : "Select an organisation to view its employees."}
              </div>
            </div>
            <div className="w-full md:w-80">
              {/* B2B admin: read-only label showing their org name */}
              {isB2BAdmin === true ? (
                <div className="rounded-md border bg-muted/30 px-3 py-2 text-sm font-medium text-foreground">
                  {resolvedOrg?.name ?? "Loading…"}
                </div>
              ) : isB2BAdmin === false ? (
                /* Super admin: full org dropdown */
                <Select value={selectedOrgId} onValueChange={setSelectedOrgId}>
                  <SelectTrigger>
                    <SelectValue placeholder="Select organisation" />
                  </SelectTrigger>
                  <SelectContent>
                    {organisations.map((org) => (
                      <SelectItem key={org.id} value={org.id}>
                        {org.name}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              ) : null /* loading — render nothing yet */}
            </div>
          </div>
        </div>
        <DataTable
          columns={columns}
          data={data}
          loading={loading}
          onRefresh={reload}
          organisationId={selectedOrgId}
        />
      </div>
    </DashboardLayout>
  );
}
