using Microsoft.Extensions.DependencyInjection;
using PoliSync.RulesEngine.Services;
using RulesEngine.Models;
using RulesEngine.Actions;
using RulesEngineAlias = RulesEngine.RulesEngine;

namespace PoliSync.RulesEngine;

public static class DependencyInjection
{
    public static IServiceCollection AddRulesEngineServices(this IServiceCollection services)
    {
        // Register Microsoft RulesEngine
        services.AddSingleton<RulesEngineAlias>(provider =>
        {
            var reSettings = new ReSettings
            {
                CustomTypes = new[] { typeof(Math), typeof(Convert) },
                CustomActions = new Dictionary<string, Func<ActionBase>>()
            };
            return new RulesEngineAlias(Array.Empty<Workflow>(), reSettings);
        });

        // Register PoliSync RulesEngine services
        services.AddSingleton<IBusinessWorkflowService, BusinessWorkflowService>();
        services.AddSingleton<IBusinessWorkflowRepository, BusinessWorkflowRepository>();
        services.AddSingleton<IBusinessRuleEvaluationService, BusinessRuleEvaluationService>();

        return services;
    }
}
