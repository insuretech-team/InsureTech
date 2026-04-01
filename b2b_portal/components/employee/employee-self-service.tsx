"use client";

import { useEffect, useState } from "react";

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { bffClient } from "@lib/sdk/b2b-sdk-client";

type ViewState = {
  loading: boolean;
  error: string;
  profile: Record<string, unknown> | null;
  coverage: Record<string, unknown> | null;
};

function formatMoney(value: unknown): string {
  if (!value || typeof value !== "object") {
    return "Not assigned";
  }
  const amount = Number((value as Record<string, unknown>).decimal_amount ?? 0);
  if (!Number.isFinite(amount) || amount <= 0) {
    return "Not assigned";
  }
  return `BDT ${Math.round(amount).toLocaleString("en-US")}`;
}

function prettifyInsuranceCategory(value: unknown): string {
  if (typeof value !== "string" || !value) {
    return "Not assigned";
  }
  return value
    .replace(/^INSURANCE_TYPE_/, "")
    .toLowerCase()
    .replace(/\b\w/g, (letter) => letter.toUpperCase());
}

function fieldValue(value: unknown, fallback = "Not available"): string {
  return typeof value === "string" && value.trim() ? value : fallback;
}

export default function EmployeeSelfService() {
  const [state, setState] = useState<ViewState>({
    loading: true,
    error: "",
    profile: null,
    coverage: null,
  });

  useEffect(() => {
    let cancelled = false;

    Promise.all([bffClient.employeeSelf.getProfile(), bffClient.employeeSelf.getCoverage()])
      .then(([profileResult, coverageResult]) => {
        if (cancelled) {
          return;
        }
        if (!profileResult.ok) {
          setState({
            loading: false,
            error: profileResult.message ?? "Unable to load your employee profile.",
            profile: null,
            coverage: null,
          });
          return;
        }
        if (!coverageResult.ok) {
          setState({
            loading: false,
            error: coverageResult.message ?? "Unable to load your coverage.",
            profile: profileResult.profile ?? null,
            coverage: null,
          });
          return;
        }
        setState({
          loading: false,
          error: "",
          profile: profileResult.profile ?? null,
          coverage: coverageResult.coverage ?? null,
        });
      })
      .catch((error) => {
        if (!cancelled) {
          setState({
            loading: false,
            error: error instanceof Error ? error.message : "Unable to load your employee access.",
            profile: null,
            coverage: null,
          });
        }
      });

    return () => {
      cancelled = true;
    };
  }, []);

  if (state.loading) {
    return (
      <div className="grid gap-6 lg:grid-cols-2">
        <Card><CardHeader><CardTitle>Loading profile...</CardTitle></CardHeader></Card>
        <Card><CardHeader><CardTitle>Loading coverage...</CardTitle></CardHeader></Card>
      </div>
    );
  }

  if (state.error) {
    return (
      <Card className="border-destructive/30">
        <CardHeader>
          <CardTitle>Employee Access</CardTitle>
          <CardDescription>{state.error}</CardDescription>
        </CardHeader>
      </Card>
    );
  }

  const profile = state.profile ?? {};
  const coverage = state.coverage ?? {};

  return (
    <div className="space-y-6">
      <div className="space-y-2">
        <p className="text-sm font-semibold uppercase tracking-[0.16em] text-[rgb(var(--brand-cold-rgb))]">
          Employee Self Service
        </p>
        <h1 className="text-3xl font-semibold tracking-tight text-foreground">
          Your profile and coverage
        </h1>
        <p className="max-w-2xl text-sm text-muted-foreground">
          Review your employee record and the insurance coverage provided by your organisation.
        </p>
      </div>

      <div className="grid gap-6 lg:grid-cols-2">
        <Card className="border-border/70 bg-card/90">
          <CardHeader>
            <CardTitle>Employee Profile</CardTitle>
            <CardDescription>Your organisation-linked identity details.</CardDescription>
          </CardHeader>
          <CardContent className="grid gap-4 sm:grid-cols-2">
            <div>
              <p className="text-xs uppercase tracking-wide text-muted-foreground">Name</p>
              <p className="mt-1 text-sm font-medium text-foreground">{fieldValue(profile.name)}</p>
            </div>
            <div>
              <p className="text-xs uppercase tracking-wide text-muted-foreground">Employee ID</p>
              <p className="mt-1 text-sm font-medium text-foreground">{fieldValue(profile.employee_id)}</p>
            </div>
            <div>
              <p className="text-xs uppercase tracking-wide text-muted-foreground">Email</p>
              <p className="mt-1 text-sm font-medium text-foreground">{fieldValue(profile.email)}</p>
            </div>
            <div>
              <p className="text-xs uppercase tracking-wide text-muted-foreground">Mobile</p>
              <p className="mt-1 text-sm font-medium text-foreground">{fieldValue(profile.mobile_number)}</p>
            </div>
            <div>
              <p className="text-xs uppercase tracking-wide text-muted-foreground">Department</p>
              <p className="mt-1 text-sm font-medium text-foreground">
                {fieldValue(profile.department_name ?? profile.department_id)}
              </p>
            </div>
            <div>
              <p className="text-xs uppercase tracking-wide text-muted-foreground">Joining Date</p>
              <p className="mt-1 text-sm font-medium text-foreground">{fieldValue(profile.date_of_joining)}</p>
            </div>
          </CardContent>
        </Card>

        <Card className="border-border/70 bg-card/90">
          <CardHeader>
            <CardTitle>Insurance Coverage</CardTitle>
            <CardDescription>Your currently assigned plan and limits.</CardDescription>
          </CardHeader>
          <CardContent className="grid gap-4 sm:grid-cols-2">
            <div>
              <p className="text-xs uppercase tracking-wide text-muted-foreground">Organisation</p>
              <p className="mt-1 text-sm font-medium text-foreground">{fieldValue(coverage.organisation_name)}</p>
            </div>
            <div>
              <p className="text-xs uppercase tracking-wide text-muted-foreground">Plan</p>
              <p className="mt-1 text-sm font-medium text-foreground">{fieldValue(coverage.assigned_plan_name)}</p>
            </div>
            <div>
              <p className="text-xs uppercase tracking-wide text-muted-foreground">Insurance Type</p>
              <p className="mt-1 text-sm font-medium text-foreground">
                {prettifyInsuranceCategory(coverage.insurance_category)}
              </p>
            </div>
            <div>
              <p className="text-xs uppercase tracking-wide text-muted-foreground">Coverage Amount</p>
              <p className="mt-1 text-sm font-medium text-foreground">{formatMoney(coverage.coverage_amount)}</p>
            </div>
            <div>
              <p className="text-xs uppercase tracking-wide text-muted-foreground">Premium</p>
              <p className="mt-1 text-sm font-medium text-foreground">{formatMoney(coverage.premium_amount)}</p>
            </div>
            <div>
              <p className="text-xs uppercase tracking-wide text-muted-foreground">Dependents</p>
              <p className="mt-1 text-sm font-medium text-foreground">
                {typeof coverage.number_of_dependent === "number" ? coverage.number_of_dependent : 0}
              </p>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
