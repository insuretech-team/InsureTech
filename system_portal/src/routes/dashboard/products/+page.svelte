<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Badge } from '$lib/components/ui/badge';
	import * as Card from '$lib/components/ui/card';
	import * as Table from '$lib/components/ui/table';
	import { Search, Plus, Filter } from 'lucide-svelte';
	import { productsDemo, formatBDT, getActiveProducts, searchProducts } from '$lib/data_detailed/products_demo';
	import {  ProductStatus, ProductCategory  } from '$lib/types';

	let searchQuery = '';
	let selectedCategory = 'ALL';
	let selectedStatus = ProductStatus.PRODUCT_STATUS_ACTIVE;

	$: filteredProducts = searchQuery
		? searchProducts(searchQuery)
		: selectedCategory === 'ALL'
			? productsDemo.filter(p => p.status === selectedStatus)
			: productsDemo.filter(p => p.category === parseInt(selectedCategory) && p.status === selectedStatus);

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

	function getStatusBadge(status: ProductStatus): string {
		switch (status) {
			case ProductStatus.PRODUCT_STATUS_ACTIVE:
				return 'default';
			case ProductStatus.PRODUCT_STATUS_INACTIVE:
				return 'secondary';
			case ProductStatus.PRODUCT_STATUS_DRAFT:
				return 'outline';
			default:
				return 'secondary';
		}
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
</script>

<div class="space-y-6">
	<!-- Header -->
	<div class="flex items-center justify-between">
		<div>
			<h1 class="text-3xl font-bold tracking-tight">Products</h1>
			<p class="text-muted-foreground">Manage insurance products and their configurations</p>
		</div>
		<Button>
			<Plus class="mr-2 h-4 w-4" />
			Add Product
		</Button>
	</div>

	<!-- Stats Cards -->
	<div class="grid gap-4 md:grid-cols-4">
		<Card.Root>
			<Card.Header class="flex flex-row items-center justify-between space-y-0 pb-2">
				<Card.Title class="text-sm font-medium">Total Products</Card.Title>
			</Card.Header>
			<Card.Content>
				<div class="text-2xl font-bold">{productsDemo.length}</div>
				<p class="text-xs text-muted-foreground">Across all categories</p>
			</Card.Content>
		</Card.Root>

		<Card.Root>
			<Card.Header class="flex flex-row items-center justify-between space-y-0 pb-2">
				<Card.Title class="text-sm font-medium">Active Products</Card.Title>
			</Card.Header>
			<Card.Content>
				<div class="text-2xl font-bold">
					{productsDemo.filter((p) => p.status === ProductStatus.PRODUCT_STATUS_ACTIVE).length}
				</div>
				<p class="text-xs text-muted-foreground">Available for purchase</p>
			</Card.Content>
		</Card.Root>

		<Card.Root>
			<Card.Header class="flex flex-row items-center justify-between space-y-0 pb-2">
				<Card.Title class="text-sm font-medium">Health Products</Card.Title>
			</Card.Header>
			<Card.Content>
				<div class="text-2xl font-bold">
					{productsDemo.filter((p) => p.category === ProductCategory.PRODUCT_CATEGORY_HEALTH).length}
				</div>
				<p class="text-xs text-muted-foreground">Most popular category</p>
			</Card.Content>
		</Card.Root>

		<Card.Root>
			<Card.Header class="flex flex-row items-center justify-between space-y-0 pb-2">
				<Card.Title class="text-sm font-medium">With Riders</Card.Title>
			</Card.Header>
			<Card.Content>
				<div class="text-2xl font-bold">
					{productsDemo.filter((p) => p.available_riders && p.available_riders.length > 0).length}
				</div>
				<p class="text-xs text-muted-foreground">Products with add-ons</p>
			</Card.Content>
		</Card.Root>
	</div>

	<!-- Filters and Search -->
	<Card.Root>
		<Card.Header>
			<div class="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
				<div class="flex gap-2">
					<div class="relative w-full md:w-80">
						<Search class="absolute left-2 top-2.5 h-4 w-4 text-muted-foreground" />
						<Input
							placeholder="Search products..."
							class="pl-8"
							bind:value={searchQuery}
						/>
					</div>
				</div>

				<div class="flex gap-2">
					<select
						bind:value={selectedCategory}
						class="flex h-10 rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background"
					>
						<option value="ALL">All Categories</option>
						<option value="{ProductCategory.PRODUCT_CATEGORY_MOTOR}">Motor</option>
						<option value="{ProductCategory.PRODUCT_CATEGORY_HEALTH}">Health</option>
						<option value="{ProductCategory.PRODUCT_CATEGORY_TRAVEL}">Travel</option>
						<option value="{ProductCategory.PRODUCT_CATEGORY_HOME}">Home</option>
						<option value="{ProductCategory.PRODUCT_CATEGORY_DEVICE}">Device</option>
						<option value="{ProductCategory.PRODUCT_CATEGORY_AGRICULTURAL}">Agricultural</option>
						<option value="{ProductCategory.PRODUCT_CATEGORY_LIFE}">Life</option>
					</select>

					<select
						bind:value={selectedStatus}
						class="flex h-10 rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background"
					>
						<option value={ProductStatus.PRODUCT_STATUS_ACTIVE}>Active</option>
						<option value={ProductStatus.PRODUCT_STATUS_INACTIVE}>Inactive</option>
						<option value={ProductStatus.PRODUCT_STATUS_DRAFT}>Draft</option>
					</select>
				</div>
			</div>
		</Card.Header>

		<Card.Content>
			<Table.Root>
				<Table.Header>
					<Table.Row>
						<Table.Head>Product Code</Table.Head>
						<Table.Head>Product Name</Table.Head>
						<Table.Head>Category</Table.Head>
						<Table.Head>Base Premium</Table.Head>
						<Table.Head>Coverage Range</Table.Head>
						<Table.Head>Tenure</Table.Head>
						<Table.Head>Riders</Table.Head>
						<Table.Head>Status</Table.Head>
						<Table.Head class="text-right">Actions</Table.Head>
					</Table.Row>
				</Table.Header>
				<Table.Body>
					{#each filteredProducts as product}
						<Table.Row>
							<Table.Cell class="font-medium">{product.product_code}</Table.Cell>
							<Table.Cell>
								<div class="flex flex-col">
									<span class="font-medium">{product.product_name}</span>
									<span class="text-xs text-muted-foreground">{product.description?.substring(0, 50)}...</span>
								</div>
							</Table.Cell>
							<Table.Cell>
								<Badge variant="outline">{getCategoryName(product.category)}</Badge>
							</Table.Cell>
							<Table.Cell>{formatBDT(product.base_premium)}</Table.Cell>
							<Table.Cell>
								<div class="text-sm">
									{formatBDT(product.min_sum_insured)} - {formatBDT(product.max_sum_insured)}
								</div>
							</Table.Cell>
							<Table.Cell>
								{product.min_tenure_months} - {product.max_tenure_months} months
							</Table.Cell>
							<Table.Cell>
								{#if product.available_riders && product.available_riders.length > 0}
									<Badge variant="secondary">{product.available_riders.length} riders</Badge>
								{:else}
									<span class="text-muted-foreground text-sm">None</span>
								{/if}
							</Table.Cell>
							<Table.Cell>
								<Badge variant={getStatusBadge(product.status)}>
									{getStatusName(product.status)}
								</Badge>
							</Table.Cell>
							<Table.Cell class="text-right">
								<Button variant="ghost" size="sm" href="/dashboard/products/{product.product_id}">
									View
								</Button>
							</Table.Cell>
						</Table.Row>
					{:else}
						<Table.Row>
							<Table.Cell colspan={9} class="text-center text-muted-foreground">
								No products found
							</Table.Cell>
						</Table.Row>
					{/each}
				</Table.Body>
			</Table.Root>
		</Card.Content>
	</Card.Root>
</div>
