// PoliSync.VehicleInsurance - Assembly Marker
namespace PoliSync.VehicleInsurance;

public sealed class AssemblyMarker
{
    public static readonly string AssemblyName = typeof(AssemblyMarker).Assembly.GetName().Name!;
}
