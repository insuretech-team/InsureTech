import { NextResponse } from "next/server";

import { unwrapGateway, type GatewayResponse } from "@lib/sdk/shared";

function getApiBaseUrl(): string {
  return (
    process.env.INSURETECH_API_BASE_URL ??
    process.env.NEXT_PUBLIC_INSURETECH_API_BASE_URL ??
    "http://localhost:8080"
  );
}

export async function GET(request: Request) {
  const { searchParams } = new URL(request.url);
  const query = searchParams.get("q")?.trim() ?? "";

  if (query.length < 2) {
    return NextResponse.json({ ok: true, organisations: [] });
  }

  try {
    const params = new URLSearchParams({
      q: query,
      page_size: "8",
    });

    const response = await fetch(`${getApiBaseUrl()}/v1/b2b/organisations:employee-login?${params}`, {
      method: "GET",
      cache: "no-store",
    });
    const body = (await response.json().catch(() => null)) as GatewayResponse<{
      organisations?: Array<{
        organisation_id?: string;
        organisation_name?: string;
        organisation_code?: string;
      }>;
    }> | null;

    if (!body) {
      return NextResponse.json(
        { ok: false, message: "Unable to load organisations", organisations: [] },
        { status: 502 }
      );
    }

    const unwrapped = unwrapGateway(body, response.status);
    if (!unwrapped.ok) {
      return NextResponse.json(
        { ok: false, message: unwrapped.message, organisations: [] },
        { status: unwrapped.status }
      );
    }

    return NextResponse.json({
      ok: true,
      organisations: (unwrapped.data.organisations ?? []).map((item) => ({
        organisationId: item.organisation_id ?? "",
        organisationName: item.organisation_name ?? "",
        organisationCode: item.organisation_code ?? "",
      })),
    });
  } catch (error) {
    return NextResponse.json(
      {
        ok: false,
        message: error instanceof Error ? error.message : "Unable to load organisations",
        organisations: [],
      },
      { status: 502 }
    );
  }
}
