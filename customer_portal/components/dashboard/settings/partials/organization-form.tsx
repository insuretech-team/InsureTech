import React from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Field, FieldGroup } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Button } from "@/components/ui/button";

const focusPurple =
  "focus-visible:ring-[#8C34C7] focus-visible:border-[#8C34C7] focus-visible:ring-2";

const OrganizationForm = () => {
  return (
    <Card>
      <form className="py-3">
        <CardHeader>
          <CardTitle>Organization Info.</CardTitle>
        </CardHeader>
        <CardContent className="text-muted-foreground text-sm">
          <FieldGroup className="space-y-4 gap-0">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <Field>
                <Label htmlFor="organizationName" className="sr-only">
                  Organization Name
                </Label>
                <Input
                  id="organizationName"
                  name="organizationName"
                  placeholder="Organization name*"
                  className={focusPurple}
                  required
                />
              </Field>

              <Field>
                <Label htmlFor="registrationNo" className="sr-only">
                  Registration No.
                </Label>
                <Input
                  id="registrationNo"
                  name="registrationNo"
                  placeholder="Registration No.*"
                  className={focusPurple}
                  required
                />
              </Field>
            </div>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <Field>
                {/* <Label htmlFor="taxId" className="sr-only">
                  Organization Name
                </Label> */}
                <Input
                  id="taxId"
                  name="taxId"
                  placeholder="Tax ID*"
                  className={focusPurple}
                  required
                />
              </Field>
            </div>
          </FieldGroup>
        </CardContent>
        <CardHeader className="mt-4">
          <CardTitle>Primary Contact</CardTitle>
          <FieldGroup className="space-y-4 gap-0">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
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

              <Field>
                <Label htmlFor="email" className="sr-only">
                  Email.
                </Label>
                <Input
                  id="email"
                  name="email"
                  placeholder="Email"
                  className={focusPurple}
                />
              </Field>
            </div>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <Field>
                <Label htmlFor="phone" className="sr-only">
                  Phone
                </Label>
                <Input
                  id="phone"
                  name="phone"
                  placeholder="Phone*"
                  className={focusPurple}
                  required
                />
              </Field>
              <Field>
                <Label htmlFor="department" className="sr-only">
                  Department
                </Label>
                <Select name="department" required>
                  <SelectTrigger
                    id="department"
                    className={`w-full ${focusPurple} hover:text-purple-300`}
                  >
                    <SelectValue placeholder="Department*" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="tech" className="focus:text-[#2b2b2b]">
                      Tech
                    </SelectItem>
                    <SelectItem
                      value="finance"
                      className="focus:text-[#2b2b2b]"
                    >
                      Finance
                    </SelectItem>
                    <SelectItem value="hr" className="focus:text-[#2b2b2b]">
                      HR
                    </SelectItem>
                  </SelectContent>
                </Select>
              </Field>
            </div>

            <div className="flex items-center justify-end">
              <Button
                variant="default"
                className="bg-[#8C34C7] hover:bg-[#7f20be]"
              >
                Save Changes
              </Button>
            </div>
          </FieldGroup>
        </CardHeader>
      </form>
    </Card>
  );
};

export default OrganizationForm;
