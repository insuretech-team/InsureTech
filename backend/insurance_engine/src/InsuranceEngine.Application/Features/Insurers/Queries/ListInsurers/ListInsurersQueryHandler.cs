using System.Collections.Generic;
using System.Threading;
using System.Threading.Tasks;
using MediatR;
using InsuranceEngine.Application.Interfaces;
using InsuranceEngine.Domain.Entities;

namespace InsuranceEngine.Application.Features.Insurers.Queries.ListInsurers;

public class ListInsurersQueryHandler : IRequestHandler<ListInsurersQuery, List<Insurer>>
{
    private readonly IInsurerRepository _insurerRepository;

    public ListInsurersQueryHandler(IInsurerRepository insurerRepository)
    {
        _insurerRepository = insurerRepository;
    }

    public async Task<List<Insurer>> Handle(ListInsurersQuery request, CancellationToken cancellationToken)
    {
        return await _insurerRepository.ListAsync();
    }
}
