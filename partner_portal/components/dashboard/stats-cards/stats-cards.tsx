import { statsData } from "@/lib/stats-cards";
import StatsCard from "./card/stats-card";

const StatsCards = () => {
  return (
    <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
      {statsData.map((stat) => (
        <StatsCard key={stat.id} {...stat} />
      ))}
    </div>
  );
};

export default StatsCards;
