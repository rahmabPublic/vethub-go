/**
 * Specialty API Controller
 *
 * Thin API layer for veterinary specialty operations.
 * Handles HTTP communication only - no business logic.
 */
import { client } from '$lib/api/client';
import type {
	SpecialtyResponse,
	CreateSpecialtyRequest,
	UpdateSpecialtyRequest
} from '$lib/api/models';

/**
 * Get all specialties
 */
export async function getSpecialties(): Promise<SpecialtyResponse[]> {
	const { data, error } = await client.GET('/v1/specialties');
	if (error) throw error;
	return data ?? [];
}

/**
 * Get specialty by ID
 */
export async function getSpecialtyById(id: number): Promise<SpecialtyResponse> {
	const { data, error } = await client.GET('/v1/specialties/{id}', {
		params: { path: { id } }
	});
	if (error) throw error;
	if (!data) throw new Error('Specialty not found');
	return data;
}

/**
 * Create a new specialty
 */
export async function createSpecialty(request: CreateSpecialtyRequest): Promise<SpecialtyResponse> {
	const { data, error } = await client.POST('/v1/specialties', {
		body: request
	});
	if (error) throw error;
	if (!data) throw new Error('Failed to create specialty');
	return data;
}

/**
 * Update an existing specialty
 */
export async function updateSpecialty(id: number, request: UpdateSpecialtyRequest): Promise<SpecialtyResponse> {
	const { data, error } = await client.PUT('/v1/specialties/{id}', {
		params: { path: { id } },
		body: request
	});
	if (error) throw error;
	if (!data) throw new Error('Failed to update specialty');
	return data;
}

/**
 * Delete a specialty
 */
export async function deleteSpecialty(id: number): Promise<void> {
	const { error } = await client.DELETE('/v1/specialties/{id}', {
		params: { path: { id } }
	});
	if (error) throw error;
}
