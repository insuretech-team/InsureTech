namespace InsuranceEngine.SharedKernel.Domain.ValueObjects;

/// <summary>
/// Value object for address information.
/// </summary>
public record Address
{
    public string? AddressLine1 { get; init; }
    public string? AddressLine2 { get; init; }
    public string? City { get; init; }
    public string? District { get; init; }
    public string? Division { get; init; }
    public string? PostalCode { get; init; }
    public string Country { get; init; } = "Bangladesh";
    public double? Latitude { get; init; }
    public double? Longitude { get; init; }

    public Address() { }

    public Address(string addressLine1, string city, string? district = null, string country = "Bangladesh")
    {
        AddressLine1 = addressLine1;
        City = city;
        District = district;
        Country = country;
    }
}
