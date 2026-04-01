import DashboardLayout from "@/components/dashboard/dashboard-layout";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Heart, Car, Shield } from "lucide-react";

const upcomingPayments = [
  {
    policyName: "Health Insurance",
    type: "Health",
    icon: Heart,
    amount: "৳14,500",
    dueDate: "22-12-2025",
    status: "Due Soon",
  },
  {
    policyName: "Auto Insurance",
    type: "Vehicle",
    icon: Car,
    amount: "৳8,500",
    dueDate: "27-12-2025",
    status: "Due Soon",
  },
  {
    policyName: "Life Insurance",
    type: "Life",
    icon: Shield,
    amount: "৳11,250",
    dueDate: "02-12-2026",
    status: "Upcoming",
  },
];

const paymentHistory = [
  {
    policyName: "Health Insurance",
    amount: "৳14,500",
    date: "22-11-2024",
    status: "Paid",
  },
  {
    policyName: "Auto Insurance",
    amount: "৳8,500",
    date: "27-10-2024",
    status: "Paid",
  },
  {
    policyName: "Life Insurance",
    amount: "৳11,250",
    date: "02-10-2024",
    status: "Paid",
  },
];

export default function PaymentsPage() {
  return (
    <DashboardLayout>
      <div className="space-y-6">
        <div>
          <h1 className="text-2xl font-bold text-foreground">Payments</h1>
          <p className="mt-1 text-sm text-muted-foreground">
            Manage your insurance premium payments
          </p>
        </div>

        <Card>
          <CardHeader>
            <CardTitle className="text-base">Upcoming Payments</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              {upcomingPayments.map((payment, index) => (
                <div
                  key={index}
                  className="flex flex-col gap-3 rounded-lg border bg-background p-4 sm:flex-row sm:items-center sm:justify-between"
                >
                  <div className="flex items-center gap-3">
                    <div className="flex size-10 shrink-0 items-center justify-center rounded-full bg-[#EEE1F7]">
                      <payment.icon className="size-5 text-[#8C34C7]" />
                    </div>
                    <div>
                      <h3 className="font-semibold text-foreground">
                        {payment.policyName}
                      </h3>
                      <p className="mt-0.5 text-sm text-muted-foreground">
                        {payment.dueDate}
                      </p>
                    </div>
                  </div>
                  <div className="flex items-center gap-3 sm:flex-row-reverse">
                    <Button
                      size="sm"
                      className="sm:w-auto bg-[#8C34C7] hover:bg-[#8C34C7]/90 hover:cursor-pointer"
                    >
                      Pay Now
                    </Button>
                    <div className="text-right sm:text-left">
                      <p className="text-lg font-bold text-[#8C34C7]">
                        {payment.amount}
                      </p>
                      <Badge
                        variant="outline"
                        className="mt-1 bg-[#8C34C7]/10 text-[#8C34C7] hover:bg-[#8C34C7]/20"
                      >
                        {payment.status}
                      </Badge>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="text-base">Payment History</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="divide-y">
              {paymentHistory.map((payment, index) => (
                <div
                  key={index}
                  className="flex items-center justify-between py-4 first:pt-0 last:pb-0"
                >
                  <div>
                    <h3 className="font-medium text-foreground">
                      {payment.policyName}
                    </h3>
                    <p className="mt-1 text-sm text-muted-foreground">
                      {payment.date}
                    </p>
                  </div>
                  <div className="text-right">
                    <p className="font-semibold text-foreground">
                      {payment.amount}
                    </p>
                    <Badge
                      variant="default"
                      className="mt-1 bg-success/10 text-success hover:bg-success/20"
                    >
                      {payment.status}
                    </Badge>
                  </div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      </div>
    </DashboardLayout>
  );
}
