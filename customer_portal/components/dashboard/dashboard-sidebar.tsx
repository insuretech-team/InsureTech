import { cn } from "@/lib/utils";
import Image from "next/image";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { navigation } from "@/lib/navigation";

const DashboardSidebar = () => {
  const pathname = usePathname();
  return (
    <div className="flex h-full flex-col">
      <div className="flex h-16 items-center border-b px-6">
        <Image
          src="logos/Insuretech.svg"
          alt="Logo"
          width={140}
          height={140}
          className="object-contain"
        />
      </div>

      {/* navigation */}
      <nav className="flex-1 space-y-1 px-0 ml-3 py-4">
        {navigation.map((item) => {
          const isActive = pathname === item.href;
          return (
            <Link
              key={item.name}
              href={item.href}
              className={cn(
                "flex items-center gap-3  px-4 py-2.5 text-sm font-medium transition-colors rounded-l-full",
                isActive
                  ? "text-[#FFFFFF] bg-gradient-to-r from-[#8C34C7] to-[#702A9F]"
                  : "text-sidebar-foreground hover:bg-sidebar-accent hover:text-sidebar-accent-foreground",
              )}
            >
              <Image
                src={item.icon}
                width={16}
                height={16}
                alt=""
                className={cn(
                  "size-5 shrink-0",
                  isActive && "invert brightness-0",
                )}
              />
              {item.name}
            </Link>
          );
        })}
      </nav>
    </div>
  );
};

export default DashboardSidebar;
