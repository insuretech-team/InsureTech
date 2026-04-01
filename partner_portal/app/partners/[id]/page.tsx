'use client';

import { use } from 'react';
import { useRouter } from 'next/navigation';
import DashboardLayout from '@/components/dashboard/dashboard-layout';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { ArrowLeft, Building2, Phone, Mail, MapPin, CreditCard, TrendingUp } from 'lucide-react';
import { getPartnerById } from '@/lib/demo-data/partners';
import { getAgentsByPartnerId } from '@/lib/demo-data/agents';
import type { PartnerStatus } from '@/lib/types';

const statusColors: Record<PartnerStatus, string> = {
  ACTIVE: 'bg-green-100 text-green-800',
  PENDING_VERIFICATION: 'bg-yellow-100 text-yellow-800',
  SUSPENDED: 'bg-red-100 text-red-800',
  TERMINATED: 'bg-gray-100 text-gray-800',
};

export default function PartnerDetailPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = use(params);
  const router = useRouter();
  const partner = getPartnerById(id);
  const agents = getAgentsByPartnerId(id);

  if (!partner) {
    return (
      <DashboardLayout>
        <div className="flex flex-col items-center justify-center py-12">
          <h2 className="text-2xl font-semibold mb-2">Partner Not Found</h2>
          <p className="text-muted-foreground mb-4">The partner you're looking for doesn't exist.</p>
          <Button onClick={() => router.push('/partners')}>
            <ArrowLeft className="mr-2 h-4 w-4" />
            Back to Partners
          </Button>
        </div>
      </DashboardLayout>
    );
  }

  return (
    <DashboardLayout>
      <div className="flex flex-col gap-6">
        {/* Header */}
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-4">
            <Button variant="ghost" size="sm" onClick={() => router.push('/partners')}>
              <ArrowLeft className="h-4 w-4" />
            </Button>
            <div>
              <h1 className="text-2xl font-semibold">{partner.organizationName}</h1>
              <p className="text-muted-foreground">{partner.partnerType.replace(/_/g, ' ')}</p>
            </div>
          </div>
          <Badge className={statusColors[partner.status]} variant="secondary">
            {partner.status.replace(/_/g, ' ')}
          </Badge>
        </div>

        {/* Tabs */}
        <Tabs defaultValue="overview" className="w-full">
          <TabsList>
            <TabsTrigger value="overview">Overview</TabsTrigger>
            <TabsTrigger value="documents">Documents</TabsTrigger>
            <TabsTrigger value="performance">Performance</TabsTrigger>
            <TabsTrigger value="agents">Agents ({agents.length})</TabsTrigger>
          </TabsList>

          {/* Overview Tab */}
          <TabsContent value="overview" className="space-y-6">
            <div className="grid gap-6 md:grid-cols-2">
              {/* Organization Information */}
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <Building2 className="h-5 w-5" />
                    Organization Information
                  </CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                  <div>
                    <div className="text-sm text-muted-foreground">Trade License</div>
                    <div className="font-medium">{partner.tradeLicense}</div>
                  </div>
                  <div>
                    <div className="text-sm text-muted-foreground">TIN Number</div>
                    <div className="font-medium">{partner.tin}</div>
                  </div>
                  <div>
                    <div className="text-sm text-muted-foreground">Partner Type</div>
                    <div className="font-medium">{partner.partnerType.replace(/_/g, ' ')}</div>
                  </div>
                </CardContent>
              </Card>

              {/* Contact Information */}
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <Phone className="h-5 w-5" />
                    Contact Information
                  </CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                  <div className="flex items-center gap-2">
                    <Mail className="h-4 w-4 text-muted-foreground" />
                    <div className="font-medium">{partner.contactEmail}</div>
                  </div>
                  <div className="flex items-center gap-2">
                    <Phone className="h-4 w-4 text-muted-foreground" />
                    <div className="font-medium">{partner.contactPhone}</div>
                  </div>
                  <div className="flex items-start gap-2">
                    <MapPin className="h-4 w-4 text-muted-foreground mt-1" />
                    <div>
                      {partner.nationwideCoverage ? (
                        <span className="font-medium text-green-600">Nationwide Coverage</span>
                      ) : (
                        <div className="font-medium">{partner.serviceLocations.join(', ')}</div>
                      )}
                    </div>
                  </div>
                </CardContent>
              </Card>

              {/* Bank Information */}
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <CreditCard className="h-5 w-5" />
                    Bank Information
                  </CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                  <div>
                    <div className="text-sm text-muted-foreground">Bank Name</div>
                    <div className="font-medium">{partner.bankName}</div>
                  </div>
                  <div>
                    <div className="text-sm text-muted-foreground">Branch</div>
                    <div className="font-medium">{partner.bankBranch || 'N/A'}</div>
                  </div>
                  <div>
                    <div className="text-sm text-muted-foreground">Account Number</div>
                    <div className="font-medium">****{partner.bankAccount.slice(-4)}</div>
                  </div>
                </CardContent>
              </Card>

              {/* Commission Structure */}
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <TrendingUp className="h-5 w-5" />
                    Commission Structure
                  </CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                  <div>
                    <div className="text-sm text-muted-foreground">Acquisition Rate</div>
                    <div className="font-medium">{partner.acquisitionRate}%</div>
                  </div>
                  <div>
                    <div className="text-sm text-muted-foreground">Renewal Rate</div>
                    <div className="font-medium">{partner.renewalRate}%</div>
                  </div>
                  {partner.claimsAssistanceRate && (
                    <div>
                      <div className="text-sm text-muted-foreground">Claims Assistance Rate</div>
                      <div className="font-medium">{partner.claimsAssistanceRate}%</div>
                    </div>
                  )}
                </CardContent>
              </Card>
            </div>

            {/* Policy Types and Service Model */}
            <Card>
              <CardHeader>
                <CardTitle>Policy Types & Service Model</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div>
                  <div className="text-sm text-muted-foreground mb-2">Assigned Policy Types</div>
                  <div className="flex flex-wrap gap-2">
                    {partner.policyTypes.map((type) => (
                      <Badge key={type} variant="secondary">
                        {type}
                      </Badge>
                    ))}
                  </div>
                </div>
                <div className="grid gap-4 md:grid-cols-2">
                  <div>
                    <div className="text-sm text-muted-foreground mb-2">Cashless Service</div>
                    {partner.cashlessEnabled ? (
                      <div className="space-y-1">
                        <div className="text-green-600 font-medium">Enabled</div>
                        <div className="text-sm">
                          Limit: ৳{partner.cashlessLimit?.toLocaleString()}
                        </div>
                        {partner.autoApprovalThreshold && (
                          <div className="text-sm">
                            Auto-approval: ৳{partner.autoApprovalThreshold.toLocaleString()}
                          </div>
                        )}
                      </div>
                    ) : (
                      <div className="text-muted-foreground">Not enabled</div>
                    )}
                  </div>
                  <div>
                    <div className="text-sm text-muted-foreground mb-2">Discount Service</div>
                    {partner.discountEnabled ? (
                      <div className="space-y-1">
                        <div className="text-blue-600 font-medium">Enabled</div>
                        <div className="text-sm">
                          Discount: {partner.discountPercentage}%
                        </div>
                      </div>
                    ) : (
                      <div className="text-muted-foreground">Not enabled</div>
                    )}
                  </div>
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          {/* Documents Tab */}
          <TabsContent value="documents">
            <Card>
              <CardHeader>
                <CardTitle>KYB Documents</CardTitle>
                <CardDescription>
                  Know Your Business verification documents
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="text-muted-foreground">
                  Document management will be implemented in the Documents module
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          {/* Performance Tab */}
          <TabsContent value="performance">
            <div className="grid gap-6 md:grid-cols-3">
              <Card>
                <CardHeader>
                  <CardTitle>Active Policies</CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="text-3xl font-bold">{partner.activePolicies}</div>
                  <div className="text-sm text-muted-foreground">Total policies serviced</div>
                </CardContent>
              </Card>
              <Card>
                <CardHeader>
                  <CardTitle>Total Claims</CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="text-3xl font-bold">{partner.totalClaims}</div>
                  <div className="text-sm text-muted-foreground">Claims submitted</div>
                </CardContent>
              </Card>
              <Card>
                <CardHeader>
                  <CardTitle>Conversion Rate</CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="text-3xl font-bold">{partner.conversionRate}%</div>
                  <div className="text-sm text-muted-foreground">Referral to policy</div>
                </CardContent>
              </Card>
            </div>
          </TabsContent>

          {/* Agents Tab */}
          <TabsContent value="agents">
            <Card>
              <CardHeader>
                <CardTitle>Agents ({agents.length})</CardTitle>
                <CardDescription>
                  Insurance agents under this partner organization
                </CardDescription>
              </CardHeader>
              <CardContent>
                {agents.length === 0 ? (
                  <div className="text-center py-8 text-muted-foreground">
                    No agents registered under this partner
                  </div>
                ) : (
                  <div className="space-y-4">
                    {agents.slice(0, 10).map((agent) => (
                      <div
                        key={agent.id}
                        className="flex items-center justify-between p-4 border rounded-lg"
                      >
                        <div>
                          <div className="font-medium">{agent.fullName}</div>
                          <div className="text-sm text-muted-foreground">{agent.phone}</div>
                        </div>
                        <div className="text-right">
                          <div className="text-sm font-medium">{agent.policiesSold} policies</div>
                          <div className="text-sm text-muted-foreground">
                            ৳{agent.commissionEarned.toLocaleString()} earned
                          </div>
                        </div>
                      </div>
                    ))}
                    {agents.length > 10 && (
                      <Button
                        variant="outline"
                        className="w-full"
                        onClick={() => router.push(`/agents?partner=${id}`)}
                      >
                        View All {agents.length} Agents
                      </Button>
                    )}
                  </div>
                )}
              </CardContent>
            </Card>
          </TabsContent>
        </Tabs>
      </div>
    </DashboardLayout>
  );
}
