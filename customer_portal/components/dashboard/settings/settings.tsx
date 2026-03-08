import React from "react";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import DashboardLayout from "../dashboard-layout";
import OrganizationForm from "./partials/organization-form";
import WorkflowForm from "./partials/workflow-form";
import NotificationForm from "./partials/notification-form";

const Settings = () => {
  return (
    <DashboardLayout>
      <Tabs defaultValue="organization" className="max-w-full">
        <TabsList className="bg-[#FFFFFF]">
          <TabsTrigger
            className="data-[state=active]:bg-[#8C34C7] 
             data-[state=active]:text-white"
            value="organization"
          >
            Organization Profile
          </TabsTrigger>
          <TabsTrigger
            className="data-[state=active]:bg-[#8C34C7] 
             data-[state=active]:text-white"
            value="workflow"
          >
            Approval Workflows{" "}
          </TabsTrigger>
          <TabsTrigger
            className="data-[state=active]:bg-[#8C34C7] 
             data-[state=active]:text-white"
            value="notification"
          >
            Notification Preferences
          </TabsTrigger>
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
