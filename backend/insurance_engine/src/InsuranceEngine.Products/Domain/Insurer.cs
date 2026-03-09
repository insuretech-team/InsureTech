using System;
using System.Collections.Generic;

namespace InsuranceEngine.Products.Domain;

public class Insurer
{
    public Guid Id { get; set; }
    public string Name { get; set; } = string.Empty;
    public string Code { get; set; } = string.Empty;
    public Guid TenantId { get; set; }
}
