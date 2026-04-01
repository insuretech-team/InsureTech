import PlanDetailModal from "@/components/modals/plan-detail-modal";
import { Button } from "@/components/ui/button";
import Image from "next/image";
import React, { useState } from "react";

type Plan = {
  name: string;
  coverage: string;
  premium: string;
  duration: string;
  enrolled: string;
};

type PlanCardProps = {
  plan: Plan;
  onView?: (plan: Plan) => void;
};

const PlanCard: React.FC<PlanCardProps> = ({ plan }) => {
  const [isPlanDetailModalOpen, setIsPlanDetailModalOpen] = useState(false);
  return (
    <>
      <div className="rounded-md border-2 border-gray-50 bg-white shadow-sm">
        <div className="px-4 pt-3 pb-2">
          <p className="text-md font-semibold text-gray-800">
            Plan name: <span className="font-semibold">{plan.name}</span>
          </p>

          <p className="mt-2 text-sm text-gray-400">Health coverage up to</p>

          <div className="mt-1 flex items-center gap-2">
            <span className="text-sm font-semibold text-[var(--primary-deep)]">
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
                <span className="font-semibold text-[var(--primary-deep)]">
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
                <span className="font-semibold text-[var(--primary-deep)]">
                  {plan.duration}
                </span>
              </span>
            </div>
          </div>
        </div>

        <div
          className="flex items-center justify-between gap-3 rounded-b-md px-4 mt-2 py-2"
          style={{ backgroundColor: "var(--bg-soft-primary)" }}
        >
          <div className="flex flex-col">
            <div className="flex items-center gap-2 text-[10px] text-gray-500">
              <Image
                src="./insurance-plans/employees.svg"
                width={16}
                height={16}
                alt="Employees"
              />
              <span className="text-sm text-[#2b2b2b] font-medium">
                Enrolled Employees
              </span>
            </div>
            <span
              className="ml-12 mt-0.5 text-md font-semibold text-[var(--primary-deep)]"
              style={{ marginLeft: "22px" }}
            >
              {plan.enrolled}
            </span>
          </div>

          <Button
            onClick={() => setIsPlanDetailModalOpen(true)}
            variant="outline"
            className="bg-[var(--primary-deep)] text-white hover:bg-[var(--primary-deep)]"
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
