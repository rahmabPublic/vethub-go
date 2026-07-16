<script lang="ts">
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { getVetById, updateVet } from '$lib/api/vet/VetController';
	import { getSpecialties } from '$lib/api/specialty/SpecialtyController';
	import type { VetResponse, SpecialtyResponse, UpdateVetRequest } from '$lib/api/models';
	import VetForm from '$lib/components/vets/VetForm.svelte';
	import { Button } from '$lib/components/ui/button';
	import { ArrowLeft, Stethoscope } from 'lucide-svelte';
	import { toast } from 'svelte-sonner';

	let vet = $state<VetResponse | null>(null);
	let specialties = $state<SpecialtyResponse[]>([]);
	let loading = $state(true);

	const vetId = $derived(Number($page.params.id));

	async function loadData() {
		loading = true;
		try {
			// Load vet and specialties in parallel
			const [vetData, specialtiesData] = await Promise.all([
				getVetById(vetId),
				getSpecialties()
			]);

			vet = vetData;
			specialties = specialtiesData;
		} catch (err) {
			toast.error('Failed to load data');
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
			const request: UpdateVetRequest = {
				firstName: data.firstName,
				lastName: data.lastName,
				specialtyIds: data.specialtyIds.length > 0 ? data.specialtyIds : undefined
			};

			await updateVet(vetId, request);
			toast.success('Veterinarian updated successfully');
			goto(`/vets/${vetId}`);
		} catch (err) {
			toast.error('Failed to update veterinarian');
			console.error('Error:', err);
		}
	}

	// Load data on mount
	$effect(() => {
		if (vetId) {
			loadData();
		}
	});
</script>

<svelte:head>
	<title>Edit {vet ? `Dr. ${vet.firstName} ${vet.lastName}` : 'Veterinarian'} | VetHub</title>
</svelte:head>

<div class="container mx-auto max-w-2xl px-4 py-8">
	<!-- Back button -->
	<Button variant="ghost" href="/vets/{vetId}" class="mb-6 gap-2">
		<ArrowLeft class="h-4 w-4" />
		Back to Veterinarian
	</Button>

	<div class="mb-6">
		<h1 class="text-2xl font-bold">Edit Veterinarian</h1>
		<p class="text-muted-foreground">Update veterinarian information</p>
	</div>

	{#if loading}
		<div class="card p-12 text-center">
			<div class="mx-auto mb-4 h-8 w-8 animate-spin rounded-full border-4 border-primary border-t-transparent"></div>
			<p class="text-muted-foreground">Loading veterinarian...</p>
		</div>
	{:else if !vet}
		<div class="card p-12 text-center">
			<Stethoscope class="mx-auto mb-4 h-12 w-12 text-muted-foreground/50" />
			<p class="text-muted-foreground">Veterinarian not found</p>
		</div>
	{:else}
		<div class="card p-6">
			<VetForm
				firstName={vet.firstName}
				lastName={vet.lastName}
				selectedSpecialtyIds={vet.specialties?.map(s => s.id) ?? []}
				{specialties}
				onSubmit={handleSubmit}
				submitLabel="Save Changes"
			/>
		</div>
	{/if}
</div>
