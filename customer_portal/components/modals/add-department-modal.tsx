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

type AddDepartmentModalProps = {
  open: boolean;
  onOpenChange: (open: boolean) => void;
};

const focusPurple =
  "focus-visible:ring-[#8C34C7] focus-visible:border-[#8C34C7] focus-visible:ring-2";

const AddDepartmentModal = ({
  open,
  onOpenChange,
}: AddDepartmentModalProps) => {
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
            Add Department
          </DialogTitle>
        </DialogHeader>

        <form onSubmit={onSubmit} className="px-6 py-6">
          <FieldGroup className="space-y-4 gap-0">
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              {/* Name */}
              <Field>
                <Label htmlFor="name" className="sr-only">
                  Name
                </Label>
                <Input
                  id="name"
                  name="name"
                  placeholder="Department Name*"
                  className={focusPurple}
                  required
                />
              </Field>

              {/* Number of Employee */}
              <Field>
                <Label htmlFor="numberOfEmployee" className="sr-only">
                  Number of Employee
                </Label>
                <Input
                  id="numberOfEmployee"
                  name="numberOfEmployee"
                  placeholder="Number of employee*"
                  className={focusPurple}
                  required
                />
              </Field>

              {/* Total premium */}
              <Field>
                <Label htmlFor="totalPremium" className="sr-only">
                  Total premium
                </Label>
                <Input
                  id="totalPremium"
                  name="totalPremium"
                  placeholder="Total premium*"
                  className={focusPurple}
                  required
                />
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

export default AddDepartmentModal;
