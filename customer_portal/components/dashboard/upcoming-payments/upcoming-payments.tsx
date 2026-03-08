import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import {
  LuArrowRight,
  LuHandHeart,
  LuCar,
  LuBriefcaseMedical,
} from "react-icons/lu";

const payments = [
  {
    type: "Health",
    icon: LuBriefcaseMedical,
    date: "22-12-2025",
    amount: "৳14,500",
    iconBg: "bg-success/10",
    iconColor: "text-success",
  },
  {
    type: "Auto",
    icon: LuCar,
    date: "27-12-2025",
    amount: "৳8,500",
    iconBg: "bg-chart-2/10",
    iconColor: "text-chart-2",
  },
  {
    type: "Life",
    icon: LuHandHeart,
    date: "02-12-2026",
    amount: "৳11,250",
    iconBg: "bg-chart-3/10",
    iconColor: "text-chart-3",
  },
];

const UpcomingPayments = () => {
  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-4">
        <CardTitle className="text-lg font-semibold">
          Upcoming Payments
        </CardTitle>
        <Button
          variant="link"
          size="sm"
          className="gap-1 text-[#8C34C7] hover:text-[#8C34C7]/90"
        >
          View All
          <LuArrowRight className="size-4" />
        </Button>
      </CardHeader>
      <CardContent className="space-y-3">
        {payments.map((payment, index) => (
          <div
            key={index}
            className="flex items-center justify-between rounded-lg border bg-card p-4 transition-colors hover:bg-gray-200"
          >
            <div className="flex items-center gap-3">
              <div
                className={`flex size-10 items-center justify-center rounded-lg ${payment.iconBg}`}
              >
                <payment.icon className={`size-5 ${payment.iconColor}`} />
              </div>
              <div>
                <p className="font-semibold text-card-foreground">
                  {payment.type}
                </p>
                <p className="text-sm text-muted-foreground">{payment.date}</p>
              </div>
            </div>
            <div className="text-right">
              <p className="text-lg font-bold text-card-foreground">
                {payment.amount}
              </p>
              <Button
                variant="ghost"
                size="sm"
                className="h-auto p-0 text-[#8C34C7] hover:text-[#8C34C7]/90"
              >
                Pay Now
              </Button>
            </div>
          </div>
        ))}
      </CardContent>
    </Card>
  );
};

export default UpcomingPayments;
