<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Badge } from '$lib/components/ui/badge';
	import * as Card from '$lib/components/ui/card';
	import * as Table from '$lib/components/ui/table';
	import { Search, Plus, FileText } from 'lucide-svelte';
	import { policiesDemo, formatBDT, getActivePolicies, getStatusColor, searchPolicies } from '$lib/data_detailed/policies_demo';
	import {  PolicyStatus  } from '$lib/types';

	let searchQuery = '';
	let selectedStatus = PolicyStatus.POLICY_STATUS_ACTIVE;

	$: filteredPolicies = searchQuery
		? searchPolicies(searchQuery)
		: policiesDemo.filter(p => p.status === selectedStatus);

	function getStatusName(status: PolicyStatus): string {
		const names: Record<number, string> = {
			1: 'Pending Payment',
			2: 'Active',
			3: 'Grace Period',
			4: 'Lapsed',
			5: 'Suspended',
			6: 'Cancelled',
			7: 'Expired'
		};
		return names[status] || 'Unknown';
	}

	function formatDate(dateStr: string): string {
		return new Date(dateStr).toLocaleDateString('en-BD', {
			year: 'numeric',
			month: 'short',
			day: 'numeric'
		});
	}
</script>

<div class="space-y-6">
	<!-- Header -->
	<div class="flex items-center justify-between">
		<div>
			<h1 class="text-3xl font-bold tracking-tight">Policies</h1>
			<p class="text-muted-foreground">Manage customer insurance policies</p>
		</div>
		<Button>
			<Plus class="mr-2 h-4 w-4" />
			Issue Policy
		</Button>
	</div>

	<!-- Stats Cards -->
	<div class="grid gap-4 md:grid-cols-4">
		<Card.Root>
			<Card.Header class="flex flex-row items-center justify-between space-y-0 pb-2">
				<Card.Title class="text-sm font-medium">Total Policies</Card.Title>
			</Card.Header>
			<Card.Content>
				<div class="text-2xl font-bold">{policiesDemo.length}</div>
				<p class="text-xs text-muted-foreground">All time</p>
			</Card.Content>
		</Card.Root>

		<Card.Root>
			<Card.Header class="flex flex-row items-center justify-between space-y-0 pb-2">
				<Card.Title class="text-sm font-medium">Active Policies</Card.Title>
			</Card.Header>
			<Card.Content>
				<div class="text-2xl font-bold">
					{policiesDemo.filter((p) => p.status === PolicyStatus.POLICY_STATUS_ACTIVE).length}
				</div>
				<p class="text-xs text-muted-foreground">Currently in force</p>
			</Card.Content>
		</Card.Root>

		<Card.Root>
			<Card.Header class="flex flex-row items-center justify-between space-y-0 pb-2">
				<Card.Title class="text-sm font-medium">Total Premium</Card.Title>
			</Card.Header>
			<Card.Content>
				<div class="text-2xl font-bold">
					{formatBDT(
						policiesDemo
							.filter((p) => p.status === PolicyStatus.POLICY_STATUS_ACTIVE)
							.reduce((sum, p) => sum + p.premium_amount, BigInt(0))
					)}
				</div>
				<p class="text-xs text-muted-foreground">Active policies</p>
			</Card.Content>
		</Card.Root>

		<Card.Root>
			<Card.Header class="flex flex-row items-center justify-between space-y-0 pb-2">
				<Card.Title class="text-sm font-medium">Total Coverage</Card.Title>
			</Card.Header>
			<Card.Content>
				<div class="text-2xl font-bold">
					{formatBDT(
						policiesDemo
							.filter((p) => p.status === PolicyStatus.POLICY_STATUS_ACTIVE)
							.reduce((sum, p) => sum + p.sumInsured, BigInt(0))
					)}
				</div>
				<p class="text-xs text-muted-foreground">Sum insured</p>
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
						placeholder="Search by policy number..."
						class="pl-8"
						bind:value={searchQuery}
					/>
				</div>

				<div class="flex gap-2">
					<select
						bind:value={selectedStatus}
						class="flex h-10 rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background"
					>
						<option value={PolicyStatus.POLICY_STATUS_ACTIVE}>Active</option>
						<option value={PolicyStatus.POLICY_STATUS_PENDING_PAYMENT}>Pending Payment</option>
						<option value={PolicyStatus.POLICY_STATUS_GRACE_PERIOD}>Grace Period</option>
						<option value={PolicyStatus.POLICY_STATUS_EXPIRED}>Expired</option>
						<option value={PolicyStatus.POLICY_STATUS_CANCELLED}>Cancelled</option>
					</select>
				</div>
			</div>
		</Card.Header>

		<Card.Content>
			<Table.Root>
				<Table.Header>
					<Table.Row>
						<Table.Head>Policy Number</Table.Head>
						<Table.Head>Customer ID</Table.Head>
						<Table.Head>Premium</Table.Head>
						<Table.Head>Sum Insured</Table.Head>
						<Table.Head>Start Date</Table.Head>
						<Table.Head>End Date</Table.Head>
						<Table.Head>Nominees</Table.Head>
						<Table.Head>Status</Table.Head>
						<Table.Head class="text-right">Actions</Table.Head>
					</Table.Row>
				</Table.Header>
				<Table.Body>
					{#each filteredPolicies as policy}
						<Table.Row>
							<Table.Cell class="font-medium">{policy.policyNumber}</Table.Cell>
							<Table.Cell>{policy.customer_id}</Table.Cell>
							<Table.Cell>{formatBDT(policy.premium_amount)}</Table.Cell>
							<Table.Cell>{formatBDT(policy.sumInsured)}</Table.Cell>
							<Table.Cell>{formatDate(policy.startDate)}</Table.Cell>
							<Table.Cell>{formatDate(policy.endDate)}</Table.Cell>
							<Table.Cell>
								{#if policy.nominees && policy.nominees.length > 0}
									<Badge variant="secondary">{policy.nominees.length} nominees</Badge>
								{:else}
									<span class="text-muted-foreground text-sm">None</span>
								{/if}
							</Table.Cell>
							<Table.Cell>
								<Badge variant={getStatusColor(policy.status)}>
									{getStatusName(policy.status)}
								</Badge>
							</Table.Cell>
							<Table.Cell class="text-right">
								<div class="flex justify-end gap-2">
									<Button variant="ghost" size="sm" href="/dashboard/policies/{policy.policy_id}">
										View
									</Button>
									{#if policy.policyDocumentUrl}
										<Button variant="ghost" size="sm">
											<FileText class="h-4 w-4" />
										</Button>
									{/if}
								</div>
							</Table.Cell>
						</Table.Row>
					{:else}
						<Table.Row>
							<Table.Cell colspan={9} class="text-center text-muted-foreground">
								No policies found
							</Table.Cell>
						</Table.Row>
					{/each}
				</Table.Body>
			</Table.Root>
		</Card.Content>
	</Card.Root>
</div>
