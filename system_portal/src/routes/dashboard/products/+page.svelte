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
		title="Products"
		description="System-wide product catalog from the generated SDK. This view is backed by the product service instead of local demo fixtures."
		meta={`${data.products.length} records`}
	/>

	<Card class="rounded-[28px] border-white/60 bg-white/82">
		<CardContent class="p-0">
			<Table.Root>
				<Table.Header>
					<Table.Row>
						<Table.Head class="pl-6">Product</Table.Head>
						<Table.Head>Category</Table.Head>
						<Table.Head>Status</Table.Head>
						<Table.Head>Premium</Table.Head>
						<Table.Head>Coverage band</Table.Head>
						<Table.Head class="pr-6">Updated</Table.Head>
					</Table.Row>
				</Table.Header>
				<Table.Body>
					{#each data.products as product}
						<Table.Row>
							<Table.Cell class="pl-6">
								<a href={`/dashboard/products/${product.id}`} class="font-medium text-primary hover:text-accent">
									{product.name}
								</a>
								<div class="text-xs text-slate-500">{product.code || product.id}</div>
							</Table.Cell>
							<Table.Cell>{product.category}</Table.Cell>
							<Table.Cell>
								<Badge variant={statusTone(product.status)}>{humanizeStatus(product.status)}</Badge>
							</Table.Cell>
							<Table.Cell>{formatCurrency(product.basePremium)}</Table.Cell>
							<Table.Cell>
								{formatCurrency(product.minSumInsured)} to {formatCurrency(product.maxSumInsured)}
							</Table.Cell>
							<Table.Cell class="pr-6">{formatDate(product.updatedAt || product.createdAt)}</Table.Cell>
						</Table.Row>
					{/each}
				</Table.Body>
			</Table.Root>
		</CardContent>
	</Card>
</div>
