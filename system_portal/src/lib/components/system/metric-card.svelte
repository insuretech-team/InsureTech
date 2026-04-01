<script lang="ts">
	import { ArrowUpRight } from 'lucide-svelte';
	import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '$lib/components/ui/card';
	import { Badge } from '$lib/components/ui/badge';
	import type { SystemMetric } from '$lib/system/types';

	let { metric }: { metric: SystemMetric } = $props();
</script>

<Card class="border-white/60 bg-white/85 shadow-[0_24px_60px_-40px_rgba(18,63,80,0.45)] backdrop-blur">
	<CardHeader class="space-y-3">
		<div class="flex items-center justify-between gap-3">
			<Badge variant={metric.tone ?? 'outline'}>{metric.label}</Badge>
			{#if metric.href}
				<a
					href={metric.href}
					class="inline-flex items-center gap-1 text-xs font-medium text-primary transition hover:text-accent"
				>
					Open
					<ArrowUpRight class="h-3.5 w-3.5" />
				</a>
			{/if}
		</div>
		<div>
			<CardTitle class="text-2xl font-semibold tracking-tight text-slate-900">
				{metric.value}
			</CardTitle>
			<CardDescription class="mt-2 text-sm leading-6 text-slate-600">
				{metric.description}
			</CardDescription>
		</div>
	</CardHeader>
	<CardContent class="pt-0" />
</Card>
