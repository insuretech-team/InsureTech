/**
 * POST /api/organisations/[id]/admins
 *
 * Two modes:
 *
 * Mode A — Create new B2B admin user (POST /api/organisations/[id]/admins):
 *   body: { email, password, mobileNumber, fullName? }
 *   → forwards to backend POST /v1/b2b/organisations/{id}/admins
 *     which calls AuthN.RegisterEmailUser + AssignRole(B2B_ORG_ADMIN)
 *
 * Mode B — Assign an existing user/member as admin
 *   POST /api/organisations/[id]/admins?action=assign
 *   body: { memberId } or { userId, temporaryPassword }
 *   → forwards to backend POST /v1/b2b/organisations/{id}/admins:assign
 */
import { NextResponse } from "next/server";
import { makeDirectHttp, makeSdkClient } from "@lib/sdk/b2b-sdk-client";
import { sdkErrorMessage } from "@lib/sdk/api-helpers";
import { resolvePortalHeaders } from "@lib/sdk/session-headers";
import {
  getBangladeshMobileValidationMessage,
  normalizeBangladeshMobile,
} from "@/src/lib/utils/bd-mobile";
import { getPasswordValidationMessage } from "@/src/lib/utils/password";

type RouteContext = { params: Promise<{ id: string }> };

export async function POST(request: Request, { params }: RouteContext) {
  const { id } = await params;
  if (!id) {
    return NextResponse.json({ ok: false, message: "organisation_id required" }, { status: 400 });
  }

  try {
    const hdrs = await resolvePortalHeaders(request);
    if (!hdrs) {
      return NextResponse.json({ ok: false, message: "Unauthorized" }, { status: 401 });
    }
    const body = (await request.json()) as Record<string, unknown>;
    const action = new URL(request.url).searchParams.get("action");

    if (action === "assign") {
      const memberId = typeof body.memberId === "string" ? body.memberId.trim() : "";
      if (memberId) {
        const result = await makeSdkClient(request, hdrs).assignOrgAdmin({
          path: { organisation_id: id },
          body: { organisation_id: id, member_id: memberId },
        });
        if (!result.response.ok) {
          return NextResponse.json({ ok: false, message: sdkErrorMessage(result) }, { status: result.response.status });
        }
        return NextResponse.json({ ok: true, message: "Admin assigned successfully" });
      }

      const userId = typeof body.userId === "string" ? body.userId.trim() : "";
      const temporaryPassword = typeof body.temporaryPassword === "string" ? body.temporaryPassword : "";
      if (!userId) {
        return NextResponse.json({ ok: false, message: "memberId or userId is required" }, { status: 400 });
      }
      if (!temporaryPassword.trim()) {
        return NextResponse.json({ ok: false, message: "Temporary password is required" }, { status: 400 });
      }
      const passwordError = getPasswordValidationMessage(temporaryPassword, "Temporary password");
      if (passwordError) {
        return NextResponse.json({ ok: false, message: passwordError }, { status: 400 });
      }

      const result = await makeDirectHttp(request, hdrs).post(`/v1/b2b/organisations/${id}/admins:assign`, {
        userId,
        temporaryPassword,
      });
      return NextResponse.json(
        { ok: result.ok, message: result.data?.message ?? (result.ok ? "Admin assigned successfully" : "Failed to assign admin") },
        { status: result.ok ? 200 : result.status }
      );
    }

    // Mode A: create new admin user — requires email + password + mobileNumber
    const email = typeof body.email === "string" ? body.email.trim() : "";
    const password = typeof body.password === "string" ? body.password : "";
    const mobileNumber = typeof body.mobileNumber === "string" ? body.mobileNumber.trim() : "";

    if (!email || !password || !mobileNumber) {
      return NextResponse.json(
        { ok: false, message: "email, password, and mobileNumber are required to create a new admin" },
        { status: 400 }
      );
    }
    const normalizedMobileNumber = normalizeBangladeshMobile(mobileNumber);
    if (!normalizedMobileNumber) {
      return NextResponse.json(
        { ok: false, message: getBangladeshMobileValidationMessage("Admin mobile number") },
        { status: 400 }
      );
    }

    // Backend assignOrgAdminPayload struct uses camelCase JSON tags:
    // json:"email", json:"password", json:"mobileNumber", json:"fullName"
    const result = await makeDirectHttp(request, hdrs).post(`/v1/b2b/organisations/${id}/admins`, {
      email,
      password,
      mobileNumber: normalizedMobileNumber,
      fullName: typeof body.fullName === "string" ? body.fullName.trim() : undefined,
    });

    return NextResponse.json(
      { ok: result.ok, message: result.data?.message ?? (result.ok ? "Admin created" : "Failed to create admin") },
      { status: result.ok ? 201 : result.status }
    );
  } catch (err) {
    return NextResponse.json({ ok: false, message: err instanceof Error ? err.message : "Error" }, { status: 502 });
  }
}
