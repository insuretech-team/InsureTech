<script lang="ts">
	import DashboardLayout from '$lib/components/dashboard-layout.svelte';

	let { data, children } = $props();
</script>

<DashboardLayout user={data.user}>
	{#snippet children()}
		{@render children()}
	{/snippet}
</DashboardLayout>
