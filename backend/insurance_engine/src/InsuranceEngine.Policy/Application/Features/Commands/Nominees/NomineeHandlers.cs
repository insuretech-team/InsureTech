using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using MediatR;
using Microsoft.EntityFrameworkCore;
using InsuranceEngine.Policy.Application.DTOs;
using InsuranceEngine.Policy.Application.Interfaces;
using InsuranceEngine.Policy.Domain.Entities;
using InsuranceEngine.Policy.Infrastructure.Persistence;
using InsuranceEngine.SharedKernel.CQRS;
using InsuranceEngine.SharedKernel.Interfaces;

namespace InsuranceEngine.Policy.Application.Features.Commands.Nominees;

public class AddNomineeCommandHandler : IRequestHandler<AddNomineeCommand, Result<Guid>>
{
    private readonly IPolicyRepository _repo;

    public AddNomineeCommandHandler(IPolicyRepository repo)
    {
        _repo = repo;
    }

    public async Task<Result<Guid>> Handle(AddNomineeCommand request, CancellationToken cancellationToken)
    {
        var policy = await _repo.GetByIdWithNomineesAsync(request.PolicyId);
        if (policy == null) return Result<Guid>.Fail(Error.NotFound("Policy", request.PolicyId.ToString()));

        var result = policy.AddNominee(request.BeneficiaryId, request.Relationship, request.SharePercentage);
        if (result.IsFailure) return Result<Guid>.Fail(result.Error!);

        await _repo.UpdateAsync(policy);
        
        // Return the id of the last added nominee
        var nominee = policy.Nominees.Last();
        return Result<Guid>.Ok(nominee.Id);
    }
}

public class UpdateNomineeCommandHandler : IRequestHandler<UpdateNomineeCommand, Result>
{
    private readonly IPolicyRepository _repo;

    public UpdateNomineeCommandHandler(IPolicyRepository repo)
    {
        _repo = repo;
    }

    public async Task<Result> Handle(UpdateNomineeCommand request, CancellationToken cancellationToken)
    {
        var policy = await _repo.GetByIdWithNomineesAsync(request.PolicyId);
        if (policy == null) return Result.Fail(Error.NotFound("Policy", request.PolicyId.ToString()));

        var result = policy.UpdateNominee(request.NomineeId, request.Relationship, request.SharePercentage);
        if (result.IsFailure) return result;

        await _repo.UpdateAsync(policy);
        return Result.Ok();
    }
}

public class DeleteNomineeCommandHandler : IRequestHandler<DeleteNomineeCommand, Result>
{
    private readonly IPolicyRepository _repo;

    public DeleteNomineeCommandHandler(IPolicyRepository repo) => _repo = repo;

    public async Task<Result> Handle(DeleteNomineeCommand request, CancellationToken cancellationToken)
    {
        var policy = await _repo.GetByIdWithNomineesAsync(request.PolicyId);
        if (policy == null) return Result.Fail(Error.NotFound("Policy", request.PolicyId.ToString()));

        var result = policy.RemoveNominee(request.NomineeId);
        if (result.IsFailure) return result;

        await _repo.UpdateAsync(policy);
        return Result.Ok();
    }
}

public class ListNomineesQueryHandler : IRequestHandler<ListNomineesQuery, List<NomineeDto>>
{
    private readonly PolicyDbContext _context;

    public ListNomineesQueryHandler(PolicyDbContext context)
    {
        _context = context;
    }

    public async Task<List<NomineeDto>> Handle(ListNomineesQuery request, CancellationToken cancellationToken)
    {
        var nominees = await _context.Nominees
            .Include(n => n.Beneficiary)
            .ThenInclude(b => b.IndividualDetails)
            .Where(n => n.PolicyId == request.PolicyId && !n.IsDeleted)
            .ToListAsync(cancellationToken);

        return nominees.Select(n => new NomineeDto(
            Id: n.Id,
            BeneficiaryId: n.BeneficiaryId,
            Relationship: n.Relationship,
            SharePercentage: n.SharePercentage
        )).ToList();
    }
}
