import type { LucideIcon } from "lucide-react";
import {
  ActivitySquare,
  BookOpenText,
  BriefcaseBusiness,
  ClipboardCheck,
  ClipboardList,
  Files,
  FolderOpen,
  LayoutDashboard,
  Layers3,
  PhoneCall,
  Settings,
  Shield,
  Users,
} from "lucide-react";

export interface PortalTabNavItem {
  href: string;
  label: string;
  icon: LucideIcon;
}

export const portalShellNavItems: PortalTabNavItem[] = [
  { href: "/", label: "Dashboard", icon: LayoutDashboard },
  { href: "/proposals", label: "Proposals", icon: Files },
  { href: "/documents", label: "Documents", icon: FolderOpen },
  { href: "/enrollment-census", label: "Enrollment & Census", icon: Users },
  { href: "/pricing-commercials", label: "Pricing & Commercials", icon: BriefcaseBusiness },
  { href: "/claim-settlement", label: "Claim Settlement", icon: ClipboardCheck },
  { href: "/claim-checklists", label: "Claims Checklists", icon: ClipboardList },
  { href: "/surveyor-desk", label: "Surveyor Desk", icon: ActivitySquare },
  { href: "/travel-assistance", label: "Travel Assistance", icon: PhoneCall },
  { href: "/tpa-claim-matrix", label: "TPA & Claim Matrix", icon: ActivitySquare },
  { href: "/knowledge-center", label: "Knowledge Center", icon: BookOpenText },
  { href: "/policy-categories", label: "Policy Categories", icon: Layers3 },
  { href: "/plan-templates", label: "Plan Templates", icon: Shield },
  { href: "/settings", label: "Settings", icon: Settings },
];

export function getPortalShellHeading(pathname: string) {
  const matched = portalShellNavItems.find((item) => (item.href === "/" ? pathname === "/" : pathname.startsWith(item.href)));
  return matched?.label ?? "Insurer Portal";
}

export const portalShellCopy = {
  brand: {
    name: "LabaidInsuretech",
    title: "Insurer Dashboard",
  },
  sidebar: {
    currentInsurerLabel: "Current insurer",
    loadingInsurer: "Loading insurer...",
    fallbackBusinessModel: "Partner-facing operations workspace",
    workspaceLabel: "Workspace",
    dataModePrefix: "Data mode:",
    signedInLabel: "Signed in as",
    fallbackEmail: "portal@insuretech",
    fallbackRole: "Partner User",
    signOutIdle: "Sign out",
    signOutBusy: "Signing out...",
  },
  topbar: {
    syncedSuffix: "workspace synced",
    closeNavAriaLabel: "Close navigation",
  },
} as const;
