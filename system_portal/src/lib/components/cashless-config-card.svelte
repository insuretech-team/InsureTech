<script lang="ts">
	import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '$lib/components/ui/card';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { Button } from '$lib/components/ui/button';
	import { Badge } from '$lib/components/ui/badge';
	import { Checkbox } from '$lib/components/ui/checkbox';
	import { CreditCard, Save, X, CheckCircle2, AlertCircle } from 'lucide-svelte';

	let {
		partnerId = '',
		partnerName = '',
		cashlessEnabled = $bindable(false),
		cashlessLimit = $bindable(0),
		autoApprovalThreshold = $bindable(0),
		preAuthRequired = $bindable(false),
		authValidityDays = $bindable(30),
		requiredDocuments = $bindable<string[]>([]),
		onSave = () => {},
		onCancel = () => {}
	} = $props();

	let editing = $state(false);

	const documentOptions = [
		'National ID Card',
		'Policy Certificate',
		'Medical Prescription',
		'Hospital Admission Form',
		'Treatment Estimate',
		'Previous Medical Records',
		'Insurance Card'
	];

	function handleToggle() {
		cashlessEnabled = !cashlessEnabled;
		editing = true;
	}

	function formatCurrency(amount: number): string {
		return new Intl.NumberFormat('en-BD', {
			style: 'currency',
			currency: 'BDT',
			minimumFractionDigits: 0
		}).format(amount / 100);
	}

	function handleSave() {
		onSave({
			partnerId,
			cashlessEnabled,
			cashlessLimit,
			autoApprovalThreshold,
			preAuthRequired,
			authValidityDays,
			requiredDocuments
		});
		editing = false;
	}

	function handleCancel() {
		editing = false;
		onCancel();
	}

	function toggleDocument(doc: string) {
		if (requiredDocuments.includes(doc)) {
			requiredDocuments = requiredDocuments.filter((d) => d !== doc);
		} else {
			requiredDocuments = [...requiredDocuments, doc];
		}
	}
</script>

<Card class="border-2 {cashlessEnabled ? 'border-green-200 bg-green-50/50 dark:border-green-900 dark:bg-green-950/20' : ''}">
	<CardHeader>
		<div class="flex items-center justify-between">
			<div class="flex items-center gap-3">
				<div class="rounded-full bg-green-100 p-2 text-green-700 dark:bg-green-900 dark:text-green-300">
					<CreditCard class="h-5 w-5" />
				</div>
				<div>
					<CardTitle class="text-lg">Cashless Configuration</CardTitle>
					<CardDescription>Configure direct billing for {partnerName}</CardDescription>
				</div>
			</div>
			<div class="flex items-center gap-2">
				{#if cashlessEnabled}
					<Badge class="bg-green-600">Enabled</Badge>
				{:else}
					<Badge variant="secondary">Disabled</Badge>
				{/if}
				<Button
					variant={cashlessEnabled ? 'outline' : 'default'}
					size="sm"
					onclick={handleToggle}
				>
					{cashlessEnabled ? 'Disable' : 'Enable'}
				</Button>
			</div>
		</div>
	</CardHeader>
	<CardContent class="space-y-4">
		{#if cashlessEnabled}
			<div class="grid gap-4 md:grid-cols-2">
				<!-- Cashless Limit -->
				<div class="space-y-2">
					<Label for="cashless-limit">
						Maximum Cashless Limit
						<span class="text-destructive">*</span>
					</Label>
					<div class="relative">
						<span class="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground">
							৳
						</span>
						<Input
							id="cashless-limit"
							type="number"
							min="0"
							step="100"
							bind:value={cashlessLimit}
							placeholder="Enter amount"
							class="pl-8"
							disabled={!editing && cashlessEnabled}
						/>
					</div>
					<p class="text-xs text-muted-foreground">
						Current limit: <strong class="text-green-600">{formatCurrency(cashlessLimit)}</strong>
					</p>
				</div>

				<!-- Auto Approval Threshold -->
				<div class="space-y-2">
					<Label for="auto-approval">
						Auto-Approval Threshold
					</Label>
					<div class="relative">
						<span class="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground">
							৳
						</span>
						<Input
							id="auto-approval"
							type="number"
							min="0"
							max={cashlessLimit}
							step="100"
							bind:value={autoApprovalThreshold}
							placeholder="Auto-approve below"
							class="pl-8"
							disabled={!editing && cashlessEnabled}
						/>
					</div>
					<p class="text-xs text-muted-foreground">
						Claims below {formatCurrency(autoApprovalThreshold)} are auto-approved
					</p>
				</div>

				<!-- Authorization Validity -->
				<div class="space-y-2">
					<Label for="auth-validity">Authorization Validity (Days)</Label>
					<Input
						id="auth-validity"
						type="number"
						min="1"
						max="365"
						bind:value={authValidityDays}
						placeholder="Days"
						disabled={!editing && cashlessEnabled}
					/>
					<p class="text-xs text-muted-foreground">
						Pre-authorization valid for {authValidityDays} days
					</p>
				</div>

				<!-- Pre-Authorization Toggle -->
				<div class="space-y-2">
					<Label>Pre-Authorization Required</Label>
					<div class="flex items-center space-x-2 rounded-lg border p-3">
						<Checkbox
							id="pre-auth"
							bind:checked={preAuthRequired}
							disabled={!editing && cashlessEnabled}
						/>
						<Label
							for="pre-auth"
							class="cursor-pointer text-sm font-normal leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
						>
							Require approval before treatment
						</Label>
					</div>
					{#if preAuthRequired}
						<div class="flex items-center gap-2 text-xs text-orange-600">
							<AlertCircle class="h-3 w-3" />
							Manual pre-authorization required
						</div>
					{:else}
						<div class="flex items-center gap-2 text-xs text-green-600">
							<CheckCircle2 class="h-3 w-3" />
							Instant cashless approval
						</div>
					{/if}
				</div>
			</div>

			<!-- Limit Visual -->
			<div class="rounded-lg border bg-card p-4">
				<p class="mb-2 text-sm font-medium">Approval Threshold</p>
				<div class="space-y-2">
					<div class="h-3 w-full rounded-full bg-secondary">
						<div
							class="h-3 rounded-full bg-gradient-to-r from-green-500 to-orange-500"
							style="width: {(autoApprovalThreshold / cashlessLimit) * 100}%"
						></div>
					</div>
					<div class="flex justify-between text-xs">
						<span class="text-green-600">
							<CheckCircle2 class="inline h-3 w-3" />
							Auto: {formatCurrency(autoApprovalThreshold)}
						</span>
						<span class="text-orange-600">
							<AlertCircle class="inline h-3 w-3" />
							Manual: {formatCurrency(cashlessLimit - autoApprovalThreshold)}
						</span>
						<span class="text-muted-foreground">Max: {formatCurrency(cashlessLimit)}</span>
					</div>
				</div>
			</div>

			<!-- Required Documents -->
			<div class="space-y-2">
				<Label>Required Documents for Cashless</Label>
				<div class="grid gap-2 rounded-lg border p-3 md:grid-cols-2">
					{#each documentOptions as doc}
						<div class="flex items-center space-x-2">
							<Checkbox
								id={doc}
								checked={requiredDocuments.includes(doc)}
								onCheckedChange={() => toggleDocument(doc)}
								disabled={!editing && cashlessEnabled}
							/>
							<Label
								for={doc}
								class="cursor-pointer text-sm font-normal leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
							>
								{doc}
							</Label>
						</div>
					{/each}
				</div>
				<p class="text-xs text-muted-foreground">
					{requiredDocuments.length} document(s) required for cashless approval
				</p>
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
			{:else if cashlessEnabled}
				<Button variant="outline" size="sm" onclick={() => (editing = true)} class="w-full">
					Edit Cashless Settings
				</Button>
			{/if}
		{:else}
			<div class="rounded-lg border border-dashed p-8 text-center">
				<CreditCard class="mx-auto mb-3 h-12 w-12 text-muted-foreground" />
				<p class="text-sm text-muted-foreground">
					Cashless facility is currently disabled for this partner.
				</p>
				<p class="mt-2 text-xs text-muted-foreground">
					Click "Enable" to configure cashless billing.
				</p>
			</div>
		{/if}
	</CardContent>
</Card>
