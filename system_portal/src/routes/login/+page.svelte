<script lang="ts">
	import { enhance } from '$app/forms';
	import { Shield, Mail, ArrowRight, KeyRound } from 'lucide-svelte';
	import BrandLogo from '$lib/components/system/brand-logo.svelte';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '$lib/components/ui/card';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import type { ActionData } from './$types';

	let { form }: { form: ActionData } = $props();

	const step = $derived(form?.step ?? 'request');
	let email = $state('');
	let code = $state('');

	$effect(() => {
		if (form?.email) {
			email = form.email;
		}
	});
</script>

<svelte:head>
	<title>System Portal Login</title>
</svelte:head>

<div class="min-h-screen bg-[radial-gradient(circle_at_top_left,_rgba(3,167,101,0.18),_transparent_24%),linear-gradient(135deg,_#0f3140_0%,_#123f50_36%,_#0d2430_100%)] text-white">
	<div class="mx-auto grid min-h-screen max-w-[1440px] gap-10 px-6 py-8 lg:grid-cols-[1.15fr_0.85fr] lg:px-10 lg:py-10">
		<section class="relative overflow-hidden rounded-[36px] border border-white/10 bg-white/6 p-8 shadow-[0_36px_70px_-44px_rgba(0,0,0,0.55)] backdrop-blur lg:p-12">
			<div class="absolute inset-0 bg-[radial-gradient(circle_at_top_right,_rgba(255,255,255,0.18),_transparent_28%)]"></div>
			<div class="relative flex h-full flex-col justify-between gap-10">
				<div class="space-y-8">
					<Badge variant="secondary" class="w-fit bg-white/12 text-white">InsureTech Platform</Badge>
					<BrandLogo class="h-12 brightness-[2] saturate-0" />
					<div class="space-y-4">
						<h1 class="max-w-xl text-4xl font-semibold tracking-tight sm:text-5xl">
							System command center for products, partners, policies, claims, and control.
						</h1>
						<p class="max-w-2xl text-base leading-7 text-white/76">
							This portal now uses the generated local SDK and the real system-user email OTP flow. Sign in with an operational email account to access the backend-wired console.
						</p>
					</div>
				</div>

				<div class="grid gap-4 sm:grid-cols-3">
					<div class="rounded-[28px] border border-white/10 bg-white/8 p-5">
						<Shield class="h-5 w-5 text-[#7ae3be]" />
						<p class="mt-4 text-sm font-semibold">Server-side sessions</p>
						<p class="mt-2 text-sm leading-6 text-white/68">Gateway-authenticated web sessions with CSRF protection.</p>
					</div>
					<div class="rounded-[28px] border border-white/10 bg-white/8 p-5">
						<Mail class="h-5 w-5 text-[#7ae3be]" />
						<p class="mt-4 text-sm font-semibold">Email OTP login</p>
						<p class="mt-2 text-sm leading-6 text-white/68">Aligned to the documented `SYSTEM_USER` authentication flow.</p>
					</div>
					<div class="rounded-[28px] border border-white/10 bg-white/8 p-5">
						<KeyRound class="h-5 w-5 text-[#7ae3be]" />
						<p class="mt-4 text-sm font-semibold">Real brand assets</p>
						<p class="mt-2 text-sm leading-6 text-white/68">Updated with the actual InsureTech logo and console styling.</p>
					</div>
				</div>
			</div>
		</section>

		<section class="flex items-center">
			<Card class="w-full rounded-[32px] border-white/60 bg-white/96 text-slate-950 shadow-[0_44px_85px_-54px_rgba(0,0,0,0.4)]">
				<CardHeader class="space-y-4 p-8">
					<Badge variant="outline" class="w-fit border-primary/20 bg-primary/5 text-primary">
						System access
					</Badge>
					<div class="space-y-2">
						<CardTitle class="text-3xl font-semibold tracking-tight">Sign in to the system portal</CardTitle>
						<CardDescription class="text-sm leading-6 text-slate-600">
							{#if step === 'verify'}
								Enter the one-time password sent to {form?.maskedEmail ?? email}.
							{:else}
								Request an OTP using the email assigned to the system user account.
							{/if}
						</CardDescription>
					</div>
				</CardHeader>

				<CardContent class="space-y-6 p-8 pt-0">
					{#if form?.error}
						<div class="rounded-2xl border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700">
							{form.error}
						</div>
					{/if}

					{#if form?.message}
						<div class="rounded-2xl border border-emerald-200 bg-emerald-50 px-4 py-3 text-sm text-emerald-700">
							{form.message}
						</div>
					{/if}

					{#if step === 'verify'}
						<form method="POST" action="?/verifyOtp" class="space-y-5" use:enhance>
							<input type="hidden" name="email" value={form?.email ?? email} />
							<input type="hidden" name="otpId" value={form?.otpId ?? ''} />

							<div class="space-y-2">
								<Label for="code">One-time password</Label>
								<Input
									id="code"
									name="code"
									type="text"
									inputmode="numeric"
									pattern="[0-9]{6}"
									maxlength={6}
									placeholder="Enter 6-digit OTP"
									bind:value={code}
									class="h-12 rounded-2xl"
									required
								/>
							</div>

							<Button type="submit" class="h-12 w-full rounded-2xl">
								Verify and continue
								<ArrowRight class="ml-2 h-4 w-4" />
							</Button>
						</form>

						<form method="POST" action="?/sendOtp" class="space-y-4">
							<input type="hidden" name="email" value={form?.email ?? email} />
							<Button type="submit" variant="outline" class="h-12 w-full rounded-2xl">
								Resend OTP
							</Button>
						</form>
					{:else}
						<form method="POST" action="?/sendOtp" class="space-y-5" use:enhance>
							<div class="space-y-2">
								<Label for="email">System user email</Label>
								<Input
									id="email"
									name="email"
									type="email"
									placeholder="operations@insuretech.example"
									bind:value={email}
									class="h-12 rounded-2xl"
									required
								/>
							</div>

							<Button type="submit" class="h-12 w-full rounded-2xl">
								Send OTP
								<ArrowRight class="ml-2 h-4 w-4" />
							</Button>
						</form>
					{/if}

					<div class="rounded-[28px] border border-slate-200 bg-slate-50 p-5">
						<p class="text-sm font-semibold text-slate-900">Access notes</p>
						<ul class="mt-3 space-y-2 text-sm leading-6 text-slate-600">
							<li>Only `SYSTEM_USER` accounts should sign in from this portal.</li>
							<li>OTP delivery and session creation are executed against the local generated SDK dependency.</li>
							<li>Successful login provisions the same portal metadata cookies used by the backend session wiring.</li>
						</ul>
					</div>
				</CardContent>
			</Card>
		</section>
	</div>
</div>
