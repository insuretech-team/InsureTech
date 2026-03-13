package com.labaid.insuretech.ui.home

import android.os.Bundle
import android.os.Handler
import android.os.Looper
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import androidx.core.content.ContextCompat
import androidx.fragment.app.Fragment
import androidx.navigation.fragment.findNavController
import androidx.recyclerview.widget.GridLayoutManager
import androidx.viewpager2.widget.ViewPager2.OnPageChangeCallback
import com.labaid.insuretech.R
import com.labaid.insuretech.adapter.FeaturedPoliciesAdapter
import com.labaid.insuretech.adapter.HomeSliderAdapter
import com.labaid.insuretech.adapter.InsuranceCardsAdapter
import com.labaid.insuretech.adapter.LatestUpdatesAdapter
import com.labaid.insuretech.data.model.featured_policy.FeaturedPolicy
import com.labaid.insuretech.data.model.insurance_card.InsuranceCard
import com.labaid.insuretech.data.model.latest_update.LatestUpdate
import com.labaid.insuretech.databinding.FragmentHomeBinding
import com.zhpan.indicator.enums.IndicatorSlideMode
import com.zhpan.indicator.enums.IndicatorStyle
import dagger.hilt.android.AndroidEntryPoint

@AndroidEntryPoint
class HomeFragment : Fragment() {

    private lateinit var binding: FragmentHomeBinding
    private var totalBannerPage = 0
    private var totalFeaturedPage = 0
    private var totalUpdatesPage = 0

    private val handler: Handler by lazy { Handler(Looper.getMainLooper()) }
    private val bannerRunnable = Runnable {
        autoScrollBanner()
    }
    private val featuredRunnable = Runnable {
        autoScrollFeatured()
    }
    private val updatesRunnable = Runnable {
        autoScrollUpdates()
    }

    override fun onCreateView(
        inflater: LayoutInflater, container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View {
        binding = FragmentHomeBinding.inflate(inflater, container, false)
        return binding.root
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)

        setupClickListeners()
        setupInsuranceCards()
        setupBannerSlider()
        setupFeaturedPolicies()
        setupLatestUpdates()
    }

    private fun setupClickListeners() {
        binding.ivAvatar.setOnClickListener {
            findNavController().navigate(R.id.action_homeFragment_to_profileFragment)
        }

        binding.tvUserName.setOnClickListener {
            findNavController().navigate(R.id.action_homeFragment_to_profileFragment)
        }
    }

    private fun setupInsuranceCards() {
        val insuranceCards = listOf(
            InsuranceCard(
                id = 1,
                name = "Health",
                iconResId = R.drawable.ic_health
            ),
            InsuranceCard(
                id = 2,
                name = "Auto",
                iconResId = R.drawable.ic_auto
            ),
            InsuranceCard(
                id = 3,
                name = "Travel",
                iconResId = R.drawable.ic_travel
            ),
            InsuranceCard(
                id = 4,
                name = "Life",
                iconResId = R.drawable.ic_life
            )
        )

        val insuranceAdapter = InsuranceCardsAdapter(insuranceCards) { card ->
            // Handle card click
            when (card.id) {
                1 -> findNavController().navigate(R.id.action_homeFragment_to_planFragment)
                // Add other navigation cases as needed
            }
        }

        binding.rvInsuranceCards.apply {
            layoutManager = GridLayoutManager(requireContext(), 2)
            adapter = insuranceAdapter
            setHasFixedSize(true)
        }
    }

    private fun setupBannerSlider() {
        val bannerImages = listOf(
            R.drawable.labaid_insuretech_cover,
            R.drawable.labaid_insuretech_cover,
            R.drawable.labaid_insuretech_cover
        )

        val bannerAdapter = HomeSliderAdapter(bannerImages)
        binding.homeSliderVp.adapter = bannerAdapter
        totalBannerPage = bannerImages.size

        binding.indicator.apply {
            setPageSize(bannerImages.size)
            setSliderColor(
                ContextCompat.getColor(requireActivity(), R.color.slider_unselected_color),
                ContextCompat.getColor(requireActivity(), R.color.band_color)
            )
            setSliderWidth(30f)
            setSlideMode(IndicatorSlideMode.WORM)
            setIndicatorStyle(IndicatorStyle.CIRCLE)
            setupWithViewPager(binding.homeSliderVp)
        }

        binding.homeSliderVp.registerOnPageChangeCallback(object : OnPageChangeCallback() {
            override fun onPageSelected(position: Int) {
                handler.removeCallbacks(bannerRunnable)
                handler.postDelayed(bannerRunnable, 3000)
            }
        })

        if (totalBannerPage > 1) {
            handler.postDelayed(bannerRunnable, 3000)
        }
    }

    private fun setupFeaturedPolicies() {
        val policies = listOf(
            FeaturedPolicy(
                name = "Seba",
                imageResId = R.drawable.plan_image_placeholder,
                coverage = listOf(
                    "Emergency medical expenses",
                    "Adventure/Sports coverage"
                ),
                poweredByLogoResId = R.drawable.ic_charmed_life_logo
            ),
            FeaturedPolicy(
                name = "Travel Plus",
                imageResId = R.drawable.plan_image_placeholder,
                coverage = listOf(
                    "Trip cancellation",
                    "Lost baggage protection"
                ),
                poweredByLogoResId = R.drawable.ic_charmed_life_logo
            )
        )

        val featuredAdapter = FeaturedPoliciesAdapter(policies)
        binding.featuredPoliciesVp.adapter = featuredAdapter
        totalFeaturedPage = policies.size

        binding.featuredIndicator.apply {
            setPageSize(policies.size)
            setSliderColor(
                ContextCompat.getColor(requireActivity(), R.color.slider_unselected_color),
                ContextCompat.getColor(requireActivity(), R.color.band_color)
            )
            setSliderWidth(30f)
            setSlideMode(IndicatorSlideMode.WORM)
            setIndicatorStyle(IndicatorStyle.CIRCLE)
            setupWithViewPager(binding.featuredPoliciesVp)
        }

        binding.featuredPoliciesVp.registerOnPageChangeCallback(object : OnPageChangeCallback() {
            override fun onPageSelected(position: Int) {
                handler.removeCallbacks(featuredRunnable)
                handler.postDelayed(featuredRunnable, 4000)
            }
        })

        if (totalFeaturedPage > 1) {
            handler.postDelayed(featuredRunnable, 4000)
        }
    }

    private fun setupLatestUpdates() {
        val updates = listOf(
            LatestUpdate(
                title = "Don't Let One Accident Stop Your Journey",
                description = "From medical support to quick claims and easy access through phone, we help you get back on the road...",
                imageResId = R.drawable.plan_image_placeholder
            ),
            LatestUpdate(
                title = "New Health Insurance Benefits",
                description = "Discover our enhanced health coverage with additional benefits for your family's wellbeing...",
                imageResId = R.drawable.plan_image_placeholder
            )
        )

        val updatesAdapter = LatestUpdatesAdapter(updates)
        binding.latestUpdatesVp.adapter = updatesAdapter
        totalUpdatesPage = updates.size

        binding.latestIndicator.apply {
            setPageSize(updates.size)
            setSliderColor(
                ContextCompat.getColor(requireActivity(), R.color.slider_unselected_color),
                ContextCompat.getColor(requireActivity(), R.color.band_color)
            )
            setSliderWidth(30f)
            setSlideMode(IndicatorSlideMode.WORM)
            setIndicatorStyle(IndicatorStyle.CIRCLE)
            setupWithViewPager(binding.latestUpdatesVp)
        }

        binding.latestUpdatesVp.registerOnPageChangeCallback(object : OnPageChangeCallback() {
            override fun onPageSelected(position: Int) {
                handler.removeCallbacks(updatesRunnable)
                handler.postDelayed(updatesRunnable, 4000)
            }
        })

        if (totalUpdatesPage > 1) {
            handler.postDelayed(updatesRunnable, 4000)
        }
    }

    private fun autoScrollBanner() {
        var currentPage = binding.homeSliderVp.currentItem
        if (currentPage == totalBannerPage - 1) {
            currentPage = 0
        } else {
            currentPage++
        }
        binding.homeSliderVp.currentItem = currentPage
    }

    private fun autoScrollFeatured() {
        var currentPage = binding.featuredPoliciesVp.currentItem
        if (currentPage == totalFeaturedPage - 1) {
            currentPage = 0
        } else {
            currentPage++
        }
        binding.featuredPoliciesVp.currentItem = currentPage
    }

    private fun autoScrollUpdates() {
        var currentPage = binding.latestUpdatesVp.currentItem
        if (currentPage == totalUpdatesPage - 1) {
            currentPage = 0
        } else {
            currentPage++
        }
        binding.latestUpdatesVp.currentItem = currentPage
    }

    override fun onDestroyView() {
        super.onDestroyView()
        handler.removeCallbacks(bannerRunnable)
        handler.removeCallbacks(featuredRunnable)
        handler.removeCallbacks(updatesRunnable)
    }
}