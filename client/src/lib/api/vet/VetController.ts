/**
 * Vet API Controller
 *
 * Thin API layer for veterinarian-related operations.
 * Handles HTTP communication only - no business logic.
 */
import { client } from '$lib/api/client';
import type {
	VetResponse,
	CreateVetRequest,
	UpdateVetRequest
} from '$lib/api/models';

/**
 * Get all vets
 */
export async function getVets(): Promise<VetResponse[]> {
	const { data, error } = await client.GET('/v1/vets');
	if (error) throw error;
	return data ?? [];
}

/**
 * Get vet by ID
 */
export async function getVetById(id: number): Promise<VetResponse> {
	const { data, error } = await client.GET('/v1/vets/{id}', {
		params: { path: { id } }
	});
	if (error) throw error;
	if (!data) throw new Error('Vet not found');
	return data;
}

/**
 * Create a new vet
 */
export async function createVet(request: CreateVetRequest): Promise<VetResponse> {
	const { data, error } = await client.POST('/v1/vets', {
		body: request
	});
	if (error) throw error;
	if (!data) throw new Error('Failed to create vet');
	return data;
}

/**
 * Update an existing vet
 */
export async function updateVet(id: number, request: UpdateVetRequest): Promise<VetResponse> {
	const { data, error } = await client.PUT('/v1/vets/{id}', {
		params: { path: { id } },
		body: request
	});
	if (error) throw error;
	if (!data) throw new Error('Failed to update vet');
	return data;
}

/**
 * Delete a vet
 */
export async function deleteVet(id: number): Promise<void> {
	const { error } = await client.DELETE('/v1/vets/{id}', {
		params: { path: { id } }
	});
	if (error) throw error;
}
