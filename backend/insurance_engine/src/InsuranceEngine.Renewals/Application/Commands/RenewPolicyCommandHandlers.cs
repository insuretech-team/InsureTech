using MediatR;
using Microsoft.Extensions.Logging;
using Insuretech.Policy.Services.V1;
using Insuretech.Common.V1;
using InsuranceEngine.SharedKernel.Persistence;
using InsuranceEngine.SharedKernel.Persistence.Entities;
using InsuranceEngine.SharedKernel.Infrastructure;
using InsuranceEngine.SharedKernel.Domain.Events;
using Microsoft.EntityFrameworkCore;
using System;
using System.Threading;
using System.Threading.Tasks;

namespace InsuranceEngine.Renewals.Application.Commands;

public sealed class RenewPolicyCommandHandler : IRequestHandler<RenewPolicyCommand, RenewPolicyTenureResponse>
{
    private readonly IRepository<PolicyEntity> _repository;
    private readonly IRepository<PolicyNomineeEntity> _nomineeRepository;
    private readonly InsuranceDbContext _dbContext;
    private readonly ILogger<RenewPolicyCommandHandler> _logger;
    private readonly IKafkaPublisher _kafkaPublisher;

    public RenewPolicyCommandHandler(
        IRepository<PolicyEntity> repository,
        IRepository<PolicyNomineeEntity> nomineeRepository,
        InsuranceDbContext dbContext,
        ILogger<RenewPolicyCommandHandler> logger,
        IKafkaPublisher kafkaPublisher)
    {
        _repository = repository;
        _nomineeRepository = nomineeRepository;
        _dbContext = dbContext;
        _logger = logger;
        _kafkaPublisher = kafkaPublisher;
    }

    public async Task<RenewPolicyTenureResponse> Handle(RenewPolicyCommand request, CancellationToken cancellationToken)
    {
        try
        {
            var existingPolicy = await _repository.GetByIdAsync(Guid.Parse(request.PolicyId), cancellationToken);
            if (existingPolicy == null)
            {
                return new RenewPolicyTenureResponse
                {
                    Error = new Error { Code = "POLICY_NOT_FOUND", Message = "Policy not found" }
                };
            }

            if (existingPolicy.Status != "ACTIVE" && existingPolicy.Status != "EXPIRED")
            {
                return new RenewPolicyTenureResponse
                {
                    Error = new Error { Code = "INVALID_STATUS", Message = $"Policy cannot be renewed from status '{existingPolicy.Status}'" }
                };
            }

            // FR-047/048: 30-day grace period
            var expiryDate = existingPolicy.EndDate;
            var now = DateTime.UtcNow;
            var daysFromExpiry = (now - expiryDate).TotalDays;

            if (daysFromExpiry > 30)
            {
                // FR-069: Reinstatement requires medical underwriting and special approval
                return new RenewPolicyTenureResponse
                {
                    Error = new Error { Code = "GRACE_PERIOD_EXPIRED", Message = "Policy is past the 30-day grace period. Reinstatement required (Underwriting module)." }
                };
            }

            if (daysFromExpiry > 0 && daysFromExpiry <= 30)
            {
                 _logger.LogInformation("Policy {PolicyId} is in Grace Period ({Days} days past expiry).", request.PolicyId, (int)daysFromExpiry);
            }

            // Get new sequence number
            var connection = _dbContext.Database.GetDbConnection();
            if (connection.State != System.Data.ConnectionState.Open)
                await connection.OpenAsync(cancellationToken);

            using var cmd = connection.CreateCommand();
            cmd.CommandText = "SELECT nextval('insurance_schema.policy_number_seq')";
            var seqResult = await cmd.ExecuteScalarAsync(cancellationToken);
            var sequenceNumber = Convert.ToInt64(seqResult);

            // Derive product code for new policy number
            var productCode = existingPolicy.PolicyNumber.Split('-').Length >= 3 
                ? existingPolicy.PolicyNumber.Split('-')[2] 
                : "0001";

            var year = DateTime.UtcNow.Year;
            var seq = sequenceNumber.ToString().PadLeft(6, '0');
            var policyNumber = $"LBT-{year}-{productCode}-{seq}";

            // Create renewed policy
            var newStartDate = existingPolicy.EndDate > DateTime.UtcNow ? existingPolicy.EndDate : DateTime.UtcNow;
            var renewedPolicy = new PolicyEntity
            {
                PolicyId = Guid.NewGuid(),
                PolicyNumber = policyNumber,
                ProductId = existingPolicy.ProductId,
                CustomerId = existingPolicy.CustomerId,
                PartnerId = existingPolicy.PartnerId,
                AgentId = existingPolicy.AgentId,
                Status = "PENDING_PAYMENT",
                PremiumAmount = existingPolicy.PremiumAmount,
                PremiumCurrency = existingPolicy.PremiumCurrency,
                SumInsured = existingPolicy.SumInsured,
                SumInsuredCurrency = existingPolicy.SumInsuredCurrency,
                TenureMonths = request.TenureMonths,
                StartDate = newStartDate,
                EndDate = newStartDate.AddMonths(request.TenureMonths),
                CreatedAt = DateTime.UtcNow,
                UpdatedAt = DateTime.UtcNow
            };

            await _repository.AddAsync(renewedPolicy, cancellationToken);

            // Copy nominees from existing policy or use new ones
            if (request.UpdateNominees && request.Nominees != null)
            {
                foreach (var nominee in request.Nominees)
                {
                    await _nomineeRepository.AddAsync(new PolicyNomineeEntity
                    {
                        NomineeId = Guid.NewGuid(),
                        PolicyId = renewedPolicy.PolicyId,
                        FullName = nominee.FullName,
                        Relationship = nominee.Relationship,
                        SharePercentage = nominee.SharePercentage,
                        DateOfBirth = nominee.DateOfBirth?.ToDateTime() ?? DateTime.UtcNow,
                        NidNumber = nominee.NidNumber,
                        PhoneNumber = nominee.PhoneNumber,
                        CreatedAt = DateTime.UtcNow,
                        UpdatedAt = DateTime.UtcNow
                    }, cancellationToken);
                }
            }
            else
            {
                // Copy existing nominees
                var existingNominees = await _nomineeRepository.FindAsync(n => n.PolicyId == existingPolicy.PolicyId, cancellationToken);
                foreach (var nominee in existingNominees)
                {
                    await _nomineeRepository.AddAsync(new PolicyNomineeEntity
                    {
                        NomineeId = Guid.NewGuid(),
                        PolicyId = renewedPolicy.PolicyId,
                        FullName = nominee.FullName,
                        Relationship = nominee.Relationship,
                        SharePercentage = nominee.SharePercentage,
                        DateOfBirth = nominee.DateOfBirth,
                        NidNumber = nominee.NidNumber,
                        PhoneNumber = nominee.PhoneNumber,
                        CreatedAt = DateTime.UtcNow,
                        UpdatedAt = DateTime.UtcNow
                    }, cancellationToken);
                }
            }

            _logger.LogInformation("Policy renewed: {OldPolicyNumber} → {NewPolicyNumber}", existingPolicy.PolicyNumber, renewedPolicy.PolicyNumber);

            // FR-047: Kafka Event for Renewal
            var renewalEvent = new PolicyRenewedEvent(
                renewedPolicy.PolicyId,
                renewedPolicy.PolicyNumber,
                renewedPolicy.CustomerId,
                renewedPolicy.PremiumAmount,
                renewedPolicy.EndDate.ToString("yyyy-MM-dd"),
                renewedPolicy.PartnerId,
                renewedPolicy.AgentId
            );
            await _kafkaPublisher.PublishAsync("insurance.policy.renewed", renewalEvent);

            return new RenewPolicyTenureResponse
            {
                NewPolicyId = renewedPolicy.PolicyId.ToString(),
                NewPolicyNumber = renewedPolicy.PolicyNumber,
                PremiumAmount = new Money { Amount = renewedPolicy.PremiumAmount, Currency = "BDT" },
                Message = "Policy renewed successfully"
            };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to renew policy {PolicyId}", request.PolicyId);
            return new RenewPolicyTenureResponse
            {
                Error = new Error { Code = "RENEW_FAILED", Message = ex.Message }
            };
        }
    }
}
