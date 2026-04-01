import { Checkbox } from "@/components/ui/checkbox";

import { Card, CardContent } from "@/components/ui/card";
import {
  Field,
  FieldContent,
  FieldDescription,
  FieldGroup,
  FieldLabel,
  FieldTitle,
} from "@/components/ui/field";
import { Button } from "@/components/ui/button";

import { workflows } from "@/lib/workflows";

const WorkflowForm = () => {
  return (
    <Card>
      <form className="py-3">
        <CardContent className="text-muted-foreground text-sm">
          <FieldGroup className="gap-3">
            {workflows.map((item) => {
              return (
                <FieldLabel key={item.id}>
                  <Field orientation="horizontal">
                    <Checkbox id="toggle-checkbox-2" name="toggle-checkbox-2" />
                    <FieldContent>
                      <FieldTitle className="text-md font-semibold text-[#2b2b2b]">
                        {item.title}
                      </FieldTitle>
                      <FieldDescription>{item.description}</FieldDescription>
                    </FieldContent>
                  </Field>
                </FieldLabel>
              );
            })}
          </FieldGroup>
          <div className="flex items-center justify-end mt-4">
            <Button variant="default" className="bg-[var(--primary-deep)] ">
              Save Changes
            </Button>
          </div>
        </CardContent>
      </form>
    </Card>
  );
};

export default WorkflowForm;
