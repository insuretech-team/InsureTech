package com.labaid.insuretech.data.model.login


import com.google.gson.annotations.SerializedName

data class Error(
    @SerializedName("code")
    val code: String?,
    @SerializedName("message")
    val message: String?,
    @SerializedName("details")
    val details: Details?,
    @SerializedName("field_violations")
    val fieldViolations: List<FieldViolation?>?,
    @SerializedName("retryable")
    val retryable: Boolean?,
    @SerializedName("retry_after_seconds")
    val retryAfterSeconds: Int?,
    @SerializedName("http_status_code")
    val httpStatusCode: Int?,
    @SerializedName("error_id")
    val errorId: String?,
    @SerializedName("documentation_url")
    val documentationUrl: String?
)