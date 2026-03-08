using System.Globalization;
using Google.Protobuf.WellKnownTypes;
using Insuretech.Underwriting.Entity.V1;
using PoliSync.SharedKernel.CQRS;
using PoliSync.SharedKernel.Domain;

namespace PoliSync.Underwriting.Domain;

public sealed class HealthDeclarationAggregate
{
    private readonly HealthDeclaration _declaration;
    private readonly List<DomainEvent> _domainEvents = [];

    private HealthDeclarationAggregate(HealthDeclaration declaration)
    {
        _declaration = declaration;
    }

    public HealthDeclaration Declaration => _declaration;
    public IReadOnlyCollection<DomainEvent> DomainEvents => _domainEvents.AsReadOnly();

    public void ClearDomainEvents() => _domainEvents.Clear();

    public static Result<HealthDeclarationAggregate> Create(
        string quoteId,
        int applicantAge,
        int heightCm,
        string weightKg,
        bool hasPreExistingConditions,
        string preExistingConditions,
        bool smoker,
        bool alcoholConsumer,
        string occupationRiskLevel,
        bool isCurrentlyHospitalized = false,
        bool hasFamilyHistory = false,
        string familyHistory = "")
    {
        if (string.IsNullOrWhiteSpace(quoteId))
        {
            return Result.Fail<HealthDeclarationAggregate>("INVALID_QUOTE_ID", "QuoteId is required");
        }

        if (heightCm <= 0)
        {
            return Result.Fail<HealthDeclarationAggregate>("INVALID_HEIGHT", "Height must be greater than zero");
        }

        if (!decimal.TryParse(weightKg, NumberStyles.Float, CultureInfo.InvariantCulture, out var parsedWeight) || parsedWeight <= 0)
        {
            return Result.Fail<HealthDeclarationAggregate>("INVALID_WEIGHT", "WeightKg must be a positive numeric value");
        }

        if (hasPreExistingConditions && string.IsNullOrWhiteSpace(preExistingConditions))
        {
            return Result.Fail<HealthDeclarationAggregate>(
                "INVALID_PRE_EXISTING_CONDITIONS",
                "Pre-existing conditions are required when hasPreExistingConditions is true");
        }

        var normalizedOccupation = NormalizeOccupationRiskLevel(occupationRiskLevel);
        if (normalizedOccupation is null)
        {
            return Result.Fail<HealthDeclarationAggregate>(
                "INVALID_OCCUPATION_RISK_LEVEL",
                "OccupationRiskLevel must be LOW, MEDIUM, or HIGH");
        }

        var bmi = CalculateBmi(parsedWeight, heightCm);
        var medicalExamRequired = RequiresMedicalExam(
            applicantAge,
            bmi,
            hasPreExistingConditions,
            smoker,
            isCurrentlyHospitalized,
            normalizedOccupation);

        var declaration = new HealthDeclaration
        {
            Id = Guid.NewGuid().ToString("N"),
            QuoteId = quoteId,
            HeightCm = heightCm,
            WeightKg = parsedWeight.ToString("0.##", CultureInfo.InvariantCulture),
            Bmi = bmi.ToString("0.00", CultureInfo.InvariantCulture),
            HasPreExistingConditions = hasPreExistingConditions,
            PreExistingConditions = hasPreExistingConditions ? preExistingConditions : string.Empty,
            IsCurrentlyHospitalized = isCurrentlyHospitalized,
            HasFamilyHistory = hasFamilyHistory,
            FamilyHistory = hasFamilyHistory ? familyHistory : string.Empty,
            Smoker = smoker,
            AlcoholConsumer = alcoholConsumer,
            OccupationRiskLevel = normalizedOccupation,
            MedicalExamRequired = medicalExamRequired,
            MedicalExamCompleted = false,
            MedicalExamResults = string.Empty,
            MedicalDocuments = string.Empty
        };

        var aggregate = new HealthDeclarationAggregate(declaration);
        aggregate._domainEvents.Add(new HealthDeclarationSubmittedDomainEvent(declaration.Id, declaration.QuoteId));
        return Result.Ok(aggregate);
    }

    public Result CompleteMedicalExam(DateTime examDateUtc, string examResults, string medicalDocuments)
    {
        if (!_declaration.MedicalExamRequired)
        {
            return Result.Fail("MEDICAL_EXAM_NOT_REQUIRED", "Medical exam is not required for this declaration");
        }

        if (string.IsNullOrWhiteSpace(examResults))
        {
            return Result.Fail("INVALID_MEDICAL_EXAM_RESULTS", "Medical exam results are required");
        }

        _declaration.MedicalExamCompleted = true;
        _declaration.MedicalExamResults = examResults;
        _declaration.MedicalDocuments = medicalDocuments ?? string.Empty;
        _declaration.MedicalExamDate = Timestamp.FromDateTime(examDateUtc.ToUniversalTime());

        _domainEvents.Add(new MedicalExamCompletedDomainEvent(_declaration.Id, _declaration.QuoteId));
        return Result.Ok();
    }

    private static decimal CalculateBmi(decimal weightKg, int heightCm)
    {
        var heightM = heightCm / 100m;
        if (heightM <= 0)
        {
            return 0m;
        }

        return weightKg / (heightM * heightM);
    }

    private static bool RequiresMedicalExam(
        int applicantAge,
        decimal bmi,
        bool hasPreExistingConditions,
        bool smoker,
        bool isCurrentlyHospitalized,
        string occupationRiskLevel)
        => hasPreExistingConditions
           || isCurrentlyHospitalized
           || applicantAge >= 55
           || smoker
           || bmi >= 30m
           || string.Equals(occupationRiskLevel, "HIGH", StringComparison.Ordinal);

    private static string? NormalizeOccupationRiskLevel(string occupationRiskLevel)
    {
        if (string.IsNullOrWhiteSpace(occupationRiskLevel))
        {
            return "LOW";
        }

        var normalized = occupationRiskLevel.Trim().ToUpperInvariant();
        return normalized is "LOW" or "MEDIUM" or "HIGH"
            ? normalized
            : null;
    }
}

public sealed record HealthDeclarationSubmittedDomainEvent(string DeclarationId, string QuoteId) : DomainEvent;

public sealed record MedicalExamCompletedDomainEvent(string DeclarationId, string QuoteId) : DomainEvent;
