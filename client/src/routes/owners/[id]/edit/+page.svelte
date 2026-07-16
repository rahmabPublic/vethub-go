<script lang="ts">
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { getOwnerById, updateOwner } from '$lib/api/owner/OwnerController';
	import type { OwnerResponse, UpdateOwnerRequest } from '$lib/api/models';
	import OwnerForm from '$lib/components/owners/OwnerForm.svelte';
	import { Button } from '$lib/components/ui/button';
	import { ArrowLeft, User } from 'lucide-svelte';
	import { toast } from 'svelte-sonner';

	let owner = $state<OwnerResponse | null>(null);
	let loading = $state(true);

	const ownerId = $derived(Number($page.params.id));

	async function loadOwner() {
		loading = true;
		try {
			owner = await getOwnerById(ownerId);
		} catch (err) {
			toast.error('Failed to load owner');
			console.error('Error:', err);
		} finally {
			loading = false;
		}
	}

	async function handleSubmit(data: {
		firstName: string;
		lastName: string;
		address: string;
		city: string;
		telephone: string;
	}) {
		try {
			const request: UpdateOwnerRequest = {
				firstName: data.firstName,
				lastName: data.lastName,
				address: data.address || undefined,
				city: data.city || undefined,
				telephone: data.telephone || undefined
			};

			await updateOwner(ownerId, request);
			toast.success('Owner updated successfully');
			goto(`/owners/${ownerId}`);
		} catch (err) {
			toast.error('Failed to update owner');
			console.error('Error:', err);
		}
	}

	// Load owner on mount
	$effect(() => {
		if (ownerId) {
			loadOwner();
		}
	});
</script>

<svelte:head>
	<title>Edit {owner ? `${owner.firstName} ${owner.lastName}` : 'Owner'} | VetHub</title>
</svelte:head>

<div class="container mx-auto max-w-2xl px-4 py-8">
	<!-- Back button -->
	<Button variant="ghost" href="/owners/{ownerId}" class="mb-6 gap-2">
		<ArrowLeft class="h-4 w-4" />
		Back to Owner
	</Button>

	{#if loading}
		<div class="card p-12 text-center">
			<div class="mx-auto mb-4 h-8 w-8 animate-spin rounded-full border-4 border-primary border-t-transparent"></div>
			<p class="text-muted-foreground">Loading owner...</p>
		</div>
	{:else if !owner}
		<div class="card p-12 text-center">
			<User class="mx-auto mb-4 h-12 w-12 text-muted-foreground/50" />
			<p class="text-muted-foreground">Owner not found</p>
		</div>
	{:else}
		<OwnerForm
			initialData={{
				firstName: owner.firstName ?? '',
				lastName: owner.lastName ?? '',
				address: owner.address ?? '',
				city: owner.city ?? '',
				telephone: owner.telephone ?? ''
			}}
			onSubmit={handleSubmit}
			submitLabel="Save Changes"
			isEdit
		/>
	{/if}
</div>
