namespace InsuranceEngine.SharedKernel.Domain.ValueObjects;

/// <summary>
/// Value object for contact information.
/// </summary>
public record ContactInfo
{
    public string? MobileNumber { get; init; }
    public string? AlternateMobile { get; init; }
    public string? Email { get; init; }
    public string? Landline { get; init; }

    public ContactInfo() { }

    public ContactInfo(string mobileNumber, string? email = null)
    {
        MobileNumber = mobileNumber;
        Email = email;
    }
}
