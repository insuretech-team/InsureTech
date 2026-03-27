using MediatR;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Logging;
using Dapper;
using Insuretech.Policy.Services.V1;

namespace InsuranceEngine.Policy.Application.Commands;

public sealed class IssuePolicyCommandHandler : IRequestHandler<IssuePolicyCommand, IssuePolicyResponse>
{
    private readonly DbContext _dbContext;
    private readonly ILogger<IssuePolicyCommandHandler> _logger;

    public IssuePolicyCommandHandler(DbContext dbContext, ILogger<IssuePolicyCommandHandler> logger)
    {
        _dbContext = dbContext;
        _logger = logger;
    }

    public async Task<IssuePolicyResponse> Handle(IssuePolicyCommand request, CancellationToken cancellationToken)
    {
        try
        {
            var sql = @"
                UPDATE insurance_schema.policies
                SET status = 'ISSUED', issued_at = @IssuedAt, updated_at = @UpdatedAt
                WHERE policy_id = @PolicyId::uuid AND deleted_at IS NULL";

            using var connection = _dbContext.Database.GetDbConnection();
            await connection.OpenAsync(cancellationToken);
            
            var rows = await connection.ExecuteAsync(sql, new
            {
                PolicyId = request.PolicyId,
                IssuedAt = DateTime.UtcNow,
                UpdatedAt = DateTime.UtcNow
            });

            if (rows == 0) throw new Exception("Policy not found or already issued");

            return new IssuePolicyResponse
            {
                Message = "Policy issued successfully"
            };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to issue policy {PolicyId}", request.PolicyId);
            throw;
        }
    }
}
