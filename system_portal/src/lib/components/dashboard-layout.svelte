<script lang="ts">
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import {
		Home,
		Users,
		Building2,
		Package,
		Settings,
		FileText,
		CreditCard,
		AlertCircle,
		BarChart3,
		Menu,
		LogOut,
		User,
		Bell,
		Search,
		ChevronDown
	} from 'lucide-svelte';
	import { Button } from '$lib/components/ui/button';
	import { Badge } from '$lib/components/ui/badge';
	import * as Avatar from '$lib/components/ui/avatar';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import * as Sheet from '$lib/components/ui/sheet';
	import { Separator } from '$lib/components/ui/separator';
	import { Input } from '$lib/components/ui/input';

	let { children } = $props();

	let sidebarOpen = $state(false);

	function handleLogout() {
		// Clear any stored session data
		localStorage.removeItem('user');
		sessionStorage.clear();
		
		// Redirect to login
		goto('/login');
	}

	const navigation = [
		{
			title: 'Overview',
			items: [
				{ name: 'Dashboard', href: '/dashboard', icon: Home },
				{ name: 'Analytics', href: '/dashboard/analytics', icon: BarChart3 }
			]
		},
		{
			title: 'Insurance',
			items: [
				{ name: 'Products', href: '/dashboard/products', icon: Package },
				{ name: 'Policies', href: '/dashboard/policies', icon: FileText },
				{ name: 'Claims', href: '/dashboard/claims', icon: AlertCircle }
			]
		},
		{
			title: 'Partners',
			items: [
				{ name: 'Life Partners', href: '/dashboard/partners/life', icon: Building2 },
				{ name: 'Non-Life Partners', href: '/dashboard/partners/non-life', icon: Building2 },
				{ name: 'Agents', href: '/dashboard/agents', icon: Users }
			]
		},
		{
			title: 'Finance',
			items: [
				{ name: 'Payments', href: '/dashboard/payments', icon: CreditCard },
				{ name: 'Commissions', href: '/dashboard/commissions', icon: CreditCard }
			]
		},
		{
			title: 'System',
			items: [
				{ name: 'Users', href: '/dashboard/users', icon: Users },
				{ name: 'Settings', href: '/dashboard/settings', icon: Settings }
			]
		}
	];

	const lifePartnerTypes = [
		{ label: 'Hospitals', count: 24, variant: 'success' },
		{ label: 'Pharmacies', count: 156, variant: 'info' },
		{ label: 'Doctors', count: 89, variant: 'warning' },
		{ label: 'Ambulances', count: 12, variant: 'destructive' }
	];

	const nonLifePartnerTypes = [
		{ label: 'Auto Repair', count: 45, variant: 'success' },
		{ label: 'Laptop Repair', count: 32, variant: 'info' },
		{ label: 'Mobile Repair', count: 78, variant: 'warning' }
	];
</script>

<!-- Mobile sidebar -->
<Sheet.Root bind:open={sidebarOpen}>
	<Sheet.Content side="left" class="w-80 p-0">
		<div class="flex h-full flex-col">
			<!-- Logo -->
			<div class="flex h-16 items-center border-b border-white/10 px-6 bg-gradient-to-r from-primary/20 via-primary/10 to-transparent">
				<img src="/logo.svg" alt="LabAid InsureTech" class="h-8" />
			</div>

			<!-- Navigation -->
			<nav class="flex-1 space-y-2 overflow-y-auto p-4">
				{#each navigation as section}
					<div class="space-y-1">
						<h3 class="mb-2 px-3 text-xs font-semibold uppercase text-muted-foreground">
							{section.title}
						</h3>
						{#each section.items as item}
							{@const isActive = $page.url.pathname === item.href}
							<a
								href={item.href}
								class="flex items-center gap-3 rounded-lg px-3 py-2 text-sm transition-all hover:bg-white/10 {isActive
									? 'bg-primary text-primary-foreground shadow-md border-l-2 border-l-accent'
									: 'text-white/70 hover:text-white hover:translate-x-0.5'}"
								onclick={() => (sidebarOpen = false)}
							>
								{@const Icon = item.icon}
								<Icon class="h-4 w-4" />
								{item.name}
							</a>
						{/each}
					</div>
					<Separator class="my-2" />
				{/each}
			</nav>
		</div>
	</Sheet.Content>
</Sheet.Root>

<div class="flex h-screen overflow-hidden bg-background">
	<!-- Desktop Sidebar -->
	<aside class="hidden w-64 flex-col border-r bg-card lg:flex">
		<!-- Logo -->
		<div class="flex h-16 items-center border-b border-white/10 px-6 bg-gradient-to-r from-primary/20 via-primary/10 to-transparent">
			<img src="/logo.svg" alt="LabAid InsureTech" class="h-8" />
		</div>

		<!-- Navigation -->
		<nav class="flex-1 space-y-2 overflow-y-auto p-4">
			{#each navigation as section}
				<div class="space-y-1">
					<h3 class="mb-2 px-3 text-xs font-semibold uppercase text-muted-foreground">
						{section.title}
					</h3>
					{#each section.items as item}
						{@const isActive = $page.url.pathname === item.href}
						<a
							href={item.href}
							class="flex items-center gap-3 rounded-lg px-3 py-2 text-sm transition-all hover:bg-white/10 {isActive
								? 'bg-primary text-primary-foreground shadow-md border-l-2 border-l-accent'
								: 'text-white/70 hover:text-white hover:translate-x-0.5'}"
						>
							{@const Icon = item.icon}
							<Icon class="h-4 w-4" />
							{item.name}
						</a>
					{/each}
				</div>
				<Separator class="my-2" />
			{/each}
		</nav>

		<!-- User section -->
		<div class="border-t p-4">
			<DropdownMenu.Root>
				<DropdownMenu.Trigger asChild>
					{#snippet child({ props })}
						<Button
							{...props}
							variant="ghost"
							class="w-full justify-start gap-3 px-3"
						>
						<Avatar.Root class="h-8 w-8">
							<Avatar.Fallback class="bg-primary text-primary-foreground">AD</Avatar.Fallback>
						</Avatar.Root>
						<div class="flex-1 text-left">
							<p class="text-sm font-medium">Admin User</p>
							<p class="text-xs text-muted-foreground">admin@labaid.com</p>
						</div>
						<ChevronDown class="h-4 w-4 text-muted-foreground" />
						</Button>
					{/snippet}
				</DropdownMenu.Trigger>
				<DropdownMenu.Content align="end" class="w-56">
					<DropdownMenu.Item>
						<User class="mr-2 h-4 w-4" />
						Profile
					</DropdownMenu.Item>
					<DropdownMenu.Item>
						<Settings class="mr-2 h-4 w-4" />
						Settings
					</DropdownMenu.Item>
					<DropdownMenu.Separator />
					<DropdownMenu.Item class="text-destructive" onclick={() => handleLogout()}>
						<LogOut class="mr-2 h-4 w-4" />
						Logout
					</DropdownMenu.Item>
				</DropdownMenu.Content>
			</DropdownMenu.Root>
		</div>
	</aside>

	<!-- Main Content -->
	<div class="flex flex-1 flex-col overflow-hidden">
		<!-- Header -->
		<header class="flex h-16 items-center gap-4 border-b border-white/10 bg-card/80 backdrop-blur-md px-6 shadow-lg">
			<!-- Mobile menu button -->
			<Button
				variant="ghost"
				size="icon"
				class="lg:hidden"
				onclick={() => (sidebarOpen = true)}
			>
				<Menu class="h-5 w-5" />
			</Button>

			<!-- Search -->
			<div class="flex-1 lg:max-w-md">
				<div class="relative">
					<Search class="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
					<Input type="search" placeholder="Search..." class="pl-10" />
				</div>
			</div>

			<!-- Right section -->
			<div class="flex items-center gap-2">
				<!-- Notifications -->
				<Button variant="ghost" size="icon" class="relative">
					<Bell class="h-5 w-5" />
					<span
						class="absolute right-1.5 top-1.5 h-2 w-2 rounded-full bg-destructive ring-2 ring-card"
					></span>
				</Button>

				<!-- User menu (mobile) -->
				<div class="lg:hidden">
					<DropdownMenu.Root>
						<DropdownMenu.Trigger asChild>
							{#snippet child({ props })}
								<Button {...props} variant="ghost" size="icon">
								<Avatar.Root class="h-8 w-8">
									<Avatar.Fallback class="bg-primary text-primary-foreground">AD</Avatar.Fallback>
								</Avatar.Root>
								</Button>
							{/snippet}
						</DropdownMenu.Trigger>
						<DropdownMenu.Content align="end" class="w-56">
							<DropdownMenu.Label>
								<div class="flex flex-col space-y-1">
									<p class="text-sm font-medium">Admin User</p>
									<p class="text-xs text-muted-foreground">admin@labaid.com</p>
								</div>
							</DropdownMenu.Label>
							<DropdownMenu.Separator />
							<DropdownMenu.Item>
								<User class="mr-2 h-4 w-4" />
								Profile
							</DropdownMenu.Item>
							<DropdownMenu.Item>
								<Settings class="mr-2 h-4 w-4" />
								Settings
							</DropdownMenu.Item>
							<DropdownMenu.Separator />
							<DropdownMenu.Item class="text-destructive" onclick={() => handleLogout()}>
								<LogOut class="mr-2 h-4 w-4" />
								Logout
							</DropdownMenu.Item>
						</DropdownMenu.Content>
					</DropdownMenu.Root>
				</div>
			</div>
		</header>

		<!-- Page Content -->
		<main class="flex-1 overflow-y-auto p-6">
			{@render children()}
		</main>
	</div>
</div>
