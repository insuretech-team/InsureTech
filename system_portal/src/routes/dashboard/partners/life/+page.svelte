<script lang="ts">
	import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '$lib/components/ui/card';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import * as Table from '$lib/components/ui/table';
	import * as Tabs from '$lib/components/ui/tabs';
	import { Hospital, PillBottle, Stethoscope, Ambulance, Plus, Search, Filter, Download } from 'lucide-svelte';
	import { hospitals, pharmacies, doctors, ambulances } from '$lib/data_detailed/partners';

	const partnerTypes = [
		{ name: 'Hospitals', count: hospitals.length, icon: Hospital, color: 'bg-blue-100 text-blue-700 dark:bg-blue-900 dark:text-blue-300' },
		{ name: 'Pharmacies', count: pharmacies.length, icon: PillBottle, color: 'bg-green-100 text-green-700 dark:bg-green-900 dark:text-green-300' },
		{ name: 'Doctors', count: doctors.length, icon: Stethoscope, color: 'bg-purple-100 text-purple-700 dark:bg-purple-900 dark:text-purple-300' },
		{ name: 'Ambulances', count: ambulances.length, icon: Ambulance, color: 'bg-red-100 text-red-700 dark:bg-red-900 dark:text-red-300' }
	];
</script>

<div class="space-y-6">
	<!-- Header -->
	<div class="flex items-center justify-between">
		<div>
			<h1 class="text-3xl font-bold tracking-tight">Life Insurance Partners</h1>
			<p class="text-muted-foreground">Manage healthcare providers, pharmacies, doctors, and ambulance services</p>
		</div>
		<Button>
			<Plus class="mr-2 h-4 w-4" />
			Add Partner
		</Button>
	</div>

	<!-- Stats Grid -->
	<div class="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
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
					<p class="text-xs text-muted-foreground mt-1">Active partners</p>
				</CardContent>
			</Card>
		{/each}
	</div>

	<!-- Partner Tabs -->
	<Tabs.Root value="hospitals" class="w-full">
		<Tabs.List class="grid w-full grid-cols-4">
			<Tabs.Trigger value="hospitals">Hospitals ({hospitals.length})</Tabs.Trigger>
			<Tabs.Trigger value="pharmacies">Pharmacies ({pharmacies.length})</Tabs.Trigger>
			<Tabs.Trigger value="doctors">Doctors ({doctors.length})</Tabs.Trigger>
			<Tabs.Trigger value="ambulances">Ambulances ({ambulances.length})</Tabs.Trigger>
		</Tabs.List>

		<!-- Hospitals Tab -->
		<Tabs.Content value="hospitals" class="space-y-4">
			<Card>
				<CardHeader>
					<div class="flex items-center justify-between">
						<div>
							<CardTitle>Hospital Partners</CardTitle>
							<CardDescription>Healthcare facilities offering cashless treatment and discounts</CardDescription>
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
							<Input type="search" placeholder="Search hospitals..." class="pl-10" />
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
								<Table.Head>Status</Table.Head>
								<Table.Head class="text-right">Actions</Table.Head>
							</Table.Row>
						</Table.Header>
						<Table.Body>
							{#each hospitals as hospital}
								<Table.Row>
									<Table.Cell class="font-medium">{hospital.id}</Table.Cell>
									<Table.Cell>{hospital.name}</Table.Cell>
									<Table.Cell>{hospital.location}</Table.Cell>
									<Table.Cell>
										<Badge variant="secondary">{hospital.services} services</Badge>
									</Table.Cell>
									<Table.Cell>
										<Badge variant="outline" class="bg-blue-50 text-blue-700 dark:bg-blue-900 dark:text-blue-300">
											{hospital.discountPercentage}%
										</Badge>
									</Table.Cell>
									<Table.Cell>
										{#if hospital.cashlessEnabled}
											<Badge class="bg-green-600">Enabled</Badge>
										{:else}
											<Badge variant="secondary">Disabled</Badge>
										{/if}
									</Table.Cell>
									<Table.Cell>
										<Badge variant={hospital.status === 'ACTIVE' ? 'default' : 'secondary'}>
											{hospital.status}
										</Badge>
									</Table.Cell>
									<Table.Cell class="text-right">
										<Button variant="ghost" size="sm" href="/dashboard/partners/{hospital.id}">
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

		<!-- Pharmacies Tab -->
		<Tabs.Content value="pharmacies" class="space-y-4">
			<Card>
				<CardHeader>
					<div class="flex items-center justify-between">
						<div>
							<CardTitle>Pharmacy Partners</CardTitle>
							<CardDescription>Retail pharmacies providing medicine discounts</CardDescription>
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
							<Input type="search" placeholder="Search pharmacies..." class="pl-10" />
						</div>
					</div>

					<Table.Root>
						<Table.Header>
							<Table.Row>
								<Table.Head>Partner ID</Table.Head>
								<Table.Head>Name</Table.Head>
								<Table.Head>Coverage</Table.Head>
								<Table.Head>Outlets</Table.Head>
								<Table.Head>Discount</Table.Head>
								<Table.Head>Cashless</Table.Head>
								<Table.Head>Status</Table.Head>
								<Table.Head class="text-right">Actions</Table.Head>
							</Table.Row>
						</Table.Header>
						<Table.Body>
							{#each pharmacies as pharmacy}
								<Table.Row>
									<Table.Cell class="font-medium">{pharmacy.id}</Table.Cell>
									<Table.Cell>{pharmacy.name}</Table.Cell>
									<Table.Cell>{pharmacy.location}</Table.Cell>
									<Table.Cell>
										<Badge variant="secondary">{pharmacy.outlets} outlets</Badge>
									</Table.Cell>
									<Table.Cell>
										<Badge variant="outline" class="bg-green-50 text-green-700 dark:bg-green-900 dark:text-green-300">
											{pharmacy.discountPercentage}%
										</Badge>
									</Table.Cell>
									<Table.Cell>
										{#if pharmacy.cashlessEnabled}
											<Badge class="bg-green-600">Enabled</Badge>
										{:else}
											<Badge variant="secondary">Disabled</Badge>
										{/if}
									</Table.Cell>
									<Table.Cell>
										<Badge variant={pharmacy.status === 'ACTIVE' ? 'default' : 'secondary'}>
											{pharmacy.status}
										</Badge>
									</Table.Cell>
									<Table.Cell class="text-right">
										<Button variant="ghost" size="sm" href="/dashboard/partners/{pharmacy.id}">
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

		<!-- Doctors Tab -->
		<Tabs.Content value="doctors" class="space-y-4">
			<Card>
				<CardHeader>
					<div class="flex items-center justify-between">
						<div>
							<CardTitle>Doctor Network</CardTitle>
							<CardDescription>Individual healthcare professionals in the partner network</CardDescription>
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
							<Input type="search" placeholder="Search doctors..." class="pl-10" />
						</div>
					</div>

					<Table.Root>
						<Table.Header>
							<Table.Row>
								<Table.Head>Doctor ID</Table.Head>
								<Table.Head>Name</Table.Head>
								<Table.Head>Specialty</Table.Head>
								<Table.Head>Hospital</Table.Head>
								<Table.Head>Discount</Table.Head>
								<Table.Head>Cashless</Table.Head>
								<Table.Head>Status</Table.Head>
								<Table.Head class="text-right">Actions</Table.Head>
							</Table.Row>
						</Table.Header>
						<Table.Body>
							{#each doctors as doctor}
								<Table.Row>
									<Table.Cell class="font-medium">{doctor.id}</Table.Cell>
									<Table.Cell>{doctor.name}</Table.Cell>
									<Table.Cell>
										<Badge variant="outline">{doctor.specialty}</Badge>
									</Table.Cell>
									<Table.Cell>{doctor.location}</Table.Cell>
									<Table.Cell>
										<Badge variant="outline" class="bg-purple-50 text-purple-700 dark:bg-purple-900 dark:text-purple-300">
											{doctor.discountPercentage}%
										</Badge>
									</Table.Cell>
									<Table.Cell>
										{#if doctor.cashlessEnabled}
											<Badge class="bg-green-600">Enabled</Badge>
										{:else}
											<Badge variant="secondary">Disabled</Badge>
										{/if}
									</Table.Cell>
									<Table.Cell>
										<Badge variant={doctor.status === 'ACTIVE' ? 'default' : 'secondary'}>
											{doctor.status}
										</Badge>
									</Table.Cell>
									<Table.Cell class="text-right">
										<Button variant="ghost" size="sm" href="/dashboard/partners/{doctor.id}">
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

		<!-- Ambulances Tab -->
		<Tabs.Content value="ambulances" class="space-y-4">
			<Card>
				<CardHeader>
					<div class="flex items-center justify-between">
						<div>
							<CardTitle>Ambulance Services</CardTitle>
							<CardDescription>Emergency medical transportation partners</CardDescription>
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
							<Input type="search" placeholder="Search ambulance services..." class="pl-10" />
						</div>
					</div>

					<Table.Root>
						<Table.Header>
							<Table.Row>
								<Table.Head>Service ID</Table.Head>
								<Table.Head>Name</Table.Head>
								<Table.Head>Location</Table.Head>
								<Table.Head>Type</Table.Head>
								<Table.Head>Vehicles</Table.Head>
								<Table.Head>Cashless</Table.Head>
								<Table.Head>Status</Table.Head>
								<Table.Head class="text-right">Actions</Table.Head>
							</Table.Row>
						</Table.Header>
						<Table.Body>
							{#each ambulances as ambulance}
								<Table.Row>
									<Table.Cell class="font-medium">{ambulance.id}</Table.Cell>
									<Table.Cell>{ambulance.name}</Table.Cell>
									<Table.Cell>{ambulance.location}</Table.Cell>
									<Table.Cell>
										<Badge variant="outline">{ambulance.serviceType}</Badge>
									</Table.Cell>
									<Table.Cell>
										<Badge variant="secondary">{ambulance.vehicles} vehicles</Badge>
									</Table.Cell>
									<Table.Cell>
										{#if ambulance.cashlessEnabled}
											<Badge class="bg-green-600">Enabled</Badge>
										{:else}
											<Badge variant="secondary">Disabled</Badge>
										{/if}
									</Table.Cell>
									<Table.Cell>
										<Badge variant={ambulance.status === 'ACTIVE' ? 'default' : 'secondary'}>
											{ambulance.status}
										</Badge>
									</Table.Cell>
									<Table.Cell class="text-right">
										<Button variant="ghost" size="sm" href="/dashboard/partners/{ambulance.id}">
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
</div>
