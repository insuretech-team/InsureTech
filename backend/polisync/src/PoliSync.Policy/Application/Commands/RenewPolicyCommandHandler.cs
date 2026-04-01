using Google.Protobuf.WellKnownTypes;
using MediatR;
using Microsoft.Extensions.Logging;
using PoliSync.SharedKernel.CQRS;
using PoliSync.SharedKernel.Domain;
using PoliSync.SharedKernel.Messaging;
using PoliSync.Policy.Infrastructure;

namespace PoliSync.Policy.Application.Commands;

public sealed class RenewPolicyCommandHandler : IRequestHandler<RenewPolicyCommand, Result<string>>
{
    private readonly IPolicyDataGateway _policyDataGateway;
    private readonly IEventBus _eventBus;
    private readonly ILogger<RenewPolicyCommandHandler> _logger;

    public RenewPolicyCommandHandler(
        IPolicyDataGateway policyDataGateway,
        IEventBus eventBus,
        ILogger<RenewPolicyCommandHandler> logger)
    {
        _policyDataGateway = policyDataGateway;
        _eventBus = eventBus;
        _logger = logger;
    }

    public async Task<Result<string>> Handle(RenewPolicyCommand request, CancellationToken cancellationToken)
    {
        try
        {
            var current = await _policyDataGateway.GetPolicyAsync(request.PolicyId, cancellationToken);
            if (current is null)
                return Result.Fail<string>("POLICY_NOT_FOUND", $"Policy {request.PolicyId} not found");

            if (current.Status != Insuretech.Policy.Entity.V1.PolicyStatus.Active &&
                current.Status != Insuretech.Policy.Entity.V1.PolicyStatus.Lapsed)
                return Result.Fail<string>("INVALID_STATUS", $"Policy in status {current.Status} cannot be renewed");

            var tenureMonths = request.NewTenureMonths > 0 ? request.NewTenureMonths : current.TenureMonths;
            var startDate = current.EndDate?.ToDateTime().Date ?? DateTime.UtcNow.Date;
            var now = DateTime.UtcNow;

            var renewed = new Insuretech.Policy.Entity.V1.Policy
            {
                PolicyId = Guid.NewGuid().ToString(),
                PolicyNumber = $"LP-{now:yyyy}-{Random.Shared.Next(100000, 999999)}",
                ProductId = current.ProductId,
                CustomerId = current.CustomerId,
                PartnerId = current.PartnerId,
                AgentId = current.AgentId,
                QuoteId = current.QuoteId,
                Status = Insuretech.Policy.Entity.V1.PolicyStatus.PendingPayment,
                PremiumAmount = new Insuretech.Common.V1.Money
                {
                    Amount = request.NewPremiumAmountPaisa > 0 ? request.NewPremiumAmountPaisa : (current.PremiumAmount?.Amount ?? 0),
                    Currency = current.PremiumAmount?.Currency ?? "BDT"
                },
                SumInsured = current.SumInsured,
                TenureMonths = tenureMonths,
                StartDate = Timestamp.FromDateTime(DateTime.SpecifyKind(startDate, DateTimeKind.Utc)),
                EndDate = Timestamp.FromDateTime(DateTime.SpecifyKind(startDate.AddMonths(tenureMonths), DateTimeKind.Utc)),
                CreatedAt = Timestamp.FromDateTime(now)
            };

            if (current.Nominees.Count > 0) renewed.Nominees.AddRange(current.Nominees);
            if (current.Riders.Count > 0)   renewed.Riders.AddRange(current.Riders);

            var created = await _policyDataGateway.CreatePolicyAsync(renewed, cancellationToken);

            await _eventBus.PublishAsync(new PolicyRenewedEvent(current.PolicyId, created.PolicyId), cancellationToken);

            _logger.LogInformation("Policy {OldPolicyId} renewed as {NewPolicyId}", request.PolicyId, created.PolicyId);
            return Result.Ok(created.PolicyId);
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to renew policy {PolicyId}", request.PolicyId);
            return Result.Fail<string>("RENEW_POLICY_FAILED", ex.Message);
        }
    }
}

public sealed record PolicyRenewedEvent(string OldPolicyId, string NewPolicyId) : DomainEvent;
