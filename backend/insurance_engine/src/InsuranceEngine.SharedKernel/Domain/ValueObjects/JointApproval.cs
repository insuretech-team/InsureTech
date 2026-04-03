using System;
using System.Collections.Generic;
using System.Linq;

namespace InsuranceEngine.SharedKernel.Domain.ValueObjects;

public record ApprovalRecord(string Role, Guid ApproverId, DateTime ApprovedAt, string Notes);

public class JointApproval
{
    private readonly List<ApprovalRecord> _approvals = new();
    public IReadOnlyList<ApprovalRecord> Approvals => _approvals.AsReadOnly();

    public int RequiredApprovals { get; }
    public IReadOnlyList<string> RequiredRoles { get; }

    public JointApproval(int requiredApprovals, IEnumerable<string> requiredRoles)
    {
        RequiredApprovals = requiredApprovals;
        RequiredRoles = requiredRoles.ToList();
    }

    public void AddApproval(string role, Guid approverId, string notes)
    {
        if (!RequiredRoles.Contains(role))
            throw new InvalidOperationException($"Role '{role}' is not authorized for this approval level.");

        if (_approvals.Any(a => a.Role == role && a.ApproverId == approverId))
            throw new InvalidOperationException("This approver has already submitted an approval for this request.");

        _approvals.Add(new ApprovalRecord(role, approverId, DateTime.UtcNow, notes));
    }

    public bool IsFullyApproved()
    {
        // Check if we have the required count
        if (_approvals.Count < RequiredApprovals) return false;

        // Ensure all required roles have at least one approval (if roles are distinct)
        // For SRS v3.11 L3: Business Admin AND Focal Person
        foreach (var role in RequiredRoles)
        {
            if (!_approvals.Any(a => a.Role == role)) return false;
        }

        return true;
    }

    public string GetStatusSummary()
    {
        var approvedRoles = string.Join(", ", _approvals.Select(a => a.Role));
        return IsFullyApproved() 
            ? "FULLY_APPROVED" 
            : $"PENDING_APPROVAL (Approved by: {approvedRoles}, Missing: {string.Join(", ", RequiredRoles.Except(_approvals.Select(a => a.Role)))})";
    }
}
