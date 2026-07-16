<script lang="ts">
	import { goto } from '$app/navigation';
	import { createVet } from '$lib/api/vet/VetController';
	import { getSpecialties } from '$lib/api/specialty/SpecialtyController';
	import type { SpecialtyResponse, CreateVetRequest } from '$lib/api/models';
	import VetForm from '$lib/components/vets/VetForm.svelte';
	import { Button } from '$lib/components/ui/button';
	import { ArrowLeft, Stethoscope } from 'lucide-svelte';
	import { toast } from 'svelte-sonner';

	let specialties = $state<SpecialtyResponse[]>([]);
	let loading = $state(true);

	async function loadSpecialties() {
		loading = true;
		try {
			specialties = await getSpecialties();
		} catch (err) {
			toast.error('Failed to load specialties');
			console.error('Error:', err);
		} finally {
			loading = false;
		}
	}

	async function handleSubmit(data: {
		firstName: string;
		lastName: string;
		specialtyIds: number[];
	}) {
		try {
			const request: CreateVetRequest = {
				firstName: data.firstName,
				lastName: data.lastName,
				specialtyIds: data.specialtyIds.length > 0 ? data.specialtyIds : undefined
			};

			const vet = await createVet(request);
			toast.success('Veterinarian created successfully');
			goto(`/vets/${vet.id}`);
		} catch (err) {
			toast.error('Failed to create veterinarian');
			console.error('Error:', err);
		}
	}

	// Load specialties on mount
	$effect(() => {
		loadSpecialties();
	});
</script>

<svelte:head>
	<title>New Veterinarian | VetHub</title>
</svelte:head>

<div class="container mx-auto max-w-2xl px-4 py-8">
	<!-- Back button -->
	<Button variant="ghost" href="/vets" class="mb-6 gap-2">
		<ArrowLeft class="h-4 w-4" />
		Back to Veterinarians
	</Button>

	<div class="mb-6">
		<h1 class="text-2xl font-bold">New Veterinarian</h1>
		<p class="text-muted-foreground">Add a new veterinarian to the clinic</p>
	</div>

	{#if loading}
		<div class="card p-12 text-center">
			<div class="mx-auto mb-4 h-8 w-8 animate-spin rounded-full border-4 border-primary border-t-transparent"></div>
			<p class="text-muted-foreground">Loading specialties...</p>
		</div>
	{:else if specialties.length === 0}
		<div class="card p-12 text-center">
			<Stethoscope class="mx-auto mb-4 h-12 w-12 text-muted-foreground/50" />
			<p class="text-muted-foreground mb-4">No specialties available. You can still create a general practice veterinarian.</p>
			<VetForm
				{specialties}
				onSubmit={handleSubmit}
				submitLabel="Create Veterinarian"
			/>
		</div>
	{:else}
		<div class="card p-6">
			<VetForm
				{specialties}
				onSubmit={handleSubmit}
				submitLabel="Create Veterinarian"
			/>
		</div>
	{/if}
</div>
