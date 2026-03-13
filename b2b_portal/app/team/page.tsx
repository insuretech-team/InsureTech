import type { Metadata } from "next";
import TeamManagementPage from "@/components/team/team-management";

export const metadata: Metadata = {
  title: "Team | InsureTech B2B Portal",
};

export default function TeamPage() {
  return <TeamManagementPage />;
}
