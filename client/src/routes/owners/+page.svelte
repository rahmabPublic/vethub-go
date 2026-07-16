<script lang="ts">
	import { getOwners } from '$lib/api/owner/OwnerController';
	import type { OwnerResponse } from '$lib/api/models';
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Badge } from '$lib/components/ui/badge';
	import * as Table from '$lib/components/ui/table';
	import { Users, Plus, Search, Phone, MapPin, Mail, PawPrint } from 'lucide-svelte';
	import { toast } from 'svelte-sonner';

	let owners = $state<OwnerResponse[]>([]);
	let loading = $state(true);
	let searchQuery = $state('');

	// Filtered owners based on search query
	let filteredOwners = $derived(() => {
		if (!searchQuery.trim()) return owners;
		const query = searchQuery.toLowerCase();
		return owners.filter(
			(owner) =>
				owner.firstName?.toLowerCase().includes(query) ||
				owner.lastName?.toLowerCase().includes(query) ||
				owner.city?.toLowerCase().includes(query) ||
				owner.pets?.some((pet) => pet.name?.toLowerCase().includes(query))
		);
	});

	async function loadOwners() {
		loading = true;
		try {
			owners = await getOwners();
		} catch (err) {
			toast.error('Failed to load owners');
			console.error('Error loading owners:', err);
		} finally {
			loading = false;
		}
	}

	// Load owners on mount
	$effect(() => {
		loadOwners();
	});
</script>

<svelte:head>
	<title>Owners | VetHub</title>
</svelte:head>

<div class="container mx-auto px-4 py-8">
	<!-- Header -->
	<div class="mb-8 flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
		<div class="flex items-center gap-3">
			<div class="flex h-12 w-12 items-center justify-center rounded-lg bg-primary/10">
				<Users class="h-6 w-6 text-primary" />
			</div>
			<div>
				<h1 class="text-2xl font-bold text-foreground">Pet Owners</h1>
				<p class="text-sm text-muted-foreground">Manage pet owner information</p>
			</div>
		</div>
		<Button href="/owners/new" class="gap-2">
			<Plus class="h-4 w-4" />
			Add Owner
		</Button>
	</div>

	<!-- Search -->
	<div class="mb-6">
		<div class="relative max-w-md">
			<Search class="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
			<Input
				type="search"
				placeholder="Search by name, city, or pet..."
				bind:value={searchQuery}
				class="pl-10"
			/>
		</div>
	</div>

	<!-- Table -->
	{#if loading}
		<div class="card p-12 text-center">
			<div class="mx-auto mb-4 h-8 w-8 animate-spin rounded-full border-4 border-primary border-t-transparent"></div>
			<p class="text-muted-foreground">Loading owners...</p>
		</div>
	{:else if filteredOwners().length === 0}
		<div class="card p-12 text-center">
			<Users class="mx-auto mb-4 h-12 w-12 text-muted-foreground/50" />
			{#if searchQuery}
				<p class="text-muted-foreground">No owners found matching "{searchQuery}"</p>
			{:else}
				<p class="text-muted-foreground">No owners registered yet</p>
			{/if}
		</div>
	{:else}
		<div class="card overflow-hidden">
			<Table.Root>
				<Table.Header>
					<Table.Row>
						<Table.Head>Name</Table.Head>
						<Table.Head>Contact</Table.Head>
						<Table.Head>Location</Table.Head>
						<Table.Head>Pets</Table.Head>
						<Table.Head class="w-[100px]">Actions</Table.Head>
					</Table.Row>
				</Table.Header>
				<Table.Body>
					{#each filteredOwners() as owner (owner.id)}
						<Table.Row class="hover:bg-muted/50">
							<Table.Cell>
								<a href="/owners/{owner.id}" class="font-medium text-foreground hover:text-primary">
									{owner.firstName} {owner.lastName}
								</a>
							</Table.Cell>
							<Table.Cell>
								<div class="flex flex-col gap-1">
									<div class="flex items-center gap-2 text-sm text-muted-foreground">
										<Phone class="h-3.5 w-3.5" />
										{owner.telephone}
									</div>
									{#if owner.email}
										<div class="flex items-center gap-2 text-sm text-muted-foreground">
											<Mail class="h-3.5 w-3.5" />
											{owner.email}
										</div>
									{/if}
								</div>
							</Table.Cell>
							<Table.Cell>
								<div class="flex items-center gap-2 text-sm text-muted-foreground">
									<MapPin class="h-3.5 w-3.5" />
									{owner.city}
								</div>
							</Table.Cell>
							<Table.Cell>
								<div class="flex flex-wrap gap-1">
									{#each owner.pets ?? [] as pet (pet.id)}
										<Badge variant="secondary" class="gap-1">
											<PawPrint class="h-3 w-3" />
											{pet.name}
										</Badge>
									{/each}
									{#if !owner.pets?.length}
										<span class="text-sm text-muted-foreground">No pets</span>
									{/if}
								</div>
							</Table.Cell>
							<Table.Cell>
								<Button variant="ghost" size="sm" href="/owners/{owner.id}">
									View
								</Button>
							</Table.Cell>
						</Table.Row>
					{/each}
				</Table.Body>
			</Table.Root>
		</div>
		<p class="mt-4 text-sm text-muted-foreground">
			Showing {filteredOwners().length} of {owners.length} owners
		</p>
	{/if}
</div>
