using MediatR;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Logging;
using Dapper;
using InsuranceEngine.SharedKernel.CQRS;

namespace InsuranceEngine.Commission.Application.Queries;

public sealed class ListCommissionsQueryHandler : IRequestHandler<ListCommissionsQuery, Result<ListCommissionsResult>>
{
    private readonly DbContext _dbContext;
    private readonly ILogger<ListCommissionsQueryHandler> _logger;

    public ListCommissionsQueryHandler(DbContext dbContext, ILogger<ListCommissionsQueryHandler> logger)
    {
        _dbContext = dbContext;
        _logger = logger;
    }

    public async Task<Result<ListCommissionsResult>> Handle(ListCommissionsQuery request, CancellationToken cancellationToken)
    {
        try
        {
            var sql = @"
                SELECT commission_id, policy_id, agent_id, premium_amount, commission_rate,
                       commission_amount, status, paid_at, created_at
                FROM insurance_schema.commissions
                WHERE (@AgentId IS NULL OR agent_id = @AgentId)
                  AND (@Status IS NULL OR status = @Status)
                  AND deleted_at IS NULL
                ORDER BY created_at DESC
                LIMIT @PageSize OFFSET @Offset";

            var offset = (request.Page - 1) * request.PageSize;

            using var connection = _dbContext.Database.GetDbConnection();
            await connection.OpenAsync(cancellationToken);

            var items = await connection.QueryAsync<CommissionDto>(sql, new
            {
                AgentId = request.AgentId,
                Status = request.Status,
                PageSize = request.PageSize,
                Offset = offset
            });

            var countSql = @"
                SELECT COUNT(*) FROM insurance_schema.commissions
                WHERE (@AgentId IS NULL OR agent_id = @AgentId)
                  AND (@Status IS NULL OR status = @Status)
                  AND deleted_at IS NULL";

            var totalCount = await connection.ExecuteScalarAsync<int>(countSql, new
            {
                AgentId = request.AgentId,
                Status = request.Status
            });

            return Result<ListCommissionsResult>.Ok(new ListCommissionsResult(items.ToList(), totalCount));
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to list commissions");
            return Result<ListCommissionsResult>.Fail("COMMISSION_LIST_FAILED", ex.Message);
        }
    }
}

public sealed class GetCommissionQueryHandler : IRequestHandler<GetCommissionQuery, Result<CommissionDto?>>
{
    private readonly DbContext _dbContext;
    private readonly ILogger<GetCommissionQueryHandler> _logger;

    public GetCommissionQueryHandler(DbContext dbContext, ILogger<GetCommissionQueryHandler> logger)
    {
        _dbContext = dbContext;
        _logger = logger;
    }

    public async Task<Result<CommissionDto?>> Handle(GetCommissionQuery request, CancellationToken cancellationToken)
    {
        try
        {
            var sql = @"
                SELECT commission_id, policy_id, agent_id, premium_amount, commission_rate,
                       commission_amount, status, paid_at, created_at
                FROM insurance_schema.commissions
                WHERE commission_id = @CommissionId AND deleted_at IS NULL";

            using var connection = _dbContext.Database.GetDbConnection();
            await connection.OpenAsync(cancellationToken);

            var commission = await connection.QueryFirstOrDefaultAsync<CommissionDto>(sql, new
            {
                CommissionId = request.CommissionId
            });

            if (commission == null)
                return Result<CommissionDto?>.NotFound("COMMISSION_NOT_FOUND", "Commission not found");

            return Result<CommissionDto?>.Ok(commission);
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to get commission {CommissionId}", request.CommissionId);
            return Result<CommissionDto?>.Fail("COMMISSION_GET_FAILED", ex.Message);
        }
    }
}
