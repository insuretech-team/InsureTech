package com.labaid.insuretech.ui.plan

import android.os.Bundle
import androidx.fragment.app.Fragment
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import androidx.navigation.fragment.findNavController
import androidx.recyclerview.widget.RecyclerView
import com.example.lifeplans.ui.plan.LifePlan
import com.example.lifeplans.ui.plan.LifePlanAdapter
import com.google.android.material.bottomsheet.BottomSheetDialog
import com.labaid.insuretech.R
import com.labaid.insuretech.databinding.FragmentPlanBinding
import com.labaid.insuretech.databinding.LayoutDeclarationBottomsheetBinding
import com.lifeplus.onetapservice.utils.Extensions.hideKeyboard
import com.lifeplus.onetapservice.utils.Extensions.toast

class PlanFragment : Fragment() {

    private lateinit var binding: FragmentPlanBinding
    private lateinit var plansRecyclerView: RecyclerView
    private lateinit var planAdapter: LifePlanAdapter
    override fun onCreateView(
        inflater: LayoutInflater, container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View? {
        binding = FragmentPlanBinding.inflate(inflater,container,false)
        return binding.root
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)

        binding.planToolbar.ivBack.setOnClickListener {
            findNavController().popBackStack()
        }

        binding.planToolbar.toolbarTitle.text = "Plan"

        setupRecyclerView(view)
    }

    private fun setupRecyclerView(view: View) {
        plansRecyclerView = view.findViewById(R.id.plansRecyclerView)

        val plans = listOf(
            LifePlan(
                id = 1,
                name = "Vorosa",
                imageUrl = "https://example.com/vorosa-image.jpg",
                coverageItems = listOf(
                    "Emergency medical expenses",
                    "Adventure/Sports coverage"
                )
            ),
            LifePlan(
                id = 2,
                name = "Astha",
                imageUrl = "https://example.com/astha-image.jpg",
                coverageItems = listOf(
                    "Emergency medical expenses",
                    "Adventure/Sports coverage"
                )
            )
        )

        planAdapter = LifePlanAdapter(plans) { selectedPlan ->
            // Handle continue button click
            displayDeclarationSheet(selectedPlan)
        }

        plansRecyclerView.adapter = planAdapter
    }

    private fun displayDeclarationSheet(selectedPlan: LifePlan) {
        val bottomSheet = BottomSheetDialog(requireContext())
        val declarationBinding = LayoutDeclarationBottomsheetBinding.inflate(layoutInflater)
        bottomSheet.setContentView(declarationBinding.root)
        bottomSheet.setCanceledOnTouchOutside(true)

        declarationBinding.apply {

            btnAgree.setOnClickListener {
                bottomSheet.dismiss()
                findNavController().navigate(R.id.action_planFragment_to_planDetailsFragment)
            }

        }

        bottomSheet.show()
    }
}