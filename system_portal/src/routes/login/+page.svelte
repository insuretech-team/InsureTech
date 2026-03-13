<script lang="ts">
	import { enhance } from '$app/forms';
	import { Button } from '$lib/components/ui/button';
	import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '$lib/components/ui/card';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { Checkbox } from '$lib/components/ui/checkbox';
	import { LogIn, Eye, EyeOff } from 'lucide-svelte';
	import type { ActionData } from './$types';

	let { form }: { form: ActionData } = $props();

	let email = $state(form?.email || '');
	let password = $state('');
	let rememberMe = $state(false);
	let showPassword = $state(false);
	let loading = $state(false);
</script>

<div class="flex min-h-screen">
	<!-- Left side - Login Form -->
	<div class="flex w-full items-center justify-center bg-background p-8 lg:w-1/2">
		<div class="w-full max-w-md space-y-8">
			<!-- Logo -->
			<div class="text-center">
				<img src="/logo.svg" alt="LabAid InsureTech" class="mx-auto h-12" />
				<h1 class="mt-6 text-3xl font-bold tracking-tight">Welcome Back</h1>
				<p class="mt-2 text-sm text-muted-foreground">
					Sign in to your admin account to continue
				</p>
			</div>

			<!-- Login Card -->
			<Card>
				<CardHeader>
					<CardTitle>Sign In</CardTitle>
					<CardDescription>Enter your credentials to access the admin portal</CardDescription>
				</CardHeader>
				<CardContent>
					<form method="POST" class="space-y-4" use:enhance={() => {
						loading = true;
						return async ({ update }) => {
							await update();
							loading = false;
						};
					}}>
						<!-- Email -->
						<div class="space-y-2">
							<Label for="email">Email Address</Label>
							<Input
								id="email"
								name="email"
								type="email"
								placeholder="test@example.com"
								bind:value={email}
								required
								autocomplete="email"
							/>
						</div>

						<!-- Password -->
						<div class="space-y-2">
							<div class="flex items-center justify-between">
								<Label for="password">Password</Label>
								<a href="/forgot-password" class="text-xs text-primary hover:underline">
									Forgot password?
								</a>
							</div>
							<div class="relative">
								<Input
									id="password"
									name="password"
									type={showPassword ? 'text' : 'password'}
									placeholder="12345"
									bind:value={password}
									required
									autocomplete="current-password"
								/>
								<button
									type="button"
									class="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground"
									onclick={() => (showPassword = !showPassword)}
								>
									{#if showPassword}
										<EyeOff class="h-4 w-4" />
									{:else}
										<Eye class="h-4 w-4" />
									{/if}
								</button>
							</div>
						</div>

						<!-- Error Message -->
						{#if form?.error}
							<div class="rounded-lg bg-destructive/10 p-3 text-sm text-destructive">
								{form.error}
							</div>
						{/if}

						<!-- Remember Me -->
						<div class="flex items-center space-x-2">
							<Checkbox id="remember" bind:checked={rememberMe} />
							<Label
								for="remember"
								class="text-sm font-normal leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
							>
								Remember me for 30 days
							</Label>
						</div>

						<!-- Submit Button -->
						<Button type="submit" class="w-full" disabled={loading}>
							{#if loading}
								<span class="mr-2 h-4 w-4 animate-spin rounded-full border-2 border-background border-t-transparent"></span>
								Signing in...
							{:else}
								<LogIn class="mr-2 h-4 w-4" />
								Sign In
							{/if}
						</Button>
					</form>

					<!-- Demo Credentials -->
					<div class="mt-4 rounded-lg bg-blue-50 p-3 dark:bg-blue-950">
						<p class="text-center text-xs font-medium text-blue-900 dark:text-blue-100">
							Demo Credentials
						</p>
						<p class="mt-1 text-center text-xs text-blue-700 dark:text-blue-300">
							Email: test@example.com | Password: 12345
						</p>
					</div>
				</CardContent>
			</Card>

			<!-- Footer -->
			<p class="text-center text-xs text-muted-foreground">
				© 2024 LabAid InsureTech. All rights reserved.
			</p>
		</div>
	</div>

	<!-- Right side - Brand Showcase -->
	<div class="hidden lg:flex lg:w-1/2 lg:relative lg:overflow-hidden bg-gradient-to-br from-primary via-primary/90 to-primary/80">
		<!-- Background Image with Overlay -->
		<div class="absolute inset-0">
			<img src="/couple.png" alt="Insurance Coverage" class="h-full w-full object-cover opacity-30" />
		</div>
		
		<!-- Brand Content Overlay -->
		<div class="relative flex flex-col justify-center p-12 text-white z-10">
			<div class="space-y-6">
				<div class="inline-block">
					<img src="/logo-header.svg" alt="LabAid InsureTech" class="h-16 brightness-0 invert" />
				</div>
				<h2 class="text-4xl font-bold">Welcome to the Future of Insurance</h2>
				<p class="text-lg text-white/90 max-w-md">
					Manage policies, partners, and claims with our comprehensive admin platform
				</p>
				<div class="flex gap-4 pt-4">
					<div class="flex items-center gap-2">
						<div class="w-2 h-2 bg-accent rounded-full"></div>
						<span class="text-sm">Life Insurance</span>
					</div>
					<div class="flex items-center gap-2">
						<div class="w-2 h-2 bg-accent rounded-full"></div>
						<span class="text-sm">Non-Life Insurance</span>
					</div>
				</div>
			</div>
		</div>
	</div>
</div>
