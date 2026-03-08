using System.Reflection;
using System.Runtime.Serialization;

namespace InsuranceEngine.Domain.Enums;

public static class EnumMemberValueMapper
{
    public static string GetEnumMemberValue<TEnum>(TEnum value) where TEnum : struct, Enum
    {
        string? name = Enum.GetName(value);
        if (name is null)
        {
            return value.ToString();
        }

        MemberInfo memberInfo = typeof(TEnum).GetMember(name).First();
        EnumMemberAttribute? attribute = memberInfo.GetCustomAttribute<EnumMemberAttribute>();
        return attribute?.Value ?? name;
    }

    public static TEnum ParseEnumMemberValue<TEnum>(string value) where TEnum : struct, Enum
    {
        foreach (TEnum enumValue in Enum.GetValues<TEnum>())
        {
            if (string.Equals(GetEnumMemberValue(enumValue), value, StringComparison.OrdinalIgnoreCase))
            {
                return enumValue;
            }
        }

        throw new ArgumentException($"Unknown {typeof(TEnum).Name} value: {value}");
    }
}
