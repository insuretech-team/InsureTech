"use client";

import Image from "next/image";
import Link from "next/link";
import { usePathname } from "next/navigation";

import { cn } from "@/lib/utils";
import { b2bDashboardClient } from "@lib/sdk";
import { usePortalPrincipal } from "@lib/auth/portal-session-context";

const DashboardSidebar = () => {
  const pathname = usePathname();
  const principal = usePortalPrincipal();
  const navigation = principal ? b2bDashboardClient.getNavigation(principal.role) : [];

  return (
    <div className="flex h-full flex-col">
      <div className="flex h-20 items-center border-b px-6">
        <Image
          src="/logos/insuretech-brand.png"
          alt="Logo"
          width={220}
          height={72}
          style={{ width: "auto", height: "auto" }}
          className="object-contain"
        />
      </div>

      {principal?.organisationName ? (
        <div className="border-b bg-muted/30 px-6 py-2">
          <p className="text-xs text-muted-foreground">Organisation</p>
          <p className="truncate text-sm font-semibold text-foreground">{principal.organisationName}</p>
        </div>
      ) : null}

      <nav className="ml-3 flex-1 space-y-1 px-0 py-4">
        {navigation.map((item) => {
          const isActive = pathname === item.href;
          return (
            <Link
              key={item.name}
              href={item.href}
              className={cn(
                "portal-nav-link",
                isActive ? "portal-nav-link-active hover:text-sidebar-primary-foreground" : ""
              )}
            >
              <Image
                src={item.icon}
                width={16}
                height={16}
                alt=""
                className={cn("size-5 shrink-0", isActive && "invert brightness-0")}
              />
              {item.name}
            </Link>
          );
        })}
      </nav>
    </div>
  );
};

export default DashboardSidebar;
