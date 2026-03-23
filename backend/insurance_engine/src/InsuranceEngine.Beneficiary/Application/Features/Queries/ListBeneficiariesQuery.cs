using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using InsuranceEngine.Beneficiary.Application.DTOs;
using InsuranceEngine.Beneficiary.Application.Interfaces;
using InsuranceEngine.SharedKernel.CQRS;
using MediatR;
using InsuranceEngine.Beneficiary.Application.Features;

namespace InsuranceEngine.Beneficiary.Application.Features.Queries;

public record ListBeneficiariesQuery(
    string? Type = null,
    string? Status = null,
    int PageSize = 10,
    int PageNumber = 1
) : IRequest<Result<PaginatedResponse<BeneficiaryDto>>>;

public class ListBeneficiariesQueryHandler : IRequestHandler<ListBeneficiariesQuery, Result<PaginatedResponse<BeneficiaryDto>>>
{
    private readonly IBeneficiaryRepository _repository;

    public ListBeneficiariesQueryHandler(IBeneficiaryRepository repository)
    {
        _repository = repository;
    }

    public async Task<Result<PaginatedResponse<BeneficiaryDto>>> Handle(ListBeneficiariesQuery request, CancellationToken cancellationToken)
    {
        var items = await _repository.ListAsync(
            request.Type, 
            request.Status, 
            request.PageNumber, 
            request.PageSize);

        var total = await _repository.GetTotalCountAsync(request.Type, request.Status);

        var dtos = items.Select(b => b.ToDto()).ToList();

        return Result<PaginatedResponse<BeneficiaryDto>>.Success(new PaginatedResponse<BeneficiaryDto>(
            Items: dtos,
            TotalCount: total,
            Page: request.PageNumber,
            PageSize: request.PageSize
        ));
    }
}
