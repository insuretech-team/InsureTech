import { BadgeCheck } from "lucide-react";
import Image from "next/image";

import { LoginForm } from "@/components/login-form";

export default function LoginPage() {
  return (
    <main className="min-h-screen px-4 py-6 sm:px-6 lg:px-8">
      <div className="mx-auto grid min-h-[calc(100vh-3rem)] max-w-7xl items-stretch gap-6 lg:grid-cols-[minmax(0,1.2fr)_460px]">
        <section className="relative overflow-hidden rounded-[36px] border border-[rgb(12_91_65_/_0.1)] bg-[linear-gradient(140deg,rgb(255_251_244_/_0.92),rgb(236_249_242_/_0.9))] p-6 shadow-[0_28px_70px_rgb(24_38_29_/_0.14)] sm:p-8 lg:p-10">
          <div className="hero-grid items-start">
            <div className="space-y-6">
              <div className="inline-flex items-center gap-2 rounded-full border border-[rgb(12_91_65_/_0.12)] bg-white/75 px-4 py-2 text-sm font-medium text-[var(--brand-deep)]">
                <BadgeCheck className="h-4 w-4" />
                Gateway-backed insurer operations
              </div>

              <div className="space-y-4">
                <div className="flex items-center gap-3">
                  <div className="rounded-[20px] bg-white/88 px-3 py-2 shadow-[0_10px_24px_rgb(12_91_65_/_0.08)]">
                    <Image src="/logos/insuretech.svg" alt="LabaidInsuretech" width={180} height={40} priority />
                  </div>
                </div>
                <p className="text-sm font-semibold uppercase tracking-[0.24em] text-[var(--brand-deep)]">Insurer operations desk</p>
                <h1 className="max-w-3xl font-[family:var(--font-heading)] text-4xl font-semibold leading-tight text-[var(--text)] sm:text-5xl">
                  Insurer Dashboard
                </h1>
                <p className="max-w-2xl text-base leading-7 text-[var(--muted)] sm:text-lg">
                  LabaidInsuretech insurer operations workspace for proposals, claims, product control, and partner-facing execution.
                </p>
              </div>

            </div>

            <div className="relative overflow-hidden rounded-[30px] border border-[rgb(12_91_65_/_0.08)] bg-[rgb(255_252_247_/_0.82)] p-4">
              <div className="relative h-[260px] overflow-hidden rounded-[24px]">
                <Image
                  src="/images/travel-banner.jpg"
                  alt="Insurer portal showcase"
                  fill
                  className="object-cover"
                  priority
                />
              </div>
              <div className="mt-4 rounded-[24px] bg-[rgb(12_91_65_/_0.92)] p-5 text-white">
                <p className="text-sm font-semibold uppercase tracking-[0.16em] text-white/72">Portal focus</p>
                <p className="mt-3 font-[family:var(--font-heading)] text-2xl font-semibold">
                  One place for insurer operations.
                </p>
                <p className="mt-2 text-sm leading-6 text-white/76">
                  Review proposals, monitor claims, and manage insurer setup without leaving the dashboard.
                </p>
              </div>
            </div>
          </div>
        </section>

        <section className="flex items-center justify-center">
          <div className="w-full max-w-[460px]">
            <LoginForm />
          </div>
        </section>
      </div>
    </main>
  );
}
