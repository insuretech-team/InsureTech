using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using InsuranceEngine.Policy.Application.DTOs;
using InsuranceEngine.Policy.Application.Interfaces;
using InsuranceEngine.SharedKernel.CQRS;
using MediatR;

namespace InsuranceEngine.Policy.Application.Features.Queries.Beneficiaries;

public record ListBeneficiariesQuery(
    string? Type = null,
    string? Status = null,
    int PageSize = 10,
    int Page = 1
) : IRequest<Result<PagedList<BeneficiaryDto>>>;

public record PagedList<T>(IEnumerable<T> Items, int TotalCount, int Page, int PageSize);

public class ListBeneficiariesQueryHandler : IRequestHandler<ListBeneficiariesQuery, Result<PagedList<BeneficiaryDto>>>
{
    private readonly IBeneficiaryRepository _repository;

    public ListBeneficiariesQueryHandler(IBeneficiaryRepository repository)
    {
        _repository = repository;
    }

    public async Task<Result<PagedList<BeneficiaryDto>>> Handle(ListBeneficiariesQuery request, CancellationToken cancellationToken)
    {
        var beneficiaries = await _repository.ListAsync(request.Type, request.Status, request.Page, request.PageSize);
        var totalCount = await _repository.GetTotalCountAsync(request.Type, request.Status);

        var dtos = beneficiaries.Select(b => new BeneficiaryDto(
            b.Id,
            b.UserId,
            b.Type.ToString(),
            b.Code,
            b.Status.ToString(),
            b.KycStatus.ToString(),
            b.KycCompletedAt,
            b.RiskScore,
            b.ReferralCode
        ));

        return Result.Ok(new PagedList<BeneficiaryDto>(dtos, totalCount, request.Page, request.PageSize));
    }
}
