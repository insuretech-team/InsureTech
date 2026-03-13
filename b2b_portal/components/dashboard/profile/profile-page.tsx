"use client";

import React from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { LuUser, LuShieldCheck } from "react-icons/lu";
import DashboardLayout from "../dashboard-layout";
import ProfileForm from "../settings/partials/profile-form";
import SecurityForm from "../settings/partials/security-form";

// User-level personal settings — profile info, password, MFA, sessions.
// Organisation-level settings live at /settings.
const VALID_TABS = ["profile", "security"] as const;
type TabValue = typeof VALID_TABS[number];

const ProfilePage = () => {
  const router = useRouter();
  const searchParams = useSearchParams();
  const tabParam = searchParams.get("tab") as TabValue | null;
  const activeTab: TabValue = VALID_TABS.includes(tabParam as TabValue) ? (tabParam as TabValue) : "profile";

  const handleTabChange = (value: string) => {
    const params = new URLSearchParams(searchParams.toString());
    params.set("tab", value);
    router.replace(`/profile?${params.toString()}`);
  };

  return (
    <DashboardLayout>
      <div className="mb-4">
        <h1 className="text-xl font-semibold text-foreground">My Account</h1>
        <p className="text-sm text-muted-foreground mt-1">
          Manage your personal profile, password, two-factor authentication and active sessions.
        </p>
      </div>
      <Tabs value={activeTab} onValueChange={handleTabChange} className="max-w-full">
        <TabsList className="bg-card">
          <TabsTrigger value="profile" className="flex items-center gap-2">
            <LuUser className="size-4" />
            My Profile
          </TabsTrigger>
          <TabsTrigger value="security" className="flex items-center gap-2">
            <LuShieldCheck className="size-4" />
            Security
          </TabsTrigger>
        </TabsList>
        <TabsContent value="profile">
          <ProfileForm />
        </TabsContent>
        <TabsContent value="security">
          <SecurityForm />
        </TabsContent>
      </Tabs>
    </DashboardLayout>
  );
};

export default ProfilePage;
