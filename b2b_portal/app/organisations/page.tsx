import type { Metadata } from "next";
import Organisations from "@/components/dashboard/organisations/Organisations";
import { requireServerSessionRole } from "@lib/auth/session";

export const metadata: Metadata = {
  title: "Organisations | InsureTech B2B Portal",
};

export default async function OrganisationsPage() {
  await requireServerSessionRole(["SYSTEM_ADMIN"], "/");
  return <Organisations />;
}
