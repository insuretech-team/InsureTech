package com.labaid.insuretech.network

import android.content.Context
import com.facebook.stetho.BuildConfig
import com.facebook.stetho.okhttp3.StethoInterceptor
import com.google.gson.GsonBuilder
import com.labaid.insuretech.data.remote.ApiService
import com.labaid.insuretech.utils.Constants
import com.labaid.insuretech.utils.prefence.PreferencesHelper
import dagger.Module
import dagger.Provides
import dagger.hilt.InstallIn
import dagger.hilt.android.qualifiers.ApplicationContext
import dagger.hilt.components.SingletonComponent
import okhttp3.Interceptor
import okhttp3.OkHttpClient
import okhttp3.Request
import okhttp3.logging.HttpLoggingInterceptor
import retrofit2.Retrofit
import retrofit2.converter.gson.GsonConverterFactory
import java.util.concurrent.TimeUnit
import javax.inject.Qualifier
import javax.inject.Singleton


@InstallIn(SingletonComponent::class)
@Module
class NetworkModule {

    // Qualifiers to distinguish between different Retrofit instances
    @Qualifier
    @Retention(AnnotationRetention.BINARY)
    annotation class MainApi

    @Qualifier
    @Retention(AnnotationRetention.BINARY)
    annotation class ChatBotApi

    @Singleton
    @Provides
    fun providesOkHttpClient(
        httpLoggingInterceptor: HttpLoggingInterceptor,
        preferencesHelper: PreferencesHelper
    ): OkHttpClient {
        val interceptor = Interceptor { chain ->
            val request: Request = chain.request().newBuilder()
                .addHeader("Authorization", "Bearer " + preferencesHelper[Constants.TOKEN, ""])
                .addHeader("Accept", "application/json")
                .addHeader("Content-Type", "application/json").build()
            chain.proceed(request)
        }
        return OkHttpClient.Builder()
            .connectTimeout(2, TimeUnit.MINUTES)
            .readTimeout(2, TimeUnit.MINUTES)
            .addInterceptor(interceptor)
            .addNetworkInterceptor(HttpLoggingInterceptor().apply {
                level = HttpLoggingInterceptor.Level.BODY
            })
            .addNetworkInterceptor(StethoInterceptor())
            .build()
    }

    @Provides
    @Singleton
    fun provideHttpLoggingInterceptor(): HttpLoggingInterceptor {
        val interceptor = HttpLoggingInterceptor()
        if (BuildConfig.DEBUG)
            interceptor.level = HttpLoggingInterceptor.Level.BODY
        else
            interceptor.level = HttpLoggingInterceptor.Level.NONE
        return interceptor
    }

    // Main API Retrofit instance (for BASE_URL)
    @MainApi
    @Singleton
    @Provides
    fun providesMainRetrofit(
        client: OkHttpClient
    ): Retrofit {
        val gson = GsonBuilder().serializeNulls().create()
        return Retrofit.Builder()
            .baseUrl(Constants.BASE_URL)
            .client(client)
            .addConverterFactory(GsonConverterFactory.create(gson))
            .build()
    }


    // Main API Service (uses BASE_URL)
    @MainApi
    @Singleton
    @Provides
    fun providesMainApiService(@MainApi retrofit: Retrofit): ApiService {
        return retrofit.create(ApiService::class.java)
    }

    @Provides
    @Singleton
    fun providePreference(@ApplicationContext context: Context): PreferencesHelper {
        return PreferencesHelper(context)
    }
}