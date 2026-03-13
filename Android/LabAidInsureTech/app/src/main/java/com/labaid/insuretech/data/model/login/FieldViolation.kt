package com.labaid.insuretech.data.model.login


import com.google.gson.annotations.SerializedName

data class FieldViolation(
    @SerializedName("field")
    val `field`: String?,
    @SerializedName("code")
    val code: String?,
    @SerializedName("description")
    val description: String?,
    @SerializedName("rejected_value")
    val rejectedValue: String?
)