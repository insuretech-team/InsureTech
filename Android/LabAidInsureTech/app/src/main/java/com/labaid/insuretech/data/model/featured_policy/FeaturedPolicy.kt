package com.labaid.insuretech.data.model.featured_policy

data class FeaturedPolicy(
    val name: String,
    val imageResId: Int,
    val coverage: List<String>,
    val poweredByLogoResId: Int
)