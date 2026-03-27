using System;
using System.Collections.Generic;
using System.Threading;
using System.Threading.Tasks;
using MediatR;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Logging;
using Dapper;
using Insuretech.Policy.Services.V1;
using Insuretech.Policy.Entity.V1;
using Insuretech.Common.V1;

namespace InsuranceEngine.Policy.Application.Queries;

public sealed class ListUserPoliciesQueryHandler : IRequestHandler<ListUserPoliciesQuery, ListUserPoliciesResponse>
{
    private readonly DbContext _dbContext;
    private readonly ILogger<ListUserPoliciesQueryHandler> _logger;

    public ListUserPoliciesQueryHandler(DbContext dbContext, ILogger<ListUserPoliciesQueryHandler> logger)
    {
        _dbContext = dbContext;
        _logger = logger;
    }

    public async Task<ListUserPoliciesResponse> Handle(ListUserPoliciesQuery request, CancellationToken cancellationToken)
    {
        try
        {
            var sql = @"
                SELECT policy_id, policy_number, product_id, customer_id, partner_id, agent_id,
                       status, premium_amount, sum_insured, tenure_months, start_date, end_date,
                       issued_at, created_at
                FROM insurance_schema.policies
                WHERE (@CustomerId IS NULL OR customer_id = @CustomerId::uuid)
                  AND (@Status IS NULL OR status = @Status)
                  AND (@ProductId IS NULL OR product_id = @ProductId::uuid)
                  AND deleted_at IS NULL
                ORDER BY created_at DESC
                LIMIT @PageSize OFFSET @Offset";

            var offset = (request.Page - 1) * request.PageSize;

            using var connection = _dbContext.Database.GetDbConnection();
            if (connection.State != System.Data.ConnectionState.Open) await connection.OpenAsync(cancellationToken);

            var items = await connection.QueryAsync<dynamic>(sql, new
            {
                CustomerId = request.CustomerId,
                Status = request.Status,
                ProductId = request.ProductId,
                PageSize = request.PageSize,
                Offset = offset
            });

            var countSql = @"
                SELECT COUNT(*) FROM insurance_schema.policies
                WHERE (@CustomerId IS NULL OR customer_id = @CustomerId::uuid)
                  AND (@Status IS NULL OR status = @Status)
                  AND (@ProductId IS NULL OR product_id = @ProductId::uuid)
                  AND deleted_at IS NULL";

            var totalCount = await connection.ExecuteScalarAsync<int>(countSql, new
            {
                CustomerId = request.CustomerId,
                Status = request.Status,
                ProductId = request.ProductId
            });

            var response = new ListUserPoliciesResponse
            {
                TotalCount = totalCount
            };

            foreach (var item in items)
            {
                response.Policies.Add(MapToProto(item));
            }

            return response;
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to list policies");
            throw;
        }
    }

    private static Insuretech.Policy.Entity.V1.Policy MapToProto(dynamic item)
    {
        var policy = new Insuretech.Policy.Entity.V1.Policy
        {
            PolicyId = item.policy_id?.ToString() ?? "",
            PolicyNumber = item.policy_number?.ToString() ?? "",
            CustomerId = item.customer_id?.ToString() ?? "",
            ProductId = item.product_id?.ToString() ?? "",
            PartnerId = item.partner_id?.ToString() ?? "",
            AgentId = item.agent_id?.ToString() ?? "",
            PremiumAmount = new Money { Amount = (long)((decimal)(item.premium_amount ?? 0) * 100), Currency = "BDT" },
            SumInsured = new Money { Amount = (long)((decimal)(item.sum_insured ?? 0) * 100), Currency = "BDT" }
        };

        string statusStr = item.status?.ToString() ?? "";
        if (System.Enum.TryParse<Insuretech.Policy.Entity.V1.PolicyStatus>(statusStr, true, out var stat)) policy.Status = stat;

        if (item.created_at != null)
        {
            policy.CreatedAt = Google.Protobuf.WellKnownTypes.Timestamp.FromDateTime(DateTime.SpecifyKind((DateTime)item.created_at, DateTimeKind.Utc));
        }

        return policy;
    }
}

public sealed class GetPolicyQueryHandler : IRequestHandler<GetPolicyQuery, GetPolicyResponse>
{
    private readonly DbContext _dbContext;
    private readonly ILogger<GetPolicyQueryHandler> _logger;

    public GetPolicyQueryHandler(DbContext dbContext, ILogger<GetPolicyQueryHandler> logger)
    {
        _dbContext = dbContext;
        _logger = logger;
    }

    public async Task<GetPolicyResponse> Handle(GetPolicyQuery request, CancellationToken cancellationToken)
    {
        try
        {
            var sql = @"
                SELECT policy_id, policy_number, product_id, customer_id, partner_id, agent_id,
                       status, premium_amount, sum_insured, tenure_months, start_date, end_date,
                       issued_at, created_at
                FROM insurance_schema.policies
                WHERE policy_id = @PolicyId::uuid AND deleted_at IS NULL";

            using var connection = _dbContext.Database.GetDbConnection();
            if (connection.State != System.Data.ConnectionState.Open) await connection.OpenAsync(cancellationToken);

            var item = await connection.QueryFirstOrDefaultAsync<dynamic>(sql, new
            {
                PolicyId = request.PolicyId
            });

            if (item == null) throw new Exception("Policy not found");

            return new GetPolicyResponse
            {
                Policy = MapToProto(item)
            };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to get policy {PolicyId}", request.PolicyId);
            throw;
        }
    }

    private static Insuretech.Policy.Entity.V1.Policy MapToProto(dynamic item)
    {
        var policy = new Insuretech.Policy.Entity.V1.Policy
        {
            PolicyId = item.policy_id?.ToString() ?? "",
            PolicyNumber = item.policy_number?.ToString() ?? "",
            CustomerId = item.customer_id?.ToString() ?? "",
            ProductId = item.product_id?.ToString() ?? "",
            PartnerId = item.partner_id?.ToString() ?? "",
            AgentId = item.agent_id?.ToString() ?? "",
            PremiumAmount = new Money { Amount = (long)((decimal)(item.premium_amount ?? 0) * 100), Currency = "BDT" },
            SumInsured = new Money { Amount = (long)((decimal)(item.sum_insured ?? 0) * 100), Currency = "BDT" }
        };

        string statusStr = item.status?.ToString() ?? "";
        if (System.Enum.TryParse<Insuretech.Policy.Entity.V1.PolicyStatus>(statusStr, true, out var stat)) policy.Status = stat;

        if (item.created_at != null)
        {
            policy.CreatedAt = Google.Protobuf.WellKnownTypes.Timestamp.FromDateTime(DateTime.SpecifyKind((DateTime)item.created_at, DateTimeKind.Utc));
        }

        return policy;
    }
}
