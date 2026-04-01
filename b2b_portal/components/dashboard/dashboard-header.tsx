"use client";

import {
  LuBell,
  LuLogOut,
  LuMenu,
  LuSettings,
  LuShieldCheck,
  LuUser,
} from "react-icons/lu";
import { useEffect, useMemo, useState } from "react";
import { useRouter } from "next/navigation";

import { bffClient } from "@lib/sdk/b2b-sdk-client";
import { usePortalSession } from "@lib/auth/portal-session-context";
import type { SessionUser } from "@lib/types/b2b";
import { Avatar, AvatarFallback, AvatarImage } from "../ui/avatar";
import { Button } from "../ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "../ui/dropdown-menu";

interface DashboardHeaderProps {
  onMenuClick: () => void;
}

function getInitials(name?: string, email?: string): string {
  if (name && name.trim()) {
    const parts = name.trim().split(/\s+/);
    if (parts.length >= 2) return (parts[0][0] + parts[1][0]).toUpperCase();
    return parts[0].slice(0, 2).toUpperCase();
  }
  if (email) return email.slice(0, 2).toUpperCase();
  return "??";
}

function toSessionUser(session: ReturnType<typeof usePortalSession>["session"]): SessionUser | null {
  if (!session) {
    return null;
  }
  return {
    userId: session.principal.user.userId,
    businessId: session.principal.businessId,
    organisationName: session.principal.organisationName ?? "",
    name: session.principal.displayName,
    email: session.principal.user.email,
    role: session.principal.role,
  };
}

const DashboardHeader = ({ onMenuClick }: DashboardHeaderProps) => {
  const router = useRouter();
  const { session } = usePortalSession();
  const user = useMemo(() => toSessionUser(session), [session]);
  const [profilePhotoUrl, setProfilePhotoUrl] = useState("");
  const [fullName, setFullName] = useState("");
  const isBeneficiary = user?.role === "B2B_BENEFICIARY";

  const fetchProfile = () => {
    bffClient.auth.getProfile().then((profileRes) => {
      if (profileRes.ok && profileRes.profile) {
        const p = profileRes.profile as Record<string, unknown>;
        if (typeof p.full_name === "string" && p.full_name) setFullName(p.full_name);
        if (typeof p.profile_photo_url === "string" && p.profile_photo_url) {
          setProfilePhotoUrl(p.profile_photo_url);
        }
      }
    }).catch(() => {});
  };

  useEffect(() => {
    setFullName(session?.principal.displayName ?? "");
    setProfilePhotoUrl("");
    if (!session) {
      return;
    }

    fetchProfile();
    window.addEventListener("profile:updated", fetchProfile);
    return () => window.removeEventListener("profile:updated", fetchProfile);
  }, [session]);

  async function handleLogout() {
    try {
      await bffClient.auth.logout();
    } catch (error) {
      console.warn("Logout request failed:", error);
    } finally {
      router.replace("/login");
      router.refresh();
    }
  }

  return (
    <header className="portal-header flex items-center px-4 md:px-6 lg:px-8">
      <Button
        variant="ghost"
        size="icon"
        className="lg:hidden"
        onClick={onMenuClick}
      >
        <LuMenu className="size-5" />
      </Button>

      {user?.organisationName ? (
        <div className="ml-4 hidden items-center gap-2 md:flex">
          <span className="rounded-md border border-border bg-muted/40 px-3 py-1 text-xs font-semibold tracking-wide text-foreground">
            {user.organisationName}
          </span>
        </div>
      ) : null}

      <div className="ml-auto flex items-center gap-2">
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button
              variant="outline"
              size="icon"
              className="relative size-9 cursor-pointer border-0 bg-transparent hover:text-primary"
            >
              <LuBell className="size-4" />
              <span className="absolute right-1 top-1 flex size-2 rounded-full bg-primary" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end" className="w-80">
            <div className="px-4 py-3">
              <p className="text-sm font-semibold">Notifications</p>
            </div>
            <DropdownMenuSeparator />
            <DropdownMenuItem className="cursor-pointer p-4 hover:bg-secondary">
              <div className="flex flex-col gap-1">
                <p className="text-sm font-medium text-gray-800">Payment Due Soon</p>
                <p className="text-xs text-muted-foreground">
                  Your Health Insurance payment is due on 22-12-2025
                </p>
                <p className="text-xs text-muted-foreground">2 hours ago</p>
              </div>
            </DropdownMenuItem>
            <DropdownMenuItem className="cursor-pointer p-4 hover:bg-secondary">
              <div className="flex flex-col gap-1">
                <p className="text-sm font-medium text-gray-800">Claim Status Updated</p>
                <p className="text-xs text-muted-foreground">
                  Your claim #173782011025648 is now under review
                </p>
                <p className="text-xs text-muted-foreground">5 hours ago</p>
              </div>
            </DropdownMenuItem>
            <DropdownMenuItem className="cursor-pointer p-4 hover:bg-secondary">
              <div className="flex flex-col gap-1">
                <p className="text-sm font-medium text-gray-800">New Policy Document</p>
                <p className="text-xs text-muted-foreground">
                  Your auto insurance policy document is ready
                </p>
                <p className="text-xs text-muted-foreground">1 day ago</p>
              </div>
            </DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuItem className="cursor-pointer justify-center text-primary hover:bg-secondary hover:text-accent">
              View All Notifications
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>

        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="link" size="icon" className="size-9 rounded-full p-0">
              <Avatar className="size-9 cursor-pointer border-2 border-border">
                {profilePhotoUrl ? (
                  <AvatarImage src={profilePhotoUrl} alt={fullName || user?.name || "User"} />
                ) : null}
                <AvatarFallback className="bg-muted-foreground text-xs font-semibold text-primary-foreground">
                  {getInitials(fullName || user?.name, user?.email)}
                </AvatarFallback>
              </Avatar>
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end" className="w-56">
            <div className="flex min-w-0 items-center gap-2 px-2 py-3">
              <Avatar className="size-10">
                {profilePhotoUrl ? (
                  <AvatarImage src={profilePhotoUrl} alt={fullName || user?.name || "User"} />
                ) : null}
                <AvatarFallback className="bg-secondary text-xs font-semibold text-primary-foreground">
                  {getInitials(fullName || user?.name, user?.email)}
                </AvatarFallback>
              </Avatar>
              <div className="flex min-w-0 flex-col">
                <p className="truncate text-sm font-medium">{fullName || user?.name || "Business User"}</p>
                <p className="truncate text-xs text-muted-foreground">
                  {user?.email ?? "user@insuretech.local"}
                </p>
              </div>
            </div>
            <DropdownMenuSeparator />
            <DropdownMenuItem
              className="cursor-pointer hover:bg-secondary"
              onClick={() => router.push("/profile")}
            >
              <LuUser className="mr-2 size-4" />
              <span className="text-foreground">My Profile</span>
            </DropdownMenuItem>
            <DropdownMenuItem
              className="cursor-pointer hover:bg-secondary"
              onClick={() => router.push("/profile?tab=security")}
            >
              <LuShieldCheck className="mr-2 size-4" />
              <span className="text-foreground">Security</span>
            </DropdownMenuItem>
            {!isBeneficiary ? (
              <DropdownMenuItem
                className="cursor-pointer hover:bg-secondary"
                onClick={() => router.push("/settings")}
              >
                <LuSettings className="mr-2 size-4" />
                <span className="text-foreground">Organisation Settings</span>
              </DropdownMenuItem>
            ) : null}
            <DropdownMenuSeparator />
            <DropdownMenuItem
              className="cursor-pointer text-destructive hover:bg-secondary"
              onClick={handleLogout}
            >
              <LuLogOut className="mr-2 size-4" />
              <span className="text-destructive">Log Out</span>
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </header>
  );
};

export default DashboardHeader;
