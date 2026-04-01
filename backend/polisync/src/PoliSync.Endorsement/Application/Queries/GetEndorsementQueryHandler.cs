using MediatR;
using Microsoft.Extensions.Logging;
using PoliSync.Endorsement.Infrastructure;
using PoliSync.SharedKernel.CQRS;
using EndorsementEntity = Insuretech.Endorsement.Entity.V1.Endorsement;

namespace PoliSync.Endorsement.Application.Queries;

public sealed class GetEndorsementQueryHandler : IRequestHandler<GetEndorsementQuery, Result<EndorsementEntity>>
{
    private readonly IEndorsementDataGateway _dataGateway;
    private readonly ILogger<GetEndorsementQueryHandler> _logger;

    public GetEndorsementQueryHandler(IEndorsementDataGateway dataGateway, ILogger<GetEndorsementQueryHandler> logger)
    {
        _dataGateway = dataGateway;
        _logger = logger;
    }

    public async Task<Result<EndorsementEntity>> Handle(GetEndorsementQuery request, CancellationToken cancellationToken)
    {
        try
        {
            var endorsement = await _dataGateway.GetEndorsementAsync(request.EndorsementId, cancellationToken);
            if (endorsement is null)
                return Result.Fail<EndorsementEntity>("ENDORSEMENT_NOT_FOUND", $"Endorsement {request.EndorsementId} not found");
            return Result.Ok(endorsement);
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to get endorsement {EndorsementId}", request.EndorsementId);
            return Result.Fail<EndorsementEntity>("GET_ENDORSEMENT_FAILED", ex.Message);
        }
    }
}
