import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { ArrowRight } from "lucide-react";

const policies = [
  {
    name: "Auto Insurance",
    enrollmentId: "173782011025648",
    coverage: "৳625,000",
    nextDue: "11-01-2026",
    status: "Active",
  },
  {
    name: "Health Insurance",
    enrollmentId: "173782011025648",
    coverage: "৳225,500",
    nextDue: "28-12-2025",
    status: "Expiring",
  },
];
const PolicyOverview = () => {
  return (
    <Card className="bg-card overflow-hidden">
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-4">
        <CardTitle className="text-base font-semibold">
          Policy Overview
        </CardTitle>
        <Button
          variant="link"
          size="sm"
          className="h-auto gap-1 p-0 text-primary hover:text-primary/90"
        >
          View All
          <ArrowRight className="size-4" />
        </Button>
      </CardHeader>
      <CardContent className="space-y-4">
        {policies.map((policy, index) => (
          <div
            key={index}
            className="space-y-3 rounded-lg border bg-background p-4"
          >
            <div className="flex items-start justify-between">
              <div>
                <h3 className="font-semibold text-foreground">{policy.name}</h3>
                <p className="mt-0.5 text-xs text-muted-foreground">
                  Enrollment ID: {policy.enrollmentId}
                </p>
              </div>
              <Badge
                variant={policy.status === "Active" ? "default" : "destructive"}
                className={
                  policy.status === "Active"
                    ? "bg-success/10 text-success"
                    : "bg-destructive/10 text-destructive"
                }
              >
                {policy.status}
              </Badge>
            </div>
            <div className="flex items-center justify-between border-t pt-3 text-sm">
              <div>
                <p className="text-xs text-muted-foreground">Coverage</p>
                <p className="mt-1 font-semibold text-foreground">
                  {policy.coverage}
                </p>
              </div>
              <div className="text-right">
                <p className="text-xs text-muted-foreground">Next Due</p>
                <p className="mt-1 font-semibold text-foreground">
                  {policy.nextDue}
                </p>
              </div>
            </div>
          </div>
        ))}
      </CardContent>
    </Card>
  );
};

export default PolicyOverview;

