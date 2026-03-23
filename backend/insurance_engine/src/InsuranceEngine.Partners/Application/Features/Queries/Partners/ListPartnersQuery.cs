using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using InsuranceEngine.Partners.Application.Features.Queries;
using InsuranceEngine.Partners.Domain.Interfaces;
using InsuranceEngine.SharedKernel.CQRS;
using MediatR;

namespace InsuranceEngine.Partners.Application.Features.Queries.Partners;

public record ListPartnersQuery() : IRequest<Result<List<PartnerDto>>>;

public class ListPartnersQueryHandler : IRequestHandler<ListPartnersQuery, Result<List<PartnerDto>>>
{
    private readonly IPartnerRepository _partnerRepository;

    public ListPartnersQueryHandler(IPartnerRepository partnerRepository)
    {
        _partnerRepository = partnerRepository;
    }

    public async Task<Result<List<PartnerDto>>> Handle(ListPartnersQuery request, CancellationToken cancellationToken)
    {
        var partners = await _partnerRepository.ListAsync(cancellationToken);
        
        var dtos = partners.Select(partner => new PartnerDto(
            partner.Id,
            partner.Name,
            partner.Code,
            partner.Email,
            partner.Phone,
            partner.Address,
            partner.Status.ToString(),
            partner.Agents.Select(a => new AgentDto(
                a.Id,
                a.Name,
                a.Code,
                a.Email,
                a.Phone,
                a.Status.ToString()
            )).ToList()
        )).ToList();

        return Result<List<PartnerDto>>.Success(dtos);
    }
}
