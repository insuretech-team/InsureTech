using FluentAssertions;
using Google.Protobuf.WellKnownTypes;
using Insuretech.Common.V1;
using Insuretech.Underwriting.Entity.V1;
using Insuretech.Underwriting.Services.V1;
using MediatR;
using Microsoft.Extensions.Logging.Abstractions;
using Moq;
using PoliSync.SharedKernel.Domain;
using PoliSync.SharedKernel.Messaging;
using PoliSync.Underwriting.Domain;
using PoliSync.Underwriting.GrpcServices;
using PoliSync.Underwriting.Infrastructure;
using HealthDeclarationEntity = Insuretech.Underwriting.Entity.V1.HealthDeclaration;
using ProtoMoney = Insuretech.Common.V1.Money;
using QuotationEntity = Insuretech.Policy.Entity.V1.Quotation;
using QuotationStatus = Insuretech.Policy.Entity.V1.QuotationStatus;
using QuoteEntity = Insuretech.Underwriting.Entity.V1.Quote;
using UnderwritingDecisionEntity = Insuretech.Underwriting.Entity.V1.UnderwritingDecision;
using Xunit;

namespace PoliSync.Underwriting.Tests;

public class UnderwritingGrpcServiceTests
{
    private sealed class CapturingEventBus : IEventBus
    {
        public List<(DomainEvent Event, string Topic)> Published { get; } = new();

        public Task PublishAsync<TEvent>(TEvent @event, string topic, CancellationToken cancellationToken = default) where TEvent : DomainEvent
        {
            Published.Add((@event, topic));
            return Task.CompletedTask;
        }

        public Task PublishAsync<TEvent>(TEvent @event, CancellationToken cancellationToken = default) where TEvent : DomainEvent
            => PublishAsync(@event, @event.EventType, cancellationToken);

        public Task PublishBatchAsync<TEvent>(IEnumerable<TEvent> events, string topic, CancellationToken cancellationToken = default) where TEvent : DomainEvent
        {
            foreach (var @event in events)
            {
                Published.Add((@event, topic));
            }

            return Task.CompletedTask;
        }
    }

    private sealed class FakeUnderwritingDataGateway : IUnderwritingDataGateway
    {
        private readonly Dictionary<string, QuoteEntity> _quotes = new();
        private readonly Dictionary<string, HealthDeclarationEntity> _declarationsByQuote = new();
        private readonly Dictionary<string, UnderwritingDecisionEntity> _decisionsByQuote = new();
        private readonly Dictionary<string, QuotationEntity> _quotations = new();

        public Task<QuoteEntity> CreateQuoteAsync(QuoteEntity quote, CancellationToken cancellationToken = default)
        {
            _quotes[quote.Id] = quote;
            return Task.FromResult(quote);
        }

        public Task<QuoteEntity?> GetQuoteAsync(string quoteId, CancellationToken cancellationToken = default)
            => Task.FromResult(_quotes.TryGetValue(quoteId, out var quote) ? quote : null);

        public Task<QuoteEntity> UpdateQuoteAsync(QuoteEntity quote, CancellationToken cancellationToken = default)
        {
            _quotes[quote.Id] = quote;
            return Task.FromResult(quote);
        }

        public Task<IReadOnlyList<QuoteEntity>> ListQuotesAsync(string beneficiaryId, int page, int pageSize, CancellationToken cancellationToken = default)
        {
            var normalizedPage = page <= 0 ? 1 : page;
            var normalizedPageSize = pageSize <= 0 ? 20 : pageSize;
            var items = _quotes.Values
                .Where(x => string.IsNullOrWhiteSpace(beneficiaryId) || x.BeneficiaryId == beneficiaryId)
                .Skip((normalizedPage - 1) * normalizedPageSize)
                .Take(normalizedPageSize)
                .ToList();
            return Task.FromResult<IReadOnlyList<QuoteEntity>>(items);
        }

        public Task<HealthDeclarationEntity?> GetHealthDeclarationByQuoteAsync(string quoteId, CancellationToken cancellationToken = default)
            => Task.FromResult(_declarationsByQuote.TryGetValue(quoteId, out var declaration) ? declaration : null);

        public Task<HealthDeclarationEntity> UpsertHealthDeclarationAsync(HealthDeclarationEntity declaration, CancellationToken cancellationToken = default)
        {
            if (string.IsNullOrWhiteSpace(declaration.Id))
            {
                declaration.Id = Guid.NewGuid().ToString("N");
            }

            _declarationsByQuote[declaration.QuoteId] = declaration;
            return Task.FromResult(declaration);
        }

        public Task<UnderwritingDecisionEntity?> GetLatestDecisionByQuoteAsync(string quoteId, CancellationToken cancellationToken = default)
            => Task.FromResult(_decisionsByQuote.TryGetValue(quoteId, out var decision) ? decision : null);

        public Task<UnderwritingDecisionEntity> UpsertUnderwritingDecisionAsync(UnderwritingDecisionEntity decision, CancellationToken cancellationToken = default)
        {
            if (string.IsNullOrWhiteSpace(decision.Id))
            {
                decision.Id = Guid.NewGuid().ToString("N");
            }

            _decisionsByQuote[decision.QuoteId] = decision;
            return Task.FromResult(decision);
        }

        public Task<QuotationEntity?> GetQuotationAsync(string quotationId, CancellationToken cancellationToken = default)
            => Task.FromResult(_quotations.TryGetValue(quotationId, out var quotation) ? quotation : null);

        public Task<QuotationEntity> UpdateQuotationAsync(QuotationEntity quotation, CancellationToken cancellationToken = default)
        {
            _quotations[quotation.QuotationId] = quotation;
            return Task.FromResult(quotation);
        }

        public void SeedQuotation(QuotationEntity quotation)
        {
            _quotations[quotation.QuotationId] = quotation;
        }
    }

    private static (UnderwritingGrpcService Service, FakeUnderwritingDataGateway Gateway, CapturingEventBus EventBus) CreateService()
    {
        var gateway = new FakeUnderwritingDataGateway();
        var eventBus = new CapturingEventBus();
        var service = new UnderwritingGrpcService(
            Mock.Of<IMediator>(),
            NullLogger<UnderwritingGrpcService>.Instance,
            gateway,
            eventBus,
            new UnderwritingRiskScorer());
        return (service, gateway, eventBus);
    }

    [Fact]
    public async Task RequestApproveConvertQuote_CompletesLifecycle()
    {
        var (service, gateway, eventBus) = CreateService();

        var requested = await service.RequestQuote(new RequestQuoteRequest
        {
            BeneficiaryId = $"ben-{Guid.NewGuid():N}",
            InsurerProductId = $"prod-{Guid.NewGuid():N}",
            SumAssured = new ProtoMoney { Amount = 800_000, Currency = "BDT" },
            TermYears = 2,
            PremiumPaymentMode = "monthly",
            ApplicantAge = 34,
            Smoker = false
        }, null!);

        gateway.SeedQuotation(new QuotationEntity
        {
            QuotationId = requested.QuoteId,
            EstimatedPremium = new ProtoMoney { Amount = requested.TotalPremium.Amount, Currency = requested.TotalPremium.Currency },
            QuotedAmount = new ProtoMoney { Amount = requested.TotalPremium.Amount, Currency = requested.TotalPremium.Currency },
            Status = QuotationStatus.Submitted
        });

        var declaration = await service.SubmitHealthDeclaration(new SubmitHealthDeclarationRequest
        {
            QuoteId = requested.QuoteId,
            HeightCm = 170,
            WeightKg = "70",
            HasPreExistingConditions = false,
            Smoker = false,
            AlcoholConsumer = false,
            OccupationRiskLevel = "low"
        }, null!);

        var approved = await service.ApproveUnderwriting(new ApproveUnderwritingRequest
        {
            QuoteId = requested.QuoteId,
            UnderwriterId = "uw-1",
            RiskLevel = "RISK_LEVEL_MEDIUM",
            PremiumAdjusted = true,
            AdjustedPremium = new ProtoMoney { Amount = 44_000, Currency = "BDT" },
            Conditions = new Struct(),
            Comments = "Approved with adjustment"
        }, null!);

        var converted = await service.ConvertQuoteToPolicy(new ConvertQuoteToPolicyRequest
        {
            QuoteId = requested.QuoteId,
            PaymentMethod = "bkash",
            PaymentReference = "tx-100"
        }, null!);

        var quote = await service.GetQuote(new GetQuoteRequest { QuoteId = requested.QuoteId }, null!);

        declaration.AutoApprovalPossible.Should().BeTrue();
        approved.DecisionId.Should().NotBeNullOrWhiteSpace();
        converted.PolicyId.Should().NotBeNullOrWhiteSpace();
        quote.Quote.Status.Should().Be(QuoteStatus.Converted);

        var linkedQuotation = await gateway.GetQuotationAsync(requested.QuoteId);
        linkedQuotation.Should().NotBeNull();
        linkedQuotation!.QuotedAmount.Amount.Should().Be(44_000);
        eventBus.Published.Should().ContainSingle(x => x.Topic == "insuretech.underwriting.decision_made.v1");
    }

    [Fact]
    public async Task RejectUnderwriting_SetsRejectedDecision()
    {
        var (service, gateway, eventBus) = CreateService();

        var requested = await service.RequestQuote(new RequestQuoteRequest
        {
            BeneficiaryId = $"ben-{Guid.NewGuid():N}",
            InsurerProductId = $"prod-{Guid.NewGuid():N}",
            SumAssured = new ProtoMoney { Amount = 900_000, Currency = "BDT" },
            TermYears = 1,
            PremiumPaymentMode = "yearly",
            ApplicantAge = 59,
            Smoker = true
        }, null!);

        gateway.SeedQuotation(new QuotationEntity
        {
            QuotationId = requested.QuoteId,
            EstimatedPremium = new ProtoMoney { Amount = requested.TotalPremium.Amount, Currency = requested.TotalPremium.Currency },
            QuotedAmount = new ProtoMoney { Amount = requested.TotalPremium.Amount, Currency = requested.TotalPremium.Currency },
            Status = QuotationStatus.Submitted
        });

        var rejected = await service.RejectUnderwriting(new RejectUnderwritingRequest
        {
            QuoteId = requested.QuoteId,
            UnderwriterId = "uw-2",
            Reason = "High risk profile",
            RiskLevel = "RISK_LEVEL_HIGH",
            Comments = "Declined"
        }, null!);

        var decision = await service.GetUnderwritingDecision(new GetUnderwritingDecisionRequest
        {
            QuoteId = requested.QuoteId
        }, null!);

        rejected.DecisionId.Should().NotBeNullOrWhiteSpace();
        decision.Decision.Decision.Should().Be(DecisionType.Rejected);
        decision.Decision.RiskLevel.Should().Be(RiskLevel.High);

        var linkedQuotation = await gateway.GetQuotationAsync(requested.QuoteId);
        linkedQuotation.Should().NotBeNull();
        linkedQuotation!.Status.Should().Be(QuotationStatus.Rejected);
        linkedQuotation.RejectionReason.Should().Be("High risk profile");
        eventBus.Published.Should().ContainSingle(x => x.Topic == "insuretech.underwriting.decision_made.v1");
    }
}
