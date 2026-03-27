using System;
using System.Collections.Generic;
using System.Threading;
using System.Threading.Tasks;
using MediatR;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Logging;
using Dapper;
using Insuretech.Claims.Services.V1;
using Insuretech.Claims.Entity.V1;
using Insuretech.Common.V1;

namespace InsuranceEngine.Claims.Application.Queries;

public sealed class ListUserClaimsQueryHandler : IRequestHandler<ListUserClaimsQuery, ListUserClaimsResponse>
{
    private readonly DbContext _dbContext;
    private readonly ILogger<ListUserClaimsQueryHandler> _logger;

    public ListUserClaimsQueryHandler(DbContext dbContext, ILogger<ListUserClaimsQueryHandler> logger)
    {
        _dbContext = dbContext;
        _logger = logger;
    }

    public async Task<ListUserClaimsResponse> Handle(ListUserClaimsQuery request, CancellationToken cancellationToken)
    {
        try
        {
            var sql = @"
                SELECT claim_id, claim_number, policy_id, claim_type, claim_amount,
                       approved_amount, settlement_amount, description, status,
                       rejection_reason, settled_at, created_at
                FROM insurance_schema.claims
                WHERE (@PolicyId IS NULL OR policy_id = @PolicyId)
                  AND (@Status IS NULL OR status = @Status)
                  AND deleted_at IS NULL
                ORDER BY created_at DESC
                LIMIT @PageSize OFFSET @Offset";

            var offset = (request.Page - 1) * request.PageSize;

            using var connection = _dbContext.Database.GetDbConnection();
            if (connection.State != System.Data.ConnectionState.Open) await connection.OpenAsync(cancellationToken);

            var items = await connection.QueryAsync<dynamic>(sql, new
            {
                PolicyId = request.PolicyId,
                Status = request.Status,
                PageSize = request.PageSize,
                Offset = offset
            });

            var countSql = @"
                SELECT COUNT(*) FROM insurance_schema.claims
                WHERE (@PolicyId IS NULL OR policy_id = @PolicyId)
                  AND (@Status IS NULL OR status = @Status)
                  AND deleted_at IS NULL";

            var totalCount = await connection.ExecuteScalarAsync<int>(countSql, new
            {
                PolicyId = request.PolicyId,
                Status = request.Status
            });

            var response = new ListUserClaimsResponse
            {
                TotalCount = totalCount
            };

            foreach (var item in items)
            {
                response.Claims.Add(MapToProto(item));
            }

            return response;
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to list claims");
            throw;
        }
    }

    private static Insuretech.Claims.Entity.V1.Claim MapToProto(dynamic item)
    {
        var claim = new Insuretech.Claims.Entity.V1.Claim
        {
            ClaimId = item.claim_id?.ToString() ?? "",
            ClaimNumber = item.claim_number?.ToString() ?? "",
            PolicyId = item.policy_id?.ToString() ?? "",
            ClaimedAmount = new Money { Amount = (long)((decimal)(item.claim_amount ?? 0) * 100), Currency = "BDT" }
        };

        string statusStr = item.status?.ToString() ?? "";
        if (System.Enum.TryParse<ClaimStatus>(statusStr, true, out var stat)) claim.Status = stat;

        if (item.created_at != null)
        {
            claim.CreatedAt = Google.Protobuf.WellKnownTypes.Timestamp.FromDateTime(DateTime.SpecifyKind((DateTime)item.created_at, DateTimeKind.Utc));
        }

        return claim;
    }
}

public sealed class GetClaimQueryHandler : IRequestHandler<GetClaimQuery, GetClaimResponse>
{
    private readonly DbContext _dbContext;
    private readonly ILogger<GetClaimQueryHandler> _logger;

    public GetClaimQueryHandler(DbContext dbContext, ILogger<GetClaimQueryHandler> logger)
    {
        _dbContext = dbContext;
        _logger = logger;
    }

    public async Task<GetClaimResponse> Handle(GetClaimQuery request, CancellationToken cancellationToken)
    {
        try
        {
            var sql = @"
                SELECT claim_id, claim_number, policy_id, claim_type, claim_amount,
                       approved_amount, settlement_amount, description, status,
                       rejection_reason, settled_at, created_at
                FROM insurance_schema.claims
                WHERE claim_id = @ClaimId AND deleted_at IS NULL";

            using var connection = _dbContext.Database.GetDbConnection();
            if (connection.State != System.Data.ConnectionState.Open) await connection.OpenAsync(cancellationToken);

            var item = await connection.QueryFirstOrDefaultAsync<dynamic>(sql, new
            {
                ClaimId = request.ClaimId
            });

            if (item == null) throw new Exception("Claim not found");

            return new GetClaimResponse
            {
                Claim = MapToProto(item)
            };
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to get claim {ClaimId}", request.ClaimId);
            throw;
        }
    }

    private static Insuretech.Claims.Entity.V1.Claim MapToProto(dynamic item)
    {
        var claim = new Insuretech.Claims.Entity.V1.Claim
        {
            ClaimId = item.claim_id?.ToString() ?? "",
            ClaimNumber = item.claim_number?.ToString() ?? "",
            PolicyId = item.policy_id?.ToString() ?? "",
            ClaimedAmount = new Money { Amount = (long)((decimal)(item.claim_amount ?? 0) * 100), Currency = "BDT" }
        };

        string statusStr = item.status?.ToString() ?? "";
        if (System.Enum.TryParse<ClaimStatus>(statusStr, true, out var stat)) claim.Status = stat;

        if (item.created_at != null)
        {
            claim.CreatedAt = Google.Protobuf.WellKnownTypes.Timestamp.FromDateTime(DateTime.SpecifyKind((DateTime)item.created_at, DateTimeKind.Utc));
        }

        return claim;
    }
}
