<script lang="ts">
	import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '$lib/components/ui/card';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import * as Tabs from '$lib/components/ui/tabs';
	import {
		TrendingUp,
		Users,
		FileText,
		DollarSign,
		Activity,
		Hospital,
		PillBottle,
		Stethoscope,
		Ambulance,
		Car,
		Laptop,
		Smartphone
	} from 'lucide-svelte';

	const stats = [
		{
			title: 'Total Policies',
			value: '2,543',
			change: '+12.3%',
			icon: FileText,
			color: 'text-blue-600'
		},
		{
			title: 'Active Partners',
			value: '436',
			change: '+8.1%',
			icon: Users,
			color: 'text-green-600'
		},
		{
			title: 'Revenue (This Month)',
			value: '৳45.2M',
			change: '+23.5%',
			icon: DollarSign,
			color: 'text-purple-600'
		},
		{
			title: 'Active Claims',
			value: '156',
			change: '-5.2%',
			icon: Activity,
			color: 'text-orange-600'
		}
	];

	const lifePartners = [
		{ type: 'Hospitals', count: 24, icon: Hospital, color: 'bg-blue-100 text-blue-700', status: 'success' },
		{ type: 'Pharmacies', count: 156, icon: PillBottle, color: 'bg-green-100 text-green-700', status: 'success' },
		{ type: 'Doctors', count: 89, icon: Stethoscope, color: 'bg-purple-100 text-purple-700', status: 'info' },
		{ type: 'Ambulances', count: 12, icon: Ambulance, color: 'bg-red-100 text-red-700', status: 'warning' }
	];

	const nonLifePartners = [
		{ type: 'Auto Repair', count: 45, icon: Car, color: 'bg-indigo-100 text-indigo-700', status: 'success' },
		{ type: 'Laptop Repair', count: 32, icon: Laptop, color: 'bg-cyan-100 text-cyan-700', status: 'info' },
		{ type: 'Mobile Repair', count: 78, icon: Smartphone, color: 'bg-pink-100 text-pink-700', status: 'warning' }
	];

	const recentPolicies = [
		{ id: 'POL-001234', customer: 'Ahmed Hassan', type: 'Health', status: 'active', premium: '৳5,000' },
		{ id: 'POL-001235', customer: 'Fatima Khan', type: 'Motor', status: 'pending', premium: '৳12,500' },
		{ id: 'POL-001236', customer: 'Rashid Ali', type: 'Device', status: 'active', premium: '৳2,800' },
		{ id: 'POL-001237', customer: 'Nadia Rahman', type: 'Travel', status: 'active', premium: '৳4,200' }
	];
</script>

<div class="space-y-6">
	<!-- Header -->
	<div class="relative">
		<div class="absolute -inset-2 bg-gradient-to-r from-primary/5 via-accent/5 to-transparent rounded-lg -z-10"></div>
		<h1 class="text-3xl font-bold tracking-tight text-primary">Dashboard</h1>
		<p class="text-muted-foreground mt-1">Welcome to LabAid InsureTech Admin Portal</p>
	</div>

	<!-- Stats Grid -->
	<div class="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
		{#each stats as stat, i}
			{@const cardColors = ['bg-card-blue', 'bg-card-purple', 'bg-card-green', 'bg-card-orange']}
			{@const borderColors = ['border-info', 'border-primary', 'border-accent', 'border-warning']}
			<Card class="{cardColors[i]} border-l-4 {borderColors[i]} hover:shadow-lg transition-all">
				<CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
					<CardTitle class="text-sm font-medium">{stat.title}</CardTitle>
					<div class="rounded-full bg-white/10 p-2">
						<svelte:component this={stat.icon} class="h-4 w-4 text-white" />
					</div>
				</CardHeader>
				<CardContent>
					<div class="text-2xl font-bold">{stat.value}</div>
					<p class="text-xs text-white/70">
						<span class={stat.change.startsWith('+') ? 'text-accent font-semibold' : 'text-destructive font-semibold'}>
							{stat.change}
						</span>
						<span class="ml-1">from last month</span>
					</p>
				</CardContent>
			</Card>
		{/each}
	</div>

	<!-- Partner Overview Tabs -->
	<Tabs.Root value="life" class="w-full">
		<Tabs.List class="grid w-full grid-cols-2">
			<Tabs.Trigger value="life">Life Partners</Tabs.Trigger>
			<Tabs.Trigger value="non-life">Non-Life Partners</Tabs.Trigger>
		</Tabs.List>

		<Tabs.Content value="life" class="space-y-4">
			<div class="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
				{#each lifePartners as partner}
					<Card>
						<CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
							<CardTitle class="text-sm font-medium">{partner.type}</CardTitle>
							<div class="rounded-full p-2 {partner.color}">
								<svelte:component this={partner.icon} class="h-4 w-4" />
							</div>
						</CardHeader>
						<CardContent>
							<div class="text-2xl font-bold">{partner.count}</div>
							<div class="mt-2 flex items-center gap-2">
								<Badge variant={partner.status === 'success' ? 'default' : partner.status === 'info' ? 'secondary' : 'outline'}>
									{partner.status === 'success' ? 'Active' : partner.status === 'info' ? 'Verified' : 'Pending'}
								</Badge>
								<Button variant="ghost" size="sm" class="h-7 text-xs">View All</Button>
							</div>
						</CardContent>
					</Card>
				{/each}
			</div>

			<!-- Discount & Cashless Info -->
			<Card class="border-t-4 border-t-accent bg-card-green">
				<CardHeader>
					<CardTitle>Life Insurance Benefits</CardTitle>
					<CardDescription>Partner networks offering discounts and cashless services</CardDescription>
				</CardHeader>
				<CardContent class="space-y-4">
					<div class="flex items-center justify-between rounded-lg border border-accent/30 bg-accent/10 p-4">
						<div class="space-y-1">
							<p class="text-sm font-medium text-white/70">Cashless Claims</p>
							<p class="text-2xl font-bold text-accent">281 Partners</p>
						</div>
						<Badge variant="default" class="bg-accent hover:bg-accent/90">Active</Badge>
					</div>
					<div class="flex items-center justify-between rounded-lg border border-primary/30 bg-primary/10 p-4">
						<div class="space-y-1">
							<p class="text-sm font-medium text-white/70">Discount Available</p>
							<p class="text-2xl font-bold text-primary">5-25%</p>
						</div>
						<Badge variant="secondary" class="bg-primary/20 text-white hover:bg-primary/30">Variable</Badge>
					</div>
				</CardContent>
			</Card>
		</Tabs.Content>

		<Tabs.Content value="non-life" class="space-y-4">
			<div class="grid gap-4 md:grid-cols-3">
				{#each nonLifePartners as partner}
					<Card>
						<CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
							<CardTitle class="text-sm font-medium">{partner.type}</CardTitle>
							<div class="rounded-full p-2 {partner.color}">
								<svelte:component this={partner.icon} class="h-4 w-4" />
							</div>
						</CardHeader>
						<CardContent>
							<div class="text-2xl font-bold">{partner.count}</div>
							<div class="mt-2 flex items-center gap-2">
								<Badge variant={partner.status === 'success' ? 'default' : partner.status === 'info' ? 'secondary' : 'outline'}>
									{partner.status === 'success' ? 'Active' : partner.status === 'info' ? 'Verified' : 'Pending'}
								</Badge>
								<Button variant="ghost" size="sm" class="h-7 text-xs">View All</Button>
							</div>
						</CardContent>
					</Card>
				{/each}
			</div>

			<!-- Discount & Cashless Info -->
			<Card class="border-t-4 border-t-accent bg-card-green">
				<CardHeader>
					<CardTitle>Non-Life Insurance Benefits</CardTitle>
					<CardDescription>Service partner networks for vehicle, device, and equipment repairs</CardDescription>
				</CardHeader>
				<CardContent class="space-y-4">
					<div class="flex items-center justify-between rounded-lg border border-accent/30 bg-accent/10 p-4">
						<div class="space-y-1">
							<p class="text-sm font-medium text-white/70">Cashless Repairs</p>
							<p class="text-2xl font-bold text-accent">155 Partners</p>
						</div>
						<Badge variant="default" class="bg-accent hover:bg-accent/90">Active</Badge>
					</div>
					<div class="flex items-center justify-between rounded-lg border border-primary/30 bg-primary/10 p-4">
						<div class="space-y-1">
							<p class="text-sm font-medium text-white/70">Service Discount</p>
							<p class="text-2xl font-bold text-primary">10-30%</p>
						</div>
						<Badge variant="secondary" class="bg-primary/20 text-white hover:bg-primary/30">Variable</Badge>
					</div>
				</CardContent>
			</Card>
		</Tabs.Content>
	</Tabs.Root>

	<!-- Recent Policies -->
	<Card class="border-t-2 border-t-primary bg-card-purple">
		<CardHeader>
			<div class="flex items-center justify-between">
				<div>
					<CardTitle class="flex items-center gap-2">
						<span>Recent Policies</span>
						<div class="h-2 w-2 rounded-full bg-accent animate-pulse"></div>
					</CardTitle>
					<CardDescription>Latest policy applications and renewals</CardDescription>
				</div>
				<Button variant="outline" size="sm" class="border-white/20 hover:bg-white/10 hover:text-white">View All</Button>
			</div>
		</CardHeader>
		<CardContent>
			<div class="space-y-4">
				{#each recentPolicies as policy}
					<div class="flex items-center justify-between border-b border-white/10 pb-4 last:border-0 last:pb-0 hover:bg-white/5 -mx-2 px-2 rounded-lg transition-colors">
						<div class="space-y-1">
							<p class="font-medium">{policy.customer}</p>
							<p class="text-sm text-white/60">{policy.id} • {policy.type} Insurance</p>
						</div>
						<div class="flex items-center gap-3">
							<p class="font-semibold text-primary">{policy.premium}</p>
							<Badge variant={policy.status === 'active' ? 'default' : 'secondary'} 
								class={policy.status === 'active' ? 'bg-accent hover:bg-accent/90' : 'bg-white/20'}>
								{policy.status}
							</Badge>
						</div>
					</div>
				{/each}
			</div>
		</CardContent>
	</Card>
</div>
