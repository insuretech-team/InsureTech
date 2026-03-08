"use client";
import {
  LuMenu,
  LuGlobe,
  LuBell,
  LuUser,
  LuSettings,
  LuLogOut,
} from "react-icons/lu";
import { Button } from "../ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "../ui/dropdown-menu";
import { Avatar, AvatarFallback, AvatarImage } from "../ui/avatar";

interface DashboardHeaderProps {
  onMenuClick: () => void;
}

const DashboardHeader = ({ onMenuClick }: DashboardHeaderProps) => {
  return (
    <header className="stick top-0 z-30 flex h-16 items-center borde-b bg-white px-4 md:px-6 lg:px-8">
      <Button
        variant="ghost"
        size="icon"
        className="lg:hidden"
        onClick={onMenuClick}
      >
        <LuMenu className="size-5" />
      </Button>

      <div className="ml-auto flex items-center gap-2">
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button
              variant="outline"
              size="icon"
              className="relative size-9 hover:text-[#8C34C7] cursor-pointer border-0 bg-transparent"
            >
              <LuBell className="size-4" />
              <span className="absolute right-1 top-1 flex size-2 rounded-full bg-[#8C34C7]" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end" className="w-80">
            <div className="px-4 py-3">
              <p className="text-sm font-semibold">Notifications</p>
            </div>
            <DropdownMenuSeparator />
            <DropdownMenuItem className="cursor-pointer p-4 hover:bg-secondary">
              <div className="flex flex-col gap-1">
                <p className="text-sm font-medium text-gray-800">
                  Payment Due Soon
                </p>
                <p className="text-xs text-muted-foreground">
                  Your Health Insurance payment is due on 22-12-2025
                </p>
                <p className="text-xs text-muted-foreground">2 hours ago</p>
              </div>
            </DropdownMenuItem>
            <DropdownMenuItem className="cursor-pointer p-4 hover:bg-secondary">
              <div className="flex flex-col gap-1">
                <p className="text-sm font-medium text-gray-800">
                  Claim Status Updated
                </p>
                <p className="text-xs text-muted-foreground">
                  Your claim #173782011025648 is now under review
                </p>
                <p className="text-xs text-muted-foreground">5 hours ago</p>
              </div>
            </DropdownMenuItem>
            <DropdownMenuItem className="cursor-pointer p-4 hover:bg-secondary">
              <div className="flex flex-col gap-1">
                <p className="text-sm font-medium text-gray-800">
                  New Policy Document
                </p>
                <p className="text-xs text-muted-foreground">
                  Your auto insurance policy document is ready
                </p>
                <p className="text-xs text-muted-foreground">1 day ago</p>
              </div>
            </DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuItem className="cursor-pointer justify-center text-[#8C34C7] hover:bg-secondary hover:text-[#702A9F]">
              View All Notifications
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>

        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button
              variant="link"
              size="icon"
              className="size-9 rounded-full p-0"
            >
              <Avatar className="size-9 cursor-pointer border-2 border-border">
                <AvatarImage src="avatar.jpg" alt="User" />
                <AvatarFallback className="bg-muted-foreground text-primary-foreground">
                  SA
                </AvatarFallback>
              </Avatar>
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end" className="w-56">
            <div className="flex items-center gap-2 px-2 py-3">
              <Avatar className="size-10">
                <AvatarImage src="avatar.jpg" alt="User" />
                <AvatarFallback className="bg-secondary text-primary-foreground">
                  SA
                </AvatarFallback>
              </Avatar>
              <div className="flex flex-col">
                <p className="text-sm font-medium">John Doe</p>
                <p className="text-xs text-muted-foreground">
                  john.doe@email.com
                </p>
              </div>
            </div>
            <DropdownMenuSeparator />
            <DropdownMenuItem className="cursor-pointer hover:bg-secondary">
              <LuUser className="mr-2 size-4" />
              <span className="text-gray-800">Profile</span>
            </DropdownMenuItem>
            <DropdownMenuItem className="cursor-pointer hover:bg-secondary">
              <LuSettings className="mr-2 size-4" />
              <span className="text-gray-800">Settings</span>
            </DropdownMenuItem>
            <DropdownMenuItem className="cursor-pointer text-destructive hover:bg-secondary">
              <LuLogOut className="mr-2 size-4" />
              <span className="text-destructive">Log Out</span>
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </header>
  );
};

export default DashboardHeader;
