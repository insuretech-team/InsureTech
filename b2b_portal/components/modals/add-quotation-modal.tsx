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

type AddQuotationModalProps = {
  open: boolean;
  onOpenChange: (open: boolean) => void;
};

const focusPurple =
  "focus-visible:ring-primary focus-visible:border-primary focus-visible:ring-2";

const AddQuotationModal = ({ open, onOpenChange }: AddQuotationModalProps) => {
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
            Add Quotation Request
          </DialogTitle>
        </DialogHeader>

        <form onSubmit={onSubmit} className="px-6 py-6">
          <FieldGroup className="space-y-4">
            {/* Select insurance type */}
            <Field>
              <Label htmlFor="insuranceType" className="sr-only">
                Select insurance type
              </Label>
              <Select name="insuranceType" required>
                <SelectTrigger
                  id="insuranceType"
                  className={`w-full h-12 ${focusPurple}`}
                >
                  <SelectValue placeholder="Select insurance type*" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem className="focus:text-foreground" value="health">
                    Health
                  </SelectItem>
                  <SelectItem className="focus:text-foreground" value="life">
                    Life
                  </SelectItem>
                  <SelectItem className="focus:text-foreground" value="motor">
                    Motor
                  </SelectItem>
                  <SelectItem className="focus:text-foreground" value="travel">
                    Travel
                  </SelectItem>
                </SelectContent>
              </Select>
            </Field>

            {/* Select plan */}
            <Field>
              <Label htmlFor="plan" className="sr-only">
                Select plan
              </Label>
              <Select name="plan" required>
                <SelectTrigger
                  id="plan"
                  className={`w-full h-12 ${focusPurple}`}
                >
                  <SelectValue placeholder="Select plan*" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem className="focus:text-foreground" value="basic">
                    Basic
                  </SelectItem>
                  <SelectItem className="focus:text-foreground" value="standard">
                    Standard
                  </SelectItem>
                  <SelectItem className="focus:text-foreground" value="premium">
                    Premium
                  </SelectItem>
                </SelectContent>
              </Select>
            </Field>

            {/* Department */}
            <Field>
              <Label htmlFor="department" className="sr-only">
                Department
              </Label>
              <Select name="department" required>
                <SelectTrigger
                  id="department"
                  className={`w-full h-12 ${focusPurple}`}
                >
                  <SelectValue placeholder="Department*" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem className="focus:text-foreground" value="hr">
                    HR
                  </SelectItem>
                  <SelectItem className="focus:text-foreground" value="tech">
                    Tech
                  </SelectItem>
                  <SelectItem className="focus:text-foreground" value="accounts">
                    Accounts
                  </SelectItem>
                  <SelectItem className="focus:text-foreground" value="business">
                    Business
                  </SelectItem>
                </SelectContent>
              </Select>
            </Field>

            {/* No. of employee */}
            <Field>
              <Label htmlFor="employeeNo" className="sr-only">
                No. of employee
              </Label>
              <Input
                id="employeeNo"
                name="employeeNo"
                placeholder="No. of employee*"
                type="number"
                min={1}
                className={`h-12 ${focusPurple}`}
                required
              />
            </Field>

            {/* Quoted amount */}
            <Field>
              <Label htmlFor="quotedAmount" className="sr-only">
                Quoted amount
              </Label>
              <Input
                id="quotedAmount"
                name="quotedAmount"
                placeholder="Quoted amount*"
                className={`h-12 ${focusPurple}`}
                required
              />
            </Field>
          </FieldGroup>

          <DialogFooter className="mt-8">
            <Button
              type="submit"
              className="w-full h-12 text-white bg-gradient-to-r from-primary to-accent hover:opacity-95"
            >
              Add Quotation to List
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
};

export default AddQuotationModal;

