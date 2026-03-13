import { Card, CardContent } from "@/components/ui/card";
import Image from "next/image";

interface PurchaseOrderCardProps {
  title: string;
  value: number;
  icon: string;
  bgColor: string;
}

const PurchaseOrderCard = (props: PurchaseOrderCardProps) => {
  return (
    <Card className="bg-card overflow-hidden">
      <CardContent className="py-0 px-2 relative">
        <div style={{ backgroundColor: props.bgColor }} className="rounded-md px-3 py-3 w-max">
          <Image src={props.icon} width={24} height={24} alt={props.title} />
        </div>
        <div className="mt-2 space-y-2 relative z-10">
          <div className="flex items-start gap-3">
            <div className="text-xs text-muted-foreground">
              <div className="mb-1 font-bold text-2xl text-foreground">{props.value}</div>
              <span className="font-medium text-sm">{props.title}</span>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  );
};

export default PurchaseOrderCard;
