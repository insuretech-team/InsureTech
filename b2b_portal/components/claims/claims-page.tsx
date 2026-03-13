import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { ArrowRight } from "lucide-react";
import DashboardLayout from "@/components/dashboard/dashboard-layout";

const claims = [
  {
    id: "173782011025648",
    policyType: "Health",
    claimAmount: "৳45,000",
    filedDate: "15-12-2024",
    status: "Under Review",
    statusColor: "bg-primary",
  },
  {
    id: "173782011025647",
    policyType: "Auto",
    claimAmount: "৳125,000",
    filedDate: "10-12-2024",
    status: "Processing",
    statusColor: "bg-orange-500",
  },
  {
    id: "173782011025646",
    policyType: "Health",
    claimAmount: "৳18,500",
    filedDate: "28-11-2024",
    status: "Approved",
    statusColor: "bg-success",
  },
  {
    id: "173782011025645",
    policyType: "Life",
    claimAmount: "৳75,000",
    filedDate: "15-11-2024",
    status: "Paid",
    statusColor: "bg-success",
  },
];

export default function ClaimsPage() {
  return (
    <DashboardLayout>
      <div className="space-y-6">
        <div className="flex flex-col justify-between gap-4 sm:flex-row sm:items-center">
          <div>
            <h1 className="text-2xl font-bold text-foreground">Claims</h1>
            <p className="mt-1 text-sm text-muted-foreground">
              Track and manage your insurance claims
            </p>
          </div>
          <Button className="sm:w-auto bg-primary hover:bg-primary/90 hover:cursor-pointer">
            File New Claim
          </Button>
        </div>

        <div className="grid gap-4 md:grid-cols-2">
          {claims.map((claim, index) => (
            <Card key={index} className="py-4">
              <CardHeader className="pb-3">
                <div className="flex items-start justify-between">
                  <div>
                    <CardTitle className="text-base">
                      Claim ID: {claim.id}
                    </CardTitle>
                    <p className="mt-1 text-sm text-muted-foreground">
                      {claim.policyType} Insurance
                    </p>
                  </div>
                  <Badge
                    className={`${claim.statusColor} text-white hover:${claim.statusColor}/90`}
                  >
                    {claim.status}
                  </Badge>
                </div>
              </CardHeader>
              <CardContent className="space-y-3">
                <div className="flex items-center justify-between text-sm">
                  <div>
                    <p className="text-muted-foreground">Claim Amount</p>
                    <p className="mt-1 text-lg font-bold text-primary">
                      {claim.claimAmount}
                    </p>
                  </div>
                  <div className="text-right">
                    <p className="text-muted-foreground">Filed Date</p>
                    <p className="mt-1 font-medium">{claim.filedDate}</p>
                  </div>
                </div>
                <Button
                  variant="link"
                  size="sm"
                  className="h-auto gap-1 p-0 text-primary hover:text-primary/90"
                >
                  View Details
                  <ArrowRight className="size-4" />
                </Button>
              </CardContent>
            </Card>
          ))}
        </div>
      </div>
    </DashboardLayout>
  );
}

