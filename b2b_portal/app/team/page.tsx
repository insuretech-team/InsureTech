import type { Metadata } from "next";
import TeamManagementPage from "@/components/team/team-management";
import { requireServerSessionRole } from "@lib/auth/session";

export const metadata: Metadata = {
  title: "Team | InsureTech B2B Portal",
};

export default async function TeamPage() {
  await requireServerSessionRole(["SYSTEM_ADMIN", "B2B_ORG_ADMIN", "BUSINESS_ADMIN"], "/");
  return <TeamManagementPage />;
}
