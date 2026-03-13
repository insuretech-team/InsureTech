// Employee form value types — shared between Create and Edit modes.

export type EmployeeGender = "EMPLOYEE_GENDER_MALE" | "EMPLOYEE_GENDER_FEMALE" | "EMPLOYEE_GENDER_OTHER" | "";

export interface EmployeeFormValues {
  // Identity
  name: string;
  employeeId: string;
  businessId: string;
  email: string;
  mobileNumber: string;
  // HR
  departmentId: string;
  dateOfBirth: string;       // YYYY-MM-DD
  dateOfJoining: string;     // YYYY-MM-DD
  gender: EmployeeGender;
  // Insurance
  insuranceCategory: number; // InsuranceType enum value
  assignedPlanId: string;
  coverageAmount: string;    // decimal string e.g. "5000"
  numberOfDependent: number;
}

export const EMPTY_EMPLOYEE_FORM: EmployeeFormValues = {
  name: "",
  employeeId: "",
  businessId: "",
  email: "",
  mobileNumber: "",
  departmentId: "",
  dateOfBirth: "",
  dateOfJoining: "",
  gender: "",
  insuranceCategory: 0,
  assignedPlanId: "",
  coverageAmount: "",
  numberOfDependent: 0,
};

export type EmployeeFormMode = "create" | "edit";

export interface EmployeeFormErrors {
  name?: string;
  employeeId?: string;
  businessId?: string;
  departmentId?: string;
  dateOfJoining?: string;
  [key: string]: string | undefined;
}
