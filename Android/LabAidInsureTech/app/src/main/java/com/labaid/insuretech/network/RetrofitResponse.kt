package com.labaid.insuretech.network

import com.google.gson.annotations.Expose
import com.google.gson.annotations.SerializedName

data class RetrofitResponse<T>(
    @SerializedName("success")
    @Expose
    var success: Boolean = false,

    @SerializedName("statusCode")
    @Expose
    val statusCode: Int,

    @SerializedName("message")
    @Expose
    var statusMessage: String = "Unexpected error, please try again",

    @SerializedName("data")
    @Expose
    var data: T? = null,

    @SerializedName("errors")
    @Expose
    var errors: Errors
)

data class Errors(
    @SerializedName("mobile")
    @Expose
    var mobile: List<String>
)


