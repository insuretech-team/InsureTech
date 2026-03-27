using MediatR;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Logging;
using Dapper;
using InsuranceEngine.SharedKernel.CQRS;

namespace InsuranceEngine.Commission.Application.Commands;

public sealed class CalculateCommissionCommandHandler : IRequestHandler<CalculateCommissionCommand, Result<string>>
{
    private readonly DbContext _dbContext;
    private readonly ILogger<CalculateCommissionCommandHandler> _logger;

    public CalculateCommissionCommandHandler(DbContext dbContext, ILogger<CalculateCommissionCommandHandler> logger)
    {
        _dbContext = dbContext;
        _logger = logger;
    }

    public async Task<Result<string>> Handle(CalculateCommissionCommand request, CancellationToken cancellationToken)
    {
        try
        {
            var commissionRate = 0.15m;
            var commissionAmount = request.PremiumAmount * commissionRate;
            var commissionId = Guid.NewGuid();
            var now = DateTime.UtcNow;

            var sql = @"
                INSERT INTO insurance_schema.commissions (
                    commission_id, policy_id, agent_id, premium_amount, commission_rate,
                    commission_amount, status, created_at
                ) VALUES (
                    @CommissionId, @PolicyId, @AgentId, @PremiumAmount, @CommissionRate,
                    @CommissionAmount, @Status, @CreatedAt
                )";

            using var connection = _dbContext.Database.GetDbConnection();
            await connection.OpenAsync(cancellationToken);

            await connection.ExecuteAsync(sql, new
            {
                CommissionId = commissionId,
                PolicyId = request.PolicyId,
                AgentId = request.AgentId,
                PremiumAmount = request.PremiumAmount,
                CommissionRate = commissionRate,
                CommissionAmount = commissionAmount,
                Status = "PENDING",
                CreatedAt = now
            });

            _logger.LogInformation("Commission calculated: {CommissionId} for Policy: {PolicyId}", commissionId, request.PolicyId);
            return Result<string>.Ok(commissionId.ToString());
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to calculate commission for policy {PolicyId}", request.PolicyId);
            return Result<string>.Fail("COMMISSION_CALCULATION_FAILED", ex.Message);
        }
    }
}

public sealed class ProcessPayoutCommandHandler : IRequestHandler<ProcessPayoutCommand, Result<bool>>
{
    private readonly DbContext _dbContext;
    private readonly ILogger<ProcessPayoutCommandHandler> _logger;

    public ProcessPayoutCommandHandler(DbContext dbContext, ILogger<ProcessPayoutCommandHandler> logger)
    {
        _dbContext = dbContext;
        _logger = logger;
    }

    public async Task<Result<bool>> Handle(ProcessPayoutCommand request, CancellationToken cancellationToken)
    {
        try
        {
            var sql = @"
                UPDATE insurance_schema.commissions
                SET status = 'PAID', paid_at = @PaidAt, updated_at = @UpdatedAt
                WHERE commission_id = @CommissionId::uuid AND status = 'PENDING' AND deleted_at IS NULL";

            using var connection = _dbContext.Database.GetDbConnection();
            await connection.OpenAsync(cancellationToken);

            var rows = await connection.ExecuteAsync(sql, new
            {
                CommissionId = request.CommissionId,
                PaidAt = DateTime.UtcNow,
                UpdatedAt = DateTime.UtcNow
            });

            if (rows == 0)
                return Result<bool>.Fail("COMMISSION_NOT_FOUND_OR_ALREADY_PAID", "Commission not found or already paid");

            _logger.LogInformation("Commission payout processed: {CommissionId}", request.CommissionId);
            return Result<bool>.Ok(true);
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to process payout for commission {CommissionId}", request.CommissionId);
            return Result<bool>.Fail("COMMISSION_PAYOUT_FAILED", ex.Message);
        }
    }
}
