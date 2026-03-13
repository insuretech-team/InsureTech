package com.labaid.insuretech.adapter

import android.content.Context
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.ArrayAdapter
import android.widget.ImageView
import android.widget.TextView
import com.labaid.insuretech.R
import com.labaid.insuretech.data.model.beneficiary.Beneficiary

class BeneficiaryAdapter(
    context: Context,
    private val beneficiaries: List<Beneficiary>
) : ArrayAdapter<Beneficiary>(context, 0, beneficiaries) {

    override fun getView(position: Int, convertView: View?, parent: ViewGroup): View {
        return createView(position, convertView, parent)
    }

    override fun getDropDownView(position: Int, convertView: View?, parent: ViewGroup): View {
        return createView(position, convertView, parent)
    }

    private fun createView(position: Int, convertView: View?, parent: ViewGroup): View {
        val view = convertView ?: LayoutInflater.from(context)
            .inflate(R.layout.item_beneficiary, parent, false)

        val beneficiary = beneficiaries[position]

        val ivIcon = view.findViewById<ImageView>(R.id.ivBeneficiaryIcon)
        val tvName = view.findViewById<TextView>(R.id.tvBeneficiaryName)

        ivIcon.setImageResource(beneficiary.iconResId)
        tvName.text = beneficiary.name

        return view
    }
}