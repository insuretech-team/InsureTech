<script lang="ts">
	import { Badge } from '$lib/components/ui/badge';
	import { Card, CardContent } from '$lib/components/ui/card';
	import * as Table from '$lib/components/ui/table';
	import PageHeader from '$lib/components/system/page-header.svelte';
	import { formatDate, humanizeStatus, statusTone } from '$lib/system/format';

	let { data } = $props();
</script>

<div class="space-y-8">
	<PageHeader
		title="Audit trail"
		description="Operational events and audit records returned from the backend audit endpoint."
		meta={`${data.events.length} recent events`}
	/>

	<Card class="rounded-[28px] border-white/60 bg-white/82">
		<CardContent class="p-0">
			<Table.Root>
				<Table.Header>
					<Table.Row>
						<Table.Head class="pl-6">Action</Table.Head>
						<Table.Head>Resource</Table.Head>
						<Table.Head>Actor</Table.Head>
						<Table.Head>Status</Table.Head>
						<Table.Head class="pr-6">Timestamp</Table.Head>
					</Table.Row>
				</Table.Header>
				<Table.Body>
					{#each data.events as eventRow}
						<Table.Row>
							<Table.Cell class="pl-6 font-medium text-slate-900">{eventRow.action}</Table.Cell>
							<Table.Cell>{eventRow.resource}</Table.Cell>
							<Table.Cell>{eventRow.actor}</Table.Cell>
							<Table.Cell><Badge variant={statusTone(eventRow.status)}>{humanizeStatus(eventRow.status)}</Badge></Table.Cell>
							<Table.Cell class="pr-6">{formatDate(eventRow.timestamp)}</Table.Cell>
						</Table.Row>
					{/each}
				</Table.Body>
			</Table.Root>
		</CardContent>
	</Card>
</div>
