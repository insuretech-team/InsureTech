package com.labaid.insuretech.adapter

import android.view.LayoutInflater
import android.view.ViewGroup
import androidx.recyclerview.widget.RecyclerView
import com.labaid.insuretech.databinding.ItemLatestUpdateBinding
import com.labaid.insuretech.data.model.latest_update.LatestUpdate

class LatestUpdatesAdapter(
    private val updates: List<LatestUpdate>
) : RecyclerView.Adapter<LatestUpdatesAdapter.ViewHolder>() {

    inner class ViewHolder(val binding: ItemLatestUpdateBinding)
        : RecyclerView.ViewHolder(binding.root)

    override fun onCreateViewHolder(parent: ViewGroup, viewType: Int): ViewHolder {
        val binding = ItemLatestUpdateBinding.inflate(
            LayoutInflater.from(parent.context), parent, false
        )
        return ViewHolder(binding)
    }

    override fun onBindViewHolder(holder: ViewHolder, position: Int) {
        val update = updates[position]
        holder.binding.apply {
            tvUpdateTitle.text = update.title
            tvUpdateDescription.text = update.description
            ivUpdateImage.setImageResource(update.imageResId)
        }
    }

    override fun getItemCount(): Int = updates.size
}