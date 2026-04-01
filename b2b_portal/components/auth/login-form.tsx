"use client";

import { FormEvent, useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { CheckCircle2, XCircle } from "lucide-react";

import { bffClient, type EmployeeLoginOrganisation } from "@lib/sdk/b2b-sdk-client";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import {
  BD_MOBILE_EXAMPLES,
  getBangladeshMobileValidationMessage,
  normalizeBangladeshMobile,
  normalizeBangladeshMobileOrRaw,
} from "@/src/lib/utils/bd-mobile";
import {
  getPasswordValidationMessage,
  PASSWORD_REQUIREMENTS_HINT,
} from "@/src/lib/utils/password";

type AuthMode = "admin" | "employee";

type DialogState =
  | { open: false }
  | { open: true; kind: "success" | "error"; title: string; message: string };

function getMobileHint(value: string): string | null {
  if (!value.trim()) return null;
  const normalized = normalizeBangladeshMobile(value);
  if (normalized) return null;
  const digits = value.replace(/\D/g, "");
  if (digits.length < 7) return null;
  return getBangladeshMobileValidationMessage();
}

function normalizeOrganisationLookupValue(value: string): string {
  return value.trim().toLowerCase().replace(/[^a-z0-9]+/g, " ");
}

function findBestOrganisationMatch(
  query: string,
  organisations: EmployeeLoginOrganisation[]
): EmployeeLoginOrganisation | null {
  const normalizedQuery = normalizeOrganisationLookupValue(query);
  if (!normalizedQuery) {
    return null;
  }

  const exactName = organisations.find(
    (organisation) => normalizeOrganisationLookupValue(organisation.organisationName) === normalizedQuery
  );
  if (exactName) {
    return exactName;
  }

  const exactCode = organisations.find(
    (organisation) => normalizeOrganisationLookupValue(organisation.organisationCode) === normalizedQuery
  );
  if (exactCode) {
    return exactCode;
  }

  const prefixNameMatches = organisations.filter((organisation) =>
    normalizeOrganisationLookupValue(organisation.organisationName).startsWith(normalizedQuery)
  );
  if (prefixNameMatches.length === 1) {
    return prefixNameMatches[0];
  }

  const containingNameMatches = organisations.filter((organisation) =>
    normalizeOrganisationLookupValue(organisation.organisationName).includes(normalizedQuery)
  );
  if (containingNameMatches.length === 1) {
    return containingNameMatches[0];
  }

  return organisations.length === 1 ? organisations[0] : null;
}

export default function LoginForm() {
  const router = useRouter();

  const [authMode, setAuthMode] = useState<AuthMode>("admin");

  const [mobileNumber, setMobileNumber] = useState("");
  const [password, setPassword] = useState("");
  const [mobileTouched, setMobileTouched] = useState(false);
  const [adminLoading, setAdminLoading] = useState(false);

  const [employeeEmail, setEmployeeEmail] = useState("");
  const [employeePassword, setEmployeePassword] = useState("");
  const [employeeLoginLoading, setEmployeeLoginLoading] = useState(false);

  const [showEmployeeActivation, setShowEmployeeActivation] = useState(false);
  const [activationOrganisationQuery, setActivationOrganisationQuery] = useState("");
  const [activationOrganisationCode, setActivationOrganisationCode] = useState("");
  const [selectedActivationOrganisation, setSelectedActivationOrganisation] =
    useState<EmployeeLoginOrganisation | null>(null);
  const [activationEmployeeId, setActivationEmployeeId] = useState("");
  const [activationEmail, setActivationEmail] = useState("");
  const [activationOtpId, setActivationOtpId] = useState("");
  const [activationOtpCode, setActivationOtpCode] = useState("");
  const [activationPassword, setActivationPassword] = useState("");
  const [activationLoading, setActivationLoading] = useState(false);
  const [activationSuggestions, setActivationSuggestions] = useState<EmployeeLoginOrganisation[]>([]);
  const [activationSuggestionsLoading, setActivationSuggestionsLoading] = useState(false);
  const [activationSuggestionsError, setActivationSuggestionsError] = useState("");

  const [dialog, setDialog] = useState<DialogState>({ open: false });

  const mobileHint = mobileTouched ? getMobileHint(mobileNumber) : null;
  const activationPasswordHint =
    activationPassword.trim() && getPasswordValidationMessage(activationPassword, "New password");

  function closeDialog() {
    setDialog({ open: false });
  }

  function openError(message: string, title = "Sign-in failed") {
    setDialog({ open: true, kind: "error", title, message });
  }

  function openSuccess(message: string, title = "Success") {
    setDialog({ open: true, kind: "success", title, message });
  }

  function resetEmployeeActivation(resetEmail = false) {
    setActivationOrganisationQuery("");
    setActivationOrganisationCode("");
    setSelectedActivationOrganisation(null);
    setActivationEmployeeId("");
    if (resetEmail) {
      setActivationEmail("");
    }
    setActivationOtpId("");
    setActivationOtpCode("");
    setActivationPassword("");
    setActivationSuggestions([]);
    setActivationSuggestionsError("");
  }

  function selectActivationOrganisation(organisation: EmployeeLoginOrganisation) {
    setActivationOrganisationQuery(organisation.organisationName);
    setActivationOrganisationCode(organisation.organisationCode);
    setSelectedActivationOrganisation(organisation);
    setActivationSuggestions([]);
    setActivationSuggestionsError("");
    setActivationOtpId("");
  }

  useEffect(() => {
    if (!showEmployeeActivation) {
      return;
    }

    const query = activationOrganisationQuery.trim();
    if (query.length < 2) {
      setActivationSuggestions([]);
      setActivationSuggestionsError("");
      if (!query) {
        setActivationOrganisationCode("");
      }
      return;
    }

    if (
      selectedActivationOrganisation &&
      normalizeOrganisationLookupValue(selectedActivationOrganisation.organisationName) ===
        normalizeOrganisationLookupValue(query)
    ) {
      setActivationSuggestions([]);
      setActivationSuggestionsError("");
      return;
    }

    let cancelled = false;
    const timer = window.setTimeout(async () => {
      setActivationSuggestionsLoading(true);
      try {
        const response = await bffClient.auth.searchEmployeeOrganisations(query);
        if (cancelled) {
          return;
        }
        if (!response.ok) {
          setActivationSuggestions([]);
          setActivationSuggestionsError(response.message ?? "Unable to find matching organisations.");
          return;
        }
        setActivationSuggestions(response.organisations ?? []);
        setActivationSuggestionsError("");

        const matchedOrganisation = findBestOrganisationMatch(query, response.organisations ?? []);
        if (matchedOrganisation) {
          setSelectedActivationOrganisation((current) => current ?? matchedOrganisation);
          setActivationOrganisationCode((current) =>
            current || matchedOrganisation.organisationCode
          );
        }
      } catch (error) {
        if (cancelled) {
          return;
        }
        setActivationSuggestions([]);
        setActivationSuggestionsError(
          error instanceof Error ? error.message : "Unable to find matching organisations."
        );
      } finally {
        if (!cancelled) {
          setActivationSuggestionsLoading(false);
        }
      }
    }, 250);

    return () => {
      cancelled = true;
      window.clearTimeout(timer);
    };
  }, [activationOrganisationQuery, selectedActivationOrganisation, showEmployeeActivation]);

  async function onAdminSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setMobileTouched(true);

    const normalized = normalizeBangladeshMobile(mobileNumber);
    if (!normalized) {
      openError(getBangladeshMobileValidationMessage());
      return;
    }

    setAdminLoading(true);
    try {
      const response = await bffClient.auth.login({ mobileNumber: normalized, password });
      if (!response.ok) {
        openError(response.message ?? "Login failed. Please try again.");
        return;
      }
      openSuccess("Redirecting you to the dashboard...", "Signed in successfully");
      setTimeout(() => {
        router.replace("/");
        router.refresh();
      }, 1200);
    } catch (submitError) {
      openError(
        submitError instanceof Error
          ? submitError.message
          : "An unexpected error occurred. Please try again."
      );
    } finally {
      setAdminLoading(false);
    }
  }

  async function onEmployeeLoginSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (!employeeEmail.trim()) {
      openError("Email is required.", "Employee sign-in failed");
      return;
    }
    if (!employeePassword.trim()) {
      openError("Password is required.", "Employee sign-in failed");
      return;
    }

    setEmployeeLoginLoading(true);
    try {
      const response = await bffClient.auth.loginEmployee({
        email: employeeEmail.trim().toLowerCase(),
        password: employeePassword,
      });
      if (!response.ok) {
        openError(response.message ?? "Employee sign-in failed.", "Employee sign-in failed");
        return;
      }
      openSuccess("Redirecting you to your coverage portal...", "Signed in successfully");
      setTimeout(() => {
        router.replace("/employee");
        router.refresh();
      }, 1200);
    } catch (error) {
      openError(
        error instanceof Error ? error.message : "Employee sign-in failed.",
        "Employee sign-in failed"
      );
    } finally {
      setEmployeeLoginLoading(false);
    }
  }

  async function sendEmployeeActivationCode(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    let chosenOrganisation = selectedActivationOrganisation;
    if (!chosenOrganisation && activationOrganisationQuery.trim().length >= 2) {
      try {
        const response = await bffClient.auth.searchEmployeeOrganisations(activationOrganisationQuery.trim());
        if (response.ok) {
          const matchedOrganisation = findBestOrganisationMatch(
            activationOrganisationQuery,
            response.organisations ?? []
          );
          if (matchedOrganisation) {
            chosenOrganisation = matchedOrganisation;
            selectActivationOrganisation(matchedOrganisation);
          } else if ((response.organisations ?? []).length > 0) {
            chosenOrganisation = response.organisations![0];
            selectActivationOrganisation(chosenOrganisation);
          }
        }
      } catch {
        // Keep the existing validation error path below if the lookup fails.
      }
    }

    const organisationCode =
      chosenOrganisation?.organisationCode?.trim().toUpperCase() ||
      activationOrganisationCode.trim().toUpperCase();
    const organisationId = chosenOrganisation?.organisationId?.trim() ?? "";
    const employeeId = activationEmployeeId.trim();
    const email = activationEmail.trim().toLowerCase();

    if (!chosenOrganisation || (!organisationCode && !organisationId)) {
      openError("Select your organisation from the matching results first.", "Activation failed");
      return;
    }
    if (!employeeId) {
      openError("Employee ID is required.", "Activation failed");
      return;
    }
    if (!email) {
      openError("Work email is required.", "Activation failed");
      return;
    }

    setActivationLoading(true);
    try {
      const response = await bffClient.auth.activateEmployee({
        organisationId,
        organisationCode,
        employeeId,
        email,
      });
      if (!response.ok) {
        openError(response.message ?? "Unable to send verification code.", "Activation failed");
        return;
      }
      setActivationOrganisationQuery(chosenOrganisation.organisationName);
      setActivationOrganisationCode(organisationCode);
      setSelectedActivationOrganisation(chosenOrganisation);
      setActivationEmployeeId(employeeId);
      setActivationEmail(email);
      setActivationOtpId(response.otpId ?? "");
      openSuccess(
        response.message ?? "Verification code sent. Enter the code and set your password below.",
        "Verification code sent"
      );
    } catch (error) {
      openError(
        error instanceof Error ? error.message : "Unable to send verification code.",
        "Activation failed"
      );
    } finally {
      setActivationLoading(false);
    }
  }

  async function completeEmployeeActivation(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    if (!activationOtpId) {
      openError("Request a verification code first.", "Activation failed");
      return;
    }
    if (!activationOtpCode.trim()) {
      openError("Verification code is required.", "Activation failed");
      return;
    }

    const passwordError = getPasswordValidationMessage(activationPassword, "New password");
    if (passwordError) {
      openError(passwordError, "Activation failed");
      return;
    }

    setActivationLoading(true);
    try {
      const response = await bffClient.auth.completeEmployeeActivation({
        email: activationEmail.trim().toLowerCase(),
        otpId: activationOtpId,
        otpCode: activationOtpCode.trim(),
        newPassword: activationPassword,
      });
      if (!response.ok) {
        openError(response.message ?? "Unable to finish activation.", "Activation failed");
        return;
      }

      setEmployeeEmail(activationEmail.trim().toLowerCase());
      setEmployeePassword("");
      setShowEmployeeActivation(false);
      resetEmployeeActivation(false);
      openSuccess(
        response.message ?? "Password set successfully. You can now sign in with email and password.",
        "Activation complete"
      );
    } catch (error) {
      openError(
        error instanceof Error ? error.message : "Unable to finish activation.",
        "Activation failed"
      );
    } finally {
      setActivationLoading(false);
    }
  }

  const isError = dialog.open && dialog.kind === "error";
  const isSuccess = dialog.open && dialog.kind === "success";

  return (
    <>
      <Dialog open={dialog.open} onOpenChange={(open) => !open && closeDialog()}>
        <DialogContent showCloseButton={isError} className="sm:max-w-sm">
          {dialog.open ? (
            <>
              <DialogHeader className="items-center gap-3">
                {isSuccess ? (
                  <CheckCircle2 className="size-12 text-emerald-500" />
                ) : (
                  <XCircle className="size-12 text-destructive" />
                )}
                <DialogTitle className="text-center text-lg">{dialog.title}</DialogTitle>
                <DialogDescription className="text-center">{dialog.message}</DialogDescription>
              </DialogHeader>
              {isError ? (
                <DialogFooter className="sm:justify-center">
                  <Button variant="outline" onClick={closeDialog} className="w-full sm:w-auto">
                    Try again
                  </Button>
                </DialogFooter>
              ) : null}
            </>
          ) : null}
        </DialogContent>
      </Dialog>

      <div className="space-y-5">
        <div className="grid grid-cols-2 gap-2 rounded-2xl bg-muted/60 p-1">
          <button
            type="button"
            onClick={() => setAuthMode("admin")}
            className={`rounded-xl px-4 py-2 text-sm font-medium transition ${
              authMode === "admin"
                ? "bg-background text-foreground shadow-sm"
                : "text-muted-foreground hover:text-foreground"
            }`}
          >
            Admin Sign In
          </button>
          <button
            type="button"
            onClick={() => setAuthMode("employee")}
            className={`rounded-xl px-4 py-2 text-sm font-medium transition ${
              authMode === "employee"
                ? "bg-background text-foreground shadow-sm"
                : "text-muted-foreground hover:text-foreground"
            }`}
          >
            Employee Sign In
          </button>
        </div>

        {authMode === "admin" ? (
          <form onSubmit={onAdminSubmit} className="space-y-5">
            <div className="space-y-2">
              <label className="text-sm font-medium text-foreground" htmlFor="mobileNumber">
                Mobile Number
              </label>
              <input
                id="mobileNumber"
                type="tel"
                className={`auth-input${mobileHint ? " border-amber-500 focus:ring-amber-500" : ""}`}
                value={mobileNumber}
                onChange={(event) => setMobileNumber(event.target.value)}
                onBlur={(event) => {
                  setMobileTouched(true);
                  setMobileNumber(normalizeBangladeshMobileOrRaw(event.target.value));
                }}
                placeholder={BD_MOBILE_EXAMPLES}
                autoComplete="tel"
                required
              />
              {mobileHint ? <p className="text-xs text-amber-600">{mobileHint}</p> : null}
            </div>

            <div className="space-y-2">
              <label className="text-sm font-medium text-foreground" htmlFor="password">
                Password
              </label>
              <input
                id="password"
                type="password"
                className="auth-input"
                value={password}
                onChange={(event) => setPassword(event.target.value)}
                placeholder="Enter password"
                required
              />
            </div>

            <button type="submit" disabled={adminLoading} className="auth-submit">
              {adminLoading ? "Signing in..." : "Sign in"}
            </button>
          </form>
        ) : (
          <div className="space-y-5">
            <form onSubmit={onEmployeeLoginSubmit} className="space-y-5">
              <div className="space-y-2">
                <label className="text-sm font-medium text-foreground" htmlFor="employee-login-email">
                  Work Email
                </label>
                <input
                  id="employee-login-email"
                  type="email"
                  className="auth-input"
                  value={employeeEmail}
                  onChange={(event) => setEmployeeEmail(event.target.value)}
                  placeholder="name@company.com"
                  autoComplete="email"
                  required
                />
              </div>

              <div className="space-y-2">
                <label className="text-sm font-medium text-foreground" htmlFor="employee-login-password">
                  Password
                </label>
                <input
                  id="employee-login-password"
                  type="password"
                  className="auth-input"
                  value={employeePassword}
                  onChange={(event) => setEmployeePassword(event.target.value)}
                  placeholder="Enter password"
                  autoComplete="current-password"
                  required
                />
              </div>

              <button type="submit" disabled={employeeLoginLoading} className="auth-submit">
                {employeeLoginLoading ? "Signing in..." : "Sign in as employee"}
              </button>
            </form>

            <div className="rounded-2xl border border-border/70 bg-muted/30 p-4">
              <div className="flex items-start justify-between gap-3">
                <div className="space-y-1">
                  <h3 className="text-sm font-semibold text-foreground">First time employee?</h3>
                  <p className="text-sm text-muted-foreground">
                    Activate with your organisation name, employee ID, and work email.
                  </p>
                </div>
                <Button
                  type="button"
                  variant={showEmployeeActivation ? "outline" : "default"}
                  onClick={() => {
                    const next = !showEmployeeActivation;
                    setShowEmployeeActivation(next);
                    if (!next) {
                      resetEmployeeActivation(true);
                    }
                  }}
                >
                  {showEmployeeActivation ? "Hide" : "Activate"}
                </Button>
              </div>

              {showEmployeeActivation ? (
                <div className="mt-4 space-y-4">
                  <form onSubmit={sendEmployeeActivationCode} className="space-y-4">
                    <div className="space-y-2">
                      <label className="text-sm font-medium text-foreground" htmlFor="employee-organisation">
                        Organisation
                      </label>
                      <input
                        id="employee-organisation"
                        type="text"
                        className="auth-input"
                        value={activationOrganisationQuery}
                        onChange={(event) => {
                          setActivationOrganisationQuery(event.target.value);
                          setActivationOrganisationCode("");
                          setSelectedActivationOrganisation(null);
                          setActivationOtpId("");
                        }}
                        placeholder="Start typing your organisation name"
                        required
                      />
                      {selectedActivationOrganisation ? (
                        <p className="text-xs text-emerald-700">
                          Selected: {selectedActivationOrganisation.organisationName}
                        </p>
                      ) : null}
                      {activationSuggestionsLoading ? (
                        <p className="text-xs text-muted-foreground">Finding matching organisations...</p>
                      ) : null}
                      {activationSuggestionsError ? (
                        <p className="text-xs text-destructive">{activationSuggestionsError}</p>
                      ) : null}
                      {!activationSuggestionsLoading &&
                      activationOrganisationQuery.trim().length >= 2 &&
                      activationSuggestions.length > 0 ? (
                        <div className="rounded-2xl border border-border/70 bg-background">
                          {activationSuggestions.map((organisation) => (
                            <button
                              key={organisation.organisationId}
                              type="button"
                              className="flex w-full items-center justify-between px-4 py-3 text-left text-sm transition hover:bg-muted/50"
                              onMouseDown={(event) => {
                                event.preventDefault();
                                selectActivationOrganisation(organisation);
                              }}
                            >
                              <span className="font-medium text-foreground">{organisation.organisationName}</span>
                              <span className="text-xs text-muted-foreground">Select</span>
                            </button>
                          ))}
                        </div>
                      ) : null}
                      {!selectedActivationOrganisation &&
                      !activationSuggestionsLoading &&
                      activationOrganisationQuery.trim().length >= 2 &&
                      activationSuggestions.length === 0 &&
                      !activationSuggestionsError ? (
                        <p className="text-xs text-muted-foreground">No matching organisation found yet.</p>
                      ) : null}
                    </div>

                    <div className="space-y-2">
                      <label className="text-sm font-medium text-foreground" htmlFor="employee-id">
                        Employee ID
                      </label>
                      <input
                        id="employee-id"
                        type="text"
                        className="auth-input"
                        value={activationEmployeeId}
                        onChange={(event) => setActivationEmployeeId(event.target.value)}
                        placeholder="Enter your employee ID"
                        required
                      />
                    </div>

                    <div className="space-y-2">
                      <label className="text-sm font-medium text-foreground" htmlFor="employee-email">
                        Work Email
                      </label>
                      <input
                        id="employee-email"
                        type="email"
                        className="auth-input"
                        value={activationEmail}
                        onChange={(event) => setActivationEmail(event.target.value)}
                        placeholder="name@company.com"
                        autoComplete="email"
                        required
                      />
                    </div>

                    <button
                      type="submit"
                      disabled={activationLoading || !selectedActivationOrganisation}
                      className="auth-submit"
                    >
                      {activationLoading
                        ? "Sending code..."
                        : activationOtpId
                          ? "Resend verification code"
                          : "Send verification code"}
                    </button>
                  </form>

                  {activationOtpId ? (
                    <form
                      onSubmit={completeEmployeeActivation}
                      className="space-y-4 rounded-2xl border border-border/70 bg-background/80 p-4"
                    >
                      <div className="space-y-2">
                        <label className="text-sm font-medium text-foreground" htmlFor="employee-otp">
                          Email Verification Code
                        </label>
                        <input
                          id="employee-otp"
                          type="text"
                          className="auth-input"
                          value={activationOtpCode}
                          onChange={(event) => setActivationOtpCode(event.target.value)}
                          placeholder="Enter the code from your email"
                          required
                        />
                      </div>

                      <div className="space-y-2">
                        <label className="text-sm font-medium text-foreground" htmlFor="employee-new-password">
                          New Password
                        </label>
                        <input
                          id="employee-new-password"
                          type="password"
                          className={`auth-input${activationPasswordHint ? " border-amber-500 focus:ring-amber-500" : ""}`}
                          value={activationPassword}
                          onChange={(event) => setActivationPassword(event.target.value)}
                          placeholder="Set your new password"
                          required
                        />
                        <p className={`text-xs ${activationPasswordHint ? "text-amber-600" : "text-muted-foreground"}`}>
                          {activationPasswordHint ?? PASSWORD_REQUIREMENTS_HINT}
                        </p>
                      </div>

                      <button type="submit" disabled={activationLoading} className="auth-submit">
                        {activationLoading ? "Setting password..." : "Verify and set password"}
                      </button>
                    </form>
                  ) : null}
                </div>
              ) : null}
            </div>
          </div>
        )}
      </div>
    </>
  );
}
