using System;
using System.Threading;
using System.Threading.Tasks;
using InsuranceEngine.Partners.Domain.Interfaces;
using InsuranceEngine.SharedKernel.CQRS;
using MediatR;

namespace InsuranceEngine.Partners.Application.Features.Commands.Agents;

public record CreateAgentCommand(
    Guid PartnerId,
    string Name,
    string Code,
    string Email,
    string? Phone = null
) : IRequest<Result>;

public class CreateAgentCommandHandler : IRequestHandler<CreateAgentCommand, Result>
{
    private readonly IPartnerRepository _partnerRepository;

    public CreateAgentCommandHandler(IPartnerRepository partnerRepository)
    {
        _partnerRepository = partnerRepository;
    }

    public async Task<Result> Handle(CreateAgentCommand request, CancellationToken cancellationToken)
    {
        var partner = await _partnerRepository.GetByIdAsync(request.PartnerId, cancellationToken);
        if (partner == null)
            return Result.Fail(Error.NotFound("Partner", request.PartnerId.ToString()));

        var result = partner.AddAgent(
            request.Name,
            request.Code,
            request.Email,
            request.Phone);

        if (!result.IsSuccess) return result;

        await _partnerRepository.UpdateAsync(partner, cancellationToken);

        return Result.Success();
    }
}
