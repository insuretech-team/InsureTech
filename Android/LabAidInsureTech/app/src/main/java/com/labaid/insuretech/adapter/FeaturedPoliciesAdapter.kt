package com.labaid.insuretech.adapter

import android.view.LayoutInflater
import android.view.ViewGroup
import androidx.recyclerview.widget.RecyclerView
import com.labaid.insuretech.databinding.ItemFeaturedPolicyBinding
import com.labaid.insuretech.data.model.featured_policy.FeaturedPolicy

class FeaturedPoliciesAdapter(
    private val policies: List<FeaturedPolicy>
) : RecyclerView.Adapter<FeaturedPoliciesAdapter.ViewHolder>() {

    inner class ViewHolder(val binding: ItemFeaturedPolicyBinding)
        : RecyclerView.ViewHolder(binding.root)

    override fun onCreateViewHolder(parent: ViewGroup, viewType: Int): ViewHolder {
        val binding = ItemFeaturedPolicyBinding.inflate(
            LayoutInflater.from(parent.context), parent, false
        )
        return ViewHolder(binding)
    }

    override fun onBindViewHolder(holder: ViewHolder, position: Int) {
        val policy = policies[position]
        holder.binding.apply {
            tvPolicyName.text = policy.name
            ivPolicyImage.setImageResource(policy.imageResId)
            ivPoweredBy.setImageResource(policy.poweredByLogoResId)

            // Set coverage items
            if (policy.coverage.isNotEmpty()) {
                tvCoverageItem1.text = policy.coverage.getOrNull(0) ?: ""
            }
            if (policy.coverage.size > 1) {
                tvCoverageItem2.text = policy.coverage.getOrNull(1) ?: ""
            }
        }
    }

    override fun getItemCount(): Int = policies.size
}