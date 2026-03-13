package com.labaid.insuretech.ui.my_info

import android.os.Bundle
import androidx.fragment.app.Fragment
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import androidx.navigation.fragment.findNavController
import com.labaid.insuretech.R
import com.labaid.insuretech.databinding.FragmentMyInfoBinding
import dagger.hilt.android.AndroidEntryPoint

@AndroidEntryPoint
class MyInfoFragment : Fragment() {
    private lateinit var binding: FragmentMyInfoBinding
    override fun onCreateView(
        inflater: LayoutInflater, container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View? {
        binding = FragmentMyInfoBinding.inflate(inflater,container,false)
        return binding.root
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)

        binding.btnBack.setOnClickListener {
            findNavController().popBackStack()
        }
    }
}