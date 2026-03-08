using System.Collections.Generic;
using MediatR;
using InsuranceEngine.Domain.Entities;

namespace InsuranceEngine.Application.Features.Insurers.Queries.ListInsurers;

public record ListInsurersQuery : IRequest<List<Insurer>>;
