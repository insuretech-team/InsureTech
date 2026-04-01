import EmployeesTable from "@/components/dashboard/employees/employees-table";
import { requireServerSessionRole } from "@lib/auth/session";

const page = async () => {
  await requireServerSessionRole(
    ["SYSTEM_ADMIN", "B2B_ORG_ADMIN", "BUSINESS_ADMIN", "HR_MANAGER"],
    "/"
  );
  return <EmployeesTable />;
};

export default page;
