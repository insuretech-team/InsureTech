"use client";

import ProfilePage from "@/components/dashboard/profile/profile-page";
import { Suspense } from "react";

const Page = () => (
    <Suspense>
        <ProfilePage />
    </Suspense>
);

export default Page;
