<script lang="ts">
	import { goto } from '$app/navigation';
	import { createOwner } from '$lib/api/owner/OwnerController';
	import type { CreateOwnerRequest } from '$lib/api/models';
	import OwnerForm from '$lib/components/owners/OwnerForm.svelte';
	import { Button } from '$lib/components/ui/button';
	import { ArrowLeft } from 'lucide-svelte';
	import { toast } from 'svelte-sonner';

	async function handleSubmit(data: {
		firstName: string;
		lastName: string;
		address: string;
		city: string;
		telephone: string;
	}) {
		try {
			const request: CreateOwnerRequest = {
				firstName: data.firstName,
				lastName: data.lastName,
				address: data.address || undefined,
				city: data.city || undefined,
				telephone: data.telephone || undefined
			};

			const owner = await createOwner(request);
			toast.success('Owner created successfully');
			goto(`/owners/${owner.id}`);
		} catch (err) {
			toast.error('Failed to create owner');
			console.error('Error:', err);
		}
	}
</script>

<svelte:head>
	<title>New Owner | VetHub</title>
</svelte:head>

<div class="container mx-auto max-w-2xl px-4 py-8">
	<!-- Back button -->
	<Button variant="ghost" href="/owners" class="mb-6 gap-2">
		<ArrowLeft class="h-4 w-4" />
		Back to Owners
	</Button>

	<OwnerForm onSubmit={handleSubmit} submitLabel="Create Owner" />
</div>
