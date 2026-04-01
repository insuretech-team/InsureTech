<script lang="ts">
	import { Badge } from '$lib/components/ui/badge';
	import { Card, CardContent } from '$lib/components/ui/card';
	import * as Table from '$lib/components/ui/table';
	import PageHeader from '$lib/components/system/page-header.svelte';
	import { humanizeStatus, statusTone } from '$lib/system/format';

	let { data } = $props();
</script>

<div class="space-y-8">
	<PageHeader
		title="Reports"
		description="Report definitions visible to the system portal through the reporting service."
		meta={`${data.reports.length} report definitions`}
	/>

	<Card class="rounded-[28px] border-white/60 bg-white/82">
		<CardContent class="p-0">
			<Table.Root>
				<Table.Header>
					<Table.Row>
						<Table.Head class="pl-6">Report</Table.Head>
						<Table.Head>Code</Table.Head>
						<Table.Head>Status</Table.Head>
						<Table.Head class="pr-6">Description</Table.Head>
					</Table.Row>
				</Table.Header>
				<Table.Body>
					{#each data.reports as report}
						<Table.Row>
							<Table.Cell class="pl-6 font-medium text-slate-900">{report.name}</Table.Cell>
							<Table.Cell>{report.code || report.id}</Table.Cell>
							<Table.Cell><Badge variant={statusTone(report.status)}>{humanizeStatus(report.status)}</Badge></Table.Cell>
							<Table.Cell class="pr-6">{report.description || 'No description available'}</Table.Cell>
						</Table.Row>
					{/each}
				</Table.Body>
			</Table.Root>
		</CardContent>
	</Card>
</div>
