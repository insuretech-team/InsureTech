'use client';

import { useState } from 'react';
import DashboardLayout from "@/components/dashboard/dashboard-layout";
import { PartnersTable } from '@/components/partners/partners-table';
import { BulkUploadModal } from '@/components/partners/bulk-upload-modal';
import { demoPartners } from '@/lib/demo-data/partners';

export default function PartnersPage() {
  const [bulkUploadOpen, setBulkUploadOpen] = useState(false);

  return (
    <DashboardLayout>
      <div className="flex flex-col gap-6">
        <div>
          <h1 className="text-2xl font-semibold">Partners</h1>
          <p className="text-muted-foreground">
            Manage partner organizations and their operations
          </p>
        </div>
        
        <PartnersTable 
          partners={demoPartners} 
          onBulkUpload={() => setBulkUploadOpen(true)}
        />

        <BulkUploadModal 
          open={bulkUploadOpen}
          onOpenChange={setBulkUploadOpen}
        />
      </div>
    </DashboardLayout>
  );
}

