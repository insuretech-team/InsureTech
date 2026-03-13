package com.labaid.insuretech.ui.otp

import android.graphics.Color
import android.os.Bundle
import android.text.Editable
import android.text.TextWatcher
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.Toast
import androidx.core.content.ContextCompat
import androidx.fragment.app.Fragment
import androidx.navigation.fragment.findNavController
import com.labaid.insuretech.R
import com.labaid.insuretech.databinding.FragmentOTPBinding
import com.labaid.insuretech.utils.Constants
import com.lifeplus.onetapservice.utils.Extensions.value
import dagger.hilt.android.AndroidEntryPoint
import ir.samanjafari.easycountdowntimer.CountDownInterface

@AndroidEntryPoint
class OTPFragment : Fragment() {

    private lateinit var binding: FragmentOTPBinding
    private var phoneNumber: String? = null

    override fun onCreateView(
        inflater: LayoutInflater, container: ViewGroup?, savedInstanceState: Bundle?
    ): View? {
        binding = FragmentOTPBinding.inflate(layoutInflater, container, false)
        return binding.root
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)

        phoneNumber = arguments?.getString(Constants.PHONE_NUMBER)

        binding.tvSubtitle.text = "A 4 digit code has been sent to your no. +88 ${phoneNumber}"

        binding.btnVerify.isEnabled = false


        binding.otpView.addTextChangedListener(object : TextWatcher {
            override fun beforeTextChanged(charSequence: CharSequence, i: Int, i1: Int, i2: Int) {}
            override fun onTextChanged(charSequence: CharSequence, i: Int, i1: Int, i2: Int) {
                val isValid = charSequence.length == 4

                binding.btnVerify.isEnabled = isValid

                if (isValid) {
                    binding.btnVerify.setTextColor(
                        ContextCompat.getColor(requireContext(), R.color.white)
                    )
                } else {
                    binding.btnVerify.setTextColor(
                        ContextCompat.getColor(requireContext(), android.R.color.darker_gray)
                    )
                }
            }

            override fun afterTextChanged(editable: Editable) {}
        })


        //Back Button
        binding.btnBack.setOnClickListener {
            findNavController().popBackStack()
        }

        binding.btnVerify.setOnClickListener {
            val otp = binding.otpView.value.trim()
            if (otp == "1111") {
                // Success - keep default color or set to green
                binding.otpView.setLineColor(resources.getColor(R.color.otp_box_color))
                // Or set to green for success
                // binding.otpView.setLineColor(Color.GREEN)

                Toast.makeText(requireContext(), "OTP Verified Successfully", Toast.LENGTH_SHORT).show()
                // Your success logic here
                findNavController().navigate(R.id.action_OTPFragment_to_homeFragment)

            } else {
                // Error - set line color to red
                binding.otpView.setLineColor(Color.RED)
                Toast.makeText(requireContext(), "Invalid OTP", Toast.LENGTH_SHORT).show()
            }
        }

        binding.easyCountDownTextview.setShowDays(false)
        binding.easyCountDownTextview.setTime(0, 0, 1, 0)
        binding.easyCountDownTextview.startTimer()
        binding.easyCountDownTextview.setOnTick(object : CountDownInterface {
            override fun onTick(time: Long) {
            }

            override fun onFinish() {
                binding.easyCountDownTextview.setTime(0, 0, 0, 0)
                binding.easyCountDownTextview.stopTimer()
                binding.resend.visibility = View.VISIBLE
                binding.didNotReceiveCode.visibility = View.VISIBLE
            }
        })

        //Resend Button
        binding.resend.setOnClickListener {
            binding.easyCountDownTextview.setTime(0, 0, 5, 0)
            binding.easyCountDownTextview.startTimer()
        }

    }

}