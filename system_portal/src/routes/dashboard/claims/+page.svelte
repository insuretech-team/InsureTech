<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Badge } from '$lib/components/ui/badge';
	import * as Card from '$lib/components/ui/card';
	import * as Table from '$lib/components/ui/table';
	import { Search, Plus, AlertTriangle, CheckCircle, Clock } from 'lucide-svelte';
	import { claimsDemo, formatBDT, getPendingClaims, getClaimStatusColor, searchClaims } from '$lib/data_detailed/claims_demo';
	import {  ClaimStatus, ClaimType  } from '$lib/types';

	let searchQuery = '';
	let selectedStatus: ClaimStatus | 'ALL' = 'ALL';

	$: filteredClaims = searchQuery
		? searchClaims(searchQuery)
		: selectedStatus === 'ALL'
			? claimsDemo
			: claimsDemo.filter(c => c.status === selectedStatus);

	function getStatusName(status: ClaimStatus): string {
		const names: Record<number, string> = {
			1: 'Submitted',
			2: 'Under Review',
			3: 'Pending Documents',
			4: 'Approved',
			5: 'Rejected',
			6: 'Settled',
			7: 'Disputed'
		};
		return names[status] || 'Unknown';
	}

	function getTypeName(type: ClaimType): string {
		const names: Record<number, string> = {
			1: 'Health - Hospitalization',
			2: 'Health - Surgery',
			3: 'Motor - Accident',
			4: 'Motor - Theft',
			5: 'Travel - Medical',
			6: 'Travel - Baggage Loss',
			7: 'Device - Damage',
			8: 'Device - Theft',
			9: 'Death'
		};
		return names[type] || 'Unknown';
	}

	function formatDate(dateStr: string): string {
		return new Date(dateStr).toLocaleDateString('en-BD', {
			year: 'numeric',
			month: 'short',
			day: 'numeric'
		});
	}

	function getDaysSinceSubmission(submitted_at: string): number {
		const submitted = new Date(submitted_at);
		const now = new Date();
		const diffTime = Math.abs(now.getTime() - submitted.getTime());
		return Math.ceil(diffTime / (1000 * 60 * 60 * 24));
	}
</script>

<div class="space-y-6">
	<!-- Header -->
	<div class="flex items-center justify-between">
		<div>
			<h1 class="text-3xl font-bold tracking-tight">Claims Management</h1>
			<p class="text-muted-foreground">Review and process insurance claims</p>
		</div>
		<Button>
			<Plus class="mr-2 h-4 w-4" />
			New Claim
		</Button>
	</div>

	<!-- Stats Cards -->
	<div class="grid gap-4 md:grid-cols-5">
		<Card.Root>
			<Card.Header class="flex flex-row items-center justify-between space-y-0 pb-2">
				<Card.Title class="text-sm font-medium">Total Claims</Card.Title>
			</Card.Header>
			<Card.Content>
				<div class="text-2xl font-bold">{claimsDemo.length}</div>
				<p class="text-xs text-muted-foreground">All time</p>
			</Card.Content>
		</Card.Root>

		<Card.Root>
			<Card.Header class="flex flex-row items-center justify-between space-y-0 pb-2">
				<Card.Title class="text-sm font-medium">Pending Review</Card.Title>
				<Clock class="h-4 w-4 text-orange-500" />
			</Card.Header>
			<Card.Content>
				<div class="text-2xl font-bold text-orange-500">
					{getPendingClaims().length}
				</div>
				<p class="text-xs text-muted-foreground">Requires attention</p>
			</Card.Content>
		</Card.Root>

		<Card.Root>
			<Card.Header class="flex flex-row items-center justify-between space-y-0 pb-2">
				<Card.Title class="text-sm font-medium">Approved</Card.Title>
				<CheckCircle class="h-4 w-4 text-green-500" />
			</Card.Header>
			<Card.Content>
				<div class="text-2xl font-bold text-green-500">
					{claimsDemo.filter((c) => c.status === ClaimStatus.CLAIM_STATUS_APPROVED || c.status === ClaimStatus.CLAIM_STATUS_SETTLED).length}
				</div>
				<p class="text-xs text-muted-foreground">Approved/Settled</p>
			</Card.Content>
		</Card.Root>

		<Card.Root>
			<Card.Header class="flex flex-row items-center justify-between space-y-0 pb-2">
				<Card.Title class="text-sm font-medium">Fraud Flagged</Card.Title>
				<AlertTriangle class="h-4 w-4 text-red-500" />
			</Card.Header>
			<Card.Content>
				<div class="text-2xl font-bold text-red-500">
					{claimsDemo.filter((c) => c.fraud_check?.flagged).length}
				</div>
				<p class="text-xs text-muted-foreground">High risk detected</p>
			</Card.Content>
		</Card.Root>

		<Card.Root>
			<Card.Header class="flex flex-row items-center justify-between space-y-0 pb-2">
				<Card.Title class="text-sm font-medium">Total Claimed</Card.Title>
			</Card.Header>
			<Card.Content>
				<div class="text-2xl font-bold">
					{formatBDT(
						claimsDemo.reduce((sum, c) => sum + c.claimed_amount, BigInt(0))
					)}
				</div>
				<p class="text-xs text-muted-foreground">Total amount</p>
			</Card.Content>
		</Card.Root>
	</div>

	<!-- Filters and Search -->
	<Card.Root>
		<Card.Header>
			<div class="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
				<div class="relative w-full md:w-80">
					<Search class="absolute left-2 top-2.5 h-4 w-4 text-muted-foreground" />
					<Input
						placeholder="Search by claim number..."
						class="pl-8"
						bind:value={searchQuery}
					/>
				</div>

				<div class="flex gap-2">
					<select
						bind:value={selectedStatus}
						class="flex h-10 rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background"
					>
						<option value="ALL">All Status</option>
						<option value={ClaimStatus.CLAIM_STATUS_SUBMITTED}>Submitted</option>
						<option value={ClaimStatus.CLAIM_STATUS_UNDER_REVIEW}>Under Review</option>
						<option value={ClaimStatus.CLAIM_STATUS_PENDING_DOCUMENTS}>Pending Documents</option>
						<option value={ClaimStatus.CLAIM_STATUS_APPROVED}>Approved</option>
						<option value={ClaimStatus.CLAIM_STATUS_REJECTED}>Rejected</option>
						<option value={ClaimStatus.CLAIM_STATUS_SETTLED}>Settled</option>
					</select>
				</div>
			</div>
		</Card.Header>

		<Card.Content>
			<Table.Root>
				<Table.Header>
					<Table.Row>
						<Table.Head>Claim Number</Table.Head>
						<Table.Head>Type</Table.Head>
						<Table.Head>Customer</Table.Head>
						<Table.Head>Claimed Amount</Table.Head>
						<Table.Head>Approved Amount</Table.Head>
						<Table.Head>Incident Date</Table.Head>
						<Table.Head>Days Pending</Table.Head>
						<Table.Head>Fraud Score</Table.Head>
						<Table.Head>Status</Table.Head>
						<Table.Head class="text-right">Actions</Table.Head>
					</Table.Row>
				</Table.Header>
				<Table.Body>
					{#each filteredClaims as claim}
						<Table.Row>
							<Table.Cell class="font-medium">{claim.claim_number}</Table.Cell>
							<Table.Cell>
								<span class="text-sm">{getTypeName(claim.type)}</span>
							</Table.Cell>
							<Table.Cell>{claim.customer_id}</Table.Cell>
							<Table.Cell>{formatBDT(claim.claimed_amount)}</Table.Cell>
							<Table.Cell>
								{#if claim.approved_amount}
									{formatBDT(claim.approved_amount)}
								{:else}
									<span class="text-muted-foreground">-</span>
								{/if}
							</Table.Cell>
							<Table.Cell>{formatDate(claim.incident_date)}</Table.Cell>
							<Table.Cell>
								{#if claim.status !== ClaimStatus.CLAIM_STATUS_SETTLED && claim.status !== ClaimStatus.CLAIM_STATUS_REJECTED}
									<Badge variant="outline">
										{getDaysSinceSubmission(claim.submitted_at)} days
									</Badge>
								{:else}
									<span class="text-muted-foreground">-</span>
								{/if}
							</Table.Cell>
							<Table.Cell>
								{#if claim.fraud_check}
									<div class="flex items-center gap-2">
										<span class="text-sm">{claim.fraud_check.fraud_score}</span>
										{#if claim.fraud_check.flagged}
											<AlertTriangle class="h-4 w-4 text-red-500" />
										{/if}
									</div>
								{:else}
									<span class="text-muted-foreground">-</span>
								{/if}
							</Table.Cell>
							<Table.Cell>
								<Badge variant={getClaimStatusColor(claim.status)}>
									{getStatusName(claim.status)}
								</Badge>
							</Table.Cell>
							<Table.Cell class="text-right">
								<Button variant="ghost" size="sm" href="/dashboard/claims/{claim.claim_id}">
									Review
								</Button>
							</Table.Cell>
						</Table.Row>
					{:else}
						<Table.Row>
							<Table.Cell colspan={10} class="text-center text-muted-foreground">
								No claims found
							</Table.Cell>
						</Table.Row>
					{/each}
				</Table.Body>
			</Table.Root>
		</Card.Content>
	</Card.Root>
</div>
