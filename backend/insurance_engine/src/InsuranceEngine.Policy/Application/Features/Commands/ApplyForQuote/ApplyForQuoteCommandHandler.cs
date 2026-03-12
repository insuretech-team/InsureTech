using System;
using System.Collections.Generic;
using System.Linq;
using System.Text.Json;
using System.Threading;
using System.Threading.Tasks;
using InsuranceEngine.Policy.Application.DTOs;
using InsuranceEngine.Policy.Application.Interfaces;
using InsuranceEngine.Policy.Domain.Entities;
using InsuranceEngine.Policy.Domain.Enums;
using InsuranceEngine.Policy.Domain.Services;
using InsuranceEngine.Products.Application.Features.Commands.CalculatePremium;
using InsuranceEngine.SharedKernel.CQRS;
using InsuranceEngine.SharedKernel.Interfaces;
using MediatR;

namespace InsuranceEngine.Policy.Application.Features.Commands.ApplyForQuote;

public class ApplyForQuoteCommandHandler : IRequestHandler<ApplyForQuoteCommand, Result<QuoteDto>>
{
    private readonly IUnderwritingRepository _repository;
    private readonly IMediator _mediator;
    private readonly QuoteNumberGenerator _quoteNumberGenerator;
    private readonly IEncryptionService _encryptionService;

    public ApplyForQuoteCommandHandler(
        IUnderwritingRepository repository,
        IMediator mediator,
        QuoteNumberGenerator quoteNumberGenerator,
        IEncryptionService encryptionService)
    {
        _repository = repository;
        _mediator = mediator;
        _quoteNumberGenerator = quoteNumberGenerator;
        _encryptionService = encryptionService;
    }

    public async Task<Result<QuoteDto>> Handle(ApplyForQuoteCommand request, CancellationToken cancellationToken)
    {
        // 1. Calculate Premium using Products module
        var premiumRequest = new CalculatePremiumCommand(
            request.ProductId,
            request.SumAssuredAmount,
            request.TermYears * 12, // assuming years to months
            request.SelectedRiderIds,
            new Dictionary<string, string> 
            { 
                { "Age", request.ApplicantAge.ToString() },
                { "Smoker", request.IsSmoker.ToString() },
                { "Occupation", request.ApplicantOccupation ?? "" }
            }
        );

        var premiumResult = await _mediator.Send(premiumRequest, cancellationToken);
        if (!premiumResult.IsSuccess)
            return Result.Fail<QuoteDto>(premiumResult.Error!);

        var premiumData = premiumResult.Value!;

        // 2. Generate Quote Number
        var productCodeResult = await _mediator.Send(new InsuranceEngine.Products.Application.Features.Queries.GetProductCode.GetProductCodeQuery(request.ProductId), cancellationToken);
        var productCode = productCodeResult.Value ?? "PROD";
        var sequence = await _repository.GetNextQuoteSequenceAsync();
        var quoteNumber = _quoteNumberGenerator.Generate(productCode, sequence);

        // 3. Create Quote Entity
        var quote = new Quote
        {
            Id = Guid.NewGuid(),
            QuoteNumber = quoteNumber,
            BeneficiaryId = request.BeneficiaryId,
            InsurerProductId = request.ProductId,
            Status = QuoteStatus.Draft,
            SumAssuredAmount = request.SumAssuredAmount,
            TermYears = request.TermYears,
            PremiumPaymentMode = request.PremiumPaymentMode,
            BasePremiumAmount = premiumData.BasePremium.Amount,
            RiderPremiumAmount = premiumData.RiderPremium.Amount,
            TaxAmount = premiumData.Vat.Amount,
            TotalPremiumAmount = premiumData.TotalPremium.Amount,
            Currency = premiumData.TotalPremium.CurrencyCode,
            PremiumCalculationJson = JsonSerializer.Serialize(premiumData.Breakdown),
            SelectedRidersJson = request.SelectedRiderIds != null ? JsonSerializer.Serialize(request.SelectedRiderIds) : null,
            ApplicantAge = request.ApplicantAge,
            ApplicantOccupation = request.ApplicantOccupation,
            IsSmoker = request.IsSmoker,
            ValidUntil = DateTime.UtcNow.AddDays(30),
            CreatedAt = DateTime.UtcNow,
            UpdatedAt = DateTime.UtcNow
        };

        // 4. Create Health Declaration
        var hdDto = request.HealthDeclaration;
        
        // Encrypt sensitive JSON fields
        var preExistingJson = hdDto.PreExistingConditions != null ? JsonSerializer.Serialize(hdDto.PreExistingConditions) : null;
        var encryptedPreExisting = preExistingJson != null ? await _encryptionService.EncryptAsync(preExistingJson) : null;

        var familyHistoryJson = hdDto.FamilyHistory != null ? JsonSerializer.Serialize(hdDto.FamilyHistory) : null;
        var encryptedFamilyHistory = familyHistoryJson != null ? await _encryptionService.EncryptAsync(familyHistoryJson) : null;

        var medicalDocumentsJson = hdDto.MedicalDocuments != null ? JsonSerializer.Serialize(hdDto.MedicalDocuments) : null;
        var encryptedMedicalDocs = medicalDocumentsJson != null ? await _encryptionService.EncryptAsync(medicalDocumentsJson) : null;

        var healthDeclaration = new UnderwritingHealthDeclaration
        {
            Id = Guid.NewGuid(),
            QuoteId = quote.Id,
            HeightCm = hdDto.HeightCm,
            WeightKg = hdDto.WeightKg,
            Bmi = hdDto.Bmi,
            HasPreExistingConditions = hdDto.HasPreExistingConditions,
            PreExistingConditionsJson = encryptedPreExisting,
            IsCurrentlyHospitalized = hdDto.IsCurrentlyHospitalized,
            HasFamilyHistory = hdDto.HasFamilyHistory,
            FamilyHistoryJson = encryptedFamilyHistory,
            IsSmoker = hdDto.IsSmoker,
            IsAlcoholConsumer = hdDto.IsAlcoholConsumer,
            OccupationRiskLevel = hdDto.OccupationRiskLevel,
            IsMedicalExamRequired = hdDto.IsMedicalExamRequired,
            IsMedicalExamCompleted = hdDto.IsMedicalExamCompleted,
            MedicalExamResultsJson = null, // Initial
            MedicalExamDate = hdDto.MedicalExamDate,
            MedicalDocumentsJson = encryptedMedicalDocs,
            CreatedAt = DateTime.UtcNow,
            UpdatedAt = DateTime.UtcNow
        };

        // 5. Persist
        await _repository.AddQuoteAsync(quote);
        await _repository.AddHealthDeclarationAsync(healthDeclaration);

        // 6. Return DTO
        return Result.Ok(new QuoteDto(
            quote.Id,
            quote.QuoteNumber,
            quote.BeneficiaryId,
            quote.InsurerProductId,
            quote.Status,
            new MoneyDto(quote.SumAssuredAmount, quote.SumAssuredCurrency),
            quote.TermYears,
            quote.PremiumPaymentMode,
            new MoneyDto(quote.BasePremiumAmount, quote.Currency),
            new MoneyDto(quote.RiderPremiumAmount, quote.Currency),
            new MoneyDto(quote.TotalPremiumAmount, quote.Currency),
            quote.ApplicantAge,
            quote.ApplicantOccupation,
            quote.IsSmoker,
            quote.ValidUntil,
            quote.CreatedAt
        ));
    }
}
