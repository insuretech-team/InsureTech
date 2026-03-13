"use client";

import React from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import DashboardLayout from "../dashboard-layout";
import OrganizationForm from "./partials/organization-form";
import WorkflowForm from "./partials/workflow-form";
import NotificationForm from "./partials/notification-form";

// Organisation-level settings only.
// User-level settings (My Profile, Password, MFA, Sessions) live at /profile.
const VALID_TABS = ["organization", "workflow", "notification"] as const;
type TabValue = typeof VALID_TABS[number];

const Settings = () => {
  const router = useRouter();
  const searchParams = useSearchParams();
  const tabParam = searchParams.get("tab") as TabValue | null;
  const activeTab: TabValue = VALID_TABS.includes(tabParam as TabValue) ? (tabParam as TabValue) : "organization";

  const handleTabChange = (value: string) => {
    const params = new URLSearchParams(searchParams.toString());
    params.set("tab", value);
    router.replace(`/settings?${params.toString()}`);
  };

  return (
    <DashboardLayout>
      <div className="mb-4">
        <h1 className="text-xl font-semibold text-foreground">Organisation Settings</h1>
        <p className="text-sm text-muted-foreground mt-1">
          Manage your organisation profile, approval workflows and notification preferences.
        </p>
      </div>
      <Tabs value={activeTab} onValueChange={handleTabChange} className="max-w-full">
        <TabsList className="bg-card">
          <TabsTrigger value="organization">Organisation Profile</TabsTrigger>
          <TabsTrigger value="workflow">Approval Workflows</TabsTrigger>
          <TabsTrigger value="notification">Notification Preferences</TabsTrigger>
        </TabsList>
        <TabsContent value="organization">
          <OrganizationForm />
        </TabsContent>
        <TabsContent value="workflow">
          <WorkflowForm />
        </TabsContent>
        <TabsContent value="notification">
          <NotificationForm />
        </TabsContent>
      </Tabs>
    </DashboardLayout>
  );
};

export default Settings;

