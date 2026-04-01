"use client";

import { BadgeCheck, LoaderCircle } from "lucide-react";
import Image from "next/image";
import Link from "next/link";

import { Panel } from "@/components/panel";
import { StatusPill } from "@/components/status-pill";
import { useCurrentInsurerId } from "@/hooks/use-current-insurer-id";
import { useDashboardWorkspace } from "@/hooks/use-dashboard-workspace";
import {
  dashboardMetricCards,
  dashboardTabCopy,
  downloadDashboardProductMatrix,
} from "@/lib/tabs/dashboard";
import { formatDateTime } from "@/lib/utils";

export function DashboardOverview() {
  const { insurerId } = useCurrentInsurerId();
  const { overview, loading, error, workspace } = useDashboardWorkspace(insurerId || undefined);

  if (loading) {
    return (
      <div className="flex min-h-[300px] items-center justify-center rounded-[32px] border border-[rgb(12_91_65_/_0.08)] bg-white/65">
        <LoaderCircle className="h-5 w-5 animate-spin text-[var(--brand-deep)]" />
      </div>
    );
  }

  if (error || !overview) {
    return (
      <Panel title={dashboardTabCopy.errorPanel.title} description={dashboardTabCopy.errorPanel.description}>
        <p className="text-sm text-[var(--danger)]">{error || dashboardTabCopy.errorPanel.emptyError}</p>
      </Panel>
    );
  }

  return (
    <div className="page-dashboard space-y-6">
      <section className="hero-grid">
        <Panel className="overflow-hidden">
          <div className="grid gap-6 lg:grid-cols-[minmax(0,1.1fr)_300px]">
            <div className="space-y-5">
              <div className="inline-flex w-fit items-center gap-2 rounded-full border border-[rgb(12_91_65_/_0.1)] bg-white/75 px-4 py-2 text-sm font-medium text-[var(--brand-deep)]">
                <BadgeCheck className="h-4 w-4" />
                {overview.source} {dashboardTabCopy.hero.feedSuffix}
              </div>

              <div>
                <p className="text-sm font-semibold uppercase tracking-[0.18em] text-[var(--brand-deep)]">
                  {workspace.heroCode} {dashboardTabCopy.hero.controlRoomSuffix}
                </p>
                <h2 className="mt-3 font-[family:var(--font-heading)] text-3xl font-semibold text-[var(--text)] sm:text-4xl">
                  {workspace.heroInsurerName}
                </h2>
                <p className="mt-3 max-w-2xl text-base leading-7 text-[var(--muted)]">
                  {dashboardTabCopy.hero.description}
                </p>
              </div>

              <div className="grid gap-3 sm:grid-cols-2">
                <div className="rounded-[26px] bg-[rgb(12_91_65_/_0.92)] p-5 text-white">
                  <p className="text-sm uppercase tracking-[0.14em] text-white/70">{dashboardTabCopy.hero.approvedTitle}</p>
                  <p className="mt-3 text-4xl font-semibold">{overview.metrics.approvedProposalCount}</p>
                  <p className="mt-2 text-sm text-white/72">{dashboardTabCopy.hero.approvedDescription}</p>
                </div>
                <div className="rounded-[26px] border border-[rgb(12_91_65_/_0.08)] bg-white/72 p-5">
                  <p className="text-sm uppercase tracking-[0.14em] text-[var(--muted)]">{dashboardTabCopy.hero.claimsTitle}</p>
                  <p className="mt-3 text-4xl font-semibold text-[var(--text)]">
                    {overview.metrics.requestedClaimCount + overview.metrics.underReviewClaimCount}
                  </p>
                  <p className="mt-2 text-sm text-[var(--muted)]">{dashboardTabCopy.hero.claimsDescription}</p>
                </div>
              </div>
            </div>

            <div className="overflow-hidden rounded-[30px] border border-[rgb(12_91_65_/_0.08)] bg-[rgb(255_252_247_/_0.86)] p-4">
              <div className="relative h-[220px] overflow-hidden rounded-[24px]">
                <Image
                  src="/images/travel-banner.jpg"
                  alt={dashboardTabCopy.hero.bannerAlt}
                  fill
                  className="object-cover"
                />
              </div>
              <div className="mt-4 rounded-[24px] bg-[rgb(245_158_11_/_0.12)] p-4">
                <p className="text-sm font-semibold uppercase tracking-[0.16em] text-[#8a5200]">{dashboardTabCopy.businessModel.eyebrow}</p>
                <p className="mt-2 font-[family:var(--font-heading)] text-2xl font-semibold text-[var(--text)]">
                  {dashboardTabCopy.businessModel.title}
                </p>
                <div className="mt-3 flex flex-wrap gap-2">
                  {workspace.businessModelPillars.map((pillar) => (
                    <span
                      key={pillar}
                      className="rounded-full border border-[rgb(245_158_11_/_0.18)] bg-white/88 px-3 py-1 text-xs font-semibold text-[#8a5200]"
                    >
                      {pillar}
                    </span>
                  ))}
                </div>
                <p className="mt-3 text-sm leading-6 text-[var(--muted)]">
                  {dashboardTabCopy.businessModel.description}
                </p>
              </div>
            </div>
          </div>
        </Panel>

        <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-1">
          {dashboardMetricCards.map((card) => {
            const Icon = card.icon;
            const value = overview.metrics[card.key];

            return (
              <Panel key={card.key} className="rounded-[28px] p-5">
                <div className="flex items-start justify-between gap-4">
                  <div>
                    <p className="text-sm font-medium text-[var(--muted)]">{card.label}</p>
                    <p className="mt-3 font-[family:var(--font-heading)] text-4xl font-semibold text-[var(--text)]">
                      {value}
                    </p>
                  </div>
                  <div className="flex h-12 w-12 items-center justify-center rounded-2xl bg-[rgb(15_157_104_/_0.12)] text-[var(--brand-deep)]">
                    <Icon className="h-5 w-5" />
                  </div>
                </div>
              </Panel>
            );
          })}
        </div>
      </section>

      <div className="grid gap-6 xl:grid-cols-[minmax(0,1.2fr)_minmax(360px,0.9fr)]">
        <Panel title={dashboardTabCopy.productLineup.title} description={dashboardTabCopy.productLineup.description}>
          <div className="grid gap-4">
            {workspace.highlightedProducts.map((product) => (
              <div
                key={product.id}
                className="rounded-[24px] border border-[rgb(12_91_65_/_0.08)] bg-white/72 p-4"
              >
                <div className="flex items-start justify-between gap-4">
                  <div>
                    <div className="flex items-center gap-2">
                      <h3 className="font-[family:var(--font-heading)] text-xl font-semibold text-[var(--text)]">
                        {product.name}
                      </h3>
                      <StatusPill status={product.status} />
                    </div>
                    <p className="mt-1 text-sm text-[var(--muted)]">
                      {product.category} product • {product.code}
                    </p>
                  </div>
                </div>
                <div className="mt-4 grid gap-3 sm:grid-cols-2">
                  <div className="rounded-[20px] bg-[rgb(12_91_65_/_0.05)] p-3">
                    <p className="text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">{dashboardTabCopy.productLineup.premiumLabel}</p>
                    <p className="mt-2 text-sm font-medium text-[var(--text)]">{product.premiumRangeText}</p>
                  </div>
                  <div className="rounded-[20px] bg-[rgb(245_158_11_/_0.1)] p-3">
                    <p className="text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">{dashboardTabCopy.productLineup.coverageLabel}</p>
                    <p className="mt-2 text-sm font-medium text-[var(--text)]">{product.coverageRangeText}</p>
                  </div>
                </div>
                <div className="mt-4 flex flex-wrap gap-2">
                  {product.features.map((feature) => (
                    <span
                      key={feature}
                      className="rounded-full border border-[rgb(12_91_65_/_0.08)] bg-white px-3 py-1 text-xs font-medium text-[var(--muted)]"
                    >
                      {feature}
                    </span>
                  ))}
                </div>
              </div>
            ))}
          </div>
        </Panel>

        <div className="space-y-6">
          <Panel title={dashboardTabCopy.recentProposals.title} description={dashboardTabCopy.recentProposals.description}>
            <div className="space-y-3">
              {workspace.recentProposals.map((proposal) => (
                <div
                  key={proposal.id}
                  className="flex items-start justify-between gap-4 rounded-[22px] border border-[rgb(12_91_65_/_0.08)] bg-white/70 p-4"
                >
                  <div>
                    <p className="font-medium text-[var(--text)]">{proposal.customerName}</p>
                    <p className="mt-1 text-sm text-[var(--muted)]">
                      {proposal.planName} • {proposal.proposalNumber}
                    </p>
                    <p className="mt-2 text-xs uppercase tracking-[0.14em] text-[var(--muted)]">
                      {dashboardTabCopy.recentProposals.submittedPrefix} {formatDateTime(proposal.submittedAt)}
                    </p>
                  </div>
                  <StatusPill status={proposal.status} />
                </div>
              ))}
            </div>
          </Panel>

          <Panel title={dashboardTabCopy.recentClaims.title} description={dashboardTabCopy.recentClaims.description}>
            <div className="space-y-3">
              {workspace.recentClaims.map((claim) => (
                <div
                  key={claim.id}
                  className="flex items-start justify-between gap-4 rounded-[22px] border border-[rgb(12_91_65_/_0.08)] bg-white/70 p-4"
                >
                  <div>
                    <p className="font-medium text-[var(--text)]">{claim.insuredName}</p>
                    <p className="mt-1 text-sm text-[var(--muted)]">
                      {claim.planName} • {claim.claimNumber}
                    </p>
                    <p className="mt-2 text-xs uppercase tracking-[0.14em] text-[var(--muted)]">
                      {dashboardTabCopy.recentClaims.submittedPrefix} {formatDateTime(claim.submittedAt)}
                    </p>
                  </div>
                  <StatusPill status={claim.status} />
                </div>
              ))}
            </div>
          </Panel>
        </div>
      </div>

      <div className="grid gap-6 xl:grid-cols-[minmax(0,1.2fr)_minmax(360px,0.9fr)]">
        <Panel
          title={dashboardTabCopy.playbooks.title}
          description={dashboardTabCopy.playbooks.description}
          action={
            <button
              className="portal-btn portal-btn-secondary"
              onClick={() => downloadDashboardProductMatrix(workspace.playbooks)}
              type="button"
            >
              {dashboardTabCopy.playbooks.actionLabel}
            </button>
          }
        >
          <div className="grid gap-4 lg:grid-cols-2">
            {workspace.playbooks.map((playbook) => (
              <div
                key={playbook.code}
                className="rounded-[24px] border border-[rgb(12_91_65_/_0.08)] bg-white/72 p-4"
              >
                <div className="flex items-center justify-between gap-3">
                  <div>
                    <p className="text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">
                      {playbook.category}
                    </p>
                    <h3 className="mt-2 font-[family:var(--font-heading)] text-2xl font-semibold text-[var(--text)]">
                      {playbook.name}
                    </h3>
                  </div>
                  <span className="rounded-full bg-[rgb(245_158_11_/_0.12)] px-3 py-1 text-xs font-semibold text-[#8a5200]">
                    {playbook.policyTerm}
                  </span>
                </div>
                <p className="mt-3 text-sm leading-6 text-[var(--muted)]">{playbook.summary}</p>
                <div className="mt-4 grid gap-3 sm:grid-cols-2">
                  <div className="rounded-[20px] bg-[rgb(12_91_65_/_0.05)] p-3">
                    <p className="text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">{dashboardTabCopy.playbooks.coverageLabel}</p>
                    <p className="mt-2 text-sm text-[var(--text)]">{playbook.coverageLimitText}</p>
                  </div>
                  <div className="rounded-[20px] bg-[rgb(245_158_11_/_0.1)] p-3">
                    <p className="text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">{dashboardTabCopy.playbooks.premiumLabel}</p>
                    <p className="mt-2 text-sm text-[var(--text)]">{playbook.premiumText}</p>
                  </div>
                </div>
                <div className="mt-4 flex flex-wrap gap-2">
                  {playbook.operationalFlags.map((flag) => (
                    <span
                      key={flag}
                      className="rounded-full border border-[rgb(12_91_65_/_0.08)] bg-white px-3 py-1 text-xs font-medium text-[var(--muted)]"
                    >
                      {flag}
                    </span>
                  ))}
                </div>
              </div>
            ))}
          </div>
        </Panel>

        <Panel title={dashboardTabCopy.documentPacks.title} description={dashboardTabCopy.documentPacks.description}>
          <div className="space-y-3">
            {workspace.playbooks.map((playbook) => (
              <div
                key={playbook.code}
                className="rounded-[22px] border border-[rgb(12_91_65_/_0.08)] bg-white/72 p-4"
              >
                <div className="flex items-center justify-between gap-3">
                  <p className="font-medium text-[var(--text)]">{playbook.name}</p>
                  <span className="text-sm text-[var(--muted)]">{playbook.requiredDocuments.length} {dashboardTabCopy.documentPacks.countSuffix}</span>
                </div>
                <div className="mt-3 flex flex-wrap gap-2">
                  {playbook.requiredDocuments.slice(0, 4).map((doc) => (
                    <span
                      key={doc}
                      className="rounded-full bg-[rgb(12_91_65_/_0.06)] px-3 py-1 text-xs font-medium text-[var(--muted)]"
                    >
                      {doc}
                    </span>
                  ))}
                </div>
              </div>
            ))}
          </div>
        </Panel>
      </div>

      <div className="grid gap-6 xl:grid-cols-[minmax(0,1.1fr)_minmax(340px,0.9fr)]">
        <Panel
          title={dashboardTabCopy.snapshot.title}
          description={dashboardTabCopy.snapshot.description}
          action={
            <Link className="portal-btn portal-btn-secondary" href="/tpa-claim-matrix">
              {dashboardTabCopy.snapshot.actionLabel}
            </Link>
          }
        >
          <div className="grid gap-4 md:grid-cols-2">
            <div className="rounded-[24px] bg-[rgb(12_91_65_/_0.92)] p-5 text-white">
              <p className="text-sm uppercase tracking-[0.16em] text-white/72">{dashboardTabCopy.snapshot.hotlineLabel}</p>
              <p className="mt-3 font-[family:var(--font-heading)] text-3xl font-semibold">
                {workspace.officialSnapshot.hotline}
              </p>
              <p className="mt-2 text-sm text-white/76">
                {dashboardTabCopy.snapshot.establishedPrefix} {workspace.officialSnapshot.established} • {workspace.officialSnapshot.website}
              </p>
            </div>
            <div className="rounded-[24px] border border-[rgb(12_91_65_/_0.08)] bg-white/72 p-5">
              <p className="text-sm uppercase tracking-[0.16em] text-[var(--muted)]">{dashboardTabCopy.snapshot.headOfficeLabel}</p>
              <p className="mt-3 text-sm leading-6 text-[var(--text)]">{workspace.officialSnapshot.headquarters}</p>
              <p className="mt-4 text-sm uppercase tracking-[0.16em] text-[var(--muted)]">{dashboardTabCopy.snapshot.productFamiliesLabel}</p>
              <div className="mt-4 flex flex-wrap gap-2">
                {workspace.officialSnapshot.products.map((item) => (
                  <span
                    key={item}
                    className="rounded-full bg-[rgb(245_158_11_/_0.12)] px-3 py-1 text-xs font-semibold text-[#8a5200]"
                  >
                    {item}
                  </span>
                ))}
              </div>
            </div>
          </div>
          <div className="mt-5 grid gap-3 md:grid-cols-3">
            {workspace.businessModelNotes.map((note) => (
              <div
                key={note}
                className="rounded-[20px] border border-[rgb(12_91_65_/_0.08)] bg-white/72 px-4 py-4 text-sm leading-6 text-[var(--muted)]"
              >
                {note}
              </div>
            ))}
          </div>
          <div className="mt-5 grid gap-3">
            {workspace.officialSnapshot.recentSignals.map((signal) => (
              <div
                key={signal}
                className="rounded-[20px] border border-[rgb(12_91_65_/_0.08)] bg-[rgb(255_252_247_/_0.82)] px-4 py-3 text-sm text-[var(--muted)]"
              >
                {signal}
              </div>
            ))}
          </div>
        </Panel>

        <Panel title={dashboardTabCopy.claimsControl.title} description={dashboardTabCopy.claimsControl.description}>
          <div className="space-y-3">
            {dashboardTabCopy.claimsControl.cards.map((card) => (
              <div key={card.title} className="rounded-[22px] border border-[rgb(12_91_65_/_0.08)] bg-white/72 p-4">
                <p className="font-medium text-[var(--text)]">{card.title}</p>
                <p className="mt-2 text-sm leading-6 text-[var(--muted)]">{card.body}</p>
              </div>
            ))}
          </div>
        </Panel>
      </div>
    </div>
  );
}
