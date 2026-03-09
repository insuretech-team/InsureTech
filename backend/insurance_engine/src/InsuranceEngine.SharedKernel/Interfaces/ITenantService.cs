using System;

namespace InsuranceEngine.SharedKernel.Interfaces;

public interface ITenantService
{
    Guid GetTenantId();
}

public class DefaultTenantService : ITenantService
{
    private readonly Guid _defaultTenantId = Guid.Parse("11111111-1111-1111-1111-111111111111");
    public Guid GetTenantId() => _defaultTenantId;
}
