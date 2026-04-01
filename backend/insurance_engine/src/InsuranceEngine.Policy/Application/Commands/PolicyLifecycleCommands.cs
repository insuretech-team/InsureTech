using InsuranceEngine.SharedKernel.CQRS;
using Insuretech.Policy.Services.V1;
using MediatR;

namespace InsuranceEngine.Policy.Application.Commands;

public sealed record IssuePolicyCommand(string PolicyId, string? QuoteId = null, string? PaymentId = null) : IRequest<IssuePolicyResponse>;



public sealed record GeneratePolicyDocumentCommand(string PolicyId) : IRequest<GeneratePolicyDocumentResponse>;
