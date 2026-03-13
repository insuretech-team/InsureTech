import PlanDetailModal from "@/components/modals/plan-detail-modal";
import { Button } from "@/components/ui/button";
import Image from "next/image";
import React, { useState } from "react";

type Plan = {
  name: string;
  coverage: string;
  premium: string; // keep string since your data is string
  duration: string;
  enrolled: string; // keep string since your data is string
};

type PlanCardProps = {
  plan: Plan;
  onView?: (plan: Plan) => void;
};

const PlanCard: React.FC<PlanCardProps> = ({ plan }) => {
  const [isPlanDetailModalOpen, setIsPlanDetailModalOpen] = useState(false);
  return (
    <>
      <div className="rounded-md border border-gray-50 bg-white shadow-sm">
        <div className="px-4 pt-3 pb-2">
          <p className="text-[11px] font-semibold text-gray-800">
            Plan name: <span className="font-semibold">{plan.name}</span>
          </p>

          <p className="mt-2 text-sm text-gray-400">Health coverage up to</p>

          <div className="mt-1 flex items-center gap-2">
            <span className="text-sm font-semibold text-primary">
              ৳{plan.coverage}
            </span>
          </div>

          <div className="mt-2 space-y-2">
            <div className="flex items-center gap-2 text-sm text-gray-500">
              <span className="flex items-center">
                <Image
                  src="./insurance-plans/sparkles.svg"
                  width={16}
                  height={16}
                  alt="Sparkles"
                />
                <span className="px-2">Premium price :</span>
                <span className="font-semibold text-primary">
                  ৳{plan.premium}
                </span>
              </span>
            </div>

            <div className="flex items-center gap-2 text-sm text-gray-500">
              <span className="flex items-center">
                <Image
                  src="./insurance-plans/clock-five.svg"
                  width={16}
                  height={16}
                  alt="Clock"
                />
                <span className="px-2">Policy duration :</span>
                <span className="font-semibold text-primary">
                  {plan.duration}
                </span>
              </span>
            </div>
          </div>
        </div>

        <div className="mt-2 flex items-center justify-between gap-3 rounded-b-md bg-primary/10 px-4 py-2">
          <div className="flex flex-col">
            <div className="flex items-center gap-2 text-[10px] text-gray-500">
              <Image
                src="./insurance-plans/employees.svg"
                width={16}
                height={16}
                alt="Employees"
              />
              <span className="text-sm text-foreground font-medium">
                Enrolled Employees
              </span>
            </div>
            <span
              className="ml-12 mt-0.5 text-md font-semibold text-primary"
              style={{ marginLeft: "22px" }}
            >
              {plan.enrolled}
            </span>
          </div>

          <Button
            onClick={() => setIsPlanDetailModalOpen(true)}
            variant="outline"
            className="bg-primary text-primary-foreground hover:bg-primary/90"
          >
            <span>View Details</span>
            <Image
              src="./insurance-plans/info.svg"
              width={16}
              height={16}
              alt="Info"
            />
          </Button>
        </div>
      </div>
      {isPlanDetailModalOpen && (
        <PlanDetailModal
          open={isPlanDetailModalOpen}
          onOpenChange={setIsPlanDetailModalOpen}
        />
      )}
    </>
  );
};

export default PlanCard;
