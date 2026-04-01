using MediatR;
using Microsoft.Extensions.Logging;
using PoliSync.Claims.Infrastructure;
using PoliSync.SharedKernel.CQRS;

namespace PoliSync.Claims.Application.Queries;

public sealed class ListClaimsQueryHandler : IRequestHandler<ListClaimsQuery, Result<ClaimListResult>>
{
    private readonly IClaimDataGateway _dataGateway;
    private readonly ILogger<ListClaimsQueryHandler> _logger;

    public ListClaimsQueryHandler(IClaimDataGateway dataGateway, ILogger<ListClaimsQueryHandler> logger)
    {
        _dataGateway = dataGateway;
        _logger = logger;
    }

    public async Task<Result<ClaimListResult>> Handle(ListClaimsQuery request, CancellationToken cancellationToken)
    {
        try
        {
            var page = request.PageNumber > 0 ? request.PageNumber : 1;
            var pageSize = request.PageSize > 0 ? request.PageSize : 20;

            var claims = await _dataGateway.ListClaimsAsync(
                request.CustomerId, request.PolicyId, page, pageSize, cancellationToken);

            _logger.LogInformation("Listed {Count} claims for customer {CustomerId}", claims.Count, request.CustomerId);
            return Result.Ok(new ClaimListResult(claims.ToList(), claims.Count));
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to list claims for customer {CustomerId}", request.CustomerId);
            return Result.Fail<ClaimListResult>("LIST_CLAIMS_FAILED", ex.Message);
        }
    }
}
