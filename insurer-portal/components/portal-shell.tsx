"use client";

import {
  LoaderCircle,
  LogOut,
  Menu,
  Shield,
  X,
} from "lucide-react";
import Image from "next/image";
import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import { useEffect, useMemo, useState } from "react";

import { api } from "@/lib/browser-client";
import { setStoredCurrentInsurerId } from "@/lib/current-insurer";
import { useInsurerContext } from "@/hooks/use-insurer-context";
import { useCurrentInsurerId } from "@/hooks/use-current-insurer-id";
import { getPortalShellHeading, portalShellCopy, portalShellNavItems } from "@/lib/tabs/portal-shell";
import { cn, initialLetters } from "@/lib/utils";

export function PortalShell({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const router = useRouter();
  const [navOpen, setNavOpen] = useState(false);
  const [loggingOut, setLoggingOut] = useState(false);
  const { insurerId, setInsurerId } = useCurrentInsurerId();
  const { context, loading } = useInsurerContext(insurerId || undefined);

  const currentHeading = useMemo(() => getPortalShellHeading(pathname), [pathname]);

  useEffect(() => {
    if (!insurerId && context.currentInsurer?.id) {
      setStoredCurrentInsurerId(context.currentInsurer.id);
      setInsurerId(context.currentInsurer.id);
    }
  }, [context.currentInsurer?.id, insurerId, setInsurerId]);

  async function handleLogout() {
    setLoggingOut(true);

    try {
      await api.auth.logout();
    } finally {
      setStoredCurrentInsurerId("");
      router.replace("/login");
      router.refresh();
      setLoggingOut(false);
    }
  }

  function handleInsurerSelect(nextInsurerId: string) {
    setStoredCurrentInsurerId(nextInsurerId);
    setInsurerId(nextInsurerId);
    router.refresh();
  }

  return (
    <div className="portal-shell bg-transparent lg:grid lg:grid-cols-[290px_minmax(0,1fr)]">
      <aside
        className={cn(
          "portal-sidebar fixed inset-y-0 left-0 z-40 w-[290px] px-5 py-5 transition-transform duration-200 lg:static lg:translate-x-0",
          navOpen ? "translate-x-0" : "-translate-x-full lg:translate-x-0",
        )}
      >
        <div className="flex h-full flex-col">
          <div className="mb-6 flex items-center justify-between">
            <div className="flex items-center gap-3">
              <div className="flex h-12 w-12 items-center justify-center overflow-hidden rounded-2xl bg-white shadow-[0_16px_30px_rgb(12_91_65_/_0.16)]">
                <Image src="/logos/favicon.svg" alt="LabaidInsuretech" width={34} height={34} />
              </div>
              <div>
                <p className="text-sm font-medium uppercase tracking-[0.18em] text-[var(--brand-deep)]">{portalShellCopy.brand.name}</p>
                <p className="font-[family:var(--font-heading)] text-lg font-semibold text-[var(--text)]">
                  {portalShellCopy.brand.title}
                </p>
              </div>
            </div>
            <button
              className="rounded-full p-2 text-[var(--muted)] lg:hidden"
              onClick={() => setNavOpen(false)}
              type="button"
            >
              <X className="h-5 w-5" />
            </button>
          </div>

          <div className="portal-panel mb-6 rounded-[26px] p-4">
            <p className="text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">{portalShellCopy.sidebar.currentInsurerLabel}</p>
            <p className="mt-2 font-[family:var(--font-heading)] text-lg font-semibold text-[var(--text)]">
              {context.currentInsurer?.name ?? portalShellCopy.sidebar.loadingInsurer}
            </p>
            <p className="mt-1 text-sm text-[var(--muted)]">
              {context.currentInsurer?.businessModel ?? portalShellCopy.sidebar.fallbackBusinessModel}
            </p>

            <label className="mt-4 block space-y-2">
              <span className="text-sm font-medium text-[var(--muted)]">{portalShellCopy.sidebar.workspaceLabel}</span>
              <select
                className="portal-select"
                value={context.currentInsurer?.id ?? insurerId}
                onChange={(event) => handleInsurerSelect(event.target.value)}
              >
                {context.insurers.map((item) => (
                  <option key={item.id} value={item.id}>
                    {item.name}
                  </option>
                ))}
              </select>
            </label>

            <div className="mt-4 flex items-center gap-2 text-sm text-[var(--muted)]">
              <span className="inline-flex h-2.5 w-2.5 rounded-full bg-[var(--brand)]" />
              {portalShellCopy.sidebar.dataModePrefix} {context.source}
            </div>
          </div>

          <nav className="space-y-2">
            {portalShellNavItems.map((item) => {
              const Icon = item.icon;
              const active = item.href === "/" ? pathname === "/" : pathname.startsWith(item.href);

              return (
                <Link
                  key={item.href}
                  className={cn("portal-nav-link", active && "portal-nav-link-active")}
                  href={item.href}
                  onClick={() => setNavOpen(false)}
                >
                  <span className="portal-nav-icon">
                    <Icon className="h-4 w-4" />
                  </span>
                  <span className="font-medium">{item.label}</span>
                </Link>
              );
            })}
          </nav>

          <div className="mt-auto portal-panel rounded-[26px] p-4">
            <p className="text-xs font-semibold uppercase tracking-[0.16em] text-[var(--muted)]">{portalShellCopy.sidebar.signedInLabel}</p>
            <div className="mt-3 flex items-center gap-3">
              <div className="flex h-11 w-11 items-center justify-center rounded-full bg-[rgb(12_91_65_/_0.1)] font-semibold text-[var(--brand-deep)]">
                {initialLetters(context.session?.email || context.currentInsurer?.name)}
              </div>
              <div className="min-w-0">
                <p className="truncate font-medium text-[var(--text)]">{context.session?.email ?? portalShellCopy.sidebar.fallbackEmail}</p>
                <p className="truncate text-sm text-[var(--muted)]">{context.session?.role ?? portalShellCopy.sidebar.fallbackRole}</p>
              </div>
            </div>

            <button
              className="portal-btn portal-btn-secondary mt-4 w-full"
              disabled={loggingOut}
              onClick={handleLogout}
              type="button"
            >
              {loggingOut ? <LoaderCircle className="h-4 w-4 animate-spin" /> : <LogOut className="h-4 w-4" />}
              {loggingOut ? portalShellCopy.sidebar.signOutBusy : portalShellCopy.sidebar.signOutIdle}
            </button>
          </div>
        </div>
      </aside>

      {navOpen ? (
        <button
          aria-label={portalShellCopy.topbar.closeNavAriaLabel}
          className="fixed inset-0 z-30 bg-[rgb(28_42_31_/_0.26)] lg:hidden"
          onClick={() => setNavOpen(false)}
          type="button"
        />
      ) : null}

      <div className="min-w-0">
        <header className="portal-topbar px-4 py-4 sm:px-6 lg:px-8">
          <div className="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
            <div className="flex items-center gap-3">
              <button
                className="rounded-full border border-[rgb(12_91_65_/_0.12)] bg-white/80 p-2 text-[var(--text)] lg:hidden"
                onClick={() => setNavOpen(true)}
                type="button"
              >
                <Menu className="h-5 w-5" />
              </button>
              <div>
                <p className="text-sm font-medium uppercase tracking-[0.18em] text-[var(--brand-deep)]">
                  {portalShellCopy.brand.name}
                </p>
                <h1 className="font-[family:var(--font-heading)] text-2xl font-semibold text-[var(--text)]">
                  {currentHeading}
                </h1>
              </div>
            </div>

            <div className="flex items-center gap-3 rounded-full border border-[rgb(12_91_65_/_0.1)] bg-white/75 px-4 py-2 text-sm text-[var(--muted)]">
              {loading ? <LoaderCircle className="h-4 w-4 animate-spin" /> : <Shield className="h-4 w-4" />}
              <span>{context.currentInsurer?.code ?? "Insurer"} {portalShellCopy.topbar.syncedSuffix}</span>
            </div>
          </div>
        </header>

        <main className="px-4 py-6 sm:px-6 lg:px-8 lg:py-8">{children}</main>
      </div>
    </div>
  );
}
