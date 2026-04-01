'use client';

import { useState, useMemo } from 'react';
import { useRouter } from 'next/navigation';
import type { Agent, AgentStatus } from '@/lib/types';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Search, Eye, UserPlus, Upload } from 'lucide-react';

interface AgentsTableProps {
  agents: Agent[];
  onRegisterAgent: () => void;
  onBulkUpload: () => void;
}

const statusColors: Record<AgentStatus, string> = {
  ACTIVE: 'bg-green-100 text-green-800',
  INACTIVE: 'bg-gray-100 text-gray-800',
  SUSPENDED: 'bg-red-100 text-red-800',
};

export function AgentsTable({ agents, onRegisterAgent, onBulkUpload }: AgentsTableProps) {
  const router = useRouter();
  const [searchQuery, setSearchQuery] = useState('');
  const [statusFilter, setStatusFilter] = useState<string>('all');
  const [territoryFilter, setTerritoryFilter] = useState<string>('all');
  const [currentPage, setCurrentPage] = useState(1);
  const itemsPerPage = 50;

  const territories = useMemo(() => {
    const uniqueTerritories = new Set(agents.map(a => a.territory).filter((t): t is string => Boolean(t)));
    return Array.from(uniqueTerritories).sort();
  }, [agents]);

  const filteredAgents = useMemo(() => {
    return agents.filter((agent) => {
      const matchesSearch = searchQuery === '' || agent.fullName.toLowerCase().includes(searchQuery.toLowerCase()) || agent.phone.includes(searchQuery) || agent.nid.includes(searchQuery);
      const matchesStatus = statusFilter === 'all' || agent.status === statusFilter;
      const matchesTerritory = territoryFilter === 'all' || agent.territory === territoryFilter;
      return matchesSearch && matchesStatus && matchesTerritory;
    });
  }, [agents, searchQuery, statusFilter, territoryFilter]);

  const totalPages = Math.ceil(filteredAgents.length / itemsPerPage);
  const paginatedAgents = filteredAgents.slice((currentPage - 1) * itemsPerPage, currentPage * itemsPerPage);

  return (
    <div className="space-y-4">
      <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
        <div className="flex flex-1 gap-2">
          <div className="relative flex-1 max-w-sm">
            <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
            <Input placeholder="Search by name, phone, or NID..." value={searchQuery} onChange={(e) => setSearchQuery(e.target.value)} className="pl-9" />
          </div>
          <Select value={statusFilter} onValueChange={setStatusFilter}>
            <SelectTrigger className="w-[150px]"><SelectValue placeholder="Status" /></SelectTrigger>
            <SelectContent>
              <SelectItem value="all">All Status</SelectItem>
              <SelectItem value="ACTIVE">Active</SelectItem>
              <SelectItem value="INACTIVE">Inactive</SelectItem>
              <SelectItem value="SUSPENDED">Suspended</SelectItem>
            </SelectContent>
          </Select>
          <Select value={territoryFilter} onValueChange={setTerritoryFilter}>
            <SelectTrigger className="w-[180px]"><SelectValue placeholder="Territory" /></SelectTrigger>
            <SelectContent>
              <SelectItem value="all">All Territories</SelectItem>
              {territories.map((territory) => (<SelectItem key={territory} value={territory}>{territory}</SelectItem>))}
            </SelectContent>
          </Select>
        </div>
        <div className="flex gap-2">
          <Button onClick={onRegisterAgent} variant="outline" className="gap-2"><UserPlus className="h-4 w-4" />Register Agent</Button>
          <Button onClick={onBulkUpload} className="gap-2"><Upload className="h-4 w-4" />Bulk Upload</Button>
        </div>
      </div>
      <div className="text-sm text-muted-foreground">Showing {paginatedAgents.length} of {filteredAgents.length} agents</div>
      <div className="rounded-md border">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Name</TableHead>
              <TableHead>Phone</TableHead>
              <TableHead>Status</TableHead>
              <TableHead>Territory</TableHead>
              <TableHead className="text-right">Policies Sold</TableHead>
              <TableHead className="text-right">Commission Earned</TableHead>
              <TableHead className="text-right">Referrals</TableHead>
              <TableHead className="text-right">Conversion Rate</TableHead>
              <TableHead className="text-right">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {paginatedAgents.length === 0 ? (
              <TableRow><TableCell colSpan={9} className="text-center py-8 text-muted-foreground">No agents found</TableCell></TableRow>
            ) : (
              paginatedAgents.map((agent) => (
                <TableRow key={agent.id} className="cursor-pointer hover:bg-muted/50" onClick={() => router.push(`/agents/${agent.id}`)}>
                  <TableCell><div><div className="font-medium">{agent.fullName}</div><div className="text-sm text-muted-foreground">Commission: {agent.commissionRate}%</div></div></TableCell>
                  <TableCell>{agent.phone}</TableCell>
                  <TableCell><Badge className={statusColors[agent.status]} variant="secondary">{agent.status}</Badge></TableCell>
                  <TableCell>{agent.territory || 'N/A'}</TableCell>
                  <TableCell className="text-right font-medium">{agent.policiesSold}</TableCell>
                  <TableCell className="text-right font-medium">৳{agent.commissionEarned.toLocaleString()}</TableCell>
                  <TableCell className="text-right">{agent.referralCount}</TableCell>
                  <TableCell className="text-right">{agent.conversionRate}%</TableCell>
                  <TableCell className="text-right"><Button variant="ghost" size="sm" onClick={(e) => { e.stopPropagation(); router.push(`/agents/${agent.id}`); }}><Eye className="h-4 w-4" /></Button></TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </div>
      {totalPages > 1 && (
        <div className="flex items-center justify-between">
          <div className="text-sm text-muted-foreground">Page {currentPage} of {totalPages}</div>
          <div className="flex gap-2">
            <Button variant="outline" size="sm" onClick={() => setCurrentPage((p) => Math.max(1, p - 1))} disabled={currentPage === 1}>Previous</Button>
            <Button variant="outline" size="sm" onClick={() => setCurrentPage((p) => Math.min(totalPages, p + 1))} disabled={currentPage === totalPages}>Next</Button>
          </div>
        </div>
      )}
    </div>
  );
}
