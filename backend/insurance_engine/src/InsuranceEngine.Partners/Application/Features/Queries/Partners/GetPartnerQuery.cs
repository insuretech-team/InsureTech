using System;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using InsuranceEngine.Partners.Application.Features.Queries;
using InsuranceEngine.Partners.Domain.Interfaces;
using InsuranceEngine.SharedKernel.CQRS;
using MediatR;

namespace InsuranceEngine.Partners.Application.Features.Queries.Partners;

public record GetPartnerQuery(Guid PartnerId) : IRequest<Result<PartnerDto>>;

public class GetPartnerQueryHandler : IRequestHandler<GetPartnerQuery, Result<PartnerDto>>
{
    private readonly IPartnerRepository _partnerRepository;

    public GetPartnerQueryHandler(IPartnerRepository partnerRepository)
    {
        _partnerRepository = partnerRepository;
    }

    public async Task<Result<PartnerDto>> Handle(GetPartnerQuery request, CancellationToken cancellationToken)
    {
        var partner = await _partnerRepository.GetByIdAsync(request.PartnerId, cancellationToken);
        if (partner == null)
            return Result<PartnerDto>.Fail(Error.NotFound("Partner", request.PartnerId.ToString()));

        var dto = new PartnerDto(
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
        );

        return Result<PartnerDto>.Success(dto);
    }
}
