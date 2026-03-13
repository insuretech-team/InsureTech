<script lang="ts">
	import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '$lib/components/ui/card';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import * as Table from '$lib/components/ui/table';
	import * as Tabs from '$lib/components/ui/tabs';
	import { Car, Laptop, Smartphone, Plus, Search, Filter, Download, MapPin } from 'lucide-svelte';
	import { autoRepairs, laptopRepairs, mobileRepairs } from '$lib/data_detailed/partners';

	const partnerTypes = [
		{ name: 'Auto Repair', count: autoRepairs.length, icon: Car, color: 'bg-indigo-100 text-indigo-700 dark:bg-indigo-900 dark:text-indigo-300' },
		{ name: 'Laptop Repair', count: laptopRepairs.length, icon: Laptop, color: 'bg-cyan-100 text-cyan-700 dark:bg-cyan-900 dark:text-cyan-300' },
		{ name: 'Mobile Repair', count: mobileRepairs.length, icon: Smartphone, color: 'bg-pink-100 text-pink-700 dark:bg-pink-900 dark:text-pink-300' }
	];
</script>

<div class="space-y-6">
	<!-- Header -->
	<div class="flex items-center justify-between">
		<div>
			<h1 class="text-3xl font-bold tracking-tight">Non-Life Insurance Partners</h1>
			<p class="text-muted-foreground">Manage repair service providers for vehicles, laptops, and mobile devices</p>
		</div>
		<Button>
			<Plus class="mr-2 h-4 w-4" />
			Add Partner
		</Button>
	</div>

	<!-- Stats Grid -->
	<div class="grid gap-4 md:grid-cols-3">
		{#each partnerTypes as type}
			<Card>
				<CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
					<CardTitle class="text-sm font-medium">{type.name}</CardTitle>
					<div class="rounded-full p-2 {type.color}">
						<svelte:component this={type.icon} class="h-4 w-4" />
					</div>
				</CardHeader>
				<CardContent>
					<div class="text-2xl font-bold">{type.count}</div>
					<p class="text-xs text-muted-foreground mt-1">Active service partners</p>
				</CardContent>
			</Card>
		{/each}
	</div>

	<!-- Partner Tabs -->
	<Tabs.Root value="auto" class="w-full">
		<Tabs.List class="grid w-full grid-cols-3">
			<Tabs.Trigger value="auto">Auto Repair ({autoRepairs.length})</Tabs.Trigger>
			<Tabs.Trigger value="laptop">Laptop Repair ({laptopRepairs.length})</Tabs.Trigger>
			<Tabs.Trigger value="mobile">Mobile Repair ({mobileRepairs.length})</Tabs.Trigger>
		</Tabs.List>

		<!-- Auto Repair Tab -->
		<Tabs.Content value="auto" class="space-y-4">
			<Card>
				<CardHeader>
					<div class="flex items-center justify-between">
						<div>
							<CardTitle>Auto Repair Partners</CardTitle>
							<CardDescription>Vehicle repair and maintenance service providers</CardDescription>
						</div>
						<div class="flex gap-2">
							<Button variant="outline" size="sm">
								<Filter class="mr-2 h-4 w-4" />
								Filter
							</Button>
							<Button variant="outline" size="sm">
								<Download class="mr-2 h-4 w-4" />
								Export
							</Button>
						</div>
					</div>
				</CardHeader>
				<CardContent>
					<div class="mb-4">
						<div class="relative">
							<Search class="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
							<Input type="search" placeholder="Search auto repair shops..." class="pl-10" />
						</div>
					</div>

					<Table.Root>
						<Table.Header>
							<Table.Row>
								<Table.Head>Partner ID</Table.Head>
								<Table.Head>Name</Table.Head>
								<Table.Head>Location</Table.Head>
								<Table.Head>Services</Table.Head>
								<Table.Head>Discount</Table.Head>
								<Table.Head>Cashless</Table.Head>
								<Table.Head>Rating</Table.Head>
								<Table.Head>Status</Table.Head>
								<Table.Head class="text-right">Actions</Table.Head>
							</Table.Row>
						</Table.Header>
						<Table.Body>
							{#each autoRepairs as shop}
								<Table.Row>
									<Table.Cell class="font-medium">{shop.id}</Table.Cell>
									<Table.Cell>
										<div>
											<p class="font-medium">{shop.name}</p>
										</div>
									</Table.Cell>
									<Table.Cell>
										<div class="flex items-center gap-1 text-muted-foreground">
											<MapPin class="h-3 w-3" />
											{shop.location}
										</div>
									</Table.Cell>
									<Table.Cell>
										<Badge variant="outline">{shop.brands}</Badge>
									</Table.Cell>
									<Table.Cell>
										<Badge variant="outline" class="bg-indigo-50 text-indigo-700 dark:bg-indigo-900 dark:text-indigo-300">
											{shop.discountPercentage}%
										</Badge>
									</Table.Cell>
									<Table.Cell>
										{#if shop.cashlessEnabled}
											<Badge class="bg-green-600">Yes</Badge>
										{:else}
											<Badge variant="secondary">No</Badge>
										{/if}
									</Table.Cell>
									<Table.Cell>
										<div class="flex items-center gap-1">
											<span class="text-yellow-500">★</span>
											<span class="text-sm">{shop.rating}</span>
										</div>
									</Table.Cell>
									<Table.Cell>
										<Badge variant={shop.status === 'ACTIVE' ? 'default' : 'secondary'}>
											{shop.status}
										</Badge>
									</Table.Cell>
									<Table.Cell class="text-right">
										<Button variant="ghost" size="sm" href="/dashboard/partners/{shop.id}">
											Configure
										</Button>
									</Table.Cell>
								</Table.Row>
							{/each}
						</Table.Body>
					</Table.Root>
				</CardContent>
			</Card>
		</Tabs.Content>

		<!-- Laptop Repair Tab -->
		<Tabs.Content value="laptop" class="space-y-4">
			<Card>
				<CardHeader>
					<div class="flex items-center justify-between">
						<div>
							<CardTitle>Laptop Repair Partners</CardTitle>
							<CardDescription>Computer and laptop repair service providers</CardDescription>
						</div>
						<div class="flex gap-2">
							<Button variant="outline" size="sm">
								<Filter class="mr-2 h-4 w-4" />
								Filter
							</Button>
							<Button variant="outline" size="sm">
								<Download class="mr-2 h-4 w-4" />
								Export
							</Button>
						</div>
					</div>
				</CardHeader>
				<CardContent>
					<div class="mb-4">
						<div class="relative">
							<Search class="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
							<Input type="search" placeholder="Search laptop repair shops..." class="pl-10" />
						</div>
					</div>

					<Table.Root>
						<Table.Header>
							<Table.Row>
								<Table.Head>Partner ID</Table.Head>
								<Table.Head>Name</Table.Head>
								<Table.Head>Location</Table.Head>
								<Table.Head>Brands</Table.Head>
								<Table.Head>Discount</Table.Head>
								<Table.Head>Cashless</Table.Head>
								<Table.Head>Rating</Table.Head>
								<Table.Head>Status</Table.Head>
								<Table.Head class="text-right">Actions</Table.Head>
							</Table.Row>
						</Table.Header>
						<Table.Body>
							{#each laptopRepairs as shop}
								<Table.Row>
									<Table.Cell class="font-medium">{shop.id}</Table.Cell>
									<Table.Cell>{shop.name}</Table.Cell>
									<Table.Cell>
										<div class="flex items-center gap-1 text-muted-foreground">
											<MapPin class="h-3 w-3" />
											{shop.location}
										</div>
									</Table.Cell>
									<Table.Cell>
										<Badge variant="outline">{shop.brands}</Badge>
									</Table.Cell>
									<Table.Cell>
										<Badge variant="outline" class="bg-cyan-50 text-cyan-700 dark:bg-cyan-900 dark:text-cyan-300">
											{shop.discountPercentage}%
										</Badge>
									</Table.Cell>
									<Table.Cell>
										{#if shop.cashlessEnabled}
											<Badge class="bg-green-600">Yes</Badge>
										{:else}
											<Badge variant="secondary">No</Badge>
										{/if}
									</Table.Cell>
									<Table.Cell>
										<div class="flex items-center gap-1">
											<span class="text-yellow-500">★</span>
											<span class="text-sm">{shop.rating}</span>
										</div>
									</Table.Cell>
									<Table.Cell>
										<Badge variant={shop.status === 'ACTIVE' ? 'default' : 'secondary'}>
											{shop.status}
										</Badge>
									</Table.Cell>
									<Table.Cell class="text-right">
										<Button variant="ghost" size="sm" href="/dashboard/partners/{shop.id}">
											Configure
										</Button>
										</Table.Cell>
								</Table.Row>
							{/each}
						</Table.Body>
					</Table.Root>
				</CardContent>
			</Card>
		</Tabs.Content>

		<!-- Mobile Repair Tab -->
		<Tabs.Content value="mobile" class="space-y-4">
			<Card>
				<CardHeader>
					<div class="flex items-center justify-between">
						<div>
							<CardTitle>Mobile Repair Partners</CardTitle>
							<CardDescription>Smartphone and mobile device repair services</CardDescription>
						</div>
						<div class="flex gap-2">
							<Button variant="outline" size="sm">
								<Filter class="mr-2 h-4 w-4" />
								Filter
							</Button>
							<Button variant="outline" size="sm">
								<Download class="mr-2 h-4 w-4" />
								Export
							</Button>
						</div>
					</div>
				</CardHeader>
				<CardContent>
					<div class="mb-4">
						<div class="relative">
							<Search class="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
							<Input type="search" placeholder="Search mobile repair shops..." class="pl-10" />
						</div>
					</div>

					<Table.Root>
						<Table.Header>
							<Table.Row>
								<Table.Head>Partner ID</Table.Head>
								<Table.Head>Name</Table.Head>
								<Table.Head>Location</Table.Head>
								<Table.Head>Brands</Table.Head>
								<Table.Head>Discount</Table.Head>
								<Table.Head>Cashless</Table.Head>
								<Table.Head>Rating</Table.Head>
								<Table.Head>Status</Table.Head>
								<Table.Head class="text-right">Actions</Table.Head>
							</Table.Row>
						</Table.Header>
						<Table.Body>
							{#each mobileRepairs as shop}
								<Table.Row>
									<Table.Cell class="font-medium">{shop.id}</Table.Cell>
									<Table.Cell>{shop.name}</Table.Cell>
									<Table.Cell>
										<div class="flex items-center gap-1 text-muted-foreground">
											<MapPin class="h-3 w-3" />
											{shop.location}
										</div>
									</Table.Cell>
									<Table.Cell>
										<Badge variant="outline">{shop.brands}</Badge>
									</Table.Cell>
									<Table.Cell>
										<Badge variant="outline" class="bg-pink-50 text-pink-700 dark:bg-pink-900 dark:text-pink-300">
											{shop.discountPercentage}%
										</Badge>
									</Table.Cell>
									<Table.Cell>
										{#if shop.cashlessEnabled}
											<Badge class="bg-green-600">Yes</Badge>
										{:else}
											<Badge variant="secondary">No</Badge>
										{/if}
									</Table.Cell>
									<Table.Cell>
										<div class="flex items-center gap-1">
											<span class="text-yellow-500">★</span>
											<span class="text-sm">{shop.rating}</span>
										</div>
									</Table.Cell>
									<Table.Cell>
										<Badge variant={shop.status === 'ACTIVE' ? 'default' : 'secondary'}>
											{shop.status}
										</Badge>
									</Table.Cell>
									<Table.Cell class="text-right">
										<Button variant="ghost" size="sm" href="/dashboard/partners/{shop.id}">
											Configure
										</Button>
									</Table.Cell>
								</Table.Row>
							{/each}
						</Table.Body>
					</Table.Root>
				</CardContent>
			</Card>
		</Tabs.Content>
	</Tabs.Root>

	<!-- Discount & Cashless Summary -->
	<div class="grid gap-4 md:grid-cols-2">
		<Card>
			<CardHeader>
				<CardTitle>Cashless Service Coverage</CardTitle>
				<CardDescription>Partners offering direct cashless repair services</CardDescription>
			</CardHeader>
			<CardContent class="space-y-4">
				<div class="flex items-center justify-between rounded-lg border p-4">
					<div class="flex items-center gap-3">
						<div class="rounded-full bg-indigo-100 p-2 text-indigo-700 dark:bg-indigo-900 dark:text-indigo-300">
							<Car class="h-5 w-5" />
						</div>
						<div>
							<p class="font-medium">Auto Repair</p>
							<p class="text-sm text-muted-foreground">{autoRepairs.filter(p => p.cashlessEnabled).length} of {autoRepairs.length} partners</p>
						</div>
					</div>
					<Badge class="bg-green-600">{Math.round((autoRepairs.filter(p => p.cashlessEnabled).length / autoRepairs.length) * 100)}%</Badge>
				</div>
				<div class="flex items-center justify-between rounded-lg border p-4">
					<div class="flex items-center gap-3">
						<div class="rounded-full bg-cyan-100 p-2 text-cyan-700 dark:bg-cyan-900 dark:text-cyan-300">
							<Laptop class="h-5 w-5" />
						</div>
						<div>
							<p class="font-medium">Laptop Repair</p>
							<p class="text-sm text-muted-foreground">{laptopRepairs.filter(p => p.cashlessEnabled).length} of {laptopRepairs.length} partners</p>
						</div>
					</div>
					<Badge class="bg-green-600">{Math.round((laptopRepairs.filter(p => p.cashlessEnabled).length / laptopRepairs.length) * 100)}%</Badge>
				</div>
				<div class="flex items-center justify-between rounded-lg border p-4">
					<div class="flex items-center gap-3">
						<div class="rounded-full bg-pink-100 p-2 text-pink-700 dark:bg-pink-900 dark:text-pink-300">
							<Smartphone class="h-5 w-5" />
						</div>
						<div>
							<p class="font-medium">Mobile Repair</p>
							<p class="text-sm text-muted-foreground">{mobileRepairs.filter(p => p.cashlessEnabled).length} of {mobileRepairs.length} partners</p>
						</div>
					</div>
					<Badge class="bg-green-600">{Math.round((mobileRepairs.filter(p => p.cashlessEnabled).length / mobileRepairs.length) * 100)}%</Badge>
				</div>
			</CardContent>
		</Card>

		<Card>
			<CardHeader>
				<CardTitle>Average Discount Rates</CardTitle>
				<CardDescription>Discount percentages across partner categories</CardDescription>
			</CardHeader>
			<CardContent class="space-y-4">
				<div class="space-y-2">
					<div class="flex items-center justify-between">
						<span class="text-sm font-medium">Auto Repair Services</span>
						<span class="text-sm font-bold">{Math.round(autoRepairs.reduce((sum, p) => sum + p.discountPercentage, 0) / autoRepairs.length)}% avg</span>
					</div>
					<div class="h-2 w-full rounded-full bg-secondary">
						<div class="h-2 rounded-full bg-indigo-600" style="width: {(autoRepairs.reduce((sum, p) => sum + p.discountPercentage, 0) / autoRepairs.length)}%"></div>
					</div>
				</div>
				<div class="space-y-2">
					<div class="flex items-center justify-between">
						<span class="text-sm font-medium">Laptop Repair Services</span>
						<span class="text-sm font-bold">{Math.round(laptopRepairs.reduce((sum, p) => sum + p.discountPercentage, 0) / laptopRepairs.length)}% avg</span>
					</div>
					<div class="h-2 w-full rounded-full bg-secondary">
						<div class="h-2 rounded-full bg-cyan-600" style="width: {(laptopRepairs.reduce((sum, p) => sum + p.discountPercentage, 0) / laptopRepairs.length)}%"></div>
					</div>
				</div>
				<div class="space-y-2">
					<div class="flex items-center justify-between">
						<span class="text-sm font-medium">Mobile Repair Services</span>
						<span class="text-sm font-bold">{Math.round(mobileRepairs.reduce((sum, p) => sum + p.discountPercentage, 0) / mobileRepairs.length)}% avg</span>
					</div>
					<div class="h-2 w-full rounded-full bg-secondary">
						<div class="h-2 rounded-full bg-pink-600" style="width: {(mobileRepairs.reduce((sum, p) => sum + p.discountPercentage, 0) / mobileRepairs.length)}%"></div>
					</div>
				</div>
			</CardContent>
		</Card>
	</div>
</div>
