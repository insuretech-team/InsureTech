import { Button } from "@/components/ui/button";
import { CreditCard, FileText, HelpCircle } from "lucide-react";
import Link from "next/link";

const QuickAccess = () => {
  return (
    <div className="rounded-lg bg-gradient-to-r from-primary to-accent p-6">
      <h2 className="mb-4 text-xl font-bold text-primary-foreground">
        Quick Access
      </h2>
      <div className="grid gap-4 sm:grid-cols-3">
        <Button
          variant="secondary"
          className="h-auto justify-start gap-3 bg-white py-3 text-foreground hover:bg-white/90"
        >
          <Link href="/payments" passHref className="flex space-x-4">
            <CreditCard className="size-5 text-primary" />
            <span className="font-medium text-primary">Make Payment</span>
          </Link>
        </Button>
        <Button
          variant="secondary"
          className="h-auto justify-start gap-3 bg-white py-3 text-foreground hover:bg-white/90"
        >
          <Link href="/claims" passHref className="flex space-x-4">
            <FileText className="size-5 text-primary" />
            <span className="font-medium text-primary">Claim Policy</span>
          </Link>
        </Button>
        <Button
          variant="secondary"
          className="h-auto justify-start gap-3 bg-white py-3 text-foreground hover:bg-white/90"
        >
          <Link href="/support" passHref className="flex space-x-4">
            <HelpCircle className="size-5 text-primary" />
            <span className="font-medium text-primary">Help & Support</span>
          </Link>
        </Button>
      </div>
    </div>
  );
};

export default QuickAccess;

