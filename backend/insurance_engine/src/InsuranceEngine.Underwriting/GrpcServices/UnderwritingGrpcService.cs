using System;
using System.Linq;
using System.Threading.Tasks;
using Google.Protobuf.WellKnownTypes;
using Grpc.Core;
using InsuranceEngine.Underwriting.Application.DTOs;
using InsuranceEngine.Underwriting.Application.Features.Commands.ApplyForQuote;
using InsuranceEngine.Underwriting.Application.Features.Commands.RecordUnderwritingDecision;
using InsuranceEngine.Underwriting.Application.Features.Queries.GetQuote;
using InsuranceEngine.Underwriting.Application.Features.Queries.ListQuotes;
using InsuranceEngine.Underwriting.Application.Features.Queries.GetUnderwritingHistory;
using Insuretech.Underwriting.Services.V1;
using Insuretech.Underwriting.Entity.V1;
using MediatR;
using Microsoft.Extensions.Logging;

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
        var command = new ApplyForQuoteCommand(
            Guid.Parse(request.BeneficiaryId),
            Guid.Parse(request.InsurerProductId),
            request.SumAssured?.Amount ?? 0,
            request.TermYears,
            request.PremiumPaymentMode,
            null, // Rider codes mapping would go here
            request.ApplicantAge,
            null,
            request.Smoker,
            new UnderwritingHealthDeclarationDto(0, 0, 0, false, null, false, false, null, false, false, null, false, false, null, null) // Default empty HD
        );

        var result = await _mediator.Send(command);

        if (result.IsSuccess)
        {
            return new RequestQuoteResponse
            {
                QuoteId = result.Value.Id.ToString(),
                QuoteNumber = result.Value.QuoteNumber,
                BasePremium = new Insuretech.Common.V1.Money { Amount = result.Value.BasePremium.Amount, Currency = result.Value.BasePremium.CurrencyCode },
                TotalPremium = new Insuretech.Common.V1.Money { Amount = result.Value.TotalPremium.Amount, Currency = result.Value.TotalPremium.CurrencyCode },
                ValidUntil = result.Value.ValidUntil.ToString("O"),
                Message = "Quote requested successfully"
            };
        }

        return new RequestQuoteResponse
        {
            Error = new Insuretech.Common.V1.Error
            {
                Code = "QUOTE_FAILED",
                Message = result.Error ?? "Unknown error"
            }
        };
    }

    public override async Task<GetQuoteResponse> GetQuote(GetQuoteRequest request, ServerCallContext context)
    {
        var result = await _mediator.Send(new GetQuoteQuery(Guid.Parse(request.QuoteId)));

        if (result.IsSuccess)
        {
            return new GetQuoteResponse
            {
                Quote = MapToProtoQuote(result.Value)
            };
        }

        return new GetQuoteResponse
        {
            Error = new Insuretech.Common.V1.Error
            {
                Code = "NOT_FOUND",
                Message = result.Error ?? "Quote not found"
            }
        };
    }

    private static Quote MapToProtoQuote(QuoteDto dto)
    {
        return new Quote
        {
            QuoteId = dto.Id.ToString(),
            QuoteNumber = dto.QuoteNumber,
            BeneficiaryId = dto.BeneficiaryId.ToString(),
            InsurerProductId = dto.InsurerProductId.ToString(),
            Status = dto.Status.ToString(),
            SumAssured = new Insuretech.Common.V1.Money { Amount = dto.SumAssured.Amount, Currency = dto.SumAssured.CurrencyCode },
            TermYears = dto.TermYears,
            PremiumPaymentMode = dto.PremiumPaymentMode,
            BasePremium = new Insuretech.Common.V1.Money { Amount = dto.BasePremium.Amount, Currency = dto.BasePremium.CurrencyCode },
            TotalPremium = new Insuretech.Common.V1.Money { Amount = dto.TotalPremium.Amount, Currency = dto.TotalPremium.CurrencyCode },
            CreatedAt = Timestamp.FromDateTime(DateTime.SpecifyKind(dto.CreatedAt, DateTimeKind.Utc))
        };
    }
}
