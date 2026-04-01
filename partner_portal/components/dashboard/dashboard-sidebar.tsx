import { cn } from "@/lib/utils";
import Image from "next/image";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { navigation } from "@/lib/navigation";
import { useEffect, useState } from "react";

type SessionUser = {
  userId: string;
  partnerId: string;
  organisationName: string;
  name: string;
  email: string;
  role: string;
};

const DashboardSidebar = () => {
  const pathname = usePathname();
  const [user, setUser] = useState<SessionUser | null>(null);

  useEffect(() => {
    // Fetch session to get user role and organization info
    fetch("/api/auth/session")
      .then((res) => res.json())
      .then((data) => {
        if (data.session) {
          setUser({
            userId: data.session.principal.user?.userId ?? "",
            partnerId: data.session.principal.partnerId ?? "",
            organisationName: data.session.principal.organisationName ?? "",
            name: data.session.principal.displayName ?? "",
            email: data.session.principal.user?.email ?? "",
            role: data.session.principal.role ?? "VIEWER",
          });
        }
      })
      .catch(() => setUser(null));
  }, []);

  // Filter navigation based on role
  const filteredNavigation = navigation.filter((item) => {
    // Partners tab only for SYSTEM_ADMIN
    if (item.href === "/partners") {
      return user?.role === "SYSTEM_ADMIN";
    }
    // All other tabs visible to authenticated users
    return true;
  });

  return (
    <div className="flex h-full flex-col">
      <div className="flex h-16 items-center border-b px-6">
        <Image
          src="logos/logo.svg"
          alt="Logo"
          width={140}
          height={140}
          className="object-contain"
        />
      </div>

      {/* Organization name banner for partner users */}
      {user?.organisationName && user?.role !== "SYSTEM_ADMIN" && (
        <div className="border-b px-6 py-2 bg-muted/30">
          <p className="text-xs text-muted-foreground">Organization</p>
          <p className="text-sm font-semibold text-foreground truncate">
            {user.organisationName}
          </p>
        </div>
      )}

      {/* navigation */}
      <nav className="flex-1 space-y-1 px-0 ml-3 py-4">
        {filteredNavigation.map((item) => {
          const isActive = pathname === item.href;
          return (
            <Link
              key={item.name}
              href={item.href}
              className={cn(
                "flex items-center gap-3  px-4 py-2.5 text-sm font-medium transition-colors rounded-l-full",
                isActive
                  ? "text-[#FFFFFF] bg-gradient-to-r from-[var(--primary-deep)] to-[var(--primary-light)]"
                  : "text-sidebar-foreground hover:bg-sidebar-accent hover:text-sidebar-accent-foreground",
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
