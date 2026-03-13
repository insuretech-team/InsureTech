package com.labaid.insuretech.data.model.plan

data class Plan(
    val id: Int,
    val name: String,
    val imageUrl: String,
    val coverageItems: List<String>,
    val poweredBy: String = "Charmed Life Insurance"
)
