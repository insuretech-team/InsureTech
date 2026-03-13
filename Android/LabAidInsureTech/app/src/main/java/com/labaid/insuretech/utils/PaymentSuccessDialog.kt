package com.labaid.insuretech.utils

import android.app.Dialog
import android.graphics.Color
import android.graphics.drawable.ColorDrawable
import android.os.Bundle
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.view.Window
import androidx.fragment.app.DialogFragment
import com.labaid.insuretech.databinding.DialogPaymentSuccessBinding

class PaymentSuccessDialog : DialogFragment() {

    private lateinit var binding: DialogPaymentSuccessBinding
    private var onContinueClicked: (() -> Unit)? = null

    override fun onCreateDialog(savedInstanceState: Bundle?): Dialog {
        val dialog = super.onCreateDialog(savedInstanceState)
        dialog.window?.requestFeature(Window.FEATURE_NO_TITLE)
        return dialog
    }

    override fun onCreateView(
        inflater: LayoutInflater,
        container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View {
        binding = DialogPaymentSuccessBinding.inflate(inflater, container, false)
        return binding.root
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)

        // Make dialog not cancelable (user must click Continue)
        isCancelable = false

        // Setup continue button
        binding.btnContinue.setOnClickListener {
            onContinueClicked?.invoke()
            dismiss()
        }

        // Start animation
        binding.animationView.playAnimation()
    }

    override fun onStart() {
        super.onStart()
        dialog?.window?.apply {
            // Set the width to match parent with proper margins
            val width = resources.displayMetrics.widthPixels

            setLayout(
                width,
                ViewGroup.LayoutParams.WRAP_CONTENT
            )

            // IMPORTANT: Set transparent background to show CardView margins
            setBackgroundDrawable(ColorDrawable(Color.TRANSPARENT))
        }
    }

    fun setOnContinueClickListener(listener: () -> Unit) {
        onContinueClicked = listener
    }

    companion object {
        const val TAG = "PaymentSuccessDialog"

        fun newInstance(): PaymentSuccessDialog {
            return PaymentSuccessDialog()
        }
    }
}