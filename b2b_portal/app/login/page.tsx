import Image from "next/image";
import { redirect } from "next/navigation";

import LoginForm from "@/components/auth/login-form";
import { getServerSession } from "@lib/auth/session";

export default async function LoginPage() {
  const session = await getServerSession();
  if (session) {
    redirect("/");
  }

  return (
    <div className="relative min-h-screen overflow-hidden bg-[radial-gradient(1100px_circle_at_15%_20%,rgb(var(--brand-jungle-rgb)/0.18),transparent_55%),radial-gradient(900px_circle_at_85%_85%,rgb(var(--brand-cold-rgb)/0.18),transparent_50%)]">
      <div className="pointer-events-none absolute inset-0 opacity-30 [background-image:linear-gradient(to_right,rgb(var(--brand-cold-rgb)/0.16)_1px,transparent_1px),linear-gradient(to_bottom,rgb(var(--brand-cold-rgb)/0.16)_1px,transparent_1px)] [background-size:32px_32px]" />

      <main className="relative mx-auto flex min-h-screen w-full max-w-7xl items-center px-4 py-8 sm:px-8">
        <div className="grid w-full overflow-hidden rounded-3xl border border-border/70 bg-card/85 shadow-[0_24px_80px_rgb(var(--brand-cold-rgb)/0.22)] backdrop-blur-sm lg:grid-cols-2">
          <section className="relative p-6 sm:p-10 lg:p-12">
            <div className="mb-6 flex items-center justify-center lg:hidden">
              <Image
                src="/logos/insuretech-logo-color.png"
                alt="InsureTech"
                width={200}
                height={58}
                style={{ width: "auto", height: "auto" }}
                priority
              />
            </div>

            <div className="mb-8 space-y-3 animate-in fade-in slide-in-from-left-6 duration-700">
              <p className="text-xs font-semibold uppercase tracking-[0.2em] text-[rgb(var(--brand-cold-rgb))]">InsureTech B2B</p>
              <h1 className="text-2xl font-semibold tracking-tight text-foreground sm:text-3xl">Business Portal Login</h1>
              <p className="max-w-md text-sm text-muted-foreground">
                Sign in with your admin mobile number and password to manage employees, policies, billing, and claims.
              </p>
            </div>

            <div className="animate-in fade-in slide-in-from-bottom-4 duration-700 delay-100">
              <LoginForm />
            </div>
          </section>

          <section className="relative hidden min-h-[640px] items-center justify-center overflow-hidden lg:flex">
            <div className="absolute inset-0 bg-[linear-gradient(145deg,rgb(var(--brand-cold-rgb)/0.98),rgb(var(--brand-jungle-rgb)/0.9))]" />
            <div className="absolute -left-20 top-16 h-64 w-64 rounded-full bg-white/10 blur-2xl animate-float-y" />
            <div className="absolute bottom-10 right-10 h-52 w-52 rounded-full bg-white/10 blur-2xl animate-float-x" />
            <div className="absolute right-20 top-1/3 h-24 w-24 rounded-full border border-white/35" />

            <div className="relative z-10 mx-auto flex max-w-md flex-col items-center gap-6 px-10 text-center text-white animate-in fade-in zoom-in-95 duration-700">
              <Image
                src="/logos/insuretech-brand.png"
                alt="InsureTech Logo"
                width={320}
                height={108}
                style={{ width: "auto", height: "auto" }}
                priority
              />
              <h2 className="text-3xl font-semibold leading-tight">Secure Insurance Operations</h2>
              <p className="text-sm text-white/85">
                Unified workspace for underwriting, policy operations, and claims workflows with enterprise-grade access control.
              </p>
              <div className="grid w-full gap-3 text-left">
                <p className="rounded-xl border border-white/20 bg-white/10 px-4 py-3 text-sm">Centralized employee onboarding and coverage management</p>
                <p className="rounded-xl border border-white/20 bg-white/10 px-4 py-3 text-sm">Real-time purchase order, billing, and renewal visibility</p>
                <p className="rounded-xl border border-white/20 bg-white/10 px-4 py-3 text-sm">Documented actions with auditable access events</p>
              </div>
            </div>
          </section>
        </div>
      </main>
    </div>
  );
}
