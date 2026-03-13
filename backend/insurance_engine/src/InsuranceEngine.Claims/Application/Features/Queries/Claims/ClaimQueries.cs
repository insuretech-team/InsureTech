using MediatR;
using InsuranceEngine.Claims.Application.DTOs;
using InsuranceEngine.SharedKernel.CQRS;
using System;
using System.Collections.Generic;

namespace InsuranceEngine.Claims.Application.Features.Queries.Claims;

public record GetClaimByIdQuery(Guid Id) : IRequest<Result<ClaimResponseDto>>;

public record ListClaimsByCustomerQuery(Guid CustomerId, int Page = 1, int PageSize = 10) : IRequest<Result<List<ClaimResponseDto>>>;
