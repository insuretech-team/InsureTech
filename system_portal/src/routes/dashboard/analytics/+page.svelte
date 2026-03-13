<script lang="ts">
	import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '$lib/components/ui/card';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import * as Tabs from '$lib/components/ui/tabs';
	import * as Table from '$lib/components/ui/table';
	import { 
		TrendingUp,
		FileText, 
		Users, 
		DollarSign, 
		Activity,
		Download,
		Calendar,
		ArrowUpRight
	} from 'lucide-svelte';
	import {
		kpiMetrics,
		monthlyChartData,
		policyStats,
		claimStats,
		claimDistributionData,
		revenueStats,
		topPartnersData,
		partnerStats,
		discountImpactData,
		ageDistributionData,
		customerStats,
		formatCurrency,
		formatNumber,
		formatPercent
	} from '$lib/data_detailed/analyticsData';

	// Icon mapping
	const iconMap: Record<string, any> = {
		dollar: DollarSign,
		file: FileText,
		activity: Activity,
		users: Users
	};
</script>

<div class="space-y-6">
	<!-- Header -->
	<div class="flex items-center justify-between">
		<div>
			<h1 class="text-3xl font-bold tracking-tight">Analytics Dashboard</h1>
			<p class="text-muted-foreground">Comprehensive insights into insurance operations</p>
		</div>
		<div class="flex gap-2">
			<Button variant="outline" size="sm">
				<Calendar class="mr-2 h-4 w-4" />
				Last 30 Days
			</Button>
			<Button size="sm">
				<Download class="mr-2 h-4 w-4" />
				Export Report
			</Button>
		</div>
	</div>

	<!-- KPI Cards -->
	<div class="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
		{#each kpiMetrics as kpi}
			<Card>
				<CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
					<CardTitle class="text-sm font-medium">{kpi.title}</CardTitle>
					<svelte:component this={iconMap[kpi.icon]} class="h-4 w-4 text-muted-foreground" />
				</CardHeader>
				<CardContent>
					<div class="text-2xl font-bold">{kpi.value}</div>
					<p class="flex items-center gap-1 text-xs text-muted-foreground mt-1">
						<ArrowUpRight class="h-3 w-3 text-green-600" />
						<span class="text-green-600">{kpi.change}</span>
						from last month
					</p>
				</CardContent>
			</Card>
		{/each}
	</div>

	<!-- Tabs -->
	<Tabs.Root value="overview" class="w-full">
		<Tabs.List>
			<Tabs.Trigger value="overview">Overview</Tabs.Trigger>
			<Tabs.Trigger value="policies">Policies</Tabs.Trigger>
			<Tabs.Trigger value="claims">Claims</Tabs.Trigger>
			<Tabs.Trigger value="revenue">Revenue</Tabs.Trigger>
			<Tabs.Trigger value="partners">Partners</Tabs.Trigger>
		</Tabs.List>

		<!-- Overview Tab -->
		<Tabs.Content value="overview" class="space-y-4">
			<div class="grid gap-4 md:grid-cols-2">
				<!-- Monthly Trends -->
				<Card class="md:col-span-2">
					<CardHeader>
						<CardTitle>Monthly Performance</CardTitle>
						<CardDescription>Last 12 months trend</CardDescription>
					</CardHeader>
					<CardContent>
						<div class="space-y-3">
							{#each monthlyChartData.slice(-6) as data}
								<div class="space-y-2">
									<div class="flex items-center justify-between text-sm">
										<span class="font-medium">{data.month}</span>
										<div class="flex gap-4 text-xs text-muted-foreground">
											<span>{formatNumber(data.policies)} policies</span>
											<span>{formatNumber(data.claims)} claims</span>
											<span>{formatCurrency(data.revenue)}</span>
										</div>
									</div>
									<div class="h-2 w-full rounded-full bg-secondary">
										<div class="h-2 rounded-full bg-primary" style="width: {(data.policies / 3000) * 100}%"></div>
									</div>
								</div>
							{/each}
						</div>
					</CardContent>
				</Card>

				<!-- Cashless vs Reimbursement -->
				<Card>
					<CardHeader>
						<CardTitle>Payment Methods</CardTitle>
						<CardDescription>Claim payment distribution</CardDescription>
					</CardHeader>
					<CardContent class="space-y-4">
						<div>
							<div class="flex items-center justify-between mb-2">
								<span class="text-sm">Cashless</span>
								<Badge class="bg-green-600">{formatNumber(claimStats.cashless)}</Badge>
							</div>
							<div class="h-3 w-full rounded-full bg-secondary">
								<div class="h-3 rounded-full bg-green-600" style="width: {(claimStats.cashless / claimStats.total) * 100}%"></div>
							</div>
						</div>
						<div>
							<div class="flex items-center justify-between mb-2">
								<span class="text-sm">Reimbursement</span>
								<Badge variant="secondary">{formatNumber(claimStats.reimbursement)}</Badge>
							</div>
							<div class="h-3 w-full rounded-full bg-secondary">
								<div class="h-3 rounded-full bg-blue-600" style="width: {(claimStats.reimbursement / claimStats.total) * 100}%"></div>
							</div>
						</div>
					</CardContent>
				</Card>

				<!-- Top Partners -->
				<Card>
					<CardHeader>
						<CardTitle>Top Partners</CardTitle>
						<CardDescription>By revenue</CardDescription>
					</CardHeader>
					<CardContent>
						<div class="space-y-3">
							{#each topPartnersData.slice(0, 5) as partner, i}
								<div class="flex items-center justify-between">
									<div class="flex items-center gap-3">
										<div class="flex h-6 w-6 items-center justify-center rounded-full bg-primary text-xs font-bold text-primary-foreground">
											{i + 1}
										</div>
										<div>
											<p class="text-sm font-medium">{partner.name}</p>
											<p class="text-xs text-muted-foreground">{partner.type}</p>
										</div>
									</div>
									<span class="text-sm font-bold text-green-600">{formatCurrency(partner.revenue)}</span>
								</div>
							{/each}
						</div>
					</CardContent>
				</Card>
			</div>
		</Tabs.Content>

		<!-- Policies Tab -->
		<Tabs.Content value="policies" class="space-y-4">
			<div class="grid gap-4 md:grid-cols-3">
				<Card>
					<CardHeader>
						<CardTitle>Total Policies</CardTitle>
					</CardHeader>
					<CardContent>
						<p class="text-3xl font-bold">{formatNumber(policyStats.total)}</p>
						<p class="text-sm text-muted-foreground mt-2">
							<span class="text-green-600">+{formatPercent(policyStats.growthRate)}</span> growth
						</p>
					</CardContent>
				</Card>

				<Card>
					<CardHeader>
						<CardTitle>Life Insurance</CardTitle>
					</CardHeader>
					<CardContent>
						<p class="text-3xl font-bold text-blue-600">{formatNumber(policyStats.life)}</p>
						<p class="text-sm text-muted-foreground mt-2">
							{formatPercent((policyStats.life / policyStats.total) * 100)} of total
						</p>
					</CardContent>
				</Card>

				<Card>
					<CardHeader>
						<CardTitle>Non-Life Insurance</CardTitle>
					</CardHeader>
					<CardContent>
						<p class="text-3xl font-bold text-green-600">{formatNumber(policyStats.nonLife)}</p>
						<p class="text-sm text-muted-foreground mt-2">
							{formatPercent((policyStats.nonLife / policyStats.total) * 100)} of total
						</p>
					</CardContent>
				</Card>
			</div>

			<!-- Age Distribution -->
			<Card>
				<CardHeader>
					<CardTitle>Policy Distribution by Age</CardTitle>
					<CardDescription>Customer demographics</CardDescription>
				</CardHeader>
				<CardContent class="space-y-4">
					{#each ageDistributionData as group}
						<div>
							<div class="flex items-center justify-between mb-2">
								<span class="text-sm font-medium">{group.ageGroup} years</span>
								<span class="text-sm text-muted-foreground">
									{formatNumber(group.count)} ({formatPercent(group.percentage)})
								</span>
							</div>
							<div class="h-3 w-full rounded-full bg-secondary">
								<div class="h-3 rounded-full bg-blue-600" style="width: {group.percentage}%"></div>
							</div>
						</div>
					{/each}
				</CardContent>
			</Card>
		</Tabs.Content>

		<!-- Claims Tab -->
		<Tabs.Content value="claims" class="space-y-4">
			<div class="grid gap-4 md:grid-cols-4">
				<Card>
					<CardHeader class="pb-2">
						<CardTitle class="text-sm">Total</CardTitle>
					</CardHeader>
					<CardContent>
						<p class="text-2xl font-bold">{formatNumber(claimStats.total)}</p>
					</CardContent>
				</Card>
				<Card>
					<CardHeader class="pb-2">
						<CardTitle class="text-sm">Approved</CardTitle>
					</CardHeader>
					<CardContent>
						<p class="text-2xl font-bold text-green-600">{formatNumber(claimStats.approved)}</p>
					</CardContent>
				</Card>
				<Card>
					<CardHeader class="pb-2">
						<CardTitle class="text-sm">Pending</CardTitle>
					</CardHeader>
					<CardContent>
						<p class="text-2xl font-bold text-orange-600">{formatNumber(claimStats.pending)}</p>
					</CardContent>
				</Card>
				<Card>
					<CardHeader class="pb-2">
						<CardTitle class="text-sm">Rejected</CardTitle>
					</CardHeader>
					<CardContent>
						<p class="text-2xl font-bold text-red-600">{formatNumber(claimStats.rejected)}</p>
					</CardContent>
				</Card>
			</div>

			<!-- Claim Distribution Table -->
			<Card>
				<CardHeader>
					<CardTitle>Claim Distribution</CardTitle>
					<CardDescription>By service type</CardDescription>
				</CardHeader>
				<CardContent>
					<Table.Root>
						<Table.Header>
							<Table.Row>
								<Table.Head>Type</Table.Head>
								<Table.Head class="text-right">Count</Table.Head>
								<Table.Head class="text-right">Amount (M)</Table.Head>
								<Table.Head class="text-right">Share</Table.Head>
							</Table.Row>
						</Table.Header>
						<Table.Body>
							{#each claimDistributionData as item}
								<Table.Row>
									<Table.Cell class="font-medium">{item.type}</Table.Cell>
									<Table.Cell class="text-right">{formatNumber(item.count)}</Table.Cell>
									<Table.Cell class="text-right">{formatCurrency(item.amount)}</Table.Cell>
									<Table.Cell class="text-right">
										<Badge variant="outline">{formatPercent(item.percentage)}</Badge>
									</Table.Cell>
								</Table.Row>
							{/each}
						</Table.Body>
					</Table.Root>
				</CardContent>
			</Card>
		</Tabs.Content>

		<!-- Revenue Tab -->
		<Tabs.Content value="revenue" class="space-y-4">
			<div class="grid gap-4 md:grid-cols-3">
				<Card>
					<CardHeader>
						<CardTitle>Total Revenue</CardTitle>
					</CardHeader>
					<CardContent>
						<p class="text-3xl font-bold text-green-600">{formatCurrency(revenueStats.total)}</p>
						<p class="text-sm text-muted-foreground mt-2">
							<span class="text-green-600">+{formatPercent(revenueStats.growthRate)}</span> growth
						</p>
					</CardContent>
				</Card>

				<Card>
					<CardHeader>
						<CardTitle>Life Premiums</CardTitle>
					</CardHeader>
					<CardContent>
						<p class="text-3xl font-bold text-blue-600">{formatCurrency(revenueStats.life)}</p>
						<p class="text-sm text-muted-foreground mt-2">
							{formatPercent((revenueStats.life / revenueStats.total) * 100)}
						</p>
					</CardContent>
				</Card>

				<Card>
					<CardHeader>
						<CardTitle>Non-Life Premiums</CardTitle>
					</CardHeader>
					<CardContent>
						<p class="text-3xl font-bold text-green-600">{formatCurrency(revenueStats.nonLife)}</p>
						<p class="text-sm text-muted-foreground mt-2">
							{formatPercent((revenueStats.nonLife / revenueStats.total) * 100)}
						</p>
					</CardContent>
				</Card>
			</div>

			<!-- Revenue Breakdown -->
			<Card>
				<CardHeader>
					<CardTitle>Revenue Breakdown</CardTitle>
				</CardHeader>
				<CardContent class="space-y-4">
					<div>
						<div class="flex items-center justify-between mb-2">
							<span class="text-sm">Life Insurance</span>
							<span class="text-sm font-bold">{formatCurrency(revenueStats.life)}</span>
						</div>
						<div class="h-3 w-full rounded-full bg-secondary">
							<div class="h-3 rounded-full bg-blue-600" style="width: {(revenueStats.life / revenueStats.total) * 100}%"></div>
						</div>
					</div>
					<div>
						<div class="flex items-center justify-between mb-2">
							<span class="text-sm">Non-Life Insurance</span>
							<span class="text-sm font-bold">{formatCurrency(revenueStats.nonLife)}</span>
						</div>
						<div class="h-3 w-full rounded-full bg-secondary">
							<div class="h-3 rounded-full bg-green-600" style="width: {(revenueStats.nonLife / revenueStats.total) * 100}%"></div>
						</div>
					</div>
					<div>
						<div class="flex items-center justify-between mb-2">
							<span class="text-sm">Discounts Given</span>
							<span class="text-sm font-bold text-red-600">-{formatCurrency(revenueStats.discounts)}</span>
						</div>
						<div class="h-3 w-full rounded-full bg-secondary">
							<div class="h-3 rounded-full bg-red-600" style="width: {(revenueStats.discounts / revenueStats.total) * 100}%"></div>
						</div>
					</div>
				</CardContent>
			</Card>
		</Tabs.Content>

		<!-- Partners Tab -->
		<Tabs.Content value="partners" class="space-y-4">
			<div class="grid gap-4 md:grid-cols-4">
				<Card>
					<CardHeader class="pb-2">
						<CardTitle class="text-sm">Total Partners</CardTitle>
					</CardHeader>
					<CardContent>
						<p class="text-2xl font-bold">{formatNumber(partnerStats.total)}</p>
					</CardContent>
				</Card>
				<Card>
					<CardHeader class="pb-2">
						<CardTitle class="text-sm">Life Partners</CardTitle>
					</CardHeader>
					<CardContent>
						<p class="text-2xl font-bold text-blue-600">{formatNumber(partnerStats.life)}</p>
					</CardContent>
				</Card>
				<Card>
					<CardHeader class="pb-2">
						<CardTitle class="text-sm">Non-Life Partners</CardTitle>
					</CardHeader>
					<CardContent>
						<p class="text-2xl font-bold text-green-600">{formatNumber(partnerStats.nonLife)}</p>
					</CardContent>
				</Card>
				<Card>
					<CardHeader class="pb-2">
						<CardTitle class="text-sm">Satisfaction</CardTitle>
					</CardHeader>
					<CardContent>
						<p class="text-2xl font-bold text-yellow-600">{partnerStats.avgSatisfaction} ★</p>
					</CardContent>
				</Card>
			</div>

			<!-- Top Partners Table -->
			<Card>
				<CardHeader>
					<CardTitle>Top Performing Partners</CardTitle>
					<CardDescription>Ranked by revenue generation</CardDescription>
				</CardHeader>
				<CardContent>
					<Table.Root>
						<Table.Header>
							<Table.Row>
								<Table.Head>Rank</Table.Head>
								<Table.Head>Partner</Table.Head>
								<Table.Head>Type</Table.Head>
								<Table.Head class="text-right">Claims</Table.Head>
								<Table.Head class="text-right">Revenue</Table.Head>
								<Table.Head class="text-right">Rating</Table.Head>
							</Table.Row>
						</Table.Header>
						<Table.Body>
							{#each topPartnersData as partner, i}
								<Table.Row>
									<Table.Cell>
										<div class="flex h-6 w-6 items-center justify-center rounded-full bg-primary text-xs font-bold text-primary-foreground">
											{i + 1}
										</div>
									</Table.Cell>
									<Table.Cell class="font-medium">{partner.name}</Table.Cell>
									<Table.Cell>
										<Badge variant="outline">{partner.type}</Badge>
									</Table.Cell>
									<Table.Cell class="text-right">{formatNumber(partner.claims)}</Table.Cell>
									<Table.Cell class="text-right font-bold text-green-600">
										{formatCurrency(partner.revenue)}
									</Table.Cell>
									<Table.Cell class="text-right">
										<span class="text-yellow-600">{partner.rating} ★</span>
									</Table.Cell>
								</Table.Row>
							{/each}
						</Table.Body>
					</Table.Root>
				</CardContent>
			</Card>

			<!-- Discount Impact Table -->
			<Card>
				<CardHeader>
					<CardTitle>Discount Impact Analysis</CardTitle>
					<CardDescription>Customer savings by partner type</CardDescription>
				</CardHeader>
				<CardContent>
					<Table.Root>
						<Table.Header>
							<Table.Row>
								<Table.Head>Partner Type</Table.Head>
								<Table.Head class="text-right">Avg Discount</Table.Head>
								<Table.Head class="text-right">Claims</Table.Head>
								<Table.Head class="text-right">Savings</Table.Head>
								<Table.Head class="text-right">Satisfaction</Table.Head>
							</Table.Row>
						</Table.Header>
						<Table.Body>
							{#each discountImpactData as item}
								<Table.Row>
									<Table.Cell class="font-medium">{item.partnerType}</Table.Cell>
									<Table.Cell class="text-right">
										<Badge variant="outline">{formatPercent(item.avgDiscount)}</Badge>
									</Table.Cell>
									<Table.Cell class="text-right">{formatNumber(item.claimsCount)}</Table.Cell>
									<Table.Cell class="text-right text-green-600 font-medium">
										{formatCurrency(item.totalSavings)}
									</Table.Cell>
									<Table.Cell class="text-right text-yellow-600">{item.satisfaction} ★</Table.Cell>
								</Table.Row>
							{/each}
						</Table.Body>
					</Table.Root>
				</CardContent>
			</Card>
		</Tabs.Content>
	</Tabs.Root>
</div>
