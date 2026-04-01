<script lang="ts">
	import { Badge } from '$lib/components/ui/badge';
	import type { Snippet } from 'svelte';

	let {
		title,
		description,
		eyebrow = 'System Portal',
		meta,
		actions
	}: {
		title: string;
		description: string;
		eyebrow?: string;
		meta?: string;
		actions?: Snippet;
	} = $props();
</script>

<div class="flex flex-col gap-4 rounded-[28px] border border-white/60 bg-white/75 p-6 shadow-[0_24px_60px_-40px_rgba(18,63,80,0.45)] backdrop-blur lg:flex-row lg:items-end lg:justify-between">
	<div class="space-y-3">
		<Badge variant="outline" class="border-primary/20 bg-primary/5 text-primary">{eyebrow}</Badge>
		<div class="space-y-2">
			<h1 class="text-3xl font-semibold tracking-tight text-slate-950">{title}</h1>
			<p class="max-w-3xl text-sm leading-6 text-slate-600">{description}</p>
		</div>
		{#if meta}
			<p class="text-xs uppercase tracking-[0.24em] text-slate-500">{meta}</p>
		{/if}
	</div>

	{#if actions}
		<div class="flex items-center gap-3">
			{@render actions()}
		</div>
	{/if}
</div>
