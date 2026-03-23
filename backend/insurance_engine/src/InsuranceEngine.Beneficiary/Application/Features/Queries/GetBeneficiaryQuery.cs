using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using InsuranceEngine.Beneficiary.Application.DTOs;
using InsuranceEngine.Beneficiary.Application.Interfaces;
using InsuranceEngine.SharedKernel.CQRS;
using MediatR;

namespace InsuranceEngine.Beneficiary.Application.Features.Queries;

public record GetBeneficiaryQuery(Guid Id) : IRequest<Result<BeneficiaryDto>>;

public class GetBeneficiaryQueryHandler : IRequestHandler<GetBeneficiaryQuery, Result<BeneficiaryDto>>
{
    private readonly IBeneficiaryRepository _repository;

    public GetBeneficiaryQueryHandler(IBeneficiaryRepository repository)
    {
        _repository = repository;
    }

    public async Task<Result<BeneficiaryDto>> Handle(GetBeneficiaryQuery request, CancellationToken cancellationToken)
    {
        var b = await _repository.GetByIdAsync(request.Id);
        if (b == null) return Result<BeneficiaryDto>.Failure("Beneficiary not found");

        return Result.Ok(b.ToDto());
    }
}
