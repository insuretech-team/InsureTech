using MediatR;
using PoliSync.SharedKernel.CQRS;
using Insuretech.Policy.Entity.V1;

namespace PoliSync.Policy.Application.Queries;

public sealed record GetPolicyQuery(string PolicyId) : IQuery<Insuretech.Policy.Entity.V1.Policy>;
