'use client';

import { useState, useMemo } from 'react';
import { useRouter } from 'next/navigation';
import type { Partner, PartnerType, PartnerStatus, PolicyType } from '@/lib/types';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { Search, Eye, Upload } from 'lucide-react';

interface PartnersTableProps {
  partners: Partner[];
  onBulkUpload: () => void;
}

const statusColors: Record<PartnerStatus, string> = {
  ACTIVE: 'bg-green-100 text-green-800',
  PENDING_VERIFICATION: 'bg-yellow-100 text-yellow-800',
  SUSPENDED: 'bg-red-100 text-red-800',
  TERMINATED: 'bg-gray-100 text-gray-800',
};

const policyTypeColors: Record<PolicyType, string> = {
  HEALTH: 'bg-blue-100 text-blue-800',
  MOTOR: 'bg-purple-100 text-purple-800',
  LIFE: 'bg-pink-100 text-pink-800',
  FIRE: 'bg-orange-100 text-orange-800',
  PET: 'bg-teal-100 text-teal-800',
  DEVICE: 'bg-indigo-100 text-indigo-800',
  TRAVEL: 'bg-cyan-100 text-cyan-800',
};

export function PartnersTable({ partners, onBulkUpload }: PartnersTableProps) {
  const router = useRouter();
  const [searchQuery, setSearchQuery] = useState('');
  const [statusFilter, setStatusFilter] = useState<string>('all');
  const [typeFilter, setTypeFilter] = useState<string>('all');
  const [policyTypeFilter, setPolicyTypeFilter] = useState<string>('all');
  const [currentPage, setCurrentPage] = useState(1);
  const itemsPerPage = 50;

  // Filter partners
  const filteredPartners = useMemo(() => {
    return partners.filter((partner) => {
      const matchesSearch =
        searchQuery === '' ||
        partner.organizationName.toLowerCase().includes(searchQuery.toLowerCase()) ||
        partner.tradeLicense.toLowerCase().includes(searchQuery.toLowerCase()) ||
        partner.tin.includes(searchQuery);

      const matchesStatus = statusFilter === 'all' || partner.status === statusFilter;
      const matchesType = typeFilter === 'all' || partner.partnerType === typeFilter;
      const matchesPolicyType =
        policyTypeFilter === 'all' || partner.policyTypes.includes(policyTypeFilter as PolicyType);

      return matchesSearch && matchesStatus && matchesType && matchesPolicyType;
    });
  }, [partners, searchQuery, statusFilter, typeFilter, policyTypeFilter]);

  // Pagination
  const totalPages = Math.ceil(filteredPartners.length / itemsPerPage);
  const paginatedPartners = filteredPartners.slice(
    (currentPage - 1) * itemsPerPage,
    currentPage * itemsPerPage
  );

  const handleRowClick = (partnerId: string) => {
    router.push(`/partners/${partnerId}`);
  };

  return (
    <div className="space-y-4">
      {/* Filters and Actions */}
      <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
        <div className="flex flex-1 gap-2">
          <div className="relative flex-1 max-w-sm">
            <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
            <Input
              placeholder="Search by name, trade license, or TIN..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="pl-9"
            />
          </div>
          <Select value={statusFilter} onValueChange={setStatusFilter}>
            <SelectTrigger className="w-[180px]">
              <SelectValue placeholder="Status" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">All Status</SelectItem>
              <SelectItem value="ACTIVE">Active</SelectItem>
              <SelectItem value="PENDING_VERIFICATION">Pending</SelectItem>
              <SelectItem value="SUSPENDED">Suspended</SelectItem>
              <SelectItem value="TERMINATED">Terminated</SelectItem>
            </SelectContent>
          </Select>
          <Select value={typeFilter} onValueChange={setTypeFilter}>
            <SelectTrigger className="w-[180px]">
              <SelectValue placeholder="Partner Type" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">All Types</SelectItem>
              <SelectItem value="HOSPITAL">Hospital</SelectItem>
              <SelectItem value="CLINIC">Clinic</SelectItem>
              <SelectItem value="PHARMACY">Pharmacy</SelectItem>
              <SelectItem value="AUTO_REPAIR">Auto Repair</SelectItem>
              <SelectItem value="PET_CLINIC">Pet Clinic</SelectItem>
              <SelectItem value="FIRE_INSPECTOR">Fire Inspector</SelectItem>
              <SelectItem value="AGENT_NETWORK">Agent Network</SelectItem>
            </SelectContent>
          </Select>
          <Select value={policyTypeFilter} onValueChange={setPolicyTypeFilter}>
            <SelectTrigger className="w-[180px]">
              <SelectValue placeholder="Policy Type" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">All Policies</SelectItem>
              <SelectItem value="HEALTH">Health</SelectItem>
              <SelectItem value="MOTOR">Motor</SelectItem>
              <SelectItem value="LIFE">Life</SelectItem>
              <SelectItem value="FIRE">Fire</SelectItem>
              <SelectItem value="PET">Pet</SelectItem>
              <SelectItem value="DEVICE">Device</SelectItem>
            </SelectContent>
          </Select>
        </div>
        <Button onClick={onBulkUpload} className="gap-2">
          <Upload className="h-4 w-4" />
          Bulk Upload
        </Button>
      </div>

      {/* Results count */}
      <div className="text-sm text-muted-foreground">
        Showing {paginatedPartners.length} of {filteredPartners.length} partners
      </div>

      {/* Table */}
      <div className="rounded-md border">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Organization</TableHead>
              <TableHead>Type</TableHead>
              <TableHead>Status</TableHead>
              <TableHead>Policy Types</TableHead>
              <TableHead>Service Model</TableHead>
              <TableHead>Location</TableHead>
              <TableHead>Agents</TableHead>
              <TableHead>Policies</TableHead>
              <TableHead className="text-right">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {paginatedPartners.length === 0 ? (
              <TableRow>
                <TableCell colSpan={9} className="text-center py-8 text-muted-foreground">
                  No partners found
                </TableCell>
              </TableRow>
            ) : (
              paginatedPartners.map((partner) => (
                <TableRow
                  key={partner.id}
                  className="cursor-pointer hover:bg-muted/50"
                  onClick={() => handleRowClick(partner.id)}
                >
                  <TableCell>
                    <div>
                      <div className="font-medium">{partner.organizationName}</div>
                      <div className="text-sm text-muted-foreground">{partner.tradeLicense}</div>
                    </div>
                  </TableCell>
                  <TableCell>
                    <span className="text-sm">{partner.partnerType.replace(/_/g, ' ')}</span>
                  </TableCell>
                  <TableCell>
                    <Badge className={statusColors[partner.status]} variant="secondary">
                      {partner.status.replace(/_/g, ' ')}
                    </Badge>
                  </TableCell>
                  <TableCell>
                    <div className="flex flex-wrap gap-1">
                      {partner.policyTypes.map((type) => (
                        <Badge
                          key={type}
                          className={policyTypeColors[type]}
                          variant="secondary"
                        >
                          {type}
                        </Badge>
                      ))}
                    </div>
                  </TableCell>
                  <TableCell>
                    <div className="flex flex-col gap-1 text-sm">
                      {partner.cashlessEnabled && (
                        <span className="text-green-600">
                          Cashless (৳{partner.cashlessLimit?.toLocaleString()})
                        </span>
                      )}
                      {partner.discountEnabled && (
                        <span className="text-blue-600">
                          Discount ({partner.discountPercentage}%)
                        </span>
                      )}
                    </div>
                  </TableCell>
                  <TableCell>
                    <div className="text-sm">
                      {partner.nationwideCoverage ? (
                        <span className="text-green-600">Nationwide</span>
                      ) : (
                        <span>{partner.serviceLocations.join(', ')}</span>
                      )}
                    </div>
                  </TableCell>
                  <TableCell>{partner.agentCount}</TableCell>
                  <TableCell>{partner.activePolicies}</TableCell>
                  <TableCell className="text-right">
                    <Button
                      variant="ghost"
                      size="sm"
                      onClick={(e) => {
                        e.stopPropagation();
                        handleRowClick(partner.id);
                      }}
                    >
                      <Eye className="h-4 w-4" />
                    </Button>
                  </TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </div>

      {/* Pagination */}
      {totalPages > 1 && (
        <div className="flex items-center justify-between">
          <div className="text-sm text-muted-foreground">
            Page {currentPage} of {totalPages}
          </div>
          <div className="flex gap-2">
            <Button
              variant="outline"
              size="sm"
              onClick={() => setCurrentPage((p) => Math.max(1, p - 1))}
              disabled={currentPage === 1}
            >
              Previous
            </Button>
            <Button
              variant="outline"
              size="sm"
              onClick={() => setCurrentPage((p) => Math.min(totalPages, p + 1))}
              disabled={currentPage === totalPages}
            >
              Next
            </Button>
          </div>
        </div>
      )}
    </div>
  );
}
