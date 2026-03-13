package com.labaid.insuretech.data.model.login


import com.google.gson.annotations.SerializedName

data class LoginResponse(
    @SerializedName("user_id")
    val userId: String?,
    @SerializedName("access_token")
    val accessToken: String?,
    @SerializedName("refresh_token")
    val refreshToken: String?,
    @SerializedName("access_token_expires_in")
    val accessTokenExpiresIn: Int?,
    @SerializedName("refresh_token_expires_in")
    val refreshTokenExpiresIn: Int?,
    @SerializedName("user")
    val user: User?,
    @SerializedName("error")
    val error: Error?
)