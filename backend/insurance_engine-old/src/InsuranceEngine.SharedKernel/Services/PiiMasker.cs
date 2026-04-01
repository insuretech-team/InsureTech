namespace InsuranceEngine.SharedKernel.Services;

public static class PiiMasker
{
    public static string MaskNid(string? nid)
    {
        if (string.IsNullOrEmpty(nid)) return string.Empty;
        if (nid.Length <= 4) return "****";
        return new string('*', nid.Length - 4) + nid[^4..];
    }

    public static string MaskPhone(string? phone)
    {
        if (string.IsNullOrEmpty(phone)) return string.Empty;
        if (phone.Length <= 4) return phone.PadRight(11, '*');
        // Mask first 7 digits of a 11-digit number
        return "*******" + phone[^4..];
    }

    public static string MaskEmail(string? email)
    {
        if (string.IsNullOrEmpty(email)) return string.Empty;
        var parts = email.Split('@');
        if (parts.Length != 2) return "*****";
        var name = parts[0];
        if (name.Length <= 2) return "**@" + parts[1];
        return name[..2] + new string('*', name.Length - 2) + "@" + parts[1];
    }
}
