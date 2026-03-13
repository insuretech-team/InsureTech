<script lang="ts">
	import DashboardLayout from '$lib/components/dashboard-layout.svelte';

	let { children } = $props();
</script>

<DashboardLayout>
	{@render children()}
</DashboardLayout>
