namespace InsuranceEngine.SharedKernel.Domain.ValueObjects;

/// <summary>
/// Value object for primary contact/focal person information.
/// </summary>
public record PrimaryContact
{
    public string? Name { get; init; }
    public string? Email { get; init; }
    public string? Phone { get; init; }
    public string? Department { get; init; }

    public PrimaryContact() { }

    public PrimaryContact(string? name, string? email, string? phone)
    {
        Name = name;
        Email = email;
        Phone = phone;
    }
}
