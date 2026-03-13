package com.labaid.insuretech.adapter

import android.view.LayoutInflater
import android.view.ViewGroup
import androidx.recyclerview.widget.RecyclerView
import com.labaid.insuretech.data.model.icoverage.HospitalizationCoverage
import com.labaid.insuretech.databinding.ItemHospitalizationCoverageBinding

class HospitalizationCoveragesAdapter(
    private val coverages: List<HospitalizationCoverage>
) : RecyclerView.Adapter<HospitalizationCoveragesAdapter.ViewHolder>() {

    inner class ViewHolder(val binding: ItemHospitalizationCoverageBinding)
        : RecyclerView.ViewHolder(binding.root)

    override fun onCreateViewHolder(parent: ViewGroup, viewType: Int): ViewHolder {
        val binding = ItemHospitalizationCoverageBinding.inflate(
            LayoutInflater.from(parent.context), parent, false
        )
        return ViewHolder(binding)
    }

    override fun onBindViewHolder(holder: ViewHolder, position: Int) {
        val coverage = coverages[position]
        holder.binding.apply {
            tvCoverageName.text = coverage.name
            ivCoverageIcon.setImageResource(coverage.iconResId)
        }
    }

    override fun getItemCount(): Int = coverages.size
}