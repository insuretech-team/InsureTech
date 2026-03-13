import { Card, CardContent } from "@/components/ui/card";
import Image from "next/image";

interface StatsCardProps {
  title: string;
  value: string | number;
  icon: string;
  bgColor: string;
  bgIcon: string;
}

const StatsCard = (props: StatsCardProps) => {
  return (
    <Card className="bg-card overflow-hidden">
      <CardContent className="p-4 relative">
        <div
          style={{ backgroundColor: props.bgColor }}
          className="rounded-md px-3 py-3 w-max"
        >
          <Image src={props.icon} width={24} height={24} alt={props.title} />
        </div>
        <div className="mt-10 space-y-2 relative z-10">
          <div className="flex items-start gap-3">
            <div className="text-xs text-muted-foreground">
              <div className="mb-1 font-bold text-3xl text-foreground">
                {props.value}
              </div>
              <span className="font-medium text-sm">{props.title}</span>
            </div>
          </div>
        </div>
        <Image
          src={props.bgIcon}
          width={100}
          height={100}
          alt={`${props.title} background icon`}
          className="absolute right-0 bottom-0 pointer-events-none"
        />
      </CardContent>
    </Card>
  );
};

export default StatsCard;

