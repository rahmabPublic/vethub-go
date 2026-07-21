<script lang="ts">
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { getPetForOwner, updatePetForOwner } from '$lib/api/pet/PetController';
	import { getPetTypes } from '$lib/api/pet-type/PetTypeController';
	import type { PetResponse, PetTypeResponse, UpdatePetRequest } from '$lib/api/models';
	import PetForm from '$lib/components/pets/PetForm.svelte';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { ArrowLeft, PawPrint } from 'lucide-svelte';
	import { toast } from 'svelte-sonner';
	import {
		getPetPictureSrc,
		hasCustomPetPicture,
		removePetPicture,
		setPetPictureFromFile
	} from '$lib/utils/petPicture';

	let pet = $state<PetResponse | null>(null);
	let petTypes = $state<PetTypeResponse[]>([]);
	let loading = $state(true);
	let pictureUpdating = $state(false);

	const ownerId = $derived(Number($page.params.id));
	const petId = $derived(Number($page.params.petId));

	async function loadData() {
		loading = true;
		try {
			// Load pet and pet types in parallel
			const [petData, typesData] = await Promise.all([
				getPetForOwner(ownerId, petId),
				getPetTypes()
			]);

			pet = petData;
			petTypes = typesData;
		} catch (err) {
			toast.error('Failed to load data');
			console.error('Error:', err);
		} finally {
			loading = false;
		}
	}

	async function handleSubmit(data: {
		name: string;
		birthDate: string;
		typeId: number;
	}) {
		try {
			const request: UpdatePetRequest = {
				name: data.name,
				birthDate: data.birthDate,
				typeId: data.typeId
			};

			await updatePetForOwner(ownerId, petId, request);
			toast.success('Pet updated successfully');
			goto(`/owners/${ownerId}/pets/${petId}`);
		} catch (err) {
			toast.error('Failed to update pet');
			console.error('Error:', err);
		}
	}

	async function handlePictureChange(event: Event) {
		const target = event.currentTarget as HTMLInputElement;
		const file = target.files?.[0];
		if (!file || !pet) return;

		pictureUpdating = true;
		try {
			await setPetPictureFromFile(pet.id, file);
			toast.success('Pet picture updated');
		} catch (err) {
			toast.error(err instanceof Error ? err.message : 'Failed to update pet picture');
		} finally {
			pictureUpdating = false;
			target.value = '';
		}
	}

	function handleRemovePicture() {
		if (!pet) return;
		pictureUpdating = true;
		removePetPicture(pet.id);
		toast.success('Pet picture removed');
		pictureUpdating = false;
	}

	// Load data on mount
	$effect(() => {
		if (ownerId && petId) {
			loadData();
		}
	});
</script>

<svelte:head>
	<title>Edit {pet ? pet.name : 'Pet'} | VetHub</title>
</svelte:head>

<div class="container mx-auto max-w-2xl px-4 py-8">
	<!-- Back button -->
	<Button variant="ghost" href="/owners/{ownerId}/pets/{petId}" class="mb-6 gap-2">
		<ArrowLeft class="h-4 w-4" />
		Back to Pet
	</Button>

	<div class="mb-6">
		<h1 class="text-2xl font-bold">Edit Pet</h1>
		<p class="text-muted-foreground">Update pet information</p>
	</div>

	{#if loading}
		<div class="card p-12 text-center">
			<div class="mx-auto mb-4 h-8 w-8 animate-spin rounded-full border-4 border-primary border-t-transparent"></div>
			<p class="text-muted-foreground">Loading pet...</p>
		</div>
	{:else if !pet}
		<div class="card p-12 text-center">
			<PawPrint class="mx-auto mb-4 h-12 w-12 text-muted-foreground/50" />
			<p class="text-muted-foreground">Pet not found</p>
		</div>
	{:else}
		<div class="card mb-6 p-6">
			<div class="mb-4 flex items-center gap-4">
				{#if getPetPictureSrc(pet.id, pet.type?.name)}
					<img
						src={getPetPictureSrc(pet.id, pet.type?.name)!}
						alt={pet.name}
						class="h-20 w-20 rounded-full border object-cover"
					/>
				{:else}
					<div class="flex h-20 w-20 items-center justify-center rounded-full bg-accent/10">
						<PawPrint class="h-9 w-9 text-accent" />
					</div>
				{/if}
				<div>
					<h2 class="font-semibold">Pet Picture</h2>
					<p class="text-sm text-muted-foreground">Upload a JPG, PNG, or WebP (max 5 MB)</p>
				</div>
			</div>

			<div class="flex flex-wrap items-center gap-3">
				<Input
					type="file"
					accept="image/jpeg,image/png,image/webp"
					onchange={handlePictureChange}
					disabled={pictureUpdating}
					class="max-w-xs"
				/>
				<Button
					type="button"
					variant="outline"
					onclick={handleRemovePicture}
					disabled={pictureUpdating || !hasCustomPetPicture(pet.id)}
				>
					Remove Picture
				</Button>
			</div>
		</div>

		<div class="card p-6">
			<PetForm
				name={pet.name}
				birthDate={pet.birthDate}
				typeId={pet.type?.id}
				{petTypes}
				onSubmit={handleSubmit}
				submitLabel="Save Changes"
			/>
		</div>
	{/if}
</div>
