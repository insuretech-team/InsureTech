package com.labaid.insuretech.repo

import com.labaid.insuretech.data.model.login.LoginReq
import com.labaid.insuretech.data.remote.ApiService
import com.labaid.insuretech.network.NetworkModule
import com.labaid.insuretech.network.Resource
import com.labaid.insuretech.utils.NetworkCallParse
import com.labaid.insuretech.utils.prefence.PreferencesHelper
import javax.inject.Inject

class AuthRepository @Inject constructor(
    @NetworkModule.MainApi private val api: ApiService,
    private val parse: NetworkCallParse,
    private  val preferencesHelper: PreferencesHelper
) {

    /*suspend fun loginUser(phoneNumber: LoginReq): Resource<Data> {
        return try {
            val response = api.loginUser(phoneNumber)
            parse.responseParse(response)
        } catch (e: Exception) {
            parse.exceptionParse(e.message)
        }
    }

    suspend fun logoutUser(): Resource<LogoutResponse> {
        return try {
            val response = api.logoutUser()
            parse.responseParse(response)
        } catch (e: Exception) {
            parse.exceptionParse(e.message)
        }
    }

    suspend fun verifyOTP(otpBody: OTPBody): Resource<OTPResponse> {
        return try {
            val response = api.verifyOTP(otpBody)
            parse.responseParse(response)
        } catch (e: Exception) {
            parse.exceptionParse(e.message)
        }
    }

    fun saveToken(token: String?){
        preferencesHelper.put(Constants.TOKEN,token)
        if (token?.isNotEmpty() == true) {
            preferencesHelper.put(Constants.IS_LOGIN, true)
        }
    }*/

}