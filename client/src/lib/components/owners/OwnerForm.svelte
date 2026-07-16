<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import * as Card from '$lib/components/ui/card';
	import { User, Save, Loader2 } from 'lucide-svelte';

	interface OwnerFormData {
		firstName: string;
		lastName: string;
		address: string;
		city: string;
		telephone: string;
		email: string;
	}

	interface Props {
		initialData?: Partial<OwnerFormData>;
		onSubmit: (data: OwnerFormData) => Promise<void>;
		submitLabel?: string;
		isEdit?: boolean;
	}

	let { initialData = {}, onSubmit, submitLabel = 'Save', isEdit = false }: Props = $props();

	let formData = $state<OwnerFormData>({
		firstName: '',
		lastName: '',
		address: '',
		city: '',
		telephone: '',
		email: ''
	});

	// Initialize form data when initialData changes
	$effect(() => {
		formData = {
			firstName: initialData.firstName ?? '',
			lastName: initialData.lastName ?? '',
			address: initialData.address ?? '',
			city: initialData.city ?? '',
			telephone: initialData.telephone ?? '',
			email: initialData.email ?? ''
		};
	});

	let submitting = $state(false);
	let errors = $state<Partial<Record<keyof OwnerFormData, string>>>({});

	function validate(): boolean {
		errors = {};

		if (!formData.firstName.trim()) {
			errors.firstName = 'First name is required';
		}

		if (!formData.lastName.trim()) {
			errors.lastName = 'Last name is required';
		}

		if (formData.telephone && !/^[\d\s\-()]+$/.test(formData.telephone)) {
			errors.telephone = 'Invalid phone number format';
		}

		if (formData.email && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(formData.email)) {
			errors.email = 'Invalid email address format';
		}

		return Object.keys(errors).length === 0;
	}

	async function handleSubmit(e: Event) {
		e.preventDefault();

		if (!validate()) return;

		submitting = true;
		try {
			await onSubmit(formData);
		} finally {
			submitting = false;
		}
	}
</script>

<Card.Root>
	<Card.Header>
		<div class="flex items-center gap-3">
			<div class="flex h-12 w-12 items-center justify-center rounded-lg bg-primary/10">
				<User class="h-6 w-6 text-primary" />
			</div>
			<div>
				<Card.Title>{isEdit ? 'Edit Owner' : 'New Owner'}</Card.Title>
				<Card.Description>
					{isEdit ? 'Update owner information' : 'Register a new pet owner'}
				</Card.Description>
			</div>
		</div>
	</Card.Header>
	<Card.Content>
		<form onsubmit={handleSubmit} class="space-y-6">
			<div class="grid gap-4 md:grid-cols-2">
				<div class="space-y-2">
					<Label for="firstName">First Name *</Label>
					<Input
						id="firstName"
						bind:value={formData.firstName}
						placeholder="George"
						class={errors.firstName ? 'border-destructive' : ''}
					/>
					{#if errors.firstName}
						<p class="text-sm text-destructive">{errors.firstName}</p>
					{/if}
				</div>

				<div class="space-y-2">
					<Label for="lastName">Last Name *</Label>
					<Input
						id="lastName"
						bind:value={formData.lastName}
						placeholder="Franklin"
						class={errors.lastName ? 'border-destructive' : ''}
					/>
					{#if errors.lastName}
						<p class="text-sm text-destructive">{errors.lastName}</p>
					{/if}
				</div>
			</div>

			<div class="space-y-2">
				<Label for="address">Address</Label>
				<Input
					id="address"
					bind:value={formData.address}
					placeholder="110 W. Liberty St."
				/>
			</div>

			<div class="grid gap-4 md:grid-cols-2">
				<div class="space-y-2">
					<Label for="city">City</Label>
					<Input id="city" bind:value={formData.city} placeholder="Madison" />
				</div>

				<div class="space-y-2">
					<Label for="telephone">Telephone</Label>
					<Input
						id="telephone"
						type="tel"
						bind:value={formData.telephone}
						placeholder="6085551023"
						class={errors.telephone ? 'border-destructive' : ''}
					/>
					{#if errors.telephone}
						<p class="text-sm text-destructive">{errors.telephone}</p>
					{/if}
				</div>
			</div>

			<div class="space-y-2">
				<Label for="email">Email</Label>
				<Input
					id="email"
					type="email"
					bind:value={formData.email}
					placeholder="george.franklin@example.com"
					class={errors.email ? 'border-destructive' : ''}
				/>
				{#if errors.email}
					<p class="text-sm text-destructive">{errors.email}</p>
				{/if}
			</div>

			<div class="flex justify-end gap-3 pt-4">
				<Button type="button" variant="outline" href="/owners">Cancel</Button>
				<Button type="submit" disabled={submitting} class="gap-2">
					{#if submitting}
						<Loader2 class="h-4 w-4 animate-spin" />
						Saving...
					{:else}
						<Save class="h-4 w-4" />
						{submitLabel}
					{/if}
				</Button>
			</div>
		</form>
	</Card.Content>
</Card.Root>
