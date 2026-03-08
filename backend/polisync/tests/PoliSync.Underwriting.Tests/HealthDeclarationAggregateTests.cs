using FluentAssertions;
using PoliSync.Underwriting.Domain;
using Xunit;

namespace PoliSync.Underwriting.Tests;

public class HealthDeclarationAggregateTests
{
    [Fact]
    public void Create_WithValidInput_ProducesDeclarationAndMedicalExamFlag()
    {
        var result = HealthDeclarationAggregate.Create(
            quoteId: Guid.NewGuid().ToString("N"),
            applicantAge: 59,
            heightCm: 170,
            weightKg: "92",
            hasPreExistingConditions: true,
            preExistingConditions: "[\"diabetes\"]",
            smoker: true,
            alcoholConsumer: false,
            occupationRiskLevel: "HIGH");

        result.IsSuccess.Should().BeTrue();
        result.Value.Should().NotBeNull();
        result.Value!.Declaration.MedicalExamRequired.Should().BeTrue();
        result.Value.Declaration.Bmi.Should().NotBeNullOrWhiteSpace();
    }

    [Fact]
    public void Create_WithInvalidOccupationRiskLevel_FailsValidation()
    {
        var result = HealthDeclarationAggregate.Create(
            quoteId: Guid.NewGuid().ToString("N"),
            applicantAge: 30,
            heightCm: 170,
            weightKg: "70",
            hasPreExistingConditions: false,
            preExistingConditions: string.Empty,
            smoker: false,
            alcoholConsumer: false,
            occupationRiskLevel: "EXTREME");

        result.IsFailure.Should().BeTrue();
        result.Error!.Code.Should().Be("INVALID_OCCUPATION_RISK_LEVEL");
    }

    [Fact]
    public void Create_WithPreExistingFlagWithoutConditions_FailsValidation()
    {
        var result = HealthDeclarationAggregate.Create(
            quoteId: Guid.NewGuid().ToString("N"),
            applicantAge: 30,
            heightCm: 170,
            weightKg: "70",
            hasPreExistingConditions: true,
            preExistingConditions: string.Empty,
            smoker: false,
            alcoholConsumer: false,
            occupationRiskLevel: "LOW");

        result.IsFailure.Should().BeTrue();
        result.Error!.Code.Should().Be("INVALID_PRE_EXISTING_CONDITIONS");
    }

    [Fact]
    public void CompleteMedicalExam_WhenRequired_UpdatesMedicalExamState()
    {
        var aggregate = HealthDeclarationAggregate.Create(
            quoteId: Guid.NewGuid().ToString("N"),
            applicantAge: 56,
            heightCm: 170,
            weightKg: "90",
            hasPreExistingConditions: false,
            preExistingConditions: string.Empty,
            smoker: true,
            alcoholConsumer: false,
            occupationRiskLevel: "MEDIUM").Value!;

        var result = aggregate.CompleteMedicalExam(
            examDateUtc: DateTime.UtcNow,
            examResults: "{\"status\":\"FIT\"}",
            medicalDocuments: "[\"doc-1\"]");

        result.IsSuccess.Should().BeTrue();
        aggregate.Declaration.MedicalExamCompleted.Should().BeTrue();
        aggregate.Declaration.MedicalExamResults.Should().Contain("FIT");
        aggregate.Declaration.MedicalExamDate.Should().NotBeNull();
    }
}
