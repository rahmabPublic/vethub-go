/**
 * Pet Type API Controller
 *
 * Thin API layer for pet type operations (e.g., "Cat", "Dog", "Bird").
 * Handles HTTP communication only - no business logic.
 */
import { client } from '$lib/api/client';
import type {
	PetTypeResponse,
	CreatePetTypeRequest,
	UpdatePetTypeRequest
} from '$lib/api/models';

/**
 * Get all pet types
 */
export async function getPetTypes(): Promise<PetTypeResponse[]> {
	const { data, error } = await client.GET('/v1/pet-types');
	if (error) throw error;
	return data ?? [];
}

/**
 * Get pet type by ID
 */
export async function getPetTypeById(id: number): Promise<PetTypeResponse> {
	const { data, error } = await client.GET('/v1/pet-types/{id}', {
		params: { path: { id } }
	});
	if (error) throw error;
	if (!data) throw new Error('Pet type not found');
	return data;
}

/**
 * Create a new pet type
 */
export async function createPetType(request: CreatePetTypeRequest): Promise<PetTypeResponse> {
	const { data, error } = await client.POST('/v1/pet-types', {
		body: request
	});
	if (error) throw error;
	if (!data) throw new Error('Failed to create pet type');
	return data;
}

/**
 * Update an existing pet type
 */
export async function updatePetType(id: number, request: UpdatePetTypeRequest): Promise<PetTypeResponse> {
	const { data, error } = await client.PUT('/v1/pet-types/{id}', {
		params: { path: { id } },
		body: request
	});
	if (error) throw error;
	if (!data) throw new Error('Failed to update pet type');
	return data;
}

/**
 * Delete a pet type
 */
export async function deletePetType(id: number): Promise<void> {
	const { error } = await client.DELETE('/v1/pet-types/{id}', {
		params: { path: { id } }
	});
	if (error) throw error;
}
