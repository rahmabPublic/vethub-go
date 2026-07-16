<script lang="ts">
	import type { components } from '$lib/types/api';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { Badge } from '$lib/components/ui/badge';
	import { Loader2, X } from 'lucide-svelte';

	type Specialty = components['schemas']['SpecialtyResponse'];

	interface Props {
		firstName?: string;
		lastName?: string;
		selectedSpecialtyIds?: number[];
		specialties: Specialty[];
		onSubmit: (data: { firstName: string; lastName: string; specialtyIds: number[] }) => Promise<void>;
		submitLabel?: string;
	}

	let {
		firstName: initialFirstName = '',
		lastName: initialLastName = '',
		selectedSpecialtyIds: initialSpecialtyIds = [],
		specialties,
		onSubmit,
		submitLabel = 'Save'
	}: Props = $props();

	let firstName = $state(initialFirstName);
	let lastName = $state(initialLastName);
	let selectedIds = $state<number[]>([...initialSpecialtyIds]);
	let submitting = $state(false);

	// Get selected specialties for display
	let selectedSpecialties = $derived(
		specialties.filter((s) => selectedIds.includes(s.id))
	);

	// Get available specialties (not yet selected)
	let availableSpecialties = $derived(
		specialties.filter((s) => !selectedIds.includes(s.id))
	);

	function addSpecialty(id: number) {
		if (!selectedIds.includes(id)) {
			selectedIds = [...selectedIds, id];
		}
	}

	function removeSpecialty(id: number) {
		selectedIds = selectedIds.filter((sid) => sid !== id);
	}

	async function handleSubmit(e: Event) {
		e.preventDefault();
		submitting = true;
		try {
			await onSubmit({
				firstName: firstName.trim(),
				lastName: lastName.trim(),
				specialtyIds: selectedIds
			});
		} finally {
			submitting = false;
		}
	}
</script>

<form onsubmit={handleSubmit} class="space-y-6">
	<div class="grid gap-4 sm:grid-cols-2">
		<div class="space-y-2">
			<Label for="firstName">First Name</Label>
			<Input
				id="firstName"
				bind:value={firstName}
				placeholder="Enter first name"
				required
				disabled={submitting}
			/>
		</div>

		<div class="space-y-2">
			<Label for="lastName">Last Name</Label>
			<Input
				id="lastName"
				bind:value={lastName}
				placeholder="Enter last name"
				required
				disabled={submitting}
			/>
		</div>
	</div>

	<div class="space-y-2">
		<Label>Specialties</Label>
		
		<!-- Selected Specialties -->
		{#if selectedSpecialties.length > 0}
			<div class="flex flex-wrap gap-2 mb-3">
				{#each selectedSpecialties as specialty (specialty.id)}
					<Badge variant="secondary" class="gap-1 pr-1">
						{specialty.name}
						<button
							type="button"
							onclick={() => removeSpecialty(specialty.id)}
							class="ml-1 rounded-full p-0.5 hover:bg-muted-foreground/20"
							disabled={submitting}
						>
							<X class="h-3 w-3" />
						</button>
					</Badge>
				{/each}
			</div>
		{/if}

		<!-- Available Specialties to Add -->
		{#if availableSpecialties.length > 0}
			<div class="flex flex-wrap gap-2">
				{#each availableSpecialties as specialty (specialty.id)}
					<button
						type="button"
						onclick={() => addSpecialty(specialty.id)}
						disabled={submitting}
						class="inline-flex items-center rounded-full border border-dashed border-border px-3 py-1 text-sm text-muted-foreground transition-colors hover:border-primary hover:text-primary"
					>
						+ {specialty.name}
					</button>
				{/each}
			</div>
		{:else if selectedSpecialties.length === specialties.length && specialties.length > 0}
			<p class="text-sm text-muted-foreground">All specialties selected</p>
		{:else}
			<p class="text-sm text-muted-foreground">No specialties available</p>
		{/if}
	</div>

	<div class="flex justify-end gap-3">
		<Button type="button" variant="outline" href="/vets" disabled={submitting}>
			Cancel
		</Button>
		<Button type="submit" disabled={submitting}>
			{#if submitting}
				<Loader2 class="mr-2 h-4 w-4 animate-spin" />
			{/if}
			{submitLabel}
		</Button>
	</div>
</form>
