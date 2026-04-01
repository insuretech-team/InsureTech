<script lang="ts">
	import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '$lib/components/ui/card';
	import PageHeader from '$lib/components/system/page-header.svelte';
	import { formatCurrency, formatDate } from '$lib/system/format';

	let { data } = $props();
</script>

<div class="space-y-8">
	<PageHeader
		title="Operational analytics"
		description="A compact analytics surface built from the same backend-backed datasets that power the rest of the console."
		meta="Derived from live route data"
	/>

	<div class="grid gap-6 lg:grid-cols-3">
		<Card class="rounded-[28px] border-white/60 bg-white/82">
			<CardHeader>
				<CardTitle>Policy health mix</CardTitle>
				<CardDescription>Active versus degraded policy states.</CardDescription>
			</CardHeader>
			<CardContent class="space-y-4">
				<div class="rounded-3xl bg-emerald-50 p-5">
					<p class="text-sm text-emerald-700">Active</p>
					<p class="mt-2 text-3xl font-semibold text-emerald-950">{data.policyHealth.active}</p>
				</div>
				<div class="rounded-3xl bg-amber-50 p-5">
					<p class="text-sm text-amber-700">Lapsed</p>
					<p class="mt-2 text-3xl font-semibold text-amber-950">{data.policyHealth.lapsed}</p>
				</div>
				<div class="rounded-3xl bg-rose-50 p-5">
					<p class="text-sm text-rose-700">Cancelled</p>
					<p class="mt-2 text-3xl font-semibold text-rose-950">{data.policyHealth.cancelled}</p>
				</div>
			</CardContent>
		</Card>

		<Card class="rounded-[28px] border-white/60 bg-white/82">
			<CardHeader>
				<CardTitle>Partner composition</CardTitle>
				<CardDescription>Category split based on partner type signals.</CardDescription>
			</CardHeader>
			<CardContent class="space-y-4">
				<div class="rounded-3xl bg-slate-50 p-5">
					<p class="text-sm text-slate-500">Life partner entities</p>
					<p class="mt-2 text-3xl font-semibold text-slate-950">{data.partnerMix.life}</p>
				</div>
				<div class="rounded-3xl bg-slate-50 p-5">
					<p class="text-sm text-slate-500">Non-life partner entities</p>
					<p class="mt-2 text-3xl font-semibold text-slate-950">{data.partnerMix.nonLife}</p>
				</div>
				<div class="rounded-3xl bg-slate-50 p-5">
					<p class="text-sm text-slate-500">Other network entities</p>
					<p class="mt-2 text-3xl font-semibold text-slate-950">{data.partnerMix.other}</p>
				</div>
			</CardContent>
		</Card>

		<Card class="rounded-[28px] border-white/60 bg-white/82">
			<CardHeader>
				<CardTitle>Claim pressure</CardTitle>
				<CardDescription>Claim count and monetary exposure from the visible dataset.</CardDescription>
			</CardHeader>
			<CardContent class="space-y-4">
				<div class="rounded-3xl bg-slate-50 p-5">
					<p class="text-sm text-slate-500">Open or unsettled claims</p>
					<p class="mt-2 text-3xl font-semibold text-slate-950">{data.claimExposure.open}</p>
				</div>
				<div class="rounded-3xl bg-slate-50 p-5">
					<p class="text-sm text-slate-500">Total claim exposure</p>
					<p class="mt-2 text-3xl font-semibold text-slate-950">{formatCurrency(data.claimExposure.totalAmount)}</p>
				</div>
				<div class="rounded-3xl bg-slate-50 p-5">
					<p class="text-sm text-slate-500">Latest incident in feed</p>
					<p class="mt-2 text-xl font-semibold text-slate-950">{formatDate(data.claimExposure.latestIncident)}</p>
				</div>
			</CardContent>
		</Card>
	</div>
</div>
