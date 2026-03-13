<script lang="ts">
	import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '$lib/components/ui/card';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { Button } from '$lib/components/ui/button';
	import { Badge } from '$lib/components/ui/badge';
	import * as Select from '$lib/components/ui/select';
	import { Percent, Save, X } from 'lucide-svelte';

	let {
		partnerId = '',
		partnerName = '',
		discountEnabled = $bindable(false),
		discountPercentage = $bindable(0),
		minDiscount = $bindable(0),
		maxDiscount = $bindable(100),
		discountType = $bindable('SERVICE'),
		onSave = () => {},
		onCancel = () => {}
	} = $props();

	let editing = $state(false);

	function handleToggle() {
		discountEnabled = !discountEnabled;
		editing = true;
	}

	function handleSave() {
		onSave({
			partnerId,
			discountEnabled,
			discountPercentage,
			minDiscount,
			maxDiscount,
			discountType
		});
		editing = false;
	}

	function handleCancel() {
		editing = false;
		onCancel();
	}
</script>

<Card class="border-2 {discountEnabled ? 'border-green-200 bg-green-50/50 dark:border-green-900 dark:bg-green-950/20' : ''}">
	<CardHeader>
		<div class="flex items-center justify-between">
			<div class="flex items-center gap-3">
				<div class="rounded-full bg-blue-100 p-2 text-blue-700 dark:bg-blue-900 dark:text-blue-300">
					<Percent class="h-5 w-5" />
				</div>
				<div>
					<CardTitle class="text-lg">Discount Configuration</CardTitle>
					<CardDescription>Set discount rates for {partnerName}</CardDescription>
				</div>
			</div>
			<div class="flex items-center gap-2">
				{#if discountEnabled}
					<Badge class="bg-green-600">Active</Badge>
				{:else}
					<Badge variant="secondary">Inactive</Badge>
				{/if}
				<Button
					variant={discountEnabled ? 'outline' : 'default'}
					size="sm"
					onclick={handleToggle}
				>
					{discountEnabled ? 'Disable' : 'Enable'}
				</Button>
			</div>
		</div>
	</CardHeader>
	<CardContent class="space-y-4">
		{#if discountEnabled}
			<div class="grid gap-4 md:grid-cols-2">
				<!-- Discount Percentage -->
				<div class="space-y-2">
					<Label for="discount-percentage">
						Discount Percentage
						<span class="text-destructive">*</span>
					</Label>
					<div class="relative">
						<Input
							id="discount-percentage"
							type="number"
							min="0"
							max="100"
							step="0.1"
							bind:value={discountPercentage}
							placeholder="Enter discount %"
							class="pr-8"
							disabled={!editing && discountEnabled}
						/>
						<div class="pointer-events-none absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground">
							%
						</div>
					</div>
					<p class="text-xs text-muted-foreground">
						Current discount: <strong class="text-green-600">{discountPercentage}%</strong>
					</p>
				</div>

				<!-- Discount Type -->
				<div class="space-y-2">
					<Label for="discount-type">
						Discount Type
						<span class="text-destructive">*</span>
					</Label>
					<Select.Root bind:selected={discountType} disabled={!editing && discountEnabled}>
						<Select.Trigger id="discount-type">
							{discountType || 'Select type'}
						</Select.Trigger>
						<Select.Content>
							<Select.Item value="SERVICE">Service Discount</Select.Item>
							<Select.Item value="PRODUCT">Product Discount</Select.Item>
							<Select.Item value="CONSULTATION">Consultation Fee</Select.Item>
							<Select.Item value="MEDICATION">Medication Discount</Select.Item>
							<Select.Item value="BULK">Bulk Purchase</Select.Item>
						</Select.Content>
					</Select.Root>
				</div>

				<!-- Min Discount -->
				<div class="space-y-2">
					<Label for="min-discount">Minimum Discount %</Label>
					<div class="relative">
						<Input
							id="min-discount"
							type="number"
							min="0"
							max={maxDiscount}
							step="0.1"
							bind:value={minDiscount}
							placeholder="Min %"
							class="pr-8"
							disabled={!editing && discountEnabled}
						/>
						<div class="pointer-events-none absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground">
							%
						</div>
					</div>
				</div>

				<!-- Max Discount -->
				<div class="space-y-2">
					<Label for="max-discount">Maximum Discount %</Label>
					<div class="relative">
						<Input
							id="max-discount"
							type="number"
							min={minDiscount}
							max="100"
							step="0.1"
							bind:value={maxDiscount}
							placeholder="Max %"
							class="pr-8"
							disabled={!editing && discountEnabled}
						/>
						<div class="pointer-events-none absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground">
							%
						</div>
					</div>
				</div>
			</div>

			<!-- Discount Range Visual -->
			<div class="rounded-lg border bg-card p-4">
				<p class="mb-2 text-sm font-medium">Discount Range</p>
				<div class="flex items-center gap-4">
					<div class="flex-1">
						<div class="h-3 w-full rounded-full bg-secondary">
							<div
								class="h-3 rounded-full bg-gradient-to-r from-blue-500 via-green-500 to-green-600"
								style="width: {(discountPercentage / 100) * 100}%"
							></div>
						</div>
						<div class="mt-1 flex justify-between text-xs text-muted-foreground">
							<span>{minDiscount}%</span>
							<span class="font-semibold text-green-600">{discountPercentage}%</span>
							<span>{maxDiscount}%</span>
						</div>
					</div>
				</div>
			</div>

			<!-- Action Buttons -->
			{#if editing}
				<div class="flex justify-end gap-2">
					<Button variant="outline" size="sm" onclick={handleCancel}>
						<X class="mr-2 h-4 w-4" />
						Cancel
					</Button>
					<Button size="sm" onclick={handleSave}>
						<Save class="mr-2 h-4 w-4" />
						Save Changes
					</Button>
				</div>
			{:else if discountEnabled}
				<Button variant="outline" size="sm" onclick={() => (editing = true)} class="w-full">
					Edit Discount Settings
				</Button>
			{/if}
		{:else}
			<div class="rounded-lg border border-dashed p-8 text-center">
				<Percent class="mx-auto mb-3 h-12 w-12 text-muted-foreground" />
				<p class="text-sm text-muted-foreground">
					Discount is currently disabled for this partner.
				</p>
				<p class="mt-2 text-xs text-muted-foreground">
					Click "Enable" to configure discount rates.
				</p>
			</div>
		{/if}
	</CardContent>
</Card>
