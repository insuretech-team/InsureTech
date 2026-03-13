import type { Metadata } from "next";
import Organisations from "@/components/dashboard/organisations/Organisations";

export const metadata: Metadata = {
  title: "Organisations | InsureTech B2B Portal",
};

export default function OrganisationsPage() {
  return <Organisations />;
}
