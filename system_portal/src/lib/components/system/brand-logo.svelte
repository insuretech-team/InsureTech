<script lang="ts">
	let {
		compact = false,
		class: className = '',
		label = 'InsureTech System Portal'
	}: {
		compact?: boolean;
		class?: string;
		label?: string;
	} = $props();
</script>

{#if compact}
	<img
		src="/logos/insuretech-system-icon.svg"
		alt={label}
		class={`h-10 w-10 ${className}`.trim()}
	/>
{:else}
	<img
		src="/logos/insuretech-system-logo.svg"
		alt={label}
		class={`h-10 w-auto ${className}`.trim()}
	/>
{/if}
