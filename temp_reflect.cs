using System;
using System.Linq;
using Grpc.Core;

foreach (var ctor in typeof(ContextPropagationToken).GetConstructors(System.Reflection.BindingFlags.Public | System.Reflection.BindingFlags.NonPublic | System.Reflection.BindingFlags.Instance))
{
    Console.WriteLine(ctor.ToString());
}
