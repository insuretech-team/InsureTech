package com.labaid.insuretech.data.remote

import com.labaid.insuretech.data.model.login.LoginReq
import com.labaid.insuretech.data.model.login.LoginResponse
import com.labaid.insuretech.network.RetrofitResponse
import retrofit2.http.Body
import retrofit2.http.POST

interface ApiService {

    @POST("/auth/login")
    suspend fun loginUser(@Body phoneNumber: LoginReq): RetrofitResponse<LoginResponse>
}
