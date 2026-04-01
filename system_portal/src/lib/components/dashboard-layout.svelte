<script lang="ts">
	import { page } from '$app/state';
	import {
		AlertCircle,
		BarChart3,
		Building2,
		FileSearch,
		FileText,
		LayoutDashboard,
		LogOut,
		Menu,
		Package,
		Search,
		Shield,
		Users
	} from 'lucide-svelte';
	import type { Snippet } from 'svelte';
	import * as Avatar from '$lib/components/ui/avatar';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Separator } from '$lib/components/ui/separator';
	import * as Sheet from '$lib/components/ui/sheet';
	import BrandLogo from '$lib/components/system/brand-logo.svelte';

	type NavigationItem = {
		label: string;
		href: string;
		icon: typeof LayoutDashboard;
	};

	type NavigationGroup = {
		title: string;
		items: NavigationItem[];
	};

	const navigation: NavigationGroup[] = [
		{
			title: 'Overview',
			items: [
				{ label: 'Dashboard', href: '/dashboard', icon: LayoutDashboard },
				{ label: 'Analytics', href: '/dashboard/analytics', icon: BarChart3 }
			]
		},
		{
			title: 'Operations',
			items: [
				{ label: 'Products', href: '/dashboard/products', icon: Package },
				{ label: 'Policies', href: '/dashboard/policies', icon: FileText },
				{ label: 'Claims', href: '/dashboard/claims', icon: AlertCircle }
			]
		},
		{
			title: 'Network',
			items: [
				{ label: 'Life Partners', href: '/dashboard/partners/life', icon: Users },
				{ label: 'Non-Life Partners', href: '/dashboard/partners/non-life', icon: Building2 },
				{ label: 'Tenants', href: '/dashboard/tenants', icon: Shield }
			]
		},
		{
			title: 'Governance',
			items: [
				{ label: 'Reports', href: '/dashboard/reports', icon: BarChart3 },
				{ label: 'Audit Trail', href: '/dashboard/audit', icon: FileSearch }
			]
		}
	];

	let {
		user,
		children
	}: {
		user: NonNullable<App.Locals['user']>;
		children: Snippet;
	} = $props();

	let mobileOpen = $state(false);

	function isActive(href: string) {
		return href === '/dashboard'
			? page.url.pathname === href
			: page.url.pathname.startsWith(href);
	}

	const displayName = $derived(user.email || 'System operator');
</script>

{#snippet Navigation()}
	<nav class="space-y-6">
		{#each navigation as group}
			<div class="space-y-2">
				<p class="px-3 text-[11px] font-semibold uppercase tracking-[0.24em] text-slate-500">
					{group.title}
				</p>

				<div class="space-y-1">
					{#each group.items as item}
						{@const Icon = item.icon}
						<a
							href={item.href}
							class={`flex items-center gap-3 rounded-2xl px-3 py-2.5 text-sm font-medium transition ${
								isActive(item.href)
									? 'bg-primary text-white shadow-[0_18px_35px_-22px_rgba(18,63,80,0.8)]'
									: 'text-slate-600 hover:bg-slate-100 hover:text-slate-950'
							}`}
							onclick={() => (mobileOpen = false)}
						>
							<Icon class="h-4 w-4" />
							<span>{item.label}</span>
						</a>
					{/each}
				</div>
			</div>
		{/each}
	</nav>
{/snippet}

<Sheet.Root bind:open={mobileOpen}>
	<Sheet.Content side="left" class="w-[20rem] border-none bg-white p-0">
		<div class="flex h-full flex-col">
			<div class="border-b border-slate-200 px-5 py-5">
				<BrandLogo class="h-9" />
				<p class="mt-3 text-sm text-slate-500">System administration and operational control.</p>
			</div>
			<div class="flex-1 overflow-y-auto px-4 py-5">
				{@render Navigation()}
			</div>
		</div>
	</Sheet.Content>
</Sheet.Root>

<div class="min-h-screen bg-[radial-gradient(circle_at_top_left,_rgba(3,167,101,0.14),_transparent_26%),linear-gradient(180deg,_#f8fcfb_0%,_#eef5f6_100%)]">
	<div class="mx-auto flex min-h-screen max-w-[1600px]">
		<aside class="hidden w-[288px] shrink-0 border-r border-white/70 bg-white/70 px-5 py-6 backdrop-blur xl:flex xl:flex-col">
			<div class="space-y-5">
				<div class="rounded-[28px] border border-primary/10 bg-gradient-to-br from-primary to-[#1f5d72] p-5 text-white shadow-[0_30px_60px_-38px_rgba(18,63,80,0.75)]">
					<BrandLogo class="h-10 brightness-[1.8] saturate-0" />
					<p class="mt-4 text-sm leading-6 text-white/78">
						Single place for tenants, products, claims, policy visibility, and governance signals.
					</p>
					<div class="mt-4 flex items-center gap-2">
						<Badge variant="secondary" class="bg-white/14 text-white">System console</Badge>
						<Badge variant="outline" class="border-white/20 text-white">OTP secured</Badge>
					</div>
				</div>

				<div class="rounded-[28px] border border-white/60 bg-white/85 p-4 shadow-[0_24px_60px_-40px_rgba(18,63,80,0.45)] backdrop-blur">
					{@render Navigation()}
				</div>
			</div>

			<div class="mt-auto rounded-[28px] border border-white/60 bg-white/85 p-4 shadow-[0_24px_60px_-40px_rgba(18,63,80,0.35)] backdrop-blur">
				<div class="flex items-center gap-3">
					<Avatar.Root class="h-11 w-11 ring-2 ring-primary/10">
						<Avatar.Fallback class="bg-primary text-white">
							{displayName.slice(0, 1).toUpperCase()}
						</Avatar.Fallback>
					</Avatar.Root>
					<div class="min-w-0">
						<p class="truncate text-sm font-semibold text-slate-900">{displayName}</p>
						<p class="truncate text-xs text-slate-500">{user.role}</p>
					</div>
				</div>

				<Separator class="my-4" />

				<form method="POST" action="/logout">
					<Button type="submit" variant="outline" class="w-full justify-start gap-2 rounded-2xl">
						<LogOut class="h-4 w-4" />
						Sign out
					</Button>
				</form>
			</div>
		</aside>

		<div class="flex min-w-0 flex-1 flex-col">
			<header class="sticky top-0 z-20 border-b border-white/60 bg-white/70 px-4 py-4 backdrop-blur lg:px-8">
				<div class="flex items-center gap-3">
					<Button
						type="button"
						variant="outline"
						size="icon"
						class="rounded-2xl xl:hidden"
						onclick={() => (mobileOpen = true)}
					>
						<Menu class="h-4 w-4" />
					</Button>

					<div class="relative flex-1">
						<Search class="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-slate-400" />
						<Input
							placeholder="Search products, policies, claims, or partner IDs"
							class="h-11 rounded-2xl border-white/70 bg-white/90 pl-10 shadow-none"
						/>
					</div>

					<Badge variant="outline" class="hidden rounded-full border-primary/20 bg-primary/5 px-3 py-1 text-primary sm:inline-flex">
						{user.role}
					</Badge>
				</div>
			</header>

			<main class="min-w-0 flex-1 px-4 py-6 lg:px-8 lg:py-8">
				{@render children()}
			</main>
		</div>
	</div>
</div>
