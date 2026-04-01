import Departments from "@/components/dashboard/departments/Departments";
import { requireServerSessionRole } from "@lib/auth/session";

const page = async () => {
  await requireServerSessionRole(
    ["SYSTEM_ADMIN", "B2B_ORG_ADMIN", "BUSINESS_ADMIN", "HR_MANAGER"],
    "/"
  );
  return <Departments />;
};

export default page;
