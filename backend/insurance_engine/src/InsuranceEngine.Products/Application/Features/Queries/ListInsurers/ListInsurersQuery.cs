using System.Collections.Generic;
using MediatR;
using InsuranceEngine.Products.Domain;

namespace InsuranceEngine.Products.Application.Features.Queries.ListInsurers;

public record ListInsurersQuery : IRequest<List<Insurer>>;
