package com.labaid.insuretech.ui.my_beneficiary

import android.os.Bundle
import androidx.fragment.app.Fragment
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import com.labaid.insuretech.R
import com.labaid.insuretech.databinding.FragmentMyBeneficiaryBinding
import dagger.hilt.android.AndroidEntryPoint

@AndroidEntryPoint
class MyBeneficiaryFragment : Fragment() {

    private lateinit var binding: FragmentMyBeneficiaryBinding

    override fun onCreateView(
        inflater: LayoutInflater, container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View? {
        binding = FragmentMyBeneficiaryBinding.inflate(inflater, container, false)
        return binding.root
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)


    }

}