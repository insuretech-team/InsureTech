import { cn } from "@/lib/utils";
import Image from "next/image";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { authClient, b2bDashboardClient } from "@lib/sdk";
import { useEffect, useState } from "react";
import type { SessionUser } from "@lib/types/b2b";

const DashboardSidebar = () => {
  const pathname = usePathname();
  const [user, setUser] = useState<SessionUser | null>(null);

  useEffect(() => {
    authClient
      .getSession()
      .then((response) => {
        const session = response.session;
        if (!session) {
          setUser(null);
          return;
        }
        setUser({
          userId: session.principal.user?.userId ?? "",
          businessId: session.principal.businessId ?? "",
          organisationName: session.principal.organisationName ?? "",
          name: session.principal.displayName ?? "",
          email: session.principal.user?.email ?? "",
          role: session.principal.role ?? "BUSINESS_ADMIN",
        });
      })
      .catch(() => setUser(null));
  }, []);

  // While session is loading (user === null), show a skeleton/empty nav
  // to avoid a flash of wrong-role navigation items.
  const navigation = b2bDashboardClient.getNavigation(user?.role);
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

      {/* Org name banner — shown below logo for B2B admin users */}
      {user?.organisationName && (
        <div className="border-b px-6 py-2 bg-muted/30">
          <p className="text-xs text-muted-foreground">Organisation</p>
          <p className="text-sm font-semibold text-foreground truncate">{user.organisationName}</p>
        </div>
      )}

      {/* navigation */}
      <nav className="flex-1 space-y-1 px-0 ml-3 py-4">
        {navigation.map((item) => {
          const isActive = pathname === item.href;
          return (
            <Link
              key={item.name}
              href={item.href}
              className={cn(
                "portal-nav-link",
                isActive
                  ? "portal-nav-link-active hover:text-sidebar-primary-foreground"
                  : "",
              )}
            >
              <Image
                src={item.icon}
                width={16}
                height={16}
                alt=""
                className={cn(
                  "size-5 shrink-0",
                  isActive && "invert brightness-0",
                )}
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
