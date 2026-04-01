"use client";

import { Download, HeartPulse, ShieldCheck, Stethoscope, Waypoints } from "lucide-react";
import { useMemo } from "react";

import { Panel } from "@/components/panel";
import {
  downloadTpaClaimMatrixCsv,
  getApprovalRangeLabel,
  getTouchpointIcon,
  getTpaClaimMatrixWorkspace,
  tpaMatrixTabCopy,
} from "@/lib/tabs/tpa-claim-matrix";

export function TpaClaimMatrix() {
  const workspace = useMemo(() => getTpaClaimMatrixWorkspace(), []);

  return (
    <div className="page-tpa-matrix space-y-6">
      <div className="grid gap-6 xl:grid-cols-[minmax(0,1.2fr)_380px]">
        <Panel title={tpaMatrixTabCopy.operatingModel.title} description={tpaMatrixTabCopy.operatingModel.description}>
          <div className="grid gap-4 lg:grid-cols-2">
            <div className="rounded-[24px] bg-[rgb(12_91_65_/_0.92)] p-5 text-white">
              <p className="text-sm uppercase tracking-[0.16em] text-white/72">{tpaMatrixTabCopy.operatingModel.modelLabel}</p>
              <p className="mt-3 font-[family:var(--font-heading)] text-2xl font-semibold">{workspace.snapshot.model}</p>
              <p className="mt-3 text-sm leading-6 text-white/78">{workspace.snapshot.network}</p>
            </div>
            <div className="rounded-[24px] border border-[rgb(12_91_65_/_0.08)] bg-white/72 p-5">
              <p className="text-sm uppercase tracking-[0.16em] text-[var(--muted)]">{tpaMatrixTabCopy.operatingModel.integrationLabel}</p>
              <p className="mt-3 text-sm leading-6 text-[var(--text)]">{workspace.snapshot.integration}</p>
              <div className="mt-4 flex flex-wrap gap-2">
                {workspace.snapshot.channels.map((channel) => (
                  <span
                    key={channel}
                    className="rounded-full bg-[rgb(245_158_11_/_0.12)] px-3 py-1 text-xs font-semibold text-[#8a5200]"
                  >
                    {channel}
                  </span>
                ))}
              </div>
            </div>
          </div>

          <div className="mt-5 grid gap-4 lg:grid-cols-3">
            <div className="rounded-[22px] border border-[rgb(12_91_65_/_0.08)] bg-white/72 p-4">
              <div className="flex items-center gap-3">
                <div className="flex h-10 w-10 items-center justify-center rounded-2xl bg-[rgb(15_157_104_/_0.12)] text-[var(--brand-deep)]">
                  <HeartPulse className="h-5 w-5" />
                </div>
                <p className="font-medium text-[var(--text)]">{tpaMatrixTabCopy.operatingModel.cards[0].title}</p>
              </div>
              <div className="mt-4 grid gap-2 text-sm text-[var(--muted)]">
                {workspace.snapshot.claimModes.map((item) => (
                  <div key={item} className="rounded-[18px] bg-[rgb(12_91_65_/_0.05)] px-3 py-2">
                    {item}
                  </div>
                ))}
              </div>
            </div>

            <div className="rounded-[22px] border border-[rgb(12_91_65_/_0.08)] bg-white/72 p-4">
              <div className="flex items-center gap-3">
                <div className="flex h-10 w-10 items-center justify-center rounded-2xl bg-[rgb(15_157_104_/_0.12)] text-[var(--brand-deep)]">
                  <Stethoscope className="h-5 w-5" />
                </div>
                <p className="font-medium text-[var(--text)]">{tpaMatrixTabCopy.operatingModel.cards[1].title}</p>
              </div>
              <div className="mt-4 grid gap-2 text-sm text-[var(--muted)]">
                {workspace.snapshot.fallbacks.map((item) => (
                  <div key={item} className="rounded-[18px] bg-[rgb(12_91_65_/_0.05)] px-3 py-2">
                    {item}
                  </div>
                ))}
              </div>
            </div>

            <div className="rounded-[22px] border border-[rgb(12_91_65_/_0.08)] bg-white/72 p-4">
              <div className="flex items-center gap-3">
                <div className="flex h-10 w-10 items-center justify-center rounded-2xl bg-[rgb(15_157_104_/_0.12)] text-[var(--brand-deep)]">
                  <ShieldCheck className="h-5 w-5" />
                </div>
                <p className="font-medium text-[var(--text)]">{tpaMatrixTabCopy.operatingModel.cards[2].title}</p>
              </div>
              <div className="mt-4 grid gap-2 text-sm text-[var(--muted)]">
                {workspace.snapshot.operatingRules.map((item) => (
                  <div key={item} className="rounded-[18px] bg-[rgb(12_91_65_/_0.05)] px-3 py-2">
                    {item}
                  </div>
                ))}
              </div>
            </div>
          </div>
        </Panel>

        <Panel title={tpaMatrixTabCopy.priorities.title} description={tpaMatrixTabCopy.priorities.description}>
          <div className="space-y-4">
            <div className="rounded-[24px] bg-[rgb(12_91_65_/_0.92)] p-5 text-white">
              <p className="text-sm uppercase tracking-[0.16em] text-white/72">{tpaMatrixTabCopy.priorities.surveyorLabel}</p>
              <p className="mt-3 font-[family:var(--font-heading)] text-2xl font-semibold">
                {workspace.surveyorLanes.length} categories require field review
              </p>
              <p className="mt-2 text-sm text-white/78">{tpaMatrixTabCopy.priorities.surveyorDescription}</p>
            </div>

            <div className="grid gap-3">
              {workspace.deskSignals.map((item) => (
                <div
                  key={item.title}
                  className="rounded-[20px] border border-[rgb(12_91_65_/_0.08)] bg-[rgb(255_252_247_/_0.82)] p-4 text-sm leading-6 text-[var(--muted)]"
                >
                  <p className="font-medium text-[var(--text)]">{item.title}</p>
                  <p className="mt-2 font-[family:var(--font-heading)] text-xl font-semibold text-[var(--text)]">
                    {item.value}
                  </p>
                  <p className="mt-2">{item.detail}</p>
                </div>
              ))}

              <div className="rounded-[20px] border border-[rgb(12_91_65_/_0.08)] bg-white/72 p-4">
                <div className="flex items-center gap-3">
                  <Waypoints className="h-5 w-5 text-[var(--brand-deep)]" />
                  <p className="font-medium text-[var(--text)]">{tpaMatrixTabCopy.priorities.ownershipTitle}</p>
                </div>
                <div className="mt-3 grid gap-2">
                  {workspace.claimMatrix.map((row) => (
                    <div
                      key={row.category}
                      className="rounded-[16px] bg-[rgb(12_91_65_/_0.05)] px-3 py-2"
                    >
                      <p className="text-sm font-medium text-[var(--text)]">{row.category}</p>
                      <p className="mt-1 text-xs text-[var(--muted)]">{row.escalationOwner}</p>
                    </div>
                  ))}
                </div>
              </div>
            </div>
          </div>
        </Panel>
      </div>

      <Panel
        title={tpaMatrixTabCopy.matrixPanel.title}
        description={tpaMatrixTabCopy.matrixPanel.description}
        action={
          <button className="portal-btn portal-btn-secondary" onClick={downloadTpaClaimMatrixCsv} type="button">
            <Download className="h-4 w-4" />
            {tpaMatrixTabCopy.matrixPanel.actionLabel}
          </button>
        }
      >
        <div className="table-wrap">
          <table className="table-base">
            <thead>
              <tr>
                {tpaMatrixTabCopy.matrixPanel.headers.map((header) => (
                  <th key={header}>{header}</th>
                ))}
              </tr>
            </thead>
            <tbody>
              {workspace.claimMatrix.map((row) => (
                <tr key={row.category}>
                  <td>
                    <p className="font-medium text-[var(--text)]">{row.category}</p>
                  </td>
                  <td className="text-sm text-[var(--muted)]">{row.planType}</td>
                  <td className="text-sm text-[var(--muted)]">{row.intakeGate}</td>
                  <td className="text-sm text-[var(--muted)]">{row.claimMode}</td>
                  <td className="text-sm text-[var(--muted)]">{row.settlementRail}</td>
                  <td className="text-sm text-[var(--muted)]">{row.typicalTat}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </Panel>

      <div className="grid gap-6 xl:grid-cols-[minmax(0,1fr)_minmax(0,1fr)]">
        <Panel title={tpaMatrixTabCopy.controlsPanel.title} description={tpaMatrixTabCopy.controlsPanel.description}>
          <div className="grid gap-4">
            {workspace.claimMatrix.map((row) => (
              <div
                key={row.category}
                className="rounded-[24px] border border-[rgb(12_91_65_/_0.08)] bg-white/72 p-4"
              >
                <p className="font-[family:var(--font-heading)] text-xl font-semibold text-[var(--text)]">{row.category}</p>
                <div className="mt-4 grid gap-4 lg:grid-cols-2">
                  <div>
                    <p className="text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">{tpaMatrixTabCopy.controlsPanel.documentsLabel}</p>
                    <div className="mt-2 flex flex-wrap gap-2">
                      {row.primaryDocuments.map((item) => (
                        <span
                          key={item}
                          className="rounded-full border border-[rgb(12_91_65_/_0.08)] bg-white px-3 py-1 text-xs font-medium text-[var(--muted)]"
                        >
                          {item}
                        </span>
                      ))}
                    </div>
                  </div>
                  <div>
                    <p className="text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">{tpaMatrixTabCopy.controlsPanel.settlementLabel}</p>
                    <p className="mt-2 text-sm leading-6 text-[var(--text)]">{row.settlementRail}</p>
                    <p className="mt-3 text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">
                      {tpaMatrixTabCopy.controlsPanel.escalationLabel}
                    </p>
                    <p className="mt-2 text-sm leading-6 text-[var(--text)]">{row.escalationOwner}</p>
                  </div>
                  <div>
                    <p className="text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">{tpaMatrixTabCopy.controlsPanel.fraudLabel}</p>
                    <div className="mt-2 flex flex-wrap gap-2">
                      {row.fraudChecks.map((item) => (
                        <span
                          key={item}
                          className="rounded-full bg-[var(--danger-soft)] px-3 py-1 text-xs font-semibold text-[var(--danger)]"
                        >
                          {item}
                        </span>
                      ))}
                    </div>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </Panel>

        <div className="space-y-6">
          <Panel title={tpaMatrixTabCopy.approvalsPanel.title} description={tpaMatrixTabCopy.approvalsPanel.description}>
            <div className="space-y-3">
              {workspace.approvals.map((row) => (
                <div
                  key={row.approvalLevel}
                  className="rounded-[22px] border border-[rgb(12_91_65_/_0.08)] bg-white/72 p-4"
                >
                  <div className="flex items-center justify-between gap-3">
                    <div>
                      <p className="font-medium text-[var(--text)]">
                        {getApprovalRangeLabel(row.min, row.max)}
                      </p>
                      <p className="mt-1 text-sm text-[var(--muted)]">{row.approvalLevel}</p>
                    </div>
                    <div className="rounded-full bg-[rgb(245_158_11_/_0.12)] px-3 py-1 text-xs font-semibold text-[#8a5200]">
                      {row.maxTat}
                    </div>
                  </div>
                  <p className="mt-3 text-sm text-[var(--muted)]">{row.approvers}</p>
                  <p className="mt-1 text-xs uppercase tracking-[0.14em] text-[var(--muted)]">{row.mode}</p>
                </div>
              ))}
            </div>
          </Panel>

          <Panel title={tpaMatrixTabCopy.touchpointsPanel.title} description={tpaMatrixTabCopy.touchpointsPanel.description}>
            <div className="grid gap-3">
              {workspace.claimMatrix.map((row) => (
                <div
                  key={row.category}
                  className="rounded-[22px] border border-[rgb(12_91_65_/_0.08)] bg-white/72 p-4"
                >
                  <div className="flex items-center gap-3">
                    <div className="flex h-10 w-10 items-center justify-center rounded-2xl bg-[rgb(15_157_104_/_0.12)] text-[var(--brand-deep)]">
                      {(() => {
                        const Icon = getTouchpointIcon(row.category);
                        return <Icon className="h-5 w-5" />;
                      })()}
                    </div>
                    <div>
                      <p className="font-medium text-[var(--text)]">{row.category}</p>
                      <p className="text-sm text-[var(--muted)]">{row.planType}</p>
                    </div>
                  </div>
                  <div className="mt-3 flex flex-wrap gap-2">
                    {row.partnerTouchpoints.map((item) => (
                      <span
                        key={item}
                        className="rounded-full border border-[rgb(12_91_65_/_0.08)] bg-[rgb(12_91_65_/_0.05)] px-3 py-1 text-xs font-medium text-[var(--muted)]"
                      >
                        {item}
                      </span>
                    ))}
                  </div>
                  <div className="mt-3 rounded-[16px] bg-[rgb(245_158_11_/_0.08)] px-3 py-3 text-sm leading-6 text-[var(--muted)]">
                    <span className="font-semibold text-[var(--text)]">{tpaMatrixTabCopy.touchpointsPanel.intakeGatePrefix}</span> {row.intakeGate}
                  </div>
                </div>
              ))}
            </div>
          </Panel>
        </div>
      </div>
    </div>
  );
}
