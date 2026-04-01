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
		title="Tenants"
		description="Tenant list backed by the tenant service in the generated SDK."
		meta={`${data.tenants.length} tenant records`}
	/>

	<Card class="rounded-[28px] border-white/60 bg-white/82">
		<CardContent class="p-0">
			<Table.Root>
				<Table.Header>
					<Table.Row>
						<Table.Head class="pl-6">Tenant</Table.Head>
						<Table.Head>Code</Table.Head>
						<Table.Head>Status</Table.Head>
						<Table.Head>Domain</Table.Head>
						<Table.Head class="pr-6">Created</Table.Head>
					</Table.Row>
				</Table.Header>
				<Table.Body>
					{#each data.tenants as tenant}
						<Table.Row>
							<Table.Cell class="pl-6 font-medium text-slate-900">{tenant.name}</Table.Cell>
							<Table.Cell>{tenant.code || tenant.id}</Table.Cell>
							<Table.Cell><Badge variant={statusTone(tenant.status)}>{humanizeStatus(tenant.status)}</Badge></Table.Cell>
							<Table.Cell>{tenant.domain || 'Not provided'}</Table.Cell>
							<Table.Cell class="pr-6">{formatDate(tenant.createdAt)}</Table.Cell>
						</Table.Row>
					{/each}
				</Table.Body>
			</Table.Root>
		</CardContent>
	</Card>
</div>
