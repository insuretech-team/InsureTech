/**
 * Organisations.tsx
 * ──────────────────
 * Full CRUD page component for Organisations (Super Admin view).
 * Row click opens OrgDetailPanel with Info / Members / Departments tabs.
 */
"use client";

import { useState } from "react";
import DashboardLayout from "../dashboard-layout";
import { useCrudList } from "@/src/hooks/useCrudList";
import { organisationClient } from "@lib/sdk/organisation-client";
import { buildOrganisationColumns } from "./data-table/columns";
import { OrganisationDataTable } from "./data-table/data-table";
import { OrgDetailPanel } from "@/components/organisations/org-detail-panel";
import { authClient } from "@lib/sdk/auth-client";
import { useEffect } from "react";
import type { Organisation } from "@lib/types/b2b";

export default function Organisations() {
  const { data, loading, reload } = useCrudList<Organisation>(
    () => organisationClient.list(),
    "organisations"
  );

  const [selectedOrg, setSelectedOrg] = useState<Organisation | null>(null);
  const [panelOpen, setPanelOpen] = useState(false);
  const [currentUserRole, setCurrentUserRole] = useState("SYSTEM_ADMIN");

  useEffect(() => {
    authClient.getSession().then((res) => {
      setCurrentUserRole(res.session?.principal.role ?? "SYSTEM_ADMIN");
    }).catch(() => {});
  }, []);

  function handleRowClick(org: Organisation) {
    setSelectedOrg(org);
    setPanelOpen(true);
  }

  const columns = buildOrganisationColumns(
    reload,
    // onView — opens detail panel
    handleRowClick,
    // onApprove — refresh + update panel if same org is currently open
    (org) => {
      if (selectedOrg?.id === org.id) {
        setSelectedOrg((prev) => prev ? { ...prev, status: "ORGANISATION_STATUS_ACTIVE" } : prev);
      }
    },
  );

  return (
    <DashboardLayout>
      <div className="space-y-4">
        <OrganisationDataTable
          columns={columns}
          data={data}
          loading={loading}
          onRefresh={reload}
        />
      </div>
      <OrgDetailPanel
        org={selectedOrg}
        currentUserRole={currentUserRole}
        open={panelOpen}
        onClose={() => setPanelOpen(false)}
        onRefresh={reload}
      />
    </DashboardLayout>
  );
}
