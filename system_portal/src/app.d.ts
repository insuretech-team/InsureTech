// See https://svelte.dev/docs/kit/types#app.d.ts
// for information about these interfaces
declare global {
	namespace App {
		// interface Error {}
		interface Locals {
			user: {
				/** User ID from the backend */
				id: string;
				/** Session ID from the backend */
				sessionId: string;
				/** Session token for API calls */
				token: string;
				/** User's portal role (SYSTEM_ADMIN, B2B_ORG_ADMIN, etc.) */
				role: string;
				/** User's email address */
				email: string;
				/** User's mobile number */
				mobileNumber: string;
				/** Session expiration timestamp */
				expiresAt: number;
			} | null;

			/** Portal context headers for backend authz */
			portalContext: {
				portal: string;
				userId: string;
				businessId: string;
				tenantId: string;
			} | null;
		}
		// interface PageData {}
		// interface PageState {}
		// interface Platform {}
	}
}

export { };
