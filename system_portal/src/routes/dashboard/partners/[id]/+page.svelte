<script lang="ts">
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { getPartnerById } from '$lib/data_detailed/partners';
	import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '$lib/components/ui/card';
	import { Badge } from '$lib/components/ui/badge';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import * as Tabs from '$lib/components/ui/tabs';
	import * as Select from '$lib/components/ui/select';
	import { Separator } from '$lib/components/ui/separator';
	import DiscountConfigCard from '$lib/components/discount-config-card.svelte';
	import CashlessConfigCard from '$lib/components/cashless-config-card.svelte';
	import { 
		ArrowLeft, 
		Building2, 
		MapPin, 
		Phone, 
		Mail, 
		Save,
		Hospital,
		PillBottle,
		Stethoscope,
		Ambulance,
		Car,
		Laptop,
		Smartphone
	} from 'lucide-svelte';

	const partnerId = $page.params.id;
	
	// Fetch partner data from dummy data
	const partnerData = getPartnerById(partnerId);
	
	if (!partnerData) {
		goto('/dashboard/partners/life');
	}

	let partner = $state(partnerData || {
		id: partnerId,
		name: 'Unknown Partner',
		type: 'HOSPITAL',
		category: 'LIFE',
		status: 'ACTIVE',
		email: '',
		phone: '',
		address: '',
		location: '',
		discountEnabled: false,
		discountPercentage: 0,
		minDiscount: 0,
		maxDiscount: 0,
		discountType: 'SERVICE',
		cashlessEnabled: false,
		cashlessLimit: 0,
		autoApprovalThreshold: 0,
		preAuthRequired: false,
		authValidityDays: 0,
		requiredDocuments: []
	});

	const partnerTypeIcons = {
		HOSPITAL: Hospital,
		PHARMACY: PillBottle,
		DOCTOR: Stethoscope,
		AMBULANCE: Ambulance,
		AUTO_REPAIR: Car,
		LAPTOP_REPAIR: Laptop,
		MOBILE_REPAIR: Smartphone
	};

	const partnerTypeColors = {
		HOSPITAL: 'bg-blue-100 text-blue-700 dark:bg-blue-900 dark:text-blue-300',
		PHARMACY: 'bg-green-100 text-green-700 dark:bg-green-900 dark:text-green-300',
		DOCTOR: 'bg-purple-100 text-purple-700 dark:bg-purple-900 dark:text-purple-300',
		AMBULANCE: 'bg-red-100 text-red-700 dark:bg-red-900 dark:text-red-300',
		AUTO_REPAIR: 'bg-indigo-100 text-indigo-700 dark:bg-indigo-900 dark:text-indigo-300',
		LAPTOP_REPAIR: 'bg-cyan-100 text-cyan-700 dark:bg-cyan-900 dark:text-cyan-300',
		MOBILE_REPAIR: 'bg-pink-100 text-pink-700 dark:bg-pink-900 dark:text-pink-300'
	};

	function handleDiscountSave(data: any) {
		console.log('Saving discount config:', data);
		// TODO: API call to save
		alert('Discount configuration saved!');
	}

	function handleCashlessSave(data: any) {
		console.log('Saving cashless config:', data);
		// TODO: API call to save
		alert('Cashless configuration saved!');
	}

	function handlePartnerUpdate() {
		console.log('Updating partner:', partner);
		// TODO: API call to update
		alert('Partner information updated!');
	}
</script>

<div class="space-y-6">
	<!-- Header -->
	<div class="flex items-center justify-between">
		<div class="flex items-center gap-4">
			<Button variant="ghost" size="icon" onclick={() => goto('/dashboard/partners/life')}>
				<ArrowLeft class="h-5 w-5" />
			</Button>
			<div>
				<div class="flex items-center gap-3">
					<div class="rounded-full p-2 {partnerTypeColors[partner.type]}">
						{@const PartnerIcon = partnerTypeIcons[partner.type]}
						<PartnerIcon class="h-6 w-6" />
					</div>
					<div>
						<h1 class="text-3xl font-bold tracking-tight">{partner.name}</h1>
						<p class="text-muted-foreground">Partner ID: {partner.id}</p>
					</div>
				</div>
			</div>
		</div>
		<div class="flex items-center gap-2">
			<Badge variant={partner.status === 'ACTIVE' ? 'default' : 'secondary'} class="text-sm">
				{partner.status}
			</Badge>
			<Badge variant="outline" class="text-sm">
				{partner.type.replace('_', ' ')}
			</Badge>
		</div>
	</div>

	<!-- Tabs -->
	<Tabs.Root value="details" class="w-full">
		<Tabs.List class="grid w-full grid-cols-3">
			<Tabs.Trigger value="details">Partner Details</Tabs.Trigger>
			<Tabs.Trigger value="discount">Discount Config</Tabs.Trigger>
			<Tabs.Trigger value="cashless">Cashless Config</Tabs.Trigger>
		</Tabs.List>

		<!-- Partner Details Tab -->
		<Tabs.Content value="details" class="space-y-4">
			<Card>
				<CardHeader>
					<CardTitle>Basic Information</CardTitle>
					<CardDescription>Manage partner's basic details and contact information</CardDescription>
				</CardHeader>
				<CardContent class="space-y-4">
					<div class="grid gap-4 md:grid-cols-2">
						<!-- Organization Name -->
						<div class="space-y-2">
							<Label for="org-name">
								Organization Name
								<span class="text-destructive">*</span>
							</Label>
							<Input id="org-name" bind:value={partner.name} placeholder="Enter organization name" />
						</div>

						<!-- Partner Type -->
						<div class="space-y-2">
							<Label for="partner-type">
								Partner Type
								<span class="text-destructive">*</span>
							</Label>
							<Select.Root bind:selected={partner.type}>
								<Select.Trigger id="partner-type">
									{partner.type || 'Select type'}
								</Select.Trigger>
								<Select.Content>
									<Select.Group>
										<Select.Label>Life Insurance Partners</Select.Label>
										<Select.Item value="HOSPITAL">Hospital</Select.Item>
										<Select.Item value="PHARMACY">Pharmacy</Select.Item>
										<Select.Item value="DOCTOR">Doctor</Select.Item>
										<Select.Item value="AMBULANCE">Ambulance</Select.Item>
									</Select.Group>
									<Select.Group>
										<Select.Label>Non-Life Partners</Select.Label>
										<Select.Item value="AUTO_REPAIR">Auto Repair</Select.Item>
										<Select.Item value="LAPTOP_REPAIR">Laptop Repair</Select.Item>
										<Select.Item value="MOBILE_REPAIR">Mobile Repair</Select.Item>
									</Select.Group>
								</Select.Content>
							</Select.Root>
						</div>

						<!-- Email -->
						<div class="space-y-2">
							<Label for="email">
								Email Address
								<span class="text-destructive">*</span>
							</Label>
							<div class="relative">
								<Mail class="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
								<Input id="email" type="email" bind:value={partner.email} placeholder="email@example.com" class="pl-10" />
							</div>
						</div>

						<!-- Phone -->
						<div class="space-y-2">
							<Label for="phone">
								Phone Number
								<span class="text-destructive">*</span>
							</Label>
							<div class="relative">
								<Phone class="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
								<Input id="phone" type="tel" bind:value={partner.phone} placeholder="+880 1XXX-XXXXXX" class="pl-10" />
							</div>
						</div>

						<!-- Trade License -->
						<div class="space-y-2">
							<Label for="trade-license">Trade License Number</Label>
							<Input id="trade-license" bind:value={partner.tradeLicense} placeholder="TL-XXX-XXXX-XXXXXX" />
						</div>

						<!-- TIN Number -->
						<div class="space-y-2">
							<Label for="tin">TIN Number</Label>
							<Input id="tin" bind:value={partner.tinNumber} placeholder="12-digit TIN" maxlength="12" />
						</div>
					</div>

					<!-- Address -->
					<div class="space-y-2">
						<Label for="address">
							Address
							<span class="text-destructive">*</span>
						</Label>
						<div class="relative">
							<MapPin class="absolute left-3 top-3 h-4 w-4 text-muted-foreground" />
							<Input id="address" bind:value={partner.address} placeholder="Full address" class="pl-10" />
						</div>
					</div>

					<Separator />

					<div class="flex justify-end">
						<Button onclick={handlePartnerUpdate}>
							<Save class="mr-2 h-4 w-4" />
							Save Changes
						</Button>
					</div>
				</CardContent>
			</Card>
		</Tabs.Content>

		<!-- Discount Config Tab -->
		<Tabs.Content value="discount">
			<DiscountConfigCard
				{partnerId}
				partnerName={partner.name}
				bind:discountEnabled={partner.discountEnabled}
				bind:discountPercentage={partner.discountPercentage}
				bind:minDiscount={partner.minDiscount}
				bind:maxDiscount={partner.maxDiscount}
				bind:discountType={partner.discountType}
				onSave={handleDiscountSave}
				onCancel={() => console.log('Discount config cancelled')}
			/>
		</Tabs.Content>

		<!-- Cashless Config Tab -->
		<Tabs.Content value="cashless">
			<CashlessConfigCard
				{partnerId}
				partnerName={partner.name}
				bind:cashlessEnabled={partner.cashlessEnabled}
				bind:cashlessLimit={partner.cashlessLimit}
				bind:autoApprovalThreshold={partner.autoApprovalThreshold}
				bind:preAuthRequired={partner.preAuthRequired}
				bind:authValidityDays={partner.authValidityDays}
				bind:requiredDocuments={partner.requiredDocuments}
				onSave={handleCashlessSave}
				onCancel={() => console.log('Cashless config cancelled')}
			/>
		</Tabs.Content>
	</Tabs.Root>
</div>
