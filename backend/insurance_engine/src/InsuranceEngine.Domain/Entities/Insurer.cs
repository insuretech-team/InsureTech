using System;
using System.Collections.Generic;

namespace InsuranceEngine.Domain.Entities;

public class Insurer
{
    public Guid Id { get; set; }
    public string Name { get; set; } = string.Empty;
    public string Code { get; set; } = string.Empty;
    
    public List<Product> Products { get; set; } = new();
}
