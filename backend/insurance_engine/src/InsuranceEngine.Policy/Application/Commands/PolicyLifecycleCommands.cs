using InsuranceEngine.SharedKernel.CQRS;
using Insuretech.Policy.Services.V1;
using MediatR;

namespace InsuranceEngine.Policy.Application.Commands;

public sealed record IssuePolicyCommand(string PolicyId, string? QuoteId = null, string? PaymentId = null) : IRequest<IssuePolicyResponse>;

public sealed record CancelPolicyCommand(string PolicyId, string Reason) : IRequest<CancelPolicyResponse>;

public sealed record RenewPolicyCommand(
    string PolicyId, 
    int TenureMonths, 
    bool UpdateNominees = false, 
    List<Insuretech.Policy.Entity.V1.Nominee>? Nominees = null) : IRequest<RenewPolicyResponse>;

public sealed record UpdatePolicyCommand(
    string PolicyId, 
    List<Insuretech.Policy.Entity.V1.Nominee>? Nominees = null,
    string? Address = null) : IRequest<UpdatePolicyResponse>;

public sealed record GeneratePolicyDocumentCommand(string PolicyId) : IRequest<GeneratePolicyDocumentResponse>;
