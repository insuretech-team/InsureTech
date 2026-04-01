using MediatR;
using Microsoft.Extensions.Logging;
using PoliSync.Claims.Infrastructure;
using PoliSync.SharedKernel.CQRS;
using ClaimEntity = Insuretech.Claims.Entity.V1.Claim;

namespace PoliSync.Claims.Application.Queries;

public sealed class GetClaimQueryHandler : IRequestHandler<GetClaimQuery, Result<ClaimEntity>>
{
    private readonly IClaimDataGateway _dataGateway;
    private readonly ILogger<GetClaimQueryHandler> _logger;

    public GetClaimQueryHandler(IClaimDataGateway dataGateway, ILogger<GetClaimQueryHandler> logger)
    {
        _dataGateway = dataGateway;
        _logger = logger;
    }

    public async Task<Result<ClaimEntity>> Handle(GetClaimQuery request, CancellationToken cancellationToken)
    {
        try
        {
            var claim = await _dataGateway.GetClaimAsync(request.ClaimId, cancellationToken);
            if (claim is null)
                return Result.Fail<ClaimEntity>("CLAIM_NOT_FOUND", $"Claim {request.ClaimId} not found");

            return Result.Ok(claim);
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to get claim {ClaimId}", request.ClaimId);
            return Result.Fail<ClaimEntity>("GET_CLAIM_FAILED", ex.Message);
        }
    }
}
