import createClient from 'openapi-fetch';
import type { paths } from '$lib/types/api';
import { SERVER_BASE_URL, API_USERNAME, API_PASSWORD } from '$lib/config/constants';
import { browser } from '$app/environment';

/**
 * Create Basic Auth header
 */
function createAuthHeader(): string {
	const credentials = `${API_USERNAME}:${API_PASSWORD}`;
	if (browser) {
		return `Basic ${btoa(credentials)}`;
	}
	// Node.js environment (SSR)
	return `Basic ${Buffer.from(credentials).toString('base64')}`;
}

/**
 * Pet Clinic API Client
 * 
 * Uses openapi-fetch for type-safe API calls.
 * Basic Auth is used for the demo backend.
 */
export const client = createClient<paths>({
	baseUrl: SERVER_BASE_URL,
	headers: {
		'Content-Type': 'application/json',
		'Authorization': createAuthHeader()
	}
});

/**
 * Error handling middleware (optional)
 * Logs errors to console for debugging
 */
client.use({
	async onResponse({ response }): Promise<Response | undefined> {
		if (!response.ok) {
			console.error(`API Error: ${response.status} ${response.statusText}`);
		}
		return response;
	}
});
