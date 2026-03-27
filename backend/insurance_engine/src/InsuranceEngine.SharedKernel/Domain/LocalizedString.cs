namespace InsuranceEngine.SharedKernel.Domain;

public class LocalizedString : ValueObject
{
    public string English { get; }
    public string Bengali { get; }

    public LocalizedString(string english, string bengali)
    {
        English = english;
        Bengali = bengali;
    }

    protected override IEnumerable<object> GetEqualityComponents()
    {
        yield return English;
        yield return Bengali;
    }

    public override string ToString() => English;
    
    public string Get(string language) => language.ToLower() == "bn" ? Bengali : English;
}
