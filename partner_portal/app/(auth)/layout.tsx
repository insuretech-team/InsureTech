"use client";
import type { ReactNode } from "react";
import Image from "next/image";
import { usePathname } from "next/navigation";

export default function AuthLayout({ children }: { children: ReactNode }) {
  const pathname = usePathname();

  const bannerImage =
    pathname === "/login" ? "/banner_login.png" : "/banner_otp.png";
  return (
    <div className="min-h-screen w-full grid grid-cols-1 lg:grid-cols-2">
      {/* Left */}
      <div className="relative hidden lg:block">
        <Image
          src={bannerImage}
          alt="Labaid Insuretech"
          fill
          priority
          className="object-cover"
        />
      </div>

      {/* Right */}
      <div className="flex items-center justify-center bg-white px-6 py-12">
        <div className="w-full max-w-[520px]">
          <div className="flex justify-center mb-10">
            <Image
              src="/logos/logo.svg"
              alt="Labaid Insuretech Company Ltd."
              width={260}
              height={70}
              priority
            />
          </div>
          {children}
        </div>
      </div>
    </div>
  );
}
