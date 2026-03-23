using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using InsuranceEngine.Commission.Application.Features.Queries;
using InsuranceEngine.Commission.Domain.Interfaces;
using InsuranceEngine.SharedKernel.CQRS;
using MediatR;

namespace InsuranceEngine.Commission.Application.Features.Queries.Commissions;

public record GetRecipientCommissionsQuery(Guid RecipientId) : IRequest<Result<List<CommissionDto>>>;

public class GetRecipientCommissionsQueryHandler : IRequestHandler<GetRecipientCommissionsQuery, Result<List<CommissionDto>>>
{
    private readonly ICommissionRepository _commissionRepository;

    public GetRecipientCommissionsQueryHandler(ICommissionRepository commissionRepository)
    {
        _commissionRepository = commissionRepository;
    }

    public async Task<Result<List<CommissionDto>>> Handle(GetRecipientCommissionsQuery request, CancellationToken cancellationToken)
    {
        var commissions = await _commissionRepository.ListByRecipientAsync(request.RecipientId, cancellationToken);
        
        var dtos = commissions.Select(c => new CommissionDto(
            c.Id,
            c.PolicyId,
            c.PartnerId,
            c.AgentId,
            c.Type.ToString(),
            c.Amount,
            c.Currency,
            c.Status.ToString(),
            c.CreatedAt
        )).ToList();

        return Result<List<CommissionDto>>.Success(dtos);
    }
}
