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
		title="Claims"
		description="Claims queue loaded from the backend claims API, normalized for the system portal."
		meta={`${data.claims.length} visible claims`}
	/>

	<Card class="rounded-[28px] border-white/60 bg-white/82">
		<CardContent class="p-0">
			<Table.Root>
				<Table.Header>
					<Table.Row>
						<Table.Head class="pl-6">Claim</Table.Head>
						<Table.Head>Claimant</Table.Head>
						<Table.Head>Policy</Table.Head>
						<Table.Head>Status</Table.Head>
						<Table.Head>Amount</Table.Head>
						<Table.Head class="pr-6">Timeline</Table.Head>
					</Table.Row>
				</Table.Header>
				<Table.Body>
					{#each data.claims as claim}
						<Table.Row>
							<Table.Cell class="pl-6 font-medium text-slate-900">{claim.claimNumber || claim.id}</Table.Cell>
							<Table.Cell>{claim.claimantName}</Table.Cell>
							<Table.Cell>{claim.policyNumber}</Table.Cell>
							<Table.Cell>
								<Badge variant={statusTone(claim.status)}>{humanizeStatus(claim.status)}</Badge>
							</Table.Cell>
							<Table.Cell>{formatCurrency(claim.amount)}</Table.Cell>
							<Table.Cell class="pr-6">
								<div>{formatDate(claim.submittedAt)}</div>
								<div class="text-xs text-slate-500">Incident {formatDate(claim.incidentDate)}</div>
							</Table.Cell>
						</Table.Row>
					{/each}
				</Table.Body>
			</Table.Root>
		</CardContent>
	</Card>
</div>
