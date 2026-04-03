using Grpc.Core;
using Grpc.Net.Client;
using Insuretech.Insurance.Services.V1;
using Insuretech.Quoting.Services.V1;
using Insuretech.Underwriting.Services.V1;
using Insuretech.Products.Services.V1;
using Insuretech.Policy.Services.V1;
using Insuretech.Claims.Services.V1;
using Insuretech.Beneficiary.Services.V1;
using Insuretech.Renewal.Services.V1;
using Insuretech.Endorsement.Services.V1;
using Insuretech.Fraud.Services.V1;
using Insuretech.Commission.Services.V1;
using Insuretech.Notification.Services.V1;
using Insuretech.Document.Services.V1;
using Insuretech.Refund.Services.V1;
using Insuretech.Payment.Services.V1;
using Microsoft.AspNetCore.Http;
using Microsoft.Extensions.Configuration;

namespace InsuranceEngine.Grpc.Clients;

/// <summary>
/// Client for calling the Go Insurance Service for database CRUD operations.
/// Includes X-* identity header forwarding for all gRPC calls (required by Go service auth).
/// </summary>
public sealed class InsuranceServiceClient : IDisposable
{
    private readonly GrpcChannel _channel;
    private readonly IHttpContextAccessor _httpContextAccessor;

    public InsuranceServiceClient(IConfiguration configuration, IHttpContextAccessor httpContextAccessor)
    {
        var address = configuration["InsuranceService:Url"] ?? "http://localhost:50115";
        _channel = GrpcChannel.ForAddress(address);
        _httpContextAccessor = httpContextAccessor;

        Insurance = new InsuranceService.InsuranceServiceClient(_channel);
        Underwriting = new UnderwritingService.UnderwritingServiceClient(_channel);
        Quoting = new QuotingService.QuotingServiceClient(_channel);
        Products = new ProductService.ProductServiceClient(_channel);
        Policies = new PolicyService.PolicyServiceClient(_channel);
        Claims = new ClaimService.ClaimServiceClient(_channel);
        Beneficiaries = new BeneficiaryService.BeneficiaryServiceClient(_channel);
        Renewals = new RenewalService.RenewalServiceClient(_channel);
        Endorsements = new EndorsementService.EndorsementServiceClient(_channel);
        Fraud = new FraudService.FraudServiceClient(_channel);
        Commissions = new CommissionService.CommissionServiceClient(_channel);
        Notifications = new NotificationService.NotificationServiceClient(_channel);
        Documents = new DocumentService.DocumentServiceClient(_channel);
        Refunds = new RefundService.RefundServiceClient(_channel);
        Payments = new PaymentService.PaymentServiceClient(_channel);
    }

    public InsuranceService.InsuranceServiceClient Insurance { get; }
    public UnderwritingService.UnderwritingServiceClient Underwriting { get; }
    public QuotingService.QuotingServiceClient Quoting { get; }
    public ProductService.ProductServiceClient Products { get; }
    public PolicyService.PolicyServiceClient Policies { get; }
    public ClaimService.ClaimServiceClient Claims { get; }
    public BeneficiaryService.BeneficiaryServiceClient Beneficiaries { get; }
    public RenewalService.RenewalServiceClient Renewals { get; }
    public EndorsementService.EndorsementServiceClient Endorsements { get; }
    public FraudService.FraudServiceClient Fraud { get; }
    public CommissionService.CommissionServiceClient Commissions { get; }
    public NotificationService.NotificationServiceClient Notifications { get; }
    public DocumentService.DocumentServiceClient Documents { get; }
    public RefundService.RefundServiceClient Refunds { get; }
    public PaymentService.PaymentServiceClient Payments { get; }

    /// <summary>
    /// Builds gRPC CallOptions with X-* identity headers forwarded from the active HTTP request.
    /// The Go insurance service requires x-user-id, x-tenant-id etc. in gRPC metadata.
    /// Call this on every gRPC invocation instead of passing cancellationToken directly.
    /// </summary>
    public CallOptions BuildCallOptions(CancellationToken ct = default)
    {
        var headers = new Metadata();
        var httpCtx = _httpContextAccessor.HttpContext;
        if (httpCtx != null)
        {
            foreach (var key in new[] {
                "x-user-id", "x-tenant-id", "x-partner-id", "x-token-id",
                "x-user-type", "x-portal", "x-roles", "x-request-id", "x-session-id" })
            {
                var val = httpCtx.Request.Headers[key].FirstOrDefault();
                if (!string.IsNullOrEmpty(val)) headers.Add(key, val);
            }
        }
        return new CallOptions(headers: headers, cancellationToken: ct);
    }

    public void Dispose()
    {
        _channel?.Dispose();
        GC.SuppressFinalize(this);
    }
}
