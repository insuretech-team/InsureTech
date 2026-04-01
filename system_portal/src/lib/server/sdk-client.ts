import {
	apiKeyServiceGenerateApiKey,
	apiKeyServiceGetApiKey,
	apiKeyServiceGetUsageStats,
	apiKeyServiceListApiKeys,
	apiKeyServiceRevokeApiKey,
	apiKeyServiceRotateApiKey,
	auditServiceGetAuditEvents,
	auditServiceGetAuditLogs,
	authServiceEmailLogin,
	authServiceGetCurrentSession,
	authServiceGetUserProfile,
	authServiceListSessions,
	authServiceLogout,
	authServiceRevokeAllSessions,
	authServiceRevokeSession,
	authServiceSendEmailOtp,
	claimServiceApproveClaim,
	claimServiceGetClaim,
	claimServiceRejectClaim,
	claimServiceSettleClaim,
	commissionServiceGetCommission,
	commissionServiceListCommissions,
	commissionServiceProcessPayout,
	createInsureTechClient,
	documentServiceDownloadDocument,
	documentServiceGenerateDocument,
	documentServiceGetDocument,
	documentServiceListDocuments,
	documentServiceListDocumentTemplates,
	fraudServiceCreateFraudRule,
	fraudServiceGetFraudAlert,
	fraudServiceListFraudAlerts,
	fraudServiceListFraudRules,
	fraudServiceUpdateFraudRule,
	insurerServiceAddInsurerProduct,
	insurerServiceGetInsurer,
	insurerServiceListInsurerProducts,
	insurerServiceListInsurers,
	insurerServiceUpdateInsurer,
	insurerServiceUpdateInsurerConfig,
	insurerServiceUpdateInsurerProduct,
	notificationServiceGetUserNotifications,
	notificationServiceMarkAsRead,
	notificationServiceSendNotification,
	partnerServiceDeletePartner,
	partnerServiceGetPartner,
	partnerServiceGetPartnerApiCredentials,
	partnerServiceGetPartnerCommission,
	partnerServiceListPartners,
	partnerServiceRotatePartnerApiKey,
	partnerServiceUpdateCommissionStructure,
	partnerServiceUpdatePartner,
	partnerServiceUpdatePartnerStatus,
	partnerServiceVerifyPartner,
	paymentServiceGetPayment,
	paymentServiceListPayments,
	productServiceCreateProduct,
	productServiceDeactivateProduct,
	productServiceGetProduct,
	productServiceListProducts,
	productServiceSearchProducts,
	productServiceUpdateProduct,
	reportServiceExecuteReport,
	reportServiceListReportDefinitions,
	reportServiceListReportExecutions,
	reportServiceListReportSchedules,
	renewalServiceGetGracePeriod,
	renewalServiceGetRenewalSchedule,
	renewalServiceListUpcomingRenewals,
	renewalServiceRenewPolicy,
	renewalServiceRevivePolicy,
	renewalServiceSendRenewalReminder,
	taskServiceAssignTask,
	taskServiceCompleteTask,
	taskServiceCreateTask,
	taskServiceGetTask,
	taskServiceListMyTasks,
	taskServiceUpdateTask,
	tenantServiceCreateTenant,
	tenantServiceGetTenant,
	tenantServiceGetTenantConfig,
	tenantServiceListTenants,
	tenantServiceUpdateTenant,
	tenantServiceUpdateTenantConfig,
	underwritingServiceApproveUnderwriting,
	underwritingServiceConvertQuoteToPolicy,
	underwritingServiceGetQuote,
	underwritingServiceListQuotes,
	underwritingServiceRejectUnderwriting,
	underwritingServiceRequestQuote,
	workflowServiceCompleteTask,
	workflowServiceGetMyTasks,
	workflowServiceGetWorkflowDefinition,
	workflowServiceGetWorkflowInstance,
	workflowServiceStartWorkflow
} from '@lifeplus/insuretech-sdk';
import type { RequestEvent } from '@sveltejs/kit';

function getBaseUrl(): string {
	return (
		process.env.INSURETECH_API_BASE_URL ??
		process.env.VITE_API_URL ??
		process.env.PUBLIC_API_URL ??
		'http://localhost:8080'
	);
}

function getApiKey(): string {
	return process.env.INSURETECH_API_KEY ?? 'system-portal';
}

function extractCookie(cookieHeader: string, name: string): string {
	const escaped = name.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
	const match = cookieHeader.match(new RegExp(`(?:^|;\\s*)${escaped}=([^;]*)`));
	return match ? decodeURIComponent(match[1]) : '';
}

export interface PortalHeaders {
	portal: string;
	userId: string;
	businessId: string;
	tenantId: string;
}

function roleToPortal(role: string): string {
	return role === 'SYSTEM_ADMIN' ? 'PORTAL_SYSTEM' : 'PORTAL_B2B';
}

export function resolvePortalHeaders(event: RequestEvent): PortalHeaders | null {
	const cookieHeader = event.request.headers.get('cookie') ?? '';
	const sessionToken = extractCookie(cookieHeader, 'session_token');

	if (!sessionToken) {
		return null;
	}

	const role = extractCookie(cookieHeader, 'portal_role') || 'SYSTEM_ADMIN';
	const userId = extractCookie(cookieHeader, 'portal_user_id');
	const businessId = extractCookie(cookieHeader, 'portal_biz_id');
	const tenantId =
		process.env.DEFAULT_TENANT_ID?.trim() ||
		extractCookie(cookieHeader, 'portal_tenant_id') ||
		'00000000-0000-0000-0000-000000000001';

	return {
		portal: roleToPortal(role),
		userId,
		businessId,
		tenantId
	};
}

function buildHeaders(event: RequestEvent, sessionOverrides?: Partial<PortalHeaders>) {
	const cookieHeader = event.request.headers.get('cookie') ?? '';
	const csrf = extractCookie(cookieHeader, 'csrf_token');
	const resolved = resolvePortalHeaders(event);

	const headers: Record<string, string> = {};
	if (cookieHeader) headers.cookie = cookieHeader;
	if (csrf) headers['X-CSRF-Token'] = csrf;

	const portal = sessionOverrides?.portal ?? resolved?.portal ?? 'PORTAL_SYSTEM';
	const userId = sessionOverrides?.userId ?? resolved?.userId ?? '';
	const businessId = sessionOverrides?.businessId ?? resolved?.businessId ?? '';
	const tenantId = sessionOverrides?.tenantId ?? resolved?.tenantId ?? '';

	if (portal) headers['x-portal'] = portal;
	if (userId) headers['x-user-id'] = userId;
	if (businessId) headers['x-business-id'] = businessId;
	if (tenantId) headers['x-tenant-id'] = tenantId;

	return headers;
}

function makeDirectHttp(event: RequestEvent, sessionOverrides?: Partial<PortalHeaders>) {
	const baseUrl = getBaseUrl();
	const baseHeaders = buildHeaders(event, sessionOverrides);

	async function request(method: string, path: string, body?: unknown) {
		const response = await fetch(`${baseUrl}${path}`, {
			method,
			cache: 'no-store',
			headers: {
				'Content-Type': 'application/json',
				...baseHeaders
			},
			body: body === undefined ? undefined : JSON.stringify(body)
		});

		const text = await response.text();
		let parsed: Record<string, unknown> = {};

		if (text) {
			try {
				parsed = JSON.parse(text) as Record<string, unknown>;
			} catch {
				parsed = {
					success: false,
					error: {
						message: text
					}
				};
			}
		}

		const success =
			typeof parsed.success === 'boolean' ? parsed.success : response.ok;
		const data = success ? (parsed.data ?? parsed) : null;
		const error =
			!success && parsed.error && typeof parsed.error === 'object'
				? (parsed.error as Record<string, unknown>)
				: null;

		return {
			ok: success,
			status: response.status,
			data,
			error,
			response
		};
	}

	return {
		get: (path: string) => request('GET', path),
		post: (path: string, body?: unknown) => request('POST', path, body),
		patch: (path: string, body?: unknown) => request('PATCH', path, body),
		put: (path: string, body?: unknown) => request('PUT', path, body),
		delete: (path: string) => request('DELETE', path)
	};
}

export function makeSdkClient(event: RequestEvent, sessionOverrides?: Partial<PortalHeaders>) {
	const client = createInsureTechClient({
		apiKey: getApiKey(),
		baseUrl: getBaseUrl(),
		headers: buildHeaders(event, sessionOverrides)
	});

	const wrap = <T extends (options?: any) => any>(fn: T) => {
		return (options?: Record<string, unknown>) =>
			fn({ client, throwOnError: false, ...(options ?? {}) });
	};

	return {
		authServiceSendEmailOtp: wrap(authServiceSendEmailOtp),
		authServiceEmailLogin: wrap(authServiceEmailLogin),
		authServiceLogout: wrap(authServiceLogout),
		authServiceGetCurrentSession: wrap(authServiceGetCurrentSession),
		authServiceGetUserProfile: wrap(authServiceGetUserProfile),
		authServiceListSessions: wrap(authServiceListSessions),
		authServiceRevokeSession: wrap(authServiceRevokeSession),
		authServiceRevokeAllSessions: wrap(authServiceRevokeAllSessions),
		tenantServiceListTenants: wrap(tenantServiceListTenants),
		tenantServiceCreateTenant: wrap(tenantServiceCreateTenant),
		tenantServiceGetTenant: wrap(tenantServiceGetTenant),
		tenantServiceUpdateTenant: wrap(tenantServiceUpdateTenant),
		tenantServiceGetTenantConfig: wrap(tenantServiceGetTenantConfig),
		tenantServiceUpdateTenantConfig: wrap(tenantServiceUpdateTenantConfig),
		insurerServiceListInsurers: wrap(insurerServiceListInsurers),
		insurerServiceGetInsurer: wrap(insurerServiceGetInsurer),
		insurerServiceUpdateInsurer: wrap(insurerServiceUpdateInsurer),
		insurerServiceListInsurerProducts: wrap(insurerServiceListInsurerProducts),
		insurerServiceAddInsurerProduct: wrap(insurerServiceAddInsurerProduct),
		insurerServiceUpdateInsurerProduct: wrap(insurerServiceUpdateInsurerProduct),
		insurerServiceUpdateInsurerConfig: wrap(insurerServiceUpdateInsurerConfig),
		partnerServiceListPartners: wrap(partnerServiceListPartners),
		partnerServiceGetPartner: wrap(partnerServiceGetPartner),
		partnerServiceUpdatePartner: wrap(partnerServiceUpdatePartner),
		partnerServiceDeletePartner: wrap(partnerServiceDeletePartner),
		partnerServiceUpdatePartnerStatus: wrap(partnerServiceUpdatePartnerStatus),
		partnerServiceVerifyPartner: wrap(partnerServiceVerifyPartner),
		partnerServiceGetPartnerCommission: wrap(partnerServiceGetPartnerCommission),
		partnerServiceUpdateCommissionStructure: wrap(partnerServiceUpdateCommissionStructure),
		partnerServiceGetPartnerApiCredentials: wrap(partnerServiceGetPartnerApiCredentials),
		partnerServiceRotatePartnerApiKey: wrap(partnerServiceRotatePartnerApiKey),
		productServiceListProducts: wrap(productServiceListProducts),
		productServiceCreateProduct: wrap(productServiceCreateProduct),
		productServiceGetProduct: wrap(productServiceGetProduct),
		productServiceUpdateProduct: wrap(productServiceUpdateProduct),
		productServiceSearchProducts: wrap(productServiceSearchProducts),
		productServiceDeactivateProduct: wrap(productServiceDeactivateProduct),
		commissionServiceListCommissions: wrap(commissionServiceListCommissions),
		commissionServiceGetCommission: wrap(commissionServiceGetCommission),
		commissionServiceProcessPayout: wrap(commissionServiceProcessPayout),
		reportServiceListReportDefinitions: wrap(reportServiceListReportDefinitions),
		reportServiceListReportExecutions: wrap(reportServiceListReportExecutions),
		reportServiceListReportSchedules: wrap(reportServiceListReportSchedules),
		reportServiceExecuteReport: wrap(reportServiceExecuteReport),
		workflowServiceGetMyTasks: wrap(workflowServiceGetMyTasks),
		workflowServiceStartWorkflow: wrap(workflowServiceStartWorkflow),
		workflowServiceGetWorkflowInstance: wrap(workflowServiceGetWorkflowInstance),
		workflowServiceGetWorkflowDefinition: wrap(workflowServiceGetWorkflowDefinition),
		workflowServiceCompleteTask: wrap(workflowServiceCompleteTask),
		taskServiceListMyTasks: wrap(taskServiceListMyTasks),
		taskServiceGetTask: wrap(taskServiceGetTask),
		taskServiceCreateTask: wrap(taskServiceCreateTask),
		taskServiceUpdateTask: wrap(taskServiceUpdateTask),
		taskServiceAssignTask: wrap(taskServiceAssignTask),
		taskServiceCompleteTask: wrap(taskServiceCompleteTask),
		documentServiceListDocuments: wrap(documentServiceListDocuments),
		documentServiceGetDocument: wrap(documentServiceGetDocument),
		documentServiceDownloadDocument: wrap(documentServiceDownloadDocument),
		documentServiceGenerateDocument: wrap(documentServiceGenerateDocument),
		documentServiceListDocumentTemplates: wrap(documentServiceListDocumentTemplates),
		claimServiceGetClaim: wrap(claimServiceGetClaim),
		claimServiceApproveClaim: wrap(claimServiceApproveClaim),
		claimServiceRejectClaim: wrap(claimServiceRejectClaim),
		claimServiceSettleClaim: wrap(claimServiceSettleClaim),
		auditServiceGetAuditLogs: wrap(auditServiceGetAuditLogs),
		auditServiceGetAuditEvents: wrap(auditServiceGetAuditEvents),
		apiKeyServiceListApiKeys: wrap(apiKeyServiceListApiKeys),
		apiKeyServiceGenerateApiKey: wrap(apiKeyServiceGenerateApiKey),
		apiKeyServiceRevokeApiKey: wrap(apiKeyServiceRevokeApiKey),
		apiKeyServiceGetApiKey: wrap(apiKeyServiceGetApiKey),
		apiKeyServiceGetUsageStats: wrap(apiKeyServiceGetUsageStats),
		apiKeyServiceRotateApiKey: wrap(apiKeyServiceRotateApiKey),
		paymentServiceListPayments: wrap(paymentServiceListPayments),
		paymentServiceGetPayment: wrap(paymentServiceGetPayment),
		notificationServiceGetUserNotifications: wrap(notificationServiceGetUserNotifications),
		notificationServiceSendNotification: wrap(notificationServiceSendNotification),
		notificationServiceMarkAsRead: wrap(notificationServiceMarkAsRead),
		underwritingServiceListQuotes: wrap(underwritingServiceListQuotes),
		underwritingServiceGetQuote: wrap(underwritingServiceGetQuote),
		underwritingServiceRequestQuote: wrap(underwritingServiceRequestQuote),
		underwritingServiceApproveUnderwriting: wrap(underwritingServiceApproveUnderwriting),
		underwritingServiceRejectUnderwriting: wrap(underwritingServiceRejectUnderwriting),
		underwritingServiceConvertQuoteToPolicy: wrap(underwritingServiceConvertQuoteToPolicy),
		renewalServiceListUpcomingRenewals: wrap(renewalServiceListUpcomingRenewals),
		renewalServiceGetRenewalSchedule: wrap(renewalServiceGetRenewalSchedule),
		renewalServiceGetGracePeriod: wrap(renewalServiceGetGracePeriod),
		renewalServiceRenewPolicy: wrap(renewalServiceRenewPolicy),
		renewalServiceSendRenewalReminder: wrap(renewalServiceSendRenewalReminder),
		renewalServiceRevivePolicy: wrap(renewalServiceRevivePolicy),
		fraudServiceListFraudAlerts: wrap(fraudServiceListFraudAlerts),
		fraudServiceGetFraudAlert: wrap(fraudServiceGetFraudAlert),
		fraudServiceListFraudRules: wrap(fraudServiceListFraudRules),
		fraudServiceCreateFraudRule: wrap(fraudServiceCreateFraudRule),
		fraudServiceUpdateFraudRule: wrap(fraudServiceUpdateFraudRule),
		_directHttp: makeDirectHttp(event, sessionOverrides)
	};
}

export type SystemSdkClient = ReturnType<typeof makeSdkClient>;
export type DirectHttpClient = ReturnType<typeof makeDirectHttp>;
