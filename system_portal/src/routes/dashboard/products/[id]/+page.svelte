<script lang="ts">
	import { page } from '$app/stores';
	import { Button } from '$lib/components/ui/button';
	import { Badge } from '$lib/components/ui/badge';
	import * as Card from '$lib/components/ui/card';
	import * as Table from '$lib/components/ui/table';
	import { ArrowLeft, Edit, Trash2 } from 'lucide-svelte';
	import { getProductById, formatBDT } from '$lib/data_detailed/products_demo';
	import {  ProductStatus, ProductCategory  } from '$lib/types';

	$: product_id = $page.params.id;
	$: product = getProductById(product_id);

	function getCategoryName(category: ProductCategory): string {
		const names: Record<number, string> = {
			1: 'Motor',
			2: 'Health',
			3: 'Travel',
			4: 'Home',
			5: 'Device',
			6: 'Agricultural',
			7: 'Life'
		};
		return names[category] || 'Unknown';
	}

	function getStatusName(status: ProductStatus): string {
		const names: Record<number, string> = {
			1: 'Draft',
			2: 'Active',
			3: 'Inactive',
			4: 'Discontinued'
		};
		return names[status] || 'Unknown';
	}

	function formatDate(dateStr: string): string {
		return new Date(dateStr).toLocaleDateString('en-BD', {
			year: 'numeric',
			month: 'long',
			day: 'numeric',
			hour: '2-digit',
			minute: '2-digit'
		});
	}
</script>

{#if product}
	<div class="space-y-6">
		<!-- Header -->
		<div class="flex items-center justify-between">
			<div class="flex items-center gap-4">
				<Button variant="ghost" size="icon" href="/dashboard/products">
					<ArrowLeft class="h-4 w-4" />
				</Button>
				<div>
					<h1 class="text-3xl font-bold tracking-tight">{product.product_name}</h1>
					<p class="text-muted-foreground">{product.product_code}</p>
				</div>
			</div>
			<div class="flex gap-2">
				<Button variant="outline">
					<Edit class="mr-2 h-4 w-4" />
					Edit
				</Button>
				<Button variant="destructive">
					<Trash2 class="mr-2 h-4 w-4" />
					Delete
				</Button>
			</div>
		</div>

		<!-- Product Overview -->
		<div class="grid gap-4 md:grid-cols-3">
			<Card.Root>
				<Card.Header>
					<Card.Title>Category</Card.Title>
				</Card.Header>
				<Card.Content>
					<Badge variant="outline" class="text-lg">{getCategoryName(product.category)}</Badge>
				</Card.Content>
			</Card.Root>

			<Card.Root>
				<Card.Header>
					<Card.Title>Status</Card.Title>
				</Card.Header>
				<Card.Content>
					<Badge class="text-lg">{getStatusName(product.status)}</Badge>
				</Card.Content>
			</Card.Root>

			<Card.Root>
				<Card.Header>
					<Card.Title>Base Premium</Card.Title>
				</Card.Header>
				<Card.Content>
					<div class="text-2xl font-bold">{formatBDT(product.base_premium)}</div>
				</Card.Content>
			</Card.Root>
		</div>

		<!-- Product Details -->
		<Card.Root>
			<Card.Header>
				<Card.Title>Product Description</Card.Title>
			</Card.Header>
			<Card.Content>
				<p class="text-muted-foreground">{product.description}</p>
			</Card.Content>
		</Card.Root>

		<!-- Coverage and Tenure -->
		<div class="grid gap-4 md:grid-cols-2">
			<Card.Root>
				<Card.Header>
					<Card.Title>Coverage Details</Card.Title>
				</Card.Header>
				<Card.Content class="space-y-4">
					<div>
						<div class="text-sm font-medium text-muted-foreground">Minimum Sum Insured</div>
						<div class="text-xl font-bold">{formatBDT(product.min_sum_insured)}</div>
					</div>
					<div>
						<div class="text-sm font-medium text-muted-foreground">Maximum Sum Insured</div>
						<div class="text-xl font-bold">{formatBDT(product.max_sum_insured)}</div>
					</div>
				</Card.Content>
			</Card.Root>

			<Card.Root>
				<Card.Header>
					<Card.Title>Tenure Details</Card.Title>
				</Card.Header>
				<Card.Content class="space-y-4">
					<div>
						<div class="text-sm font-medium text-muted-foreground">Minimum Tenure</div>
						<div class="text-xl font-bold">{product.min_tenure_months} months</div>
					</div>
					<div>
						<div class="text-sm font-medium text-muted-foreground">Maximum Tenure</div>
						<div class="text-xl font-bold">{product.max_tenure_months} months</div>
					</div>
				</Card.Content>
			</Card.Root>
		</div>

		<!-- Riders -->
		{#if product.available_riders && product.available_riders.length > 0}
			<Card.Root>
				<Card.Header>
					<Card.Title>Available Riders</Card.Title>
					<Card.Description>Add-on coverages available with this product</Card.Description>
				</Card.Header>
				<Card.Content>
					<Table.Root>
						<Table.Header>
							<Table.Row>
								<Table.Head>Rider Name</Table.Head>
								<Table.Head>Description</Table.Head>
								<Table.Head>Premium</Table.Head>
								<Table.Head>Coverage</Table.Head>
								<Table.Head>Mandatory</Table.Head>
							</Table.Row>
						</Table.Header>
						<Table.Body>
							{#each product.available_riders as rider}
								<Table.Row>
									<Table.Cell class="font-medium">{rider.rider_name}</Table.Cell>
									<Table.Cell>{rider.description}</Table.Cell>
									<Table.Cell>{formatBDT(rider.premium_amount)}</Table.Cell>
									<Table.Cell>{formatBDT(rider.coverage_amount)}</Table.Cell>
									<Table.Cell>
										{#if rider.is_mandatory}
											<Badge variant="destructive">Required</Badge>
										{:else}
											<Badge variant="secondary">Optional</Badge>
										{/if}
									</Table.Cell>
								</Table.Row>
							{/each}
						</Table.Body>
					</Table.Root>
				</Card.Content>
			</Card.Root>
		{/if}

		<!-- Exclusions -->
		{#if product.exclusions && product.exclusions.length > 0}
			<Card.Root>
				<Card.Header>
					<Card.Title>Exclusions</Card.Title>
					<Card.Description>Conditions and scenarios not covered by this product</Card.Description>
				</Card.Header>
				<Card.Content>
					<ul class="space-y-2">
						{#each product.exclusions as exclusion}
							<li class="flex items-start gap-2">
								<span class="text-destructive mt-1">•</span>
								<span>{exclusion}</span>
							</li>
						{/each}
					</ul>
				</Card.Content>
			</Card.Root>
		{/if}

		<!-- Metadata -->
		<Card.Root>
			<Card.Header>
				<Card.Title>Product Metadata</Card.Title>
			</Card.Header>
			<Card.Content>
				<div class="grid gap-4 md:grid-cols-2">
					<div>
						<div class="text-sm font-medium text-muted-foreground">Product ID</div>
						<div class="font-mono text-sm">{product.product_id}</div>
					</div>
					<div>
						<div class="text-sm font-medium text-muted-foreground">Product Code</div>
						<div class="font-mono text-sm">{product.product_code}</div>
					</div>
					<div>
						<div class="text-sm font-medium text-muted-foreground">Created At</div>
						<div class="text-sm">{formatDate(product.created_at)}</div>
					</div>
					<div>
						<div class="text-sm font-medium text-muted-foreground">Last Updated</div>
						<div class="text-sm">{formatDate(product.updated_at)}</div>
					</div>
					<div>
						<div class="text-sm font-medium text-muted-foreground">Created By</div>
						<div class="text-sm">{product.created_by}</div>
					</div>
				</div>
			</Card.Content>
		</Card.Root>
	</div>
{:else}
	<div class="flex h-[50vh] items-center justify-center">
		<Card.Root class="w-96">
			<Card.Header>
				<Card.Title>Product Not Found</Card.Title>
				<Card.Description>The product you're looking for doesn't exist.</Card.Description>
			</Card.Header>
			<Card.Content>
				<Button href="/dashboard/products">Back to Products</Button>
			</Card.Content>
		</Card.Root>
	</div>
{/if}
