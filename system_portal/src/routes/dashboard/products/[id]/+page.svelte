<script lang="ts">
	import { Badge } from '$lib/components/ui/badge';
	import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '$lib/components/ui/card';
	import PageHeader from '$lib/components/system/page-header.svelte';
	import { formatCurrency, formatDate, humanizeStatus, statusTone } from '$lib/system/format';

	let { data } = $props();
</script>

<div class="space-y-8">
	<PageHeader
		title={data.product.name}
		description={data.product.description || 'Backend product detail returned from the generated SDK.'}
		meta={data.product.code || data.product.id}
	/>

	<div class="grid gap-6 lg:grid-cols-3">
		<Card class="rounded-[28px] border-white/60 bg-white/82">
			<CardHeader>
				<CardTitle>Commercial terms</CardTitle>
				<CardDescription>Premium and coverage configuration.</CardDescription>
			</CardHeader>
			<CardContent class="space-y-4">
				<div>
					<p class="text-sm text-slate-500">Base premium</p>
					<p class="mt-2 text-3xl font-semibold text-slate-950">{formatCurrency(data.product.basePremium)}</p>
				</div>
				<div>
					<p class="text-sm text-slate-500">Coverage band</p>
					<p class="mt-2 text-lg font-semibold text-slate-950">
						{formatCurrency(data.product.minSumInsured)} to {formatCurrency(data.product.maxSumInsured)}
					</p>
				</div>
			</CardContent>
		</Card>

		<Card class="rounded-[28px] border-white/60 bg-white/82">
			<CardHeader>
				<CardTitle>Lifecycle</CardTitle>
				<CardDescription>Status and timing metadata.</CardDescription>
			</CardHeader>
			<CardContent class="space-y-4">
				<div>
					<p class="text-sm text-slate-500">Status</p>
					<div class="mt-2">
						<Badge variant={statusTone(data.product.status)}>{humanizeStatus(data.product.status)}</Badge>
					</div>
				</div>
				<div>
					<p class="text-sm text-slate-500">Created</p>
					<p class="mt-2 text-lg font-semibold text-slate-950">{formatDate(data.product.createdAt)}</p>
				</div>
				<div>
					<p class="text-sm text-slate-500">Last updated</p>
					<p class="mt-2 text-lg font-semibold text-slate-950">{formatDate(data.product.updatedAt)}</p>
				</div>
			</CardContent>
		</Card>

		<Card class="rounded-[28px] border-white/60 bg-white/82">
			<CardHeader>
				<CardTitle>Classification</CardTitle>
				<CardDescription>Product identifiers visible in this portal.</CardDescription>
			</CardHeader>
			<CardContent class="space-y-4">
				<div>
					<p class="text-sm text-slate-500">Category</p>
					<p class="mt-2 text-lg font-semibold text-slate-950">{data.product.category}</p>
				</div>
				<div>
					<p class="text-sm text-slate-500">Code</p>
					<p class="mt-2 text-lg font-semibold text-slate-950">{data.product.code || 'Not provided'}</p>
				</div>
				<div>
					<p class="text-sm text-slate-500">Product ID</p>
					<p class="mt-2 break-all text-sm font-medium text-slate-700">{data.product.id}</p>
				</div>
			</CardContent>
		</Card>
	</div>
</div>
