using System.ComponentModel.DataAnnotations;

namespace InsuranceEngine.Domain.Entities;

public class ContactInfo
{
    [MaxLength(50)]
    public string? MobileNumber { get; set; }

    [MaxLength(255)]
    public string? Email { get; set; }

    [MaxLength(50)]
    public string? AlternateMobile { get; set; }

    [MaxLength(50)]
    public string? Landline { get; set; }
}
