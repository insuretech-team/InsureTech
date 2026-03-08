using Google.Protobuf.WellKnownTypes;
using Grpc.Core;
using Insuretech.Insurance.Services.V1;
using Insuretech.Policy.Entity.V1;
using MediatR;
using Microsoft.Extensions.Logging;
using PoliSync.Quotes.Infrastructure;

namespace PoliSync.Quotes.GrpcServices;

public sealed class QuotesGrpcService : InsuranceService.InsuranceServiceBase
{
    private readonly IMediator _mediator;
    private readonly ILogger<QuotesGrpcService> _logger;
    private readonly IQuotationDataGateway _quotationDataGateway;

    public QuotesGrpcService(
        IMediator mediator,
        ILogger<QuotesGrpcService> logger,
        IQuotationDataGateway quotationDataGateway)
    {
        _mediator = mediator;
        _logger = logger;
        _quotationDataGateway = quotationDataGateway;
    }

    public override async Task<CreateQuotationResponse> CreateQuotation(CreateQuotationRequest request, ServerCallContext context)
    {
        if (request.Quotation is null)
        {
            throw new RpcException(new Status(StatusCode.InvalidArgument, "Quotation is required"));
        }

        var now = DateTime.UtcNow;
        var quotation = request.Quotation.Clone();
        quotation.QuotationId = string.IsNullOrWhiteSpace(quotation.QuotationId)
            ? Guid.NewGuid().ToString("N")
            : quotation.QuotationId;
        quotation.QuotationNumber = string.IsNullOrWhiteSpace(quotation.QuotationNumber)
            ? BuildQuotationNumber()
            : quotation.QuotationNumber;
        quotation.Status = quotation.Status == QuotationStatus.Unspecified
            ? QuotationStatus.Draft
            : quotation.Status;
        quotation.ValidUntil = HasValue(quotation.ValidUntil)
            ? quotation.ValidUntil
            : Timestamp.FromDateTime(now.AddDays(30));
        quotation.CreatedAt = Timestamp.FromDateTime(now);
        quotation.UpdatedAt = Timestamp.FromDateTime(now);

        if (quotation.Status != QuotationStatus.Draft && !HasValue(quotation.SubmissionDate))
        {
            quotation.SubmissionDate = Timestamp.FromDateTime(now);
        }

        var created = await _quotationDataGateway.CreateQuotationAsync(quotation, GetCancellationToken(context));
        _logger.LogInformation("Quotation created: {QuotationId}", created.QuotationId);

        return new CreateQuotationResponse { Quotation = created };
    }

    public override async Task<GetQuotationResponse> GetQuotation(GetQuotationRequest request, ServerCallContext context)
    {
        var quotation = await _quotationDataGateway.GetQuotationAsync(request.QuotationId, GetCancellationToken(context));
        if (quotation is null)
        {
            throw new RpcException(new Status(StatusCode.NotFound, "Quotation not found"));
        }

        return new GetQuotationResponse { Quotation = quotation };
    }

    public override async Task<UpdateQuotationResponse> UpdateQuotation(UpdateQuotationRequest request, ServerCallContext context)
    {
        if (request.Quotation is null || string.IsNullOrWhiteSpace(request.Quotation.QuotationId))
        {
            throw new RpcException(new Status(StatusCode.InvalidArgument, "Quotation with QuotationId is required"));
        }

        var existing = await _quotationDataGateway.GetQuotationAsync(request.Quotation.QuotationId, GetCancellationToken(context));
        if (existing is null)
        {
            throw new RpcException(new Status(StatusCode.NotFound, "Quotation not found"));
        }

        var updated = request.Quotation.Clone();
        ApplyLifecycleValidation(existing, updated);

        updated.CreatedAt = existing.CreatedAt;
        updated.QuotationNumber = string.IsNullOrWhiteSpace(updated.QuotationNumber)
            ? existing.QuotationNumber
            : updated.QuotationNumber;
        updated.ValidUntil = HasValue(updated.ValidUntil) ? updated.ValidUntil : existing.ValidUntil;
        updated.SubmissionDate = ResolveSubmissionDate(existing, updated);
        updated.UpdatedAt = Timestamp.FromDateTime(DateTime.UtcNow);

        var persisted = await _quotationDataGateway.UpdateQuotationAsync(updated, GetCancellationToken(context));
        return new UpdateQuotationResponse { Quotation = persisted };
    }

    public override async Task<Empty> DeleteQuotation(DeleteQuotationRequest request, ServerCallContext context)
    {
        var existing = await _quotationDataGateway.GetQuotationAsync(request.QuotationId, GetCancellationToken(context));
        if (existing is null)
        {
            throw new RpcException(new Status(StatusCode.NotFound, "Quotation not found"));
        }

        if (existing.Status == QuotationStatus.Approved)
        {
            throw new RpcException(new Status(StatusCode.FailedPrecondition, "Approved quotation cannot be deleted"));
        }

        await _quotationDataGateway.DeleteQuotationAsync(request.QuotationId, GetCancellationToken(context));
        return new Empty();
    }

    public override async Task<ListQuotationsResponse> ListQuotations(ListQuotationsRequest request, ServerCallContext context)
    {
        var page = request.Page <= 0 ? 1 : request.Page;
        var pageSize = request.PageSize <= 0 ? 20 : request.PageSize;
        var quotations = await _quotationDataGateway.ListQuotationsAsync(
            request.BusinessId,
            page,
            pageSize,
            GetCancellationToken(context));

        var ordered = quotations.OrderByDescending(q => q.UpdatedAt?.Seconds ?? 0).ToList();
        var response = new ListQuotationsResponse { Total = ordered.Count };
        response.Quotations.AddRange(ordered);
        return response;
    }

    private static void ApplyLifecycleValidation(Quotation existing, Quotation updated)
    {
        var nextStatus = updated.Status == QuotationStatus.Unspecified ? existing.Status : updated.Status;
        updated.Status = nextStatus;

        if (existing.Status is QuotationStatus.Approved or QuotationStatus.Rejected && nextStatus != existing.Status)
        {
            throw new RpcException(new Status(StatusCode.FailedPrecondition, "Terminal quotation cannot change status"));
        }

        if (nextStatus == existing.Status)
        {
            return;
        }

        var allowed = existing.Status switch
        {
            QuotationStatus.Draft => nextStatus == QuotationStatus.Submitted,
            QuotationStatus.Submitted => nextStatus is QuotationStatus.Received or QuotationStatus.Approved or QuotationStatus.Rejected,
            QuotationStatus.Received => nextStatus is QuotationStatus.Approved or QuotationStatus.Rejected,
            _ => false
        };

        if (!allowed)
        {
            throw new RpcException(new Status(StatusCode.FailedPrecondition, $"Invalid quotation transition from {existing.Status} to {nextStatus}"));
        }
    }

    private static Timestamp ResolveSubmissionDate(Quotation existing, Quotation updated)
    {
        if (HasValue(updated.SubmissionDate))
        {
            return updated.SubmissionDate;
        }

        if (HasValue(existing.SubmissionDate))
        {
            return existing.SubmissionDate;
        }

        return updated.Status == QuotationStatus.Draft
            ? new Timestamp()
            : Timestamp.FromDateTime(DateTime.UtcNow);
    }

    private static bool HasValue(Timestamp? timestamp)
        => timestamp is not null && timestamp.Seconds > 0;

    private static string BuildQuotationNumber()
        => $"QUO-{DateTime.UtcNow:yyyy}-{Random.Shared.Next(100000, 999999)}";

    private static CancellationToken GetCancellationToken(ServerCallContext? context)
        => context?.CancellationToken ?? CancellationToken.None;
}
