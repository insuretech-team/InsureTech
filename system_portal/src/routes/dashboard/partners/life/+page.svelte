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
		title="Life partners"
		description="Hospitals, pharmacies, doctors, and similar health-service entities currently visible to the partner service."
		meta={`${data.partners.length} life-network records`}
	/>

	<Card class="rounded-[28px] border-white/60 bg-white/82">
		<CardContent class="p-0">
			<Table.Root>
				<Table.Header>
					<Table.Row>
						<Table.Head class="pl-6">Partner</Table.Head>
						<Table.Head>Type</Table.Head>
						<Table.Head>Status</Table.Head>
						<Table.Head>Contact</Table.Head>
						<Table.Head class="pr-6">Joined</Table.Head>
					</Table.Row>
				</Table.Header>
				<Table.Body>
					{#each data.partners as partner}
						<Table.Row>
							<Table.Cell class="pl-6">
								<a href={`/dashboard/partners/${partner.id}`} class="font-medium text-primary hover:text-accent">
									{partner.name}
								</a>
								<div class="text-xs text-slate-500">{partner.id}</div>
							</Table.Cell>
							<Table.Cell>{partner.type}</Table.Cell>
							<Table.Cell><Badge variant={statusTone(partner.status)}>{humanizeStatus(partner.status)}</Badge></Table.Cell>
							<Table.Cell>
								<div>{partner.email || 'No email'}</div>
								<div class="text-xs text-slate-500">{partner.phone || 'No phone'}</div>
							</Table.Cell>
							<Table.Cell class="pr-6">{formatDate(partner.joinedAt)}</Table.Cell>
						</Table.Row>
					{/each}
				</Table.Body>
			</Table.Root>
		</CardContent>
	</Card>
</div>
