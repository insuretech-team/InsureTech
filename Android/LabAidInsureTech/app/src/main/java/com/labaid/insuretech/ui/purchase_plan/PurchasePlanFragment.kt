package com.labaid.insuretech.ui.purchase_plan

import android.app.Activity
import android.app.DatePickerDialog
import android.content.Intent
import android.net.Uri
import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.ArrayAdapter
import android.widget.Toast
import androidx.activity.result.contract.ActivityResultContracts
import androidx.fragment.app.Fragment
import com.labaid.insuretech.R
import com.labaid.insuretech.adapter.BeneficiaryAdapter
import com.labaid.insuretech.data.model.beneficiary.Beneficiary
import com.labaid.insuretech.databinding.FragmentPurchasePlanBinding
import com.labaid.insuretech.utils.PaymentSuccessDialog
import dagger.hilt.android.AndroidEntryPoint
import java.text.SimpleDateFormat
import java.util.*

@AndroidEntryPoint
class PurchasePlanFragment : Fragment() {

    private lateinit var binding: FragmentPurchasePlanBinding
    private var frontNidUri: Uri? = null
    private var backNidUri: Uri? = null

    private val frontNidPickerLauncher = registerForActivityResult(
        ActivityResultContracts.StartActivityForResult()
    ) { result ->
        if (result.resultCode == Activity.RESULT_OK) {
            result.data?.data?.let { uri ->
                frontNidUri = uri
                Toast.makeText(requireContext(), "Front NID uploaded", Toast.LENGTH_SHORT).show()
                // Optionally update UI to show file is uploaded
            }
        }
    }

    private val backNidPickerLauncher = registerForActivityResult(
        ActivityResultContracts.StartActivityForResult()
    ) { result ->
        if (result.resultCode == Activity.RESULT_OK) {
            result.data?.data?.let { uri ->
                backNidUri = uri
                Toast.makeText(requireContext(), "Back NID uploaded", Toast.LENGTH_SHORT).show()
                // Optionally update UI to show file is uploaded
            }
        }
    }

    override fun onCreateView(
        inflater: LayoutInflater, container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View {
        binding = FragmentPurchasePlanBinding.inflate(inflater, container, false)
        return binding.root
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)

        setupToolbar()
        setupPlanDetails()
        setupBeneficiaryDropdown()
        setupOccupationDropdown()
        setupIncomeRangeDropdown()
        setupDateOfBirth()
        setupFileUploads()
        setupContinueButton()
    }

    private fun setupToolbar() {
        binding.purchasePlanToolbar.apply {
            ivBack.setOnClickListener {
                requireActivity().onBackPressed()
            }
            toolbarTitle.text = "Purchase Plan"
        }
    }

    private fun setupPlanDetails() {
        // TODO: Get plan details from arguments or ViewModel
        binding.apply {
            tvPlanName.text = "Seba"
            tvCoverageAmount.text = "৳ 25,000"
            tvPremiumPrice.text = "৳ 800"
            tvPolicyDuration.text = "1 year"
        }
    }

    private fun setupBeneficiaryDropdown() {
        val beneficiaries = listOf(
            Beneficiary(1, "Myself", R.drawable.ic_person),
            Beneficiary(2, "Spouse", R.drawable.ic_spouse),
            Beneficiary(3, "Children", R.drawable.ic_child)
        )

        val adapter = BeneficiaryAdapter(requireContext(), beneficiaries)
        binding.actvBeneficiary.apply {
            setAdapter(adapter)
            setText(beneficiaries[1].name, false) // Default to "Spouse"
            setOnItemClickListener { _, _, position, _ ->
                val selected = beneficiaries[position]
                // Handle selection
            }
        }
    }

    private fun setupOccupationDropdown() {
        val occupations = arrayOf(
            "Business",
            "Service",
            "Professional",
            "Student",
            "Retired",
            "Unemployed",
            "Other"
        )

        val adapter = ArrayAdapter(requireContext(), R.layout.dropdown_item, occupations)
        binding.actvOccupation.apply {
            setAdapter(adapter)
            setOnItemClickListener { _, _, position, _ ->
                val selected = occupations[position]
                // Handle selection
            }
        }
    }

    private fun setupIncomeRangeDropdown() {
        val incomeRanges = arrayOf(
            "Below ৳ 10,000",
            "৳ 10,000 - ৳ 25,000",
            "৳ 25,000 - ৳ 50,000",
            "৳ 50,000 - ৳ 1,00,000",
            "৳ 1,00,000 - ৳ 2,00,000",
            "Above ৳ 2,00,000"
        )

        val adapter = ArrayAdapter(requireContext(), R.layout.dropdown_item, incomeRanges)
        binding.actvIncomeRange.apply {
            setAdapter(adapter)
            setOnItemClickListener { _, _, position, _ ->
                val selected = incomeRanges[position]
                // Handle selection
            }
        }
    }

    private fun setupDateOfBirth() {
        binding.etDateOfBirth.setOnClickListener {
            showDatePickerDialog()
        }

        binding.layoutDateOfBirth.setEndIconOnClickListener {
            showDatePickerDialog()
        }
    }

    private fun showDatePickerDialog() {
        val calendar = Calendar.getInstance()
        val year = calendar.get(Calendar.YEAR)
        val month = calendar.get(Calendar.MONTH)
        val day = calendar.get(Calendar.DAY_OF_MONTH)

        val datePickerDialog = DatePickerDialog(
            requireContext(),
            { _, selectedYear, selectedMonth, selectedDay ->
                val selectedDate = Calendar.getInstance().apply {
                    set(selectedYear, selectedMonth, selectedDay)
                }
                val dateFormat = SimpleDateFormat("dd/MM/yyyy", Locale.getDefault())
                binding.etDateOfBirth.setText(dateFormat.format(selectedDate.time))
            },
            year,
            month,
            day
        )

        // Set max date to today (can't select future dates)
        datePickerDialog.datePicker.maxDate = System.currentTimeMillis()

        datePickerDialog.show()
    }

    private fun setupFileUploads() {
        binding.cardUploadFront.setOnClickListener {
            openFilePicker(frontNidPickerLauncher)
        }

        binding.cardUploadBack.setOnClickListener {
            openFilePicker(backNidPickerLauncher)
        }
    }

    private fun openFilePicker(launcher: androidx.activity.result.ActivityResultLauncher<Intent>) {
        val intent = Intent(Intent.ACTION_GET_CONTENT).apply {
            type = "image/*"
            putExtra(Intent.EXTRA_MIME_TYPES, arrayOf("image/jpeg", "image/png"))
        }
        launcher.launch(intent)
    }

    private fun setupContinueButton() {
        binding.btnContinueToPay.setOnClickListener {
            /*if (validateForm()) {
                // Proceed to payment
                Toast.makeText(requireContext(), "Proceeding to payment...", Toast.LENGTH_SHORT).show()
                // TODO: Navigate to payment screen
            }*/

            showPaymentSuccessDialog()
        }
    }

    private fun showPaymentSuccessDialog() {
        val dialog = PaymentSuccessDialog.newInstance()

        dialog.setOnContinueClickListener {
            // Handle what happens after user clicks Continue
            // For example: navigate to home or confirmation screen
            Toast.makeText(requireContext(), "Navigating to home...", Toast.LENGTH_SHORT).show()

            // Navigate to home or other screen
            // findNavController().navigate(R.id.action_purchasePlanFragment_to_homeFragment)
        }

        dialog.show(childFragmentManager, PaymentSuccessDialog.TAG)
    }

    private fun validateForm(): Boolean {
        val name = binding.etName.text.toString().trim()
        val dateOfBirth = binding.etDateOfBirth.text.toString().trim()

        return when {
            name.isEmpty() -> {
                binding.layoutName.error = "Name is required"
                false
            }
            dateOfBirth.isEmpty() -> {
                binding.layoutDateOfBirth.error = "Date of birth is required"
                false
            }
            frontNidUri == null -> {
                Toast.makeText(requireContext(), "Please upload front side of NID", Toast.LENGTH_SHORT).show()
                false
            }
            backNidUri == null -> {
                Toast.makeText(requireContext(), "Please upload back side of NID", Toast.LENGTH_SHORT).show()
                false
            }
            else -> {
                binding.layoutName.error = null
                binding.layoutDateOfBirth.error = null
                true
            }
        }
    }
}