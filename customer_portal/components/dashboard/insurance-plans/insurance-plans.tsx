import React from "react";
import DashboardLayout from "../dashboard-layout";
import PlanCard from "./plan-card";
import Image from "next/image";

const InsurancePlans = () => {
  const healthPlans = [
    {
      id: 1,
      name: "Seba",
      coverage: "25,000",
      premium: "800",
      duration: "1 year",
      enrolled: "125",
    },
    {
      id: 2,
      name: "Surokkha",
      coverage: "25,000",
      premium: "900",
      duration: "1 year",
      enrolled: "125",
    },
    {
      id: 3,
      name: "Susastho",
      coverage: "25,000",
      premium: "800",
      duration: "1 year",
      enrolled: "125",
    },
  ];

  const lifePlans = [
    {
      id: 1,
      name: "Verosa",
      coverage: "25,000",
      premium: "800",
      duration: "1 year",
      enrolled: "125",
    },
    {
      id: 2,
      name: "Astha",
      coverage: "25,000",
      premium: "800",
      duration: "1 year",
      enrolled: "125",
    },
    {
      id: 3,
      name: "Nivroy",
      coverage: "25,000",
      premium: "800",
      duration: "1 year",
      enrolled: "125",
    },
  ];

  return (
    <DashboardLayout>
      <h4 className="text-md font-semibold text-gray-800">
        Insurance Plan List
      </h4>
      <div className="mt-3 grid grid-cols-1 gap-4 lg:grid-cols-2">
        {/* Health Insurance */}
        <div className="rounded-md bg-white p-4 shadow-sm">
          <div className="flex items-center gap-2">
            <h5 className="text-sm font-semibold text-gray-800">
              <span className="flex items-center">
                <Image
                  src="./insurance-plans/health.svg"
                  width={32}
                  height={32}
                  alt="Health"
                />
                <span className="pl-2">Health Insurance</span>
              </span>
            </h5>
          </div>

          <div className="mt-4 space-y-3">
            {healthPlans.map((p) => (
              <PlanCard key={p.id} plan={p} />
            ))}
          </div>
        </div>

        {/* Life Insurance */}
        <div className="rounded-md bg-white p-4 shadow-sm">
          <div className="flex items-center gap-2">
            <span className="flex items-center">
              <Image
                src="./insurance-plans/life.svg"
                width={32}
                height={32}
                alt="Life"
              />
              <span className="pl-2">Life Insurance</span>
            </span>
          </div>

          <div className="mt-4 space-y-3">
            {lifePlans.map((p) => (
              <PlanCard key={p.id} plan={p} />
            ))}
          </div>
        </div>
      </div>
    </DashboardLayout>
  );
};

export default InsurancePlans;
