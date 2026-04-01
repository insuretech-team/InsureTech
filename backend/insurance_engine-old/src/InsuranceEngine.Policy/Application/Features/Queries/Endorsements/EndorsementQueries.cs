using System;
using System.Collections.Generic;
using System.Threading;
using System.Threading.Tasks;
using MediatR;
using InsuranceEngine.Policy.Application.Interfaces;
using InsuranceEngine.Policy.Domain.Entities;
using InsuranceEngine.Policy.Domain.Enums;
using InsuranceEngine.SharedKernel.CQRS;

namespace InsuranceEngine.Policy.Application.Features.Queries.Endorsements;

public record GetEndorsementQuery(Guid Id) : IRequest<Result<Endorsement>>;

public record ListEndorsementsQuery(
    Guid? PolicyId = null,
    EndorsementStatus? Status = null,
    int Page = 1,
    int PageSize = 10
) : IRequest<Result<EndorsementListDto>>;

public class EndorsementListDto
{
    public List<Endorsement> Items { get; set; } = new();
    public int TotalCount { get; set; }
}

public class EndorsementQueryHandlers : 
    IRequestHandler<GetEndorsementQuery, Result<Endorsement>>,
    IRequestHandler<ListEndorsementsQuery, Result<EndorsementListDto>>
{
    private readonly IEndorsementRepository _repo;

    public EndorsementQueryHandlers(IEndorsementRepository repo)
    {
        _repo = repo;
    }

    public async Task<Result<Endorsement>> Handle(GetEndorsementQuery request, CancellationToken cancellationToken)
    {
        var item = await _repo.GetByIdAsync(request.Id);
        return item != null ? Result<Endorsement>.Success(item) : Result<Endorsement>.Failure("Endorsement not found");
    }

    public async Task<Result<EndorsementListDto>> Handle(ListEndorsementsQuery request, CancellationToken cancellationToken)
    {
        var (items, total) = await _repo.ListAsync(request.PolicyId, request.Status, request.Page, request.PageSize);
        return Result<EndorsementListDto>.Success(new EndorsementListDto { Items = items, TotalCount = total });
    }
}
