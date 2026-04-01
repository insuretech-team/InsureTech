// PoliSync.LifeInsurance - Assembly Marker
namespace PoliSync.LifeInsurance;

public sealed class AssemblyMarker
{
    public static readonly string AssemblyName = typeof(AssemblyMarker).Assembly.GetName().Name!;
}
