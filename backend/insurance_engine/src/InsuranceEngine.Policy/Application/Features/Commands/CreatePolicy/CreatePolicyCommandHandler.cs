using System;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using MediatR;
using Newtonsoft.Json;
using InsuranceEngine.Policy.Application.Interfaces;
using InsuranceEngine.Policy.Domain.Entities;
using InsuranceEngine.Policy.Domain.Enums;
using InsuranceEngine.Policy.Domain.Events;
using InsuranceEngine.Policy.Domain.Services;
using InsuranceEngine.Policy.Domain.ValueObjects;
using InsuranceEngine.SharedKernel.CQRS;
using InsuranceEngine.SharedKernel.Interfaces;

namespace InsuranceEngine.Policy.Application.Features.Commands.CreatePolicy;

public class CreatePolicyCommandHandler : IRequestHandler<CreatePolicyCommand, Result<CreatePolicyResponse>>
{
    private readonly IPolicyRepository _policyRepository;
    private readonly PolicyNumberGenerator _policyNumberGenerator;
    private readonly PolicyDuplicateDetector _duplicateDetector;
    private readonly IEventBus _eventBus;
    private readonly IEncryptionService _encryptionService;

    public CreatePolicyCommandHandler(
        IPolicyRepository policyRepository,
        PolicyNumberGenerator policyNumberGenerator,
        PolicyDuplicateDetector duplicateDetector,
        IEventBus eventBus,
        IEncryptionService encryptionService)
    {
        _policyRepository = policyRepository;
        _policyNumberGenerator = policyNumberGenerator;
        _duplicateDetector = duplicateDetector;
        _eventBus = eventBus;
        _encryptionService = encryptionService;
    }

    public async Task<Result<CreatePolicyResponse>> Handle(CreatePolicyCommand request, CancellationToken cancellationToken)
    {
        // FR-063 + FR-033: Duplicate detection & NID uniqueness
        var duplicateCheck = await _duplicateDetector.ValidateAsync(
            request.CustomerId, request.ProductId, request.Applicant.NidNumber);

        if (!duplicateCheck.IsSuccess)
            return Result<CreatePolicyResponse>.Fail(duplicateCheck.Error!);

        // Get product code for policy number generation
        var productCode = await _policyRepository.GetProductCodeAsync(request.ProductId);
        if (productCode == null)
            return Result<CreatePolicyResponse>.Fail(Error.NotFound("Product", request.ProductId.ToString()));

        // Generate policy number
        var seqNumber = await _policyRepository.GetNextSequenceNumberAsync();
        var policyNumber = _policyNumberGenerator.Generate(productCode, seqNumber);

        // Encrypt PII in applicant
        var applicant = new Applicant
        {
            FullName = request.Applicant.FullName,
            DateOfBirth = request.Applicant.DateOfBirth,
            NidNumber = !string.IsNullOrEmpty(request.Applicant.NidNumber)
                ? _encryptionService.Encrypt(request.Applicant.NidNumber) : null,
            Occupation = request.Applicant.Occupation,
            AnnualIncome = request.Applicant.AnnualIncome,
            Address = request.Applicant.Address,
            PhoneNumber = !string.IsNullOrEmpty(request.Applicant.PhoneNumber)
                ? _encryptionService.Encrypt(request.Applicant.PhoneNumber) : null,
            HealthDeclaration = request.Applicant.HealthDeclaration != null
                ? new Domain.ValueObjects.HealthDeclaration
                {
                    HasPreExistingConditions = request.Applicant.HealthDeclaration.HasPreExistingConditions,
                    Conditions = request.Applicant.HealthDeclaration.Conditions,
                    IsSmoker = request.Applicant.HealthDeclaration.IsSmoker,
                    BloodGroup = request.Applicant.HealthDeclaration.BloodGroup
                } : null
        };

        var endDate = request.StartDate.AddMonths(request.TenureMonths);

        var policy = new PolicyEntity
        {
            Id = Guid.NewGuid(),
            PolicyNumber = policyNumber,
            ProductId = request.ProductId,
            CustomerId = request.CustomerId,
            PartnerId = request.PartnerId,
            AgentId = request.AgentId,
            Status = PolicyStatus.PendingPayment,
            PremiumAmount = request.PremiumAmount,
            PremiumCurrency = "BDT",
            SumInsuredAmount = request.SumInsuredAmount,
            SumInsuredCurrency = "BDT",
            TenureMonths = request.TenureMonths,
            StartDate = request.StartDate,
            EndDate = endDate,
            ProposerDetailsJson = JsonConvert.SerializeObject(applicant),
            CreatedAt = DateTime.UtcNow,
            UpdatedAt = DateTime.UtcNow
        };

        if (request.Nominees != null)
        {
            foreach (var n in request.Nominees)
            {
                var addResult = policy.AddNominee(n.BeneficiaryId, n.FullName, n.Relationship, n.SharePercentage,
                    n.DateOfBirth, n.NidNumber, n.PhoneNumber, n.NomineeDobText);
                if (addResult.IsFailure)
                    return Result<CreatePolicyResponse>.Fail(addResult.Error!);
            }
        }

        // Add riders
        if (request.Riders != null)
        {
            foreach (var r in request.Riders)
            {
                policy.Riders.Add(new PolicyRider
                {
                    Id = Guid.NewGuid(),
                    PolicyId = policy.Id,
                    RiderName = r.RiderName,
                    PremiumAmount = r.PremiumAmount.Amount,
                    PremiumCurrency = r.PremiumAmount.CurrencyCode,
                    CoverageAmount = r.CoverageAmount.Amount,
                    CoverageCurrency = r.CoverageAmount.CurrencyCode,
                    CreatedAt = DateTime.UtcNow,
                    UpdatedAt = DateTime.UtcNow
                });
            }
        }

        var policyId = await _policyRepository.AddAsync(policy);

        await _eventBus.PublishAsync("insurance.policy.v1", new PolicyCreatedEvent(
            PolicyId: policyId,
            PolicyNumber: policyNumber,
            CustomerId: policy.CustomerId,
            ProductId: policy.ProductId,
            PremiumAmount: policy.PremiumAmount,
            StartDate: policy.StartDate,
            EndDate: policy.EndDate
        ));

        return Result<CreatePolicyResponse>.Ok(new CreatePolicyResponse(policyId, policyNumber));
    }
}
