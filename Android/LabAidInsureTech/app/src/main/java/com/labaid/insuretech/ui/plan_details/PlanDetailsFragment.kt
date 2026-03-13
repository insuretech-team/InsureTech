package com.labaid.insuretech.ui.plan_details

import android.os.Bundle
import androidx.fragment.app.Fragment
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import androidx.navigation.fragment.findNavController
import androidx.recyclerview.widget.GridLayoutManager
import com.labaid.insuretech.R
import com.labaid.insuretech.adapter.HospitalizationCoveragesAdapter
import com.labaid.insuretech.data.model.icoverage.HospitalizationCoverage
import com.labaid.insuretech.databinding.FragmentPlanDetailsBinding
import dagger.hilt.android.AndroidEntryPoint

@AndroidEntryPoint
class PlanDetailsFragment : Fragment() {

    private lateinit var binding: FragmentPlanDetailsBinding

    override fun onCreateView(
        inflater: LayoutInflater, container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View {
        binding = FragmentPlanDetailsBinding.inflate(inflater, container, false)
        return binding.root
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)

        setupToolbar()
        setupPlanDetails()
        setupHospitalizationCoverages()
        setupPurchaseButton()
    }

    private fun setupToolbar() {
        binding.planDetailsToolbar.apply {
            ivBack.setOnClickListener {
                requireActivity().onBackPressed()
            }
            toolbarTitle.text = "Plan Details"
        }
    }

    private fun setupPlanDetails() {
        // TODO: Get plan details from arguments or ViewModel
        binding.apply {
            tvPlanName.text = "Seba"
            tvCoverageAmount.text = "৳ 25,000"
            tvPremiumPrice.text = "৳ 800"
            tvPolicyDuration.text = "1 year"
            tvPolicyPurpose.text = "To provide financial protection against medical expenses arising from illness, accidents, hospitalization, and emergency medical care."
        }
    }

    private fun setupHospitalizationCoverages() {
        val coverages = listOf(
            HospitalizationCoverage(
                name = "In-patient treatment expenses",
                iconResId = R.drawable.ic_hospital_bed
            ),
            HospitalizationCoverage(
                name = "Cabin room rent",
                iconResId = R.drawable.ic_hospital_bed
            ),
            HospitalizationCoverage(
                name = "ICU / CCU charges",
                iconResId = R.drawable.ic_hospital_bed
            ),
            HospitalizationCoverage(
                name = "Doctor & specialist consultation fees",
                iconResId = R.drawable.ic_hospital_bed
            )
        )

        val adapter = HospitalizationCoveragesAdapter(coverages)
        binding.rvHospitalizationCoverages.apply {
            layoutManager = GridLayoutManager(requireContext(), 2)
            this.adapter = adapter
            setHasFixedSize(true)
        }
    }

    private fun setupPurchaseButton() {
        binding.btnPurchasePlan.setOnClickListener {
            findNavController().navigate(R.id.action_planDetailsFragment_to_purchasePlanFragment)
        }
    }
}