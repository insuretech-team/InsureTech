import DashboardLayout from "@/components/dashboard/dashboard-layout";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Heart, Car, Shield } from "lucide-react";

const policies = [
  {
    name: "Health Insurance",
    type: "Health",
    icon: Heart,
    enrollmentId: "173782011025648",
    coverage: "৳225,500",
    premium: "৳14,500",
    nextDue: "28-12-2025",
    status: "Expiring",
    provider: "Chartered Life",
  },
  {
    name: "Auto Insurance",
    type: "Vehicle",
    icon: Car,
    enrollmentId: "173782011025648",
    coverage: "৳625,000",
    premium: "৳8,500",
    nextDue: "11-01-2026",
    status: "Active",
    provider: "National Insurance",
  },
  {
    name: "Life Insurance",
    type: "Life",
    icon: Shield,
    enrollmentId: "173782011025649",
    coverage: "৳1,200,000",
    premium: "৳11,250",
    nextDue: "02-12-2026",
    status: "Active",
    provider: "MetLife Bangladesh",
  },
];

const PoliciesPage = () => {
  return (
    <DashboardLayout>
      <div className="space-y-6">
        <div>
          <h1 className="text-2xl font-bold text-foreground">My Policies</h1>
          <p className="mt-1 text-sm text-muted-foreground">
            Manage and view all your insurance policies
          </p>
        </div>

        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
          {policies.map((policy, index) => (
            <Card key={index} className="py-4">
              <CardHeader>
                <div className="flex items-start justify-between">
                  <div className="flex items-center gap-3">
                    <div className="flex size-10 items-center justify-center rounded-full bg-[#EEE1F7]">
                      <policy.icon className="size-5 text-[#8C34C7]" />
                    </div>
                    <div>
                      <CardTitle className="text-base">{policy.name}</CardTitle>
                      <p className="text-xs text-muted-foreground">
                        {policy.type}
                      </p>
                    </div>
                  </div>
                  <Badge
                    variant={
                      policy.status === "Active" ? "default" : "destructive"
                    }
                    className={
                      policy.status === "Active"
                        ? "bg-success/10 text-success hover:bg-success/20"
                        : "bg-destructive/10 text-destructive hover:bg-destructive/20"
                    }
                  >
                    {policy.status}
                  </Badge>
                </div>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="space-y-2 text-sm">
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Provider</span>
                    <span className="font-medium">{policy.provider}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Enrollment ID</span>
                    <span className="font-mono text-xs">
                      {policy.enrollmentId}
                    </span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Coverage</span>
                    <span className="font-semibold text-success">
                      {policy.coverage}
                    </span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Premium</span>
                    <span className="font-medium">{policy.premium}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Next Due</span>
                    <span className="font-medium">{policy.nextDue}</span>
                  </div>
                </div>
                <Button
                  className="w-full bg-transparent hover:text-[#8C34C7]"
                  variant="outline"
                >
                  View Details
                </Button>
              </CardContent>
            </Card>
          ))}
        </div>
      </div>
    </DashboardLayout>
  );
};

export default PoliciesPage;
