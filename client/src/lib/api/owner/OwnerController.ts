/**
 * Owner API Controller
 *
 * Thin API layer for owner-related operations.
 * Handles HTTP communication only - no business logic.
 */
import { client } from '$lib/api/client';
import type {
	OwnerResponse,
	CreateOwnerRequest,
	UpdateOwnerRequest
} from '$lib/api/models';

/**
 * Get all owners
 */
export async function getOwners(): Promise<OwnerResponse[]> {
	const { data, error } = await client.GET('/v1/owners');
	if (error) throw error;
	return data ?? [];
}

/**
 * Get owner by ID
 */
export async function getOwnerById(id: number): Promise<OwnerResponse> {
	const { data, error } = await client.GET('/v1/owners/{id}', {
		params: { path: { id } }
	});
	if (error) throw error;
	if (!data) throw new Error('Owner not found');
	return data;
}

/**
 * Create a new owner
 */
export async function createOwner(request: CreateOwnerRequest): Promise<OwnerResponse> {
	const { data, error } = await client.POST('/v1/owners', {
		body: request
	});
	if (error) throw error;
	if (!data) throw new Error('Failed to create owner');
	return data;
}

/**
 * Update an existing owner
 */
export async function updateOwner(id: number, request: UpdateOwnerRequest): Promise<OwnerResponse> {
	const { data, error } = await client.PUT('/v1/owners/{id}', {
		params: { path: { id } },
		body: request
	});
	if (error) throw error;
	if (!data) throw new Error('Failed to update owner');
	return data;
}

/**
 * Delete an owner
 */
export async function deleteOwner(id: number): Promise<void> {
	const { error } = await client.DELETE('/v1/owners/{id}', {
		params: { path: { id } }
	});
	if (error) throw error;
}
