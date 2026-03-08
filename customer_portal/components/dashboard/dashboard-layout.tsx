"use client";
import { type ReactNode, useState } from "react";
import { Button } from "../ui/button";
import { LuX } from "react-icons/lu";
import DashboardSidebar from "./dashboard-sidebar";
import DashboardHeader from "./dashboard-header";

interface DashboardLayoutProps {
  children: ReactNode;
}

const DashboardLayout = ({ children }: DashboardLayoutProps) => {
  const [sidebarOpen, setSidebarOpen] = useState(false);
  return (
    <div className="min-h-screen bg-background">
      {/* mobile sidebar backdrop */}
      {sidebarOpen && (
        <div
          className="fixed inset-0 z-40 bg-black/50 lg:hidden"
          onClick={() => setSidebarOpen(false)}
        />
      )}
      <aside
        className={`fixed left-0 top-0 z-50 h-full w-64 transform bg-sidebar transition-transform duration-200 ease-in-out lg:translate-x-0 ${
          sidebarOpen ? "translate-x-0" : "-translate-x-full"
        }`}
      >
        <DashboardSidebar />
        <Button
          variant="ghost"
          size="icon"
          className="absolute right-4 top-4 lg:hidden"
          onClick={() => setSidebarOpen(false)}
        >
          <LuX className="size-5" />
        </Button>
      </aside>
      <div className="lg:pl-64">
        <DashboardHeader onMenuClick={() => setSidebarOpen(true)} />
        <main className="p-4 md:p-6 lg:p-8">{children}</main>
      </div>
    </div>
  );
};

export default DashboardLayout;
