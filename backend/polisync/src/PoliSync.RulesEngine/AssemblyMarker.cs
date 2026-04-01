// PoliSync.RulesEngine - Assembly Marker
namespace PoliSync.RulesEngine;

public sealed class AssemblyMarker
{
    public static readonly string AssemblyName = typeof(AssemblyMarker).Assembly.GetName().Name!;
}
