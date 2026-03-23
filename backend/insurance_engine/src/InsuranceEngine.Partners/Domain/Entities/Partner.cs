using System;
using System.Collections.Generic;
using InsuranceEngine.SharedKernel.Domain;
using InsuranceEngine.Partners.Domain.Enums;
using InsuranceEngine.SharedKernel.CQRS;

namespace InsuranceEngine.Partners.Domain.Entities;

public class Partner : AggregateRoot<Guid>
{
    public string Name { get; private set; } = string.Empty;
    public string Code { get; private set; } = string.Empty;
    public string Email { get; private set; } = string.Empty;
    public string? Phone { get; private set; }
    public PartnerStatus Status { get; private set; }
    public string? Address { get; private set; }
    public DateTime CreatedAt { get; private set; }
    public DateTime UpdatedAt { get; private set; }

    public List<Agent> Agents { get; private set; } = new();

    private Partner() { }

    public static Partner Create(string name, string code, string email, string? phone = null, string? address = null)
    {
        return new Partner
        {
            Id = Guid.NewGuid(),
            Name = name,
            Code = code,
            Email = email,
            Phone = phone,
            Address = address,
            Status = PartnerStatus.Active,
            CreatedAt = DateTime.UtcNow,
            UpdatedAt = DateTime.UtcNow
        };
    }

    public Result AddAgent(string name, string code, string email, string? phone = null)
    {
        if (Agents.Any(a => a.Code == code))
            return Result.Failure("Agent code already exists");

        var agent = Agent.Create(Id, name, code, email, phone);
        Agents.Add(agent);
        UpdatedAt = DateTime.UtcNow;
        return Result.Success();
    }

    public void UpdateStatus(PartnerStatus status)
    {
        Status = status;
        UpdatedAt = DateTime.UtcNow;
    }
}

public class Agent : Entity<Guid>
{
    public Guid PartnerId { get; private set; }
    public string Name { get; private set; } = string.Empty;
    public string Code { get; private set; } = string.Empty;
    public string Email { get; private set; } = string.Empty;
    public string? Phone { get; private set; }
    public PartnerStatus Status { get; private set; }
    public DateTime CreatedAt { get; private set; }
    public DateTime UpdatedAt { get; private set; }

    private Agent() { }

    internal static Agent Create(Guid partnerId, string name, string code, string email, string? phone = null)
    {
        return new Agent
        {
            Id = Guid.NewGuid(),
            PartnerId = partnerId,
            Name = name,
            Code = code,
            Email = email,
            Phone = phone,
            Status = PartnerStatus.Active,
            CreatedAt = DateTime.UtcNow,
            UpdatedAt = DateTime.UtcNow
        };
    }
}
