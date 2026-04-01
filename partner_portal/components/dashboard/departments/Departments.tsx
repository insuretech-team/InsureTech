import React from "react";
import DashboardLayout from "../dashboard-layout";
import { DataTable } from "./data-table/data-table";
import {
  departmentColumns,
  Department,
} from "@/components/dashboard/departments/data-table/columns";

const data: Department[] = [
  {
    id: "EMP-001",
    name: "Finance",
    employeeNo: 32,
    totalPremium: "৳ 432,650",
  },
  {
    id: "EMP-002",
    name: "Accounts",
    employeeNo: 32,
    totalPremium: "৳ 432,650",
  },
  {
    id: "EMP-003",
    name: "Tech",
    employeeNo: 32,
    totalPremium: "৳ 432,650",
  },
  {
    id: "EMP-004",
    name: "Tech",
    employeeNo: 32,
    totalPremium: "৳ 432,650",
  },
  {
    id: "EMP-005",
    name: "Business",
    employeeNo: 32,
    totalPremium: "৳ 432,650",
  },
];

const Departments = () => {
  return (
    <DashboardLayout>
      <div className="space-y-4">
        <DataTable columns={departmentColumns} data={data} loading={false} />
      </div>
    </DashboardLayout>
  );
};

export default Departments;
