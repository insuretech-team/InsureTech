using Grpc.Core;
using MediatR;
using Microsoft.Extensions.Logging;
using Insuretech.Underwriting.Services.V1;
using InsuranceEngine.Underwriting.Application.Commands;
using InsuranceEngine.Underwriting.Application.Queries;

namespace InsuranceEngine.Underwriting.GrpcServices;

public sealed class UnderwritingGrpcService : UnderwritingService.UnderwritingServiceBase
{
    private readonly IMediator _mediator;
    private readonly ILogger<UnderwritingGrpcService> _logger;

    public UnderwritingGrpcService(IMediator mediator, ILogger<UnderwritingGrpcService> logger)
    {
        _mediator = mediator;
        _logger = logger;
    }

    public override async Task<RequestQuoteResponse> RequestQuote(RequestQuoteRequest request, ServerCallContext context)
    {
        if (string.IsNullOrEmpty(request.BeneficiaryId) || string.IsNullOrEmpty(request.InsurerProductId))
            throw new RpcException(new Status(StatusCode.InvalidArgument, "Beneficiary ID and Product ID are required"));

        return await _mediator.Send(new RequestQuoteCommand(
            request.BeneficiaryId, request.InsurerProductId, request.SumAssured?.Amount ?? 0,
            request.TermYears, request.PremiumPaymentMode,
            request.RiderCodes.Count > 0 ? request.RiderCodes.ToList() : null,
            request.ApplicantAge, request.Smoker), context.CancellationToken);
    }

    public override async Task<GetQuoteResponse> GetQuote(GetQuoteRequest request, ServerCallContext context)
    {
        if (string.IsNullOrEmpty(request.QuoteId))
            throw new RpcException(new Status(StatusCode.InvalidArgument, "Quote ID is required"));
        return await _mediator.Send(new GetQuoteQuery(request.QuoteId), context.CancellationToken);
    }

    public override async Task<ListQuotesResponse> ListQuotes(ListQuotesRequest request, ServerCallContext context)
    {
        if (string.IsNullOrEmpty(request.BeneficiaryId))
            throw new RpcException(new Status(StatusCode.InvalidArgument, "Beneficiary ID is required"));
        return await _mediator.Send(new ListQuotesQuery(request.BeneficiaryId, request.Status,
            request.Page <= 0 ? 1 : request.Page, request.PageSize <= 0 ? 10 : request.PageSize), context.CancellationToken);
    }

    public override async Task<SubmitHealthDeclarationResponse> SubmitHealthDeclaration(SubmitHealthDeclarationRequest request, ServerCallContext context)
    {
        if (string.IsNullOrEmpty(request.QuoteId))
            throw new RpcException(new Status(StatusCode.InvalidArgument, "Quote ID is required"));
        return await _mediator.Send(new SubmitHealthDeclarationCommand(
            request.QuoteId, request.HeightCm, request.WeightKg,
            request.HasPreExistingConditions, request.PreExistingConditions,
            request.Smoker, request.AlcoholConsumer, request.OccupationRiskLevel), context.CancellationToken);
    }

    public override async Task<GetHealthDeclarationResponse> GetHealthDeclaration(GetHealthDeclarationRequest request, ServerCallContext context)
    {
        if (string.IsNullOrEmpty(request.QuoteId))
            throw new RpcException(new Status(StatusCode.InvalidArgument, "Quote ID is required"));
        return await _mediator.Send(new GetHealthDeclarationQuery(request.QuoteId), context.CancellationToken);
    }

    public override async Task<ApproveUnderwritingResponse> ApproveUnderwriting(ApproveUnderwritingRequest request, ServerCallContext context)
    {
        if (string.IsNullOrEmpty(request.QuoteId) || string.IsNullOrEmpty(request.UnderwriterId))
            throw new RpcException(new Status(StatusCode.InvalidArgument, "Quote ID and Underwriter ID are required"));
        return await _mediator.Send(new ApproveUnderwritingCommand(
            request.QuoteId, request.UnderwriterId, request.RiskLevel,
            request.PremiumAdjusted, request.AdjustedPremium?.Amount, request.Comments), context.CancellationToken);
    }

    public override async Task<RejectUnderwritingResponse> RejectUnderwriting(RejectUnderwritingRequest request, ServerCallContext context)
    {
        if (string.IsNullOrEmpty(request.QuoteId) || string.IsNullOrEmpty(request.UnderwriterId))
            throw new RpcException(new Status(StatusCode.InvalidArgument, "Quote ID and Underwriter ID are required"));
        return await _mediator.Send(new RejectUnderwritingCommand(
            request.QuoteId, request.UnderwriterId, request.Reason, request.RiskLevel, request.Comments), context.CancellationToken);
    }

    public override async Task<ConvertQuoteToPolicyResponse> ConvertQuoteToPolicy(ConvertQuoteToPolicyRequest request, ServerCallContext context)
    {
        if (string.IsNullOrEmpty(request.QuoteId))
            throw new RpcException(new Status(StatusCode.InvalidArgument, "Quote ID is required"));
        return await _mediator.Send(new ConvertQuoteToPolicyCommand(
            request.QuoteId, request.PaymentMethod, request.PaymentReference), context.CancellationToken);
    }

    public override async Task<GetUnderwritingDecisionResponse> GetUnderwritingDecision(GetUnderwritingDecisionRequest request, ServerCallContext context)
    {
        if (string.IsNullOrEmpty(request.QuoteId))
            throw new RpcException(new Status(StatusCode.InvalidArgument, "Quote ID is required"));
        return await _mediator.Send(new GetUnderwritingDecisionQuery(request.QuoteId), context.CancellationToken);
    }
}
