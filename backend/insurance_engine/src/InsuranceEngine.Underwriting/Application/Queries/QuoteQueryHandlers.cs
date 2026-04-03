using MediatR;
using Microsoft.Extensions.Logging;
using Insuretech.Underwriting.Services.V1;
using Insuretech.Common.V1;
using InsuranceEngine.SharedKernel.Persistence;
using InsuranceEngine.SharedKernel.Persistence.Entities;
using InsuranceEngine.Products.Domain.Entities;
using Google.Protobuf.WellKnownTypes;
using System.Linq.Expressions;

namespace InsuranceEngine.Underwriting.Application.Queries;

public sealed class GetQuoteQueryHandler : IRequestHandler<GetQuoteQuery, GetQuoteResponse>
{
    private readonly IRepository<QuoteEntity> _quoteRepository;
    private readonly IRepository<HealthDeclarationEntity> _declarationRepository;
    private readonly IRepository<UnderwritingDecisionEntity> _decisionRepository;
    private readonly ILogger<GetQuoteQueryHandler> _logger;

    public GetQuoteQueryHandler(
        IRepository<QuoteEntity> quoteRepository,
        IRepository<HealthDeclarationEntity> declarationRepository,
        IRepository<UnderwritingDecisionEntity> decisionRepository,
        ILogger<GetQuoteQueryHandler> logger)
    {
        _quoteRepository = quoteRepository;
        _declarationRepository = declarationRepository;
        _decisionRepository = decisionRepository;
        _logger = logger;
    }

    public async Task<GetQuoteResponse> Handle(GetQuoteQuery request, CancellationToken cancellationToken)
    {
        try
        {
            var entity = await _quoteRepository.GetByIdAsync(Guid.Parse(request.QuoteId), cancellationToken);
            if (entity == null)
                return new GetQuoteResponse { Error = new Error { Code = "QUOTE_NOT_FOUND", Message = "Quote not found" } };

            var response = new GetQuoteResponse { Quote = MapQuoteToProto(entity) };

            var declarations = await _declarationRepository.FindAsync(d => d.QuoteId == entity.QuoteId, cancellationToken);
            if (declarations.Count > 0)
            {
                var d = declarations[0];
                response.HealthDeclaration = new Insuretech.Underwriting.Entity.V1.HealthDeclaration
                {
                    HeightCm = d.HeightCm,
                    WeightKg = d.WeightKg,
                    HasPreExistingConditions = d.HasPreExistingConditions,
                    PreExistingConditions = d.PreExistingConditions ?? "",
                    Smoker = d.Smoker,
                    AlcoholConsumer = d.AlcoholConsumer,
                    OccupationRiskLevel = d.OccupationRiskLevel ?? ""
                };
            }

            var decisionsList = await _decisionRepository.FindAsync(dd => dd.QuoteId == entity.QuoteId, cancellationToken);
            if (decisionsList.Count > 0)
            {
                var latestDecision = decisionsList.OrderByDescending(x => x.CreatedAt).First();
                var decisionProto = new Insuretech.Underwriting.Entity.V1.UnderwritingDecision
                {
                    UnderwriterId = latestDecision.UnderwriterId.ToString(),
                    UnderwriterComments = latestDecision.Comments ?? ""
                };
                if (System.Enum.TryParse<Insuretech.Underwriting.Entity.V1.RiskLevel>(latestDecision.RiskLevel ?? "", true, out var rl))
                    decisionProto.RiskLevel = rl;
                if (System.Enum.TryParse<Insuretech.Underwriting.Entity.V1.DecisionType>(latestDecision.Decision, true, out var dec))
                    decisionProto.Decision = dec;

                response.Decision = decisionProto;
            }

            return response;
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to get quote");
            throw;
        }
    }

    internal static Insuretech.Underwriting.Entity.V1.Quote MapQuoteToProto(QuoteEntity e)
    {
        var q = new Insuretech.Underwriting.Entity.V1.Quote
        {
            Id = e.QuoteId.ToString(),
            QuoteNumber = e.QuoteNumber,
            BeneficiaryId = e.BeneficiaryId.ToString(),
            InsurerProductId = e.InsurerProductId.ToString(),
            TermYears = e.TermYears,
            PremiumPaymentMode = e.PremiumPaymentMode,
            ApplicantAge = e.ApplicantAge,
            ApplicantOccupation = e.ApplicantOccupation ?? "",
            Smoker = e.Smoker,
            SelectedRiders = e.SelectedRiders ?? "",
            PremiumCalculation = e.PremiumCalculation ?? "",
            ConvertedPolicyId = e.ConvertedPolicyId?.ToString() ?? "",
            SumAssured = new Money { Amount = e.SumAssured, Currency = e.SumAssuredCurrency },
            BasePremium = new Money { Amount = e.BasePremium, Currency = e.BasePremiumCurrency },
            TotalPremium = new Money { Amount = e.TotalPremium, Currency = e.TotalPremiumCurrency }
        };

        if (System.Enum.TryParse<Insuretech.Underwriting.Entity.V1.QuoteStatus>(e.Status, true, out var s)) q.Status = s;
        q.ValidUntil = Timestamp.FromDateTime(DateTime.SpecifyKind(e.ValidUntil, DateTimeKind.Utc));
        if (e.RiderPremium.HasValue) q.RiderPremium = new Money { Amount = e.RiderPremium.Value, Currency = "BDT" };
        if (e.TaxAmount.HasValue) q.TaxAmount = new Money { Amount = e.TaxAmount.Value, Currency = "BDT" };
        if (e.ConvertedAt.HasValue) q.ConvertedAt = Timestamp.FromDateTime(DateTime.SpecifyKind(e.ConvertedAt.Value, DateTimeKind.Utc));

        return q;
    }
}

public sealed class ListQuotesQueryHandler : IRequestHandler<ListQuotesQuery, ListQuotesResponse>
{
    private readonly IRepository<QuoteEntity> _repository;
    private readonly ILogger<ListQuotesQueryHandler> _logger;

    public ListQuotesQueryHandler(IRepository<QuoteEntity> repository, ILogger<ListQuotesQueryHandler> logger)
    {
        _repository = repository;
        _logger = logger;
    }

    public async Task<ListQuotesResponse> Handle(ListQuotesQuery request, CancellationToken cancellationToken)
    {
        try
        {
            Expression<Func<QuoteEntity, bool>> predicate = q => q.DeletedAt == null;

            if (!string.IsNullOrEmpty(request.BeneficiaryId))
            {
                var beneficiaryId = Guid.Parse(request.BeneficiaryId);
                predicate = Combine(predicate, q => q.BeneficiaryId == beneficiaryId);
            }
            if (!string.IsNullOrEmpty(request.Status))
            {
                var status = request.Status;
                predicate = Combine(predicate, q => q.Status == status);
            }

            var (items, totalCount) = await _repository.GetPagedAsync(
                page: request.Page, pageSize: request.PageSize,
                predicate: predicate, orderBy: q => q.CreatedAt, descending: true,
                cancellationToken: cancellationToken);

            var response = new ListQuotesResponse { TotalCount = totalCount };
            foreach (var e in items)
                response.Quotes.Add(GetQuoteQueryHandler.MapQuoteToProto(e));
            return response;
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to list quotes");
            throw;
        }
    }

    private static Expression<Func<T, bool>> Combine<T>(Expression<Func<T, bool>> e1, Expression<Func<T, bool>> e2)
    {
        var p = Expression.Parameter(typeof(T));
        var v1 = new ReplaceVisitor(e1.Parameters[0], p);
        var v2 = new ReplaceVisitor(e2.Parameters[0], p);
        return Expression.Lambda<Func<T, bool>>(Expression.AndAlso(v1.Visit(e1.Body)!, v2.Visit(e2.Body)!), p);
    }
    private class ReplaceVisitor(Expression old, Expression @new) : ExpressionVisitor
    { public override Expression Visit(Expression? n) => n == old ? @new : base.Visit(n)!; }
}

public sealed class GetHealthDeclarationQueryHandler : IRequestHandler<GetHealthDeclarationQuery, GetHealthDeclarationResponse>
{
    private readonly IRepository<HealthDeclarationEntity> _repository;

    public GetHealthDeclarationQueryHandler(IRepository<HealthDeclarationEntity> repository)
    {
        _repository = repository;
    }

    public async Task<GetHealthDeclarationResponse> Handle(GetHealthDeclarationQuery request, CancellationToken cancellationToken)
    {
        var quoteId = Guid.Parse(request.QuoteId);
        var declarations = await _repository.FindAsync(d => d.QuoteId == quoteId, cancellationToken);
        if (declarations.Count == 0)
            return new GetHealthDeclarationResponse { Error = new Error { Code = "NOT_FOUND", Message = "No health declaration found" } };

        var d = declarations[0];
        return new GetHealthDeclarationResponse
        {
            HealthDeclaration = new Insuretech.Underwriting.Entity.V1.HealthDeclaration
            {
                HeightCm = d.HeightCm,
                WeightKg = d.WeightKg,
                HasPreExistingConditions = d.HasPreExistingConditions,
                PreExistingConditions = d.PreExistingConditions ?? "",
                Smoker = d.Smoker,
                AlcoholConsumer = d.AlcoholConsumer,
                OccupationRiskLevel = d.OccupationRiskLevel ?? ""
            }
        };
    }
}

public sealed class GetUnderwritingDecisionQueryHandler : IRequestHandler<GetUnderwritingDecisionQuery, GetUnderwritingDecisionResponse>
{
    private readonly IRepository<UnderwritingDecisionEntity> _repository;

    public GetUnderwritingDecisionQueryHandler(IRepository<UnderwritingDecisionEntity> repository)
    {
        _repository = repository;
    }

    public async Task<GetUnderwritingDecisionResponse> Handle(GetUnderwritingDecisionQuery request, CancellationToken cancellationToken)
    {
        var quoteId = Guid.Parse(request.QuoteId);
        var decisions = await _repository.FindAsync(d => d.QuoteId == quoteId, cancellationToken);
        if (decisions.Count == 0)
            return new GetUnderwritingDecisionResponse { Error = new Error { Code = "NOT_FOUND", Message = "No decision found" } };

        var dd = decisions.OrderByDescending(x => x.CreatedAt).First();
        var result = new Insuretech.Underwriting.Entity.V1.UnderwritingDecision
        {
            UnderwriterId = dd.UnderwriterId.ToString(),
            UnderwriterComments = dd.Comments ?? ""
        };
        if (System.Enum.TryParse<Insuretech.Underwriting.Entity.V1.RiskLevel>(dd.RiskLevel ?? "", true, out var rl))
            result.RiskLevel = rl;
        if (System.Enum.TryParse<Insuretech.Underwriting.Entity.V1.DecisionType>(dd.Decision, true, out var dec))
            result.Decision = dec;

        return new GetUnderwritingDecisionResponse { Decision = result };
    }
}
