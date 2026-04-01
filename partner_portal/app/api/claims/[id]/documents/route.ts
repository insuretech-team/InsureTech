import { NextResponse } from "next/server";
import { makeSdkClient } from "@lib/sdk/partner-sdk-client";
import { resolvePortalHeaders } from "@lib/sdk/session-headers";
import { sdkErrorMessage, badRequest } from "@lib/sdk/api-helpers";

const ALLOWED_FILE_TYPES = ["image/jpeg", "image/png", "application/pdf"];
const MAX_FILE_SIZE = 5 * 1024 * 1024; // 5MB

/** POST /api/claims/[id]/documents - Upload claim document */
export async function POST(
  request: Request,
  { params }: { params: { id: string } }
) {
  const hdrs = await resolvePortalHeaders(request);
  if (!hdrs) return NextResponse.json({ ok: false, message: "Unauthorized" }, { status: 401 });

  let formData: FormData;
  try {
    formData = await request.formData();
  } catch {
    return badRequest("Invalid form data");
  }

  const file = formData.get("file") as File | null;
  if (!file) {
    return badRequest("file is required");
  }

  // Validate file type
  if (!ALLOWED_FILE_TYPES.includes(file.type)) {
    return NextResponse.json(
      {
        ok: false,
        message: `Invalid file type. Allowed types: JPEG, PNG, PDF. Received: ${file.type}`,
      },
      { status: 400 }
    );
  }

  // Validate file size
  if (file.size > MAX_FILE_SIZE) {
    return NextResponse.json(
      {
        ok: false,
        message: `File size exceeds maximum allowed size of 5MB. File size: ${(file.size / 1024 / 1024).toFixed(2)}MB`,
      },
      { status: 400 }
    );
  }

  const documentType = formData.get("document_type") as string | null;
  if (!documentType) {
    return badRequest("document_type is required");
  }

  const sdk = makeSdkClient(request, hdrs);
  const result = await sdk.uploadClaimDocument({
    path: { claim_id: params.id },
    body: {
      file,
      document_type: documentType,
    },
  });

  if (!result.response.ok) {
    return NextResponse.json(
      { ok: false, message: sdkErrorMessage(result) },
      { status: result.response.status }
    );
  }

  return NextResponse.json(
    { ok: true, message: "Document uploaded successfully", data: result.data },
    { status: result.response.status || 201 }
  );
}
