package com.labaid.insuretech.data.model.login


import com.google.gson.annotations.SerializedName

data class Details(
    @SerializedName("property1")
    val property1: String?,
    @SerializedName("property2")
    val property2: String?
)