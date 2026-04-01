import { FileText, HandCoins, Package2, ShieldAlert } from "lucide-react";

import {
  pragatiBusinessModelNotes,
  pragatiBusinessModelPillars,
  pragatiOfficialSnapshot,
} from "@/lib/claims-intelligence";
import type { PortalOverview } from "@/lib/types";

export const dashboardMetricCards = [
  { key: "productCount", label: "Active products", icon: Package2 },
  { key: "proposalCount", label: "Proposal queue", icon: FileText },
  { key: "claimCount", label: "Claims in motion", icon: ShieldAlert },
  { key: "settledClaimCount", label: "Settled claims", icon: HandCoins },
] as const;

export const dashboardTabCopy = {
  errorPanel: {
    title: "Overview unavailable",
    description: "The portal could not load this insurer workspace.",
    emptyError: "No overview data was returned.",
  },
  hero: {
    feedSuffix: "workspace feed",
    controlRoomSuffix: "control room",
    bannerAlt: "Pragati dashboard banner",
    fallbackCode: "INS",
    fallbackInsurer: "Selected insurer",
    description:
      "Track proposal movement, claim activity, and insurer product footprint from one workspace that mirrors the platform auth and SDK flow.",
    approvedTitle: "Approved proposals",
    approvedDescription: "Decisions already moved beyond intake.",
    claimsTitle: "Claims needing attention",
    claimsDescription: "Pending documents plus under-review cases.",
  },
  businessModel: {
    eyebrow: "Business model",
    title: "Multi-channel non-life growth model",
    description:
      "Keep product setup, partner routing, and claims operations aligned to how Pragati distributes and services business across retail, bank, and corporate channels.",
  },
  productLineup: {
    title: "Product lineup",
    description: "Top insurer products currently visible in this workspace.",
    premiumLabel: "Premium range",
    coverageLabel: "Coverage range",
  },
  recentProposals: {
    title: "Recent proposals",
    description: "Fresh submissions moving through underwriting.",
    submittedPrefix: "Submitted",
  },
  recentClaims: {
    title: "Recent claims",
    description: "A quick view of settlement activity and backlog.",
    submittedPrefix: "Filed",
  },
  playbooks: {
    title: "Pragati non-life playbooks",
    description: "Doc-backed insurer product intelligence from the KBank non-life portfolio.",
    actionLabel: "Download product matrix",
    exportFileName: "pragati-product-matrix.csv",
    coverageLabel: "Coverage",
    premiumLabel: "Premium",
  },
  documentPacks: {
    title: "Required document packs",
    description: "A quick view of the claim and review documents expected by product.",
    countSuffix: "docs",
  },
  snapshot: {
    title: "Official Pragati insurer snapshot",
    description: "Core public insurer profile and business-positioning signals placed on the main dashboard.",
    actionLabel: "Open TPA matrix",
    hotlineLabel: "Hotline",
    establishedPrefix: "Established",
    headOfficeLabel: "Head office",
    productFamiliesLabel: "Official product families",
  },
  claimsControl: {
    title: "Claims control focus",
    description: "Use the SRS-aligned claim lanes to keep reviewers on the right path.",
    cards: [
      {
        title: "Cashless hospital lane",
        body:
          "Health claims should prefer the TPA path with LabAid hospital integration, provider validation, and manual fallback when EHR connectivity fails.",
      },
      {
        title: "Amount-based routing",
        body:
          "Small claims can flow to auto or officer review, while BDT 50K-2L requires joint approval and BDT 2L+ needs board plus insurer sign-off.",
      },
      {
        title: "Fraud checkpoints",
        body:
          "Watch for rapid policy-to-claim movement, frequent repeat types, non-network providers, and claims sitting right at the coverage limit.",
      },
    ],
  },
} as const;

interface DashboardPlaybookLike {
  code: string;
  name: string;
  category: string;
  audience: string;
  coverageLimitText: string;
  premiumText: string;
  policyTerm: string;
  ageRange: string;
  operationalFlags: string[];
  requiredDocuments: string[];
  summary: string;
}

const dashboardProductMatrixHeaders = [
  "product_name",
  "category",
  "audience",
  "coverage_limit",
  "premium",
  "policy_term",
  "age_range",
  "operational_flags",
] as const;

function buildCsvDownload(fileName: string, rows: string[][]) {
  const csv = `\uFEFF${rows
    .map((row) => row.map((cell) => `"${cell.replace(/"/g, '""')}"`).join(","))
    .join("\r\n")}`;
  const blob = new Blob([csv], { type: "text/csv;charset=utf-8" });
  const url = URL.createObjectURL(blob);
  const link = document.createElement("a");
  link.href = url;
  link.download = fileName;
  link.click();
  URL.revokeObjectURL(url);
}

export function buildDashboardProductMatrixRows(playbooks: DashboardPlaybookLike[]) {
  return [
    [...dashboardProductMatrixHeaders],
    ...playbooks.map((item) => [
      item.name,
      item.category,
      item.audience,
      item.coverageLimitText,
      item.premiumText,
      item.policyTerm,
      item.ageRange,
      item.operationalFlags.join(" | "),
    ]),
  ];
}

export function downloadDashboardProductMatrix(playbooks: DashboardPlaybookLike[]) {
  buildCsvDownload(dashboardTabCopy.playbooks.exportFileName, buildDashboardProductMatrixRows(playbooks));
}

export function getDashboardWorkspace(overview: PortalOverview | null, playbooks: DashboardPlaybookLike[]) {
  return {
    overview,
    highlightedProducts: overview?.products.slice(0, 3) ?? [],
    recentProposals: overview?.proposals.slice(0, 5) ?? [],
    recentClaims: overview?.claims.slice(0, 5) ?? [],
    playbooks,
    businessModelPillars: pragatiBusinessModelPillars,
    businessModelNotes: pragatiBusinessModelNotes,
    officialSnapshot: pragatiOfficialSnapshot,
    heroCode: overview?.currentInsurer?.code ?? dashboardTabCopy.hero.fallbackCode,
    heroInsurerName: overview?.currentInsurer?.name ?? dashboardTabCopy.hero.fallbackInsurer,
  };
}
