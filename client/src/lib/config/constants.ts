/**
 * API Configuration
 * 
 * The backend runs on localhost:8080 during development.
 * In production, this could be configured via environment variable.
 */
const baseUrl = import.meta.env.VITE_SERVER_BASE_URL || 'http://localhost:8080';
export const SERVER_BASE_URL: string = `${baseUrl}/api`;

/**
 * Demo API Credentials
 * 
 * For development, we use the Spring Security default user with a fixed password.
 * Configure these via environment variables for different environments.
 */
export const API_USERNAME = import.meta.env.VITE_API_USERNAME || 'user';
export const API_PASSWORD = import.meta.env.VITE_API_PASSWORD || 'password';

/**
 * Application metadata
 */
export const APP_NAME = 'Pet Clinic';
export const APP_VERSION = __APP_VERSION__ || '0.1.0';

/**
 * Declare the global __APP_VERSION__ that's injected by Vite
 */
declare const __APP_VERSION__: string;
