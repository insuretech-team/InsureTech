using MediatR;
using Microsoft.Extensions.Logging;
using PoliSync.SharedKernel.CQRS;
using PoliSync.Policy.Infrastructure;
using Insuretech.Policy.Entity.V1;

namespace PoliSync.Policy.Application.Queries;

public sealed class ListUserPoliciesQueryHandler : IRequestHandler<ListUserPoliciesQuery, Result<PolicyListResult>>
{
    private readonly IPolicyDataGateway _policyDataGateway;
    private readonly ILogger<ListUserPoliciesQueryHandler> _logger;

    public ListUserPoliciesQueryHandler(
        IPolicyDataGateway policyDataGateway,
        ILogger<ListUserPoliciesQueryHandler> logger)
    {
        _policyDataGateway = policyDataGateway;
        _logger = logger;
    }

    public async Task<Result<PolicyListResult>> Handle(ListUserPoliciesQuery request, CancellationToken cancellationToken)
    {
        try
        {
            var page = request.PageNumber > 0 ? request.PageNumber : 1;
            var pageSize = request.PageSize > 0 ? request.PageSize : 20;

            var policies = await _policyDataGateway.ListPoliciesAsync(
                request.UserId, page, pageSize, cancellationToken);

            _logger.LogInformation("Listed {Count} policies for user {UserId}", policies.Count, request.UserId);

            return Result.Ok(new PolicyListResult([.. policies], policies.Count));
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to list policies for user {UserId}", request.UserId);
            return Result.Fail<PolicyListResult>("LIST_POLICIES_FAILED", ex.Message);
        }
    }
}
