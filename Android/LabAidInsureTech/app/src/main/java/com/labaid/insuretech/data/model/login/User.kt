package com.labaid.insuretech.data.model.login


import com.google.gson.annotations.SerializedName

data class User(
    @SerializedName("user_id")
    val userId: String?,
    @SerializedName("mobile_number")
    val mobileNumber: String?,
    @SerializedName("email")
    val email: String?,
    @SerializedName("status")
    val status: String?,
    @SerializedName("created_at")
    val createdAt: String?,
    @SerializedName("updated_at")
    val updatedAt: String?,
    @SerializedName("last_login_at")
    val lastLoginAt: String?,
    @SerializedName("created_by")
    val createdBy: String?,
    @SerializedName("updated_by")
    val updatedBy: String?,
    @SerializedName("login_attempts")
    val loginAttempts: Int?,
    @SerializedName("locked_until")
    val lockedUntil: String?,
    @SerializedName("deleted_at")
    val deletedAt: String?
)