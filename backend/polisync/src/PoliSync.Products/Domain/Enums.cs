namespace PoliSync.Products.Domain;

/// <summary>
/// Product category — maps to proto ProductCategory enum.
/// </summary>
public enum ProductCategory
{
    Unspecified = 0,
    Motor = 1,
    Health = 2,
    Travel = 3,
    Home = 4,
    Device = 5,
    Agricultural = 6,
    Life = 7
}

/// <summary>
/// Product lifecycle status — maps to proto ProductStatus enum.
/// </summary>
public enum ProductStatus
{
    Unspecified = 0,
    Draft = 1,
    Active = 2,
    Inactive = 3,
    Discontinued = 4
}

/// <summary>
/// Pricing rule type — maps to proto RuleType enum.
/// </summary>
public enum RuleType
{
    Unspecified = 0,
    AgeBased = 1,
    LocationBased = 2,
    OccupationBased = 3,
    VehicleType = 4,
    HealthCondition = 5
}

/// <summary>
/// Pricing action type — maps to proto ActionType enum.
/// </summary>
public enum ActionType
{
    Unspecified = 0,
    IncreasePercentage = 1,
    DecreasePercentage = 2,
    FixedAmount = 3,
    Reject = 4
}
