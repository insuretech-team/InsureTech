package com.example.lifeplans.ui.plan

import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.ImageView
import android.widget.TextView
import androidx.recyclerview.widget.RecyclerView
import com.bumptech.glide.Glide
import com.google.android.material.button.MaterialButton
import com.labaid.insuretech.R

// Data class for Life Plan
data class LifePlan(
    val id: Int,
    val name: String,
    val imageUrl: String,
    val coverageItems: List<String>,
    val poweredBy: String = "Charmed Life Insurance"
)

// ViewHolder for the plan item
class LifePlanViewHolder(itemView: View) : RecyclerView.ViewHolder(itemView) {
    val ivPlanImage: ImageView = itemView.findViewById(R.id.ivPlanImage)
    val tvPlanName: TextView = itemView.findViewById(R.id.tvPlanName)
    val tvCoverageItem1: TextView = itemView.findViewById(R.id.tvCoverageItem1)
    val tvCoverageItem2: TextView = itemView.findViewById(R.id.tvCoverageItem2)
    val btnContinue: MaterialButton = itemView.findViewById(R.id.btnContinue)
}

// Adapter for the RecyclerView
class LifePlanAdapter(
    private val plans: List<LifePlan>,
    private val onContinueClick: (LifePlan) -> Unit
) : RecyclerView.Adapter<LifePlanViewHolder>() {

    override fun onCreateViewHolder(parent: ViewGroup, viewType: Int): LifePlanViewHolder {
        val view = LayoutInflater.from(parent.context)
            .inflate(R.layout.item_plan_card, parent, false)
        return LifePlanViewHolder(view)
    }

    override fun onBindViewHolder(holder: LifePlanViewHolder, position: Int) {
        val plan = plans[position]

        holder.tvPlanName.text = plan.name

        // Set coverage items (you might want to make this more dynamic)
        if (plan.coverageItems.isNotEmpty()) {
            holder.tvCoverageItem1.text = plan.coverageItems[0]
        }
        if (plan.coverageItems.size > 1) {
            holder.tvCoverageItem2.text = plan.coverageItems[1]
        }

        // Load image using your preferred image loading library (Glide, Picasso, Coil, etc.)
        /*Glide.with(holder.itemView.context)
            .load(plan.imageUrl)
           .into(holder.ivPlanImage)*/

        Glide.with(holder.itemView.context)
            .load(plan.imageUrl)
            .placeholder(R.drawable.plan_image_placeholder) // while loading
            .error(R.drawable.plan_image_placeholder)       // if failed
            .fallback(R.drawable.plan_image_placeholder)    // if url is null
            .into(holder.ivPlanImage)

        holder.btnContinue.setOnClickListener {
            onContinueClick(plan)
        }
    }

    override fun getItemCount(): Int = plans.size
}