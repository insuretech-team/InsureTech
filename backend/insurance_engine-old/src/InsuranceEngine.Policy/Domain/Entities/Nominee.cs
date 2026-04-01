using InsuranceEngine.SharedKernel.Domain;
using InsuranceEngine.SharedKernel.Domain.ValueObjects;

namespace InsuranceEngine.Policy.Domain.Entities;


public class Nominee : Entity<Guid>
{
    public Nominee(Guid id) : base(id) { }
    public Nominee() { }
    public Guid PolicyId { get; set; }
    public Guid? BeneficiaryId { get; set; }

    // --- Proto-aligned inline fields ---

    public string FullName { get; set; } = string.Empty;
    public string Relationship { get; set; } = string.Empty;
    public double SharePercentage { get; set; }
    public DateTime? DateOfBirth { get; set; }
    public string? NomineeDobText { get; set; }
    public string? NidNumber { get; set; }
    public string? PhoneNumber { get; set; }

    // Audit
    public DateTime CreatedAt { get; set; }
    public DateTime UpdatedAt { get; set; }
    public DateTime? DeletedAt { get; set; }
}

