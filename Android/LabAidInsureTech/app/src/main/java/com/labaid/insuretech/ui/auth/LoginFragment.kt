package com.labaid.insuretech.ui.auth

import android.content.Context
import android.graphics.Color
import android.os.Bundle
import android.text.Editable
import android.text.TextWatcher
import androidx.fragment.app.Fragment
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.view.inputmethod.InputMethodManager
import android.widget.Toast
import androidx.core.content.ContextCompat
import androidx.fragment.app.viewModels
import androidx.navigation.fragment.findNavController
import com.labaid.insuretech.R
import com.labaid.insuretech.databinding.FragmentLoginBinding
import com.labaid.insuretech.databinding.FragmentSplashBinding
import com.labaid.insuretech.utils.Constants
import com.labaid.insuretech.viewmodel.AuthViewModel
import com.lifeplus.onetapservice.utils.Extensions.makeGone
import com.lifeplus.onetapservice.utils.Extensions.makeVisible
import com.lifeplus.onetapservice.utils.Extensions.toast
import dagger.hilt.android.AndroidEntryPoint

@AndroidEntryPoint
class LoginFragment : Fragment() {
    private lateinit var binding: FragmentLoginBinding
    private val viewModel: AuthViewModel by viewModels()

    override fun onCreateView(
        inflater: LayoutInflater, container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View? {
        binding = FragmentLoginBinding.inflate(inflater,container,false)
        return binding.root
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)

        binding.btnGetCode.isEnabled = false

        binding.etPhone.setOnFocusChangeListener { _, hasFocus ->
            if (hasFocus) {
                // Convert 20dp to pixels dynamically based on the device's density
                val paddingStart = (35 * binding.etPhone.context.resources.displayMetrics.density).toInt()

                // Set the padding for the EditText
                binding.etPhone.setPadding(
                    paddingStart, // Start padding in pixels
                    binding.etPhone.paddingTop,
                    binding.etPhone.paddingRight,
                    binding.etPhone.paddingBottom
                )

                binding.tvBdCode.visibility = View.VISIBLE


                // Show the keyboard explicitly
                val inputMethodManager = binding.etPhone.context.getSystemService(
                    Context.INPUT_METHOD_SERVICE) as InputMethodManager
                inputMethodManager.showSoftInput(binding.etPhone, InputMethodManager.SHOW_IMPLICIT)
            }
        }

        binding.etPhone.addTextChangedListener(object : TextWatcher {
            override fun beforeTextChanged(charSequence: CharSequence, i: Int, i1: Int, i2: Int) {}
            override fun onTextChanged(charSequence: CharSequence, i: Int, i1: Int, i2: Int) {

                val isValid = charSequence.length == 11

                binding.btnGetCode.isEnabled = isValid

                if (isValid) {
                    binding.btnGetCode.setTextColor(
                        ContextCompat.getColor(requireContext(), R.color.white)
                    )
                } else {
                    binding.btnGetCode.setTextColor(
                        ContextCompat.getColor(requireContext(), android.R.color.darker_gray)
                    )
                }
            }

            override fun afterTextChanged(editable: Editable) {}
        })

        binding.btnGetCode.setOnClickListener {
            val phone = binding.etPhone.text.toString().trim()
            val corporateId = binding.etCorporateID.text.toString().trim()

            val bdPhoneRegex = Regex("^01[3-9]\\d{8}$")

            if (!bdPhoneRegex.matches(phone)) {
                toast("Enter valid Bangladeshi mobile number")
                return@setOnClickListener
            }

            /*if (corporateId.isEmpty()){
                toast("Enter corporate id")
                return@setOnClickListener
            }*/
            /*if (!bdPhoneRegex.matches(phone)) {
                // Show error message
                binding.phoneNumberValidationText.visibility = View.VISIBLE
                binding.phoneNumberValidationText.text = "Enter valid Bangladeshi phone number"

                // Change EditText border to red
                binding.etPhone.setBackgroundResource(R.drawable.edit_text_border_error)

                return@setOnClickListener
            } else {
                // Hide error message
                binding.phoneNumberValidationText.visibility = View.GONE

                // Reset EditText border to normal
                binding.etPhone.setBackgroundResource(R.drawable.button_stroke_background)
            }*/

            if (phone != null){
                val bundle = Bundle().apply {
                    putString(Constants.PHONE_NUMBER, phone)
                }
                findNavController().navigate(R.id.action_loginFragment_to_OTPFragment,bundle)
            }
        }

        binding.tvLoginAsCorporate.setOnClickListener {
            binding.tvLoginAsCorporate.makeGone()
            binding.tvAreYouACorporatePerson.makeGone()
            binding.btnBack.makeVisible()
            binding.etCorporateID.makeVisible()
            binding.tvTitle.text = "Enter your corporate ID & phone no."
            binding.etPhone.text?.clear()
        }

        binding.btnBack.setOnClickListener {
            binding.etCorporateID.makeGone()
            binding.btnBack.makeGone()
            binding.tvLoginAsCorporate.makeGone()
            binding.tvAreYouACorporatePerson.makeGone()
            binding.tvTitle.text = getString(R.string.enter_your_phone_no)
        }
    }
}