<script lang="ts">
	import { Badge } from '$lib/components/ui/badge';
	import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '$lib/components/ui/card';
	import * as Table from '$lib/components/ui/table';
	import MetricCard from '$lib/components/system/metric-card.svelte';
	import PageHeader from '$lib/components/system/page-header.svelte';
	import { formatCurrency, formatDate, humanizeStatus, statusTone } from '$lib/system/format';

	let { data } = $props();
</script>

<div class="space-y-8">
	<PageHeader
		title="System operations overview"
		description="Live operational snapshot for the InsureTech platform, wired to products, partners, policies, claims, and tenant services through the generated local SDK and backend APIs."
		meta="Session-backed system console"
	/>

	<div class="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
		{#each data.overview.metrics as metric}
			<MetricCard {metric} />
		{/each}
	</div>

	<div class="grid gap-6 xl:grid-cols-[1.2fr_0.8fr]">
		<Card class="rounded-[28px] border-white/60 bg-white/82 backdrop-blur">
			<CardHeader>
				<CardTitle>Operational pulse</CardTitle>
				<CardDescription>Counts and portfolio values visible from the current backend session.</CardDescription>
			</CardHeader>
			<CardContent class="grid gap-4 sm:grid-cols-2">
				<div class="rounded-3xl bg-slate-50 p-5">
					<p class="text-sm text-slate-500">Live premium in view</p>
					<p class="mt-3 text-3xl font-semibold text-slate-950">{formatCurrency(data.breakdown.livePremium)}</p>
					<p class="mt-2 text-sm text-slate-600">Average premium {formatCurrency(data.breakdown.averagePolicyPremium)}</p>
				</div>
				<div class="rounded-3xl bg-slate-50 p-5">
					<p class="text-sm text-slate-500">Claim exposure</p>
					<p class="mt-3 text-3xl font-semibold text-slate-950">{formatCurrency(data.breakdown.claimExposure)}</p>
					<p class="mt-2 text-sm text-slate-600">{data.claimExposure.open} claims not yet settled</p>
				</div>
				<div class="rounded-3xl bg-slate-50 p-5">
					<p class="text-sm text-slate-500">Partner network mix</p>
					<p class="mt-3 text-3xl font-semibold text-slate-950">{data.breakdown.partnerCount}</p>
					<p class="mt-2 text-sm text-slate-600">
						Life {data.partnerMix.life} • Non-life {data.partnerMix.nonLife} • Other {data.partnerMix.other}
					</p>
				</div>
				<div class="rounded-3xl bg-slate-50 p-5">
					<p class="text-sm text-slate-500">Tenant readiness</p>
					<p class="mt-3 text-3xl font-semibold text-slate-950">{data.tenantSummary.active}</p>
					<p class="mt-2 text-sm text-slate-600">Last tenant created {formatDate(data.tenantSummary.lastCreated)}</p>
				</div>
			</CardContent>
		</Card>

		<Card class="rounded-[28px] border-white/60 bg-white/82 backdrop-blur">
			<CardHeader>
				<CardTitle>Policy health</CardTitle>
				<CardDescription>Current state distribution from policy records.</CardDescription>
			</CardHeader>
			<CardContent class="space-y-4">
				<div class="rounded-3xl border border-emerald-100 bg-emerald-50 p-5">
					<p class="text-sm text-emerald-700">Active policies</p>
					<p class="mt-2 text-3xl font-semibold text-emerald-950">{data.policyHealth.active}</p>
				</div>
				<div class="rounded-3xl border border-amber-100 bg-amber-50 p-5">
					<p class="text-sm text-amber-700">Lapsed policies</p>
					<p class="mt-2 text-3xl font-semibold text-amber-950">{data.policyHealth.lapsed}</p>
				</div>
				<div class="rounded-3xl border border-rose-100 bg-rose-50 p-5">
					<p class="text-sm text-rose-700">Cancelled policies</p>
					<p class="mt-2 text-3xl font-semibold text-rose-950">{data.policyHealth.cancelled}</p>
				</div>
			</CardContent>
		</Card>
	</div>

	<div class="grid gap-6 xl:grid-cols-2">
		<Card class="rounded-[28px] border-white/60 bg-white/82 backdrop-blur">
			<CardHeader>
				<CardTitle>Latest products</CardTitle>
				<CardDescription>Catalog records returned directly from the product service.</CardDescription>
			</CardHeader>
			<CardContent>
				<Table.Root>
					<Table.Header>
						<Table.Row>
							<Table.Head>Product</Table.Head>
							<Table.Head>Status</Table.Head>
							<Table.Head>Base premium</Table.Head>
						</Table.Row>
					</Table.Header>
					<Table.Body>
						{#each data.overview.products.slice(0, 5) as product}
							<Table.Row>
								<Table.Cell>
									<a href={`/dashboard/products/${product.id}`} class="font-medium text-primary hover:text-accent">
										{product.name}
									</a>
									<div class="text-xs text-slate-500">{product.code || product.id}</div>
								</Table.Cell>
								<Table.Cell>
									<Badge variant={statusTone(product.status)}>{humanizeStatus(product.status)}</Badge>
								</Table.Cell>
								<Table.Cell>{formatCurrency(product.basePremium)}</Table.Cell>
							</Table.Row>
						{/each}
					</Table.Body>
				</Table.Root>
			</CardContent>
		</Card>

		<Card class="rounded-[28px] border-white/60 bg-white/82 backdrop-blur">
			<CardHeader>
				<CardTitle>Newest claims</CardTitle>
				<CardDescription>Recent cases from the claims API endpoint.</CardDescription>
			</CardHeader>
			<CardContent>
				<Table.Root>
					<Table.Header>
						<Table.Row>
							<Table.Head>Claim</Table.Head>
							<Table.Head>Status</Table.Head>
							<Table.Head>Amount</Table.Head>
						</Table.Row>
					</Table.Header>
					<Table.Body>
						{#each data.overview.claims.slice(0, 5) as claim}
							<Table.Row>
								<Table.Cell>
									<div class="font-medium text-slate-900">{claim.claimNumber}</div>
									<div class="text-xs text-slate-500">{claim.claimantName}</div>
								</Table.Cell>
								<Table.Cell>
									<Badge variant={statusTone(claim.status)}>{humanizeStatus(claim.status)}</Badge>
								</Table.Cell>
								<Table.Cell>{formatCurrency(claim.amount)}</Table.Cell>
							</Table.Row>
						{/each}
					</Table.Body>
				</Table.Root>
			</CardContent>
		</Card>
	</div>
</div>
