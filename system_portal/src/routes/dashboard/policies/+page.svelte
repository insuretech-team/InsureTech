<script lang="ts">
	import { Badge } from '$lib/components/ui/badge';
	import { Card, CardContent } from '$lib/components/ui/card';
	import * as Table from '$lib/components/ui/table';
	import PageHeader from '$lib/components/system/page-header.svelte';
	import { formatCurrency, formatDate, humanizeStatus, statusTone } from '$lib/system/format';

	let { data } = $props();
</script>

<div class="space-y-8">
	<PageHeader
		title="Policies"
		description="Live policy rows fetched from the policy API. This replaces the older mock policy dashboard."
		meta={`${data.policies.length} visible policies`}
	/>

	<Card class="rounded-[28px] border-white/60 bg-white/82">
		<CardContent class="p-0">
			<Table.Root>
				<Table.Header>
					<Table.Row>
						<Table.Head class="pl-6">Policy</Table.Head>
						<Table.Head>Customer</Table.Head>
						<Table.Head>Product</Table.Head>
						<Table.Head>Status</Table.Head>
						<Table.Head>Premium</Table.Head>
						<Table.Head>Coverage</Table.Head>
						<Table.Head class="pr-6">Term</Table.Head>
					</Table.Row>
				</Table.Header>
				<Table.Body>
					{#each data.policies as policy}
						<Table.Row>
							<Table.Cell class="pl-6 font-medium text-slate-900">{policy.policyNumber || policy.id}</Table.Cell>
							<Table.Cell>{policy.customerName}</Table.Cell>
							<Table.Cell>{policy.productName}</Table.Cell>
							<Table.Cell>
								<Badge variant={statusTone(policy.status)}>{humanizeStatus(policy.status)}</Badge>
							</Table.Cell>
							<Table.Cell>{formatCurrency(policy.premium)}</Table.Cell>
							<Table.Cell>{formatCurrency(policy.sumInsured)}</Table.Cell>
							<Table.Cell class="pr-6">
								{formatDate(policy.startDate)} to {formatDate(policy.endDate)}
							</Table.Cell>
						</Table.Row>
					{/each}
				</Table.Body>
			</Table.Root>
		</CardContent>
	</Card>
</div>
