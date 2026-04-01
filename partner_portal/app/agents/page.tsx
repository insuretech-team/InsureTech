'use client';

import { useState } from 'react';
import DashboardLayout from "@/components/dashboard/dashboard-layout";
import { AgentsTable } from '@/components/agents/agents-table';
import { demoAgents } from '@/lib/demo-data/agents';
import { Button } from '@/components/ui/button';
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog';

export default function AgentsPage() {
  const [registerOpen, setRegisterOpen] = useState(false);
  const [bulkUploadOpen, setBulkUploadOpen] = useState(false);

  return (
    <DashboardLayout>
      <div className="flex flex-col gap-6">
        <div>
          <h1 className="text-2xl font-semibold">Agents</h1>
          <p className="text-muted-foreground">
            Manage partner agents and their activities
          </p>
        </div>
        
        <AgentsTable 
          agents={demoAgents}
          onRegisterAgent={() => setRegisterOpen(true)}
          onBulkUpload={() => setBulkUploadOpen(true)}
        />

        {/* Register Agent Modal - Simplified */}
        <Dialog open={registerOpen} onOpenChange={setRegisterOpen}>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Register New Agent</DialogTitle>
            </DialogHeader>
            <div className="py-4 text-muted-foreground">
              Agent registration form will be implemented here
            </div>
          </DialogContent>
        </Dialog>

        {/* Bulk Upload Modal - Simplified */}
        <Dialog open={bulkUploadOpen} onOpenChange={setBulkUploadOpen}>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Bulk Upload Agents</DialogTitle>
            </DialogHeader>
            <div className="py-4 text-muted-foreground">
              Bulk upload functionality will be implemented here
            </div>
          </DialogContent>
        </Dialog>
      </div>
    </DashboardLayout>
  );
}

