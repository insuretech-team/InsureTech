package com.labaid.insuretech.adapter

import android.view.LayoutInflater
import android.view.ViewGroup
import androidx.recyclerview.widget.RecyclerView
import com.labaid.insuretech.databinding.HomeSliderLayoutBinding

class HomeSliderAdapter(
    private val images: List<Int>
) : RecyclerView.Adapter<HomeSliderAdapter.ViewHolder>() {

    inner class ViewHolder(val binding: HomeSliderLayoutBinding)
        : RecyclerView.ViewHolder(binding.root)

    override fun onCreateViewHolder(parent: ViewGroup, viewType: Int): ViewHolder {
        val binding = HomeSliderLayoutBinding.inflate(
            LayoutInflater.from(parent.context), parent, false
        )
        return ViewHolder(binding)
    }

    override fun onBindViewHolder(holder: ViewHolder, position: Int) {
        holder.binding.ivHomeSlider.setImageResource(images[position])
    }

    override fun getItemCount(): Int = images.size
}
