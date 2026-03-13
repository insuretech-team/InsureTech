package com.lifeplus.onetapservice.utils

import android.content.Context
import android.net.ConnectivityManager
import android.net.NetworkCapabilities
import android.os.Build
import dagger.hilt.android.qualifiers.ApplicationContext
import javax.inject.Inject

class NetworkUtil @Inject constructor(
    @ApplicationContext val appContext: Context
) {
    fun isConnectedToInternet(): Boolean {
        val connectivityManager =
            appContext.getSystemService(Context.CONNECTIVITY_SERVICE) as ConnectivityManager

        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.M) {
            connectivityManager.getNetworkCapabilities(connectivityManager.activeNetwork)
                ?.let { networkCapabilities ->
                    return networkCapabilities.hasTransport(NetworkCapabilities.TRANSPORT_CELLULAR) ||
                            networkCapabilities.hasTransport(NetworkCapabilities.TRANSPORT_WIFI)
                }
        } else {
            return connectivityManager.activeNetworkInfo?.isConnectedOrConnecting ?: false
        }
        return false
    }

    fun isInternetAvailable(): Boolean {
        if (isConnectedToInternet()) {
            val command = "ping -c 1 google.com"
            try {
                return Runtime.getRuntime().exec(command).waitFor() == 0
            } catch (e: Exception) {
                e.printStackTrace()
            }
        }
        return false
    }
}