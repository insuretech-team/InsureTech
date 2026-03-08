"use client";

import * as React from "react";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Field, FieldGroup } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { LuCalendarDays } from "react-icons/lu";

type AddEmployeeModalProps = {
  open: boolean;
  onOpenChange: (open: boolean) => void;
};

const focusPurple =
  "focus-visible:ring-[#8C34C7] focus-visible:border-[#8C34C7] focus-visible:ring-2";

const AddEmployeeModal = ({ open, onOpenChange }: AddEmployeeModalProps) => {
  const onSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    // TODO: handle submit (API call)
    // onOpenChange(false);
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-2xl p-0">
        <DialogHeader className="px-6 py-4 border-b">
          <DialogTitle className="text-xl font-semibold">
            Add Employee
          </DialogTitle>
        </DialogHeader>

        <form onSubmit={onSubmit} className="px-6 py-6">
          <FieldGroup className="space-y-4 gap-0">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {/* Name */}
              <Field>
                <Label htmlFor="name" className="sr-only">
                  Name
                </Label>
                <Input
                  id="name"
                  name="name"
                  placeholder="Name*"
                  className={focusPurple}
                  required
                />
              </Field>

              {/* Employee ID */}
              <Field>
                <Label htmlFor="employeeId" className="sr-only">
                  Employee ID
                </Label>
                <Input
                  id="employeeId"
                  name="employeeId"
                  placeholder="Employee ID*"
                  className={focusPurple}
                  required
                />
              </Field>
            </div>
            {/* Email */}
            <Field>
              <Label htmlFor="email" className="sr-only">
                Email
              </Label>
              <Input
                id="email"
                name="email"
                placeholder="Email"
                type="email"
                className={focusPurple}
              />
            </Field>

            {/* DOB + Gender */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <Field>
                <Label htmlFor="dob" className="sr-only">
                  Date of birth
                </Label>
                <div className="relative">
                  <Input
                    id="dob"
                    name="dob"
                    placeholder="Date of birth*"
                    type="text"
                    onFocus={(e) => (e.target.type = "date")}
                    onBlur={(e) => {
                      if (!e.target.value) e.target.type = "text";
                    }}
                    className={`${focusPurple} pr-10`}
                    required
                  />
                  <LuCalendarDays className="pointer-events-none absolute right-3 top-1/2 -translate-y-1/2 text-[#8C34C7]" />
                </div>
              </Field>

              <Field>
                <Label htmlFor="gender" className="sr-only">
                  Gender
                </Label>
                <Select name="gender" required>
                  <SelectTrigger
                    id="gender"
                    className={`w-full ${focusPurple} hover:text-purple-300`}
                  >
                    <SelectValue placeholder="Gender*" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="male" className="focus:text-[#2b2b2b]">
                      Male
                    </SelectItem>
                    <SelectItem value="female" className="focus:text-[#2b2b2b]">
                      Female
                    </SelectItem>
                    <SelectItem value="other" className="focus:text-[#2b2b2b]">
                      Other
                    </SelectItem>
                  </SelectContent>
                </Select>
              </Field>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {/* Mobile */}
              <Field>
                <Label htmlFor="mobile" className="sr-only">
                  Mobile no.
                </Label>
                <Input
                  id="mobile"
                  name="mobile"
                  placeholder="Mobile no.*"
                  type="tel"
                  className={focusPurple}
                  required
                />
              </Field>
              {/* Date of joining */}
              <Field>
                <Label htmlFor="doj" className="sr-only">
                  Date of joining
                </Label>
                <div className="relative">
                  <Input
                    id="doj"
                    name="doj"
                    placeholder="Date of joining*"
                    type="text"
                    onFocus={(e) => (e.target.type = "date")}
                    onBlur={(e) => {
                      if (!e.target.value) e.target.type = "text";
                    }}
                    className={`${focusPurple} pr-10`}
                    required
                  />
                  <LuCalendarDays className="pointer-events-none absolute right-3 top-1/2 -translate-y-1/2 text-[#8C34C7]" />
                </div>
              </Field>
            </div>

            {/* Department + Income range */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <Field>
                <Label htmlFor="department" className="sr-only">
                  Department
                </Label>
                <Select name="department" required>
                  <SelectTrigger
                    id="department"
                    className={`w-full ${focusPurple}`}
                  >
                    <SelectValue placeholder="Department*" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="hr" className="focus:text-[#2b2b2b]">
                      HR
                    </SelectItem>
                    <SelectItem
                      value="engineering"
                      className="focus:text-[#2b2b2b]"
                    >
                      Engineering
                    </SelectItem>
                    <SelectItem value="sales" className="focus:text-[#2b2b2b]">
                      Sales
                    </SelectItem>
                    <SelectItem
                      value="finance"
                      className="focus:text-[#2b2b2b]"
                    >
                      Finance
                    </SelectItem>
                  </SelectContent>
                </Select>
              </Field>

              <Field>
                <Label htmlFor="income" className="sr-only">
                  Income range
                </Label>
                <Select name="income" required>
                  <SelectTrigger
                    id="income"
                    className={`w-full ${focusPurple}`}
                  >
                    <SelectValue placeholder="Income range*" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="0-25k" className="focus:text-[#2b2b2b]">
                      $0 - $25k
                    </SelectItem>
                    <SelectItem value="25-50k" className="focus:text-[#2b2b2b]">
                      $25k - $50k
                    </SelectItem>
                    <SelectItem
                      value="50-100k"
                      className="focus:text-[#2b2b2b]"
                    >
                      $50k - $100k
                    </SelectItem>
                    <SelectItem value="100k+" className="focus:text-[#2b2b2b]">
                      $100k+
                    </SelectItem>
                  </SelectContent>
                </Select>
              </Field>
            </div>
          </FieldGroup>

          <DialogFooter className="mt-8">
            <Button
              type="submit"
              className="w-full h-12 text-white bg-gradient-to-r from-[#8C34C7] to-[#702A9F] hover:opacity-95"
            >
              Add Employee to List
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
};

export default AddEmployeeModal;
