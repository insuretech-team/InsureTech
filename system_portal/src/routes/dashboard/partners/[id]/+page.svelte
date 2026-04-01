<script lang="ts">
	import { Badge } from '$lib/components/ui/badge';
	import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '$lib/components/ui/card';
	import PageHeader from '$lib/components/system/page-header.svelte';
	import { formatDate, humanizeStatus, statusTone } from '$lib/system/format';

	let { data } = $props();
</script>

<div class="space-y-8">
	<PageHeader
		title={data.partner.name}
		description="Partner profile normalized from the backend partner detail response."
		meta={data.partner.id}
	/>

	<div class="grid gap-6 lg:grid-cols-3">
		<Card class="rounded-[28px] border-white/60 bg-white/82">
			<CardHeader>
				<CardTitle>Identity</CardTitle>
				<CardDescription>Primary partner classification.</CardDescription>
			</CardHeader>
			<CardContent class="space-y-4">
				<div>
					<p class="text-sm text-slate-500">Type</p>
					<p class="mt-2 text-lg font-semibold text-slate-950">{data.partner.type}</p>
				</div>
				<div>
					<p class="text-sm text-slate-500">Category</p>
					<p class="mt-2 text-lg font-semibold text-slate-950">{data.partner.category}</p>
				</div>
				<div>
					<p class="text-sm text-slate-500">Status</p>
					<div class="mt-2">
						<Badge variant={statusTone(data.partner.status)}>{humanizeStatus(data.partner.status)}</Badge>
					</div>
				</div>
			</CardContent>
		</Card>

		<Card class="rounded-[28px] border-white/60 bg-white/82">
			<CardHeader>
				<CardTitle>Contact</CardTitle>
				<CardDescription>Direct communication fields exposed by the partner record.</CardDescription>
			</CardHeader>
			<CardContent class="space-y-4">
				<div>
					<p class="text-sm text-slate-500">Email</p>
					<p class="mt-2 text-lg font-semibold text-slate-950">{data.partner.email || 'Not provided'}</p>
				</div>
				<div>
					<p class="text-sm text-slate-500">Phone</p>
					<p class="mt-2 text-lg font-semibold text-slate-950">{data.partner.phone || 'Not provided'}</p>
				</div>
				<div>
					<p class="text-sm text-slate-500">Address</p>
					<p class="mt-2 text-sm leading-6 text-slate-700">{data.partner.address || 'Not provided'}</p>
				</div>
			</CardContent>
		</Card>

		<Card class="rounded-[28px] border-white/60 bg-white/82">
			<CardHeader>
				<CardTitle>Lifecycle</CardTitle>
				<CardDescription>Enrollment and backend identifiers.</CardDescription>
			</CardHeader>
			<CardContent class="space-y-4">
				<div>
					<p class="text-sm text-slate-500">Joined</p>
					<p class="mt-2 text-lg font-semibold text-slate-950">{formatDate(data.partner.joinedAt)}</p>
				</div>
				<div>
					<p class="text-sm text-slate-500">Partner ID</p>
					<p class="mt-2 break-all text-sm font-medium text-slate-700">{data.partner.id}</p>
				</div>
			</CardContent>
		</Card>
	</div>
</div>
