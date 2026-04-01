import { fallbackClaims, fallbackProposals } from "@/lib/mock-data";
import type { PortalOverview } from "@/lib/types";
import { directHttp, fetchCurrentSession } from "@/lib/server/insuretech";
import { fallbackContext, mapClaim, mapConfig, mapInsurer, mapProduct, mapProposal } from "@/lib/server/mappers";

function preferPragatiInsurer<T extends { id: string; name: string }>(items: T[]) {
  return items.find((entry) => entry.name.toLowerCase().includes("pragati")) ?? items[0];
}

export async function loadContext(request: Request, insurerId?: string) {
  const session = await fetchCurrentSession(request);
  if (!session) return null;

  const listResponse = await directHttp(request, "/v1/insurers?page_size=50", { session });
  const insurers = listResponse.ok
    ? ((listResponse.data.insurers as unknown[]) ?? []).map((entry) => mapInsurer(entry)).filter((entry) => entry.id)
    : fallbackContext().insurers;

  const selectedId = insurerId || preferPragatiInsurer(insurers)?.id || fallbackContext().currentInsurer.id;

  let currentInsurer = insurers.find((entry) => entry.id === selectedId) ?? fallbackContext().currentInsurer;
  let config = fallbackContext().config;
  let products = fallbackContext().products;

  if (selectedId) {
    const [detailResponse, productsResponse] = await Promise.all([
      directHttp(request, `/v1/insurers/${selectedId}`, { session }),
      directHttp(request, `/v1/insurers/${selectedId}/products?page_size=100`, { session }),
    ]);

    if (detailResponse.ok) {
      currentInsurer = mapInsurer(detailResponse.data.insurer ?? detailResponse.data);
      config = mapConfig(selectedId, detailResponse.data.config);
    }

    if (productsResponse.ok) {
      const liveProducts = ((productsResponse.data.insurer_products as unknown[]) ?? [])
        .map((entry) => mapProduct(entry))
        .filter((entry) => entry.id);
      if (liveProducts.length) products = liveProducts;
    }
  }

  const source = listResponse.ok ? "live" : "fallback";
  return { session, insurers, currentInsurer, config, products, source } as const;
}

export async function loadOverview(request: Request, insurerId?: string): Promise<PortalOverview | null> {
  const context = await loadContext(request, insurerId);
  if (!context) return null;

  const [proposalsResponse, claimsResponse] = await Promise.all([
    directHttp(
      request,
      `/v1/insurance-proposals?insurer_id=${encodeURIComponent(context.currentInsurer.id)}&page_size=50`,
      { session: context.session },
    ),
    directHttp(request, "/v1/claims?page_size=50", { session: context.session }),
  ]);

  const proposals = proposalsResponse.ok
    ? (((proposalsResponse.data.proposals as unknown[]) ?? []).map((entry) => mapProposal(entry, "live")))
    : fallbackProposals;

  const claims = claimsResponse.ok
    ? (((claimsResponse.data.claims as unknown[]) ?? []).map((entry) => mapClaim(entry, "live")))
    : fallbackClaims;

  const source: PortalOverview["source"] =
    proposalsResponse.ok && claimsResponse.ok
      ? context.source
      : proposalsResponse.ok || claimsResponse.ok
        ? "mixed"
        : "fallback";

  return {
    currentInsurer: context.currentInsurer,
    insurers: context.insurers,
    products: context.products,
    proposals,
    claims,
    config: context.config,
    metrics: {
      insurerCount: context.insurers.length,
      productCount: context.products.length,
      proposalCount: proposals.length,
      approvedProposalCount: proposals.filter((item) => item.status === "Approved").length,
      claimCount: claims.length,
      settledClaimCount: claims.filter((item) => item.status === "Settled").length,
      requestedClaimCount: claims.filter((item) => item.status === "Pending Documents").length,
      underReviewClaimCount: claims.filter((item) => item.status === "Under Review").length,
    },
    source,
  };
}
