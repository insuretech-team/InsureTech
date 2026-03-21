using System;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using MediatR;
using Newtonsoft.Json;
using InsuranceEngine.Policy.Application.Interfaces;
using InsuranceEngine.Policy.Domain.Entities;
using InsuranceEngine.Policy.Domain.Enums;
using InsuranceEngine.Policy.Domain.ValueObjects;
using InsuranceEngine.Products.Application.Interfaces;
using InsuranceEngine.Products.Domain.Services;
using InsuranceEngine.SharedKernel.CQRS;

namespace InsuranceEngine.Policy.Application.Features.Commands.Endorsements;

public record ApproveEndorsementCommand(Guid EndorsementId, Guid ApprovedBy) : IRequest<Result>;

public class ApproveEndorsementCommandHandler : IRequestHandler<ApproveEndorsementCommand, Result>
{
    private readonly IEndorsementRepository _endorsementRepo;
    private readonly IPolicyRepository _policyRepo;
    private readonly IProductRepository _productRepo;
    private readonly PricingEngine _pricingEngine;

    public ApproveEndorsementCommandHandler(
        IEndorsementRepository endorsementRepo,
        IPolicyRepository policyRepo,
        IProductRepository productRepo,
        PricingEngine pricingEngine)
    {
        _endorsementRepo = endorsementRepo;
        _policyRepo = policyRepo;
        _productRepo = productRepo;
        _pricingEngine = pricingEngine;
    }

    public async Task<Result> Handle(ApproveEndorsementCommand request, CancellationToken cancellationToken)
    {
        var endorsement = await _endorsementRepo.GetByIdAsync(request.EndorsementId);
        if (endorsement == null) return Result.Failure("Endorsement not found");
        if (endorsement.Status != EndorsementStatus.Pending) return Result.Failure("Only pending endorsements can be approved");

        var policy = await _policyRepo.GetByIdWithNomineesAsync(endorsement.PolicyId);
        if (policy == null) return Result.Failure("Policy not found");

        // Apply changes based on type
        var result = endorsement.Type switch
        {
            EndorsementType.AddressChange or EndorsementType.ContactChange => ApplyContactChange(policy, endorsement),
            EndorsementType.NomineeChange => ApplyNomineeChange(policy, endorsement),
            EndorsementType.SumAssuredChange => await ApplySumAssuredChange(policy, endorsement),
            _ => Result.Failure($"Endorsement type '{endorsement.Type}' implementation pending")
        };

        if (!result.IsSuccess) return result;

        // Update endorsement
        endorsement.Status = EndorsementStatus.Applied;
        endorsement.ApprovedBy = request.ApprovedBy;
        endorsement.At = DateTime.UtcNow;
        endorsement.UpdatedAt = DateTime.UtcNow;

        await _policyRepo.UpdateAsync(policy);
        await _endorsementRepo.UpdateAsync(endorsement);

        return Result.Success();
    }

    private Result ApplyContactChange(PolicyEntity policy, Endorsement endorsement)
    {
        // For Address/Contact change, we just update the ProposerDetailsJson
        // The 'ChangesJson' should contain the partial/full Applicant object
        try
        {
            var newDetails = JsonConvert.DeserializeObject<Applicant>(endorsement.ChangesJson);
            if (newDetails == null) return Result.Failure("Invalid changes data");

            var currentDetails = !string.IsNullOrEmpty(policy.ProposerDetailsJson) 
                ? JsonConvert.DeserializeObject<Applicant>(policy.ProposerDetailsJson) 
                : new Applicant { FullName = "", DateOfBirth = DateTime.MinValue };

            if (currentDetails == null) currentDetails = new Applicant { FullName = "", DateOfBirth = DateTime.MinValue };

            // Merge logic using record 'with' expression
            currentDetails = currentDetails with
            {
                FullName = !string.IsNullOrEmpty(newDetails.FullName) ? newDetails.FullName : currentDetails.FullName,
                Address = !string.IsNullOrEmpty(newDetails.Address) ? newDetails.Address : currentDetails.Address,
                PhoneNumber = !string.IsNullOrEmpty(newDetails.PhoneNumber) ? newDetails.PhoneNumber : currentDetails.PhoneNumber
            };


            policy.ProposerDetailsJson = JsonConvert.SerializeObject(currentDetails);
            policy.UpdatedAt = DateTime.UtcNow;
            return Result.Success();
        }
        catch (Exception ex)
        {
            return Result.Failure($"Failed to apply contact changes: {ex.Message}");
        }
    }

    private Result ApplyNomineeChange(PolicyEntity policy, Endorsement endorsement)
    {
        try
        {
            // Expected format: { "Nominees": [ { "Action": "Add/Update/Remove", "Data": { ... } } ] }
            var changes = JsonConvert.DeserializeObject<NomineeChangeSet>(endorsement.ChangesJson);
            if (changes == null || changes.Items == null) return Result.Failure("Invalid nominee changes data");

            foreach (var item in changes.Items)
            {
                var r = item.Action.ToLower() switch
                {
                    "add" => policy.AddNominee(item.Data.BeneficiaryId, item.Data.FullName, item.Data.Relationship, item.Data.SharePercentage, item.Data.DateOfBirth, item.Data.NidNumber, item.Data.PhoneNumber, item.Data.NomineeDobText),
                    "update" => policy.UpdateNominee(item.Data.Id, item.Data.FullName, item.Data.Relationship, item.Data.SharePercentage, item.Data.DateOfBirth, item.Data.NidNumber, item.Data.PhoneNumber, item.Data.NomineeDobText),
                    "remove" => policy.RemoveNominee(item.Data.Id),
                    _ => Result.Failure($"Unknown nominee action: {item.Action}")
                };
                if (!r.IsSuccess) return r;
            }

            policy.UpdatedAt = DateTime.UtcNow;
            return Result.Success();
        }
        catch (Exception ex)
        {
            return Result.Failure($"Failed to apply nominee changes: {ex.Message}");
        }
    }

    private async Task<Result> ApplySumAssuredChange(PolicyEntity policy, Endorsement endorsement)
    {
        try
        {
            // Expected format: { "NewSumInsured": 100000 }
            var change = JsonConvert.DeserializeObject<SumAssuredChange>(endorsement.ChangesJson);
            if (change == null) return Result.Failure("Invalid sum assured changes data");

            var oldSumInsured = policy.SumInsuredAmount;
            policy.SumInsuredAmount = (long)change.NewSumInsured;

            // Recalculate premium
            var product = await _productRepo.GetByIdAsync(policy.ProductId);
            if (product == null) return Result.Failure("Product not found for premium recalculation");

            var applicantData = new Dictionary<string, string>();
            if (!string.IsNullOrEmpty(policy.ProposerDetailsJson))
            {
                var applicant = JsonConvert.DeserializeObject<Applicant>(policy.ProposerDetailsJson);
                if (applicant != null)
                {
                    applicantData["occupation"] = applicant.Occupation ?? "";
                    applicantData["income"] = applicant.AnnualIncome.ToString();
                    // Extract age if possible... (simplified for now)
                    var age = DateTime.UtcNow.Year - (applicant.DateOfBirth?.Year ?? DateTime.UtcNow.Year);
                    applicantData["age"] = age.ToString();

                }
            }

            var calcResult = _pricingEngine.Calculate(product, policy.SumInsuredAmount, policy.TenureMonths, new(), applicantData);
            
            endorsement.PremiumAdjustmentAmount = calcResult.TotalPremium.Amount - policy.TotalPayableAmount;
            endorsement.PremiumRefundRequired = endorsement.PremiumAdjustmentAmount < 0;

            policy.PremiumAmount = calcResult.BasePremium.Amount;
            policy.VatTaxAmount = calcResult.Vat.Amount;
            policy.ServiceFeeAmount = calcResult.ServiceFee.Amount;
            policy.TotalPayableAmount = calcResult.TotalPremium.Amount;
            
            policy.UpdatedAt = DateTime.UtcNow;
            return Result.Success();
        }
        catch (Exception ex)
        {
            return Result.Failure($"Failed to apply sum assured changes: {ex.Message}");
        }
    }

    private class NomineeChangeSet
    {
        public List<NomineeChangeItem> Items { get; set; } = new();
    }

    private class NomineeChangeItem
    {
        public string Action { get; set; } = string.Empty; // add, update, remove
        public NomineeChangeData Data { get; set; } = new();
    }

    private class NomineeChangeData
    {
        public Guid Id { get; set; }
        public Guid? BeneficiaryId { get; set; }
        public string FullName { get; set; } = string.Empty;
        public string Relationship { get; set; } = string.Empty;
        public double SharePercentage { get; set; }
        public DateTime? DateOfBirth { get; set; }
        public string? NidNumber { get; set; }
        public string? PhoneNumber { get; set; }
        public string? NomineeDobText { get; set; }
    }

    private class SumAssuredChange
    {
        public decimal NewSumInsured { get; set; }
    }
}
