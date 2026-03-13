"use client";

import Settings from "@/components/dashboard/settings/settings";
import { Suspense } from "react";

const page = () => {
  return (
    <Suspense>
      <Settings />
    </Suspense>
  );
};

export default page;
