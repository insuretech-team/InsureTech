package com.labaid.insuretech.utils


import com.labaid.insuretech.network.Resource
import com.labaid.insuretech.network.RetrofitResponse
import com.lifeplus.onetapservice.utils.NetworkUtil
import javax.inject.Inject

class NetworkCallParse @Inject constructor(
    private val networkUtil: NetworkUtil
) {

    /*fun <T> responseParse(response: RetrofitResponse<T>): Resource<T> {
        return if (response.statusCode == Constants.STATUS_OK && response.data != null) {
            Resource.success(response.data!!)
        } else {
            Resource.error(response.statusMessage)
        }
    }*/

    fun <T> responseParse(response: RetrofitResponse<T>): Resource<T> {
        return if (response.success && response.data != null) {
            Resource.success(response.data!!)
        } else {
            Resource.error(response.statusMessage)
        }
    }

    /* fun <T> responseParse(response: RetrofitResponse<T>): ResponseData {
         return if (response.success && response.data != null) {
             ResponseData.Success(response.data!!)
         } else {
            ResponseData.Error(response.statusMessage)
         }
     }
     fun exceptionParse(msg: String?): ResponseData {
         msg?.let {
             if (msg.contains(Constants.UNABLE_TO_RESOLVE_HOST, true)
                 && !networkUtil.isInternetAvailable()
             ) {
                 return ResponseData.Error(Constants.NO_INTERNET)
             }
         }
         return ResponseData.Error("Unexpected error, please try again")
     }*/

    fun <T> exceptionParse(msg: String?): Resource<T> {
        msg?.let {
            if (msg.contains(Constants.UNABLE_TO_RESOLVE_HOST, true)
                && !networkUtil.isInternetAvailable()
            ) {
                return Resource.error(Constants.NO_INTERNET)
            } else return Resource.error(msg)
        }
        return Resource.error("Unexpected error, please try again")
    }
}