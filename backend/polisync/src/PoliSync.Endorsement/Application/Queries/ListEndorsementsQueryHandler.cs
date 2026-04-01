using MediatR;
using Microsoft.Extensions.Logging;
using PoliSync.Endorsement.Infrastructure;
using PoliSync.SharedKernel.CQRS;

namespace PoliSync.Endorsement.Application.Queries;

public sealed class ListEndorsementsQueryHandler : IRequestHandler<ListEndorsementsQuery, Result<EndorsementListResult>>
{
    private readonly IEndorsementDataGateway _dataGateway;
    private readonly ILogger<ListEndorsementsQueryHandler> _logger;

    public ListEndorsementsQueryHandler(IEndorsementDataGateway dataGateway, ILogger<ListEndorsementsQueryHandler> logger)
    {
        _dataGateway = dataGateway;
        _logger = logger;
    }

    public async Task<Result<EndorsementListResult>> Handle(ListEndorsementsQuery request, CancellationToken cancellationToken)
    {
        try
        {
            var all = await _dataGateway.ListEndorsementsByPolicyAsync(request.PolicyId, cancellationToken);
            var page = request.PageNumber > 0 ? request.PageNumber : 1;
            var pageSize = request.PageSize > 0 ? request.PageSize : 20;
            var paged = all.Skip((page - 1) * pageSize).Take(pageSize).ToList();
            return Result.Ok(new EndorsementListResult(paged, all.Count));
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to list endorsements for policy {PolicyId}", request.PolicyId);
            return Result.Fail<EndorsementListResult>("LIST_ENDORSEMENTS_FAILED", ex.Message);
        }
    }
}
