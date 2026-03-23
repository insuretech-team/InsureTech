using System;
using System.Threading;
using System.Threading.Tasks;
using InsuranceEngine.Partners.Domain.Entities;
using InsuranceEngine.Partners.Domain.Interfaces;
using InsuranceEngine.SharedKernel.CQRS;
using MediatR;

namespace InsuranceEngine.Partners.Application.Features.Commands.Partners;

public record CreatePartnerCommand(
    string OrganizationName,
    string Code,
    string Email,
    string? Phone = null
) : IRequest<Result<Guid>>;

public class CreatePartnerCommandHandler : IRequestHandler<CreatePartnerCommand, Result<Guid>>
{
    private readonly IPartnerRepository _partnerRepository;

    public CreatePartnerCommandHandler(IPartnerRepository partnerRepository)
    {
        _partnerRepository = partnerRepository;
    }

    public async Task<Result<Guid>> Handle(CreatePartnerCommand request, CancellationToken cancellationToken)
    {
        var existing = await _partnerRepository.GetByCodeAsync(request.Code, cancellationToken);
        if (existing != null)
            return Result<Guid>.Fail(Error.Validation($"Partner with code {request.Code} already exists"));

        var partner = Partner.Create(
            request.OrganizationName,
            request.Code,
            request.Email,
            request.Phone);

        await _partnerRepository.CreateAsync(partner, cancellationToken);

        return Result<Guid>.Success(partner.Id);
    }
}
