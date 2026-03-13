package com.labaid.insuretech.utils.prefence

import android.content.Context
import android.content.SharedPreferences
import com.google.gson.Gson

class PreferencesHelper(context: Context) {
    private val preferencesHelper: SharedPreferences
    private val sslPref = "labaid_insure_tech"


    init {
        preferencesHelper = context.getSharedPreferences(sslPref, Context.MODE_PRIVATE)
    }

    fun put(key: String, value: String?) {
        preferencesHelper.edit().putString(key, value).apply()
    }

    fun put(key: String, value: Int) {
        preferencesHelper.edit().putInt(key, value).apply()
    }

    fun put(key: String, value: Long) {
        preferencesHelper.edit().putLong(key, value).apply()
    }

    fun put(key: String, value: Float) {
        preferencesHelper.edit().putFloat(key, value).apply()
    }

    fun put(key: String, value: Boolean) {
        preferencesHelper.edit().putBoolean(key, value).apply()
    }

    operator fun get(key: String, defaultValue: String): String {
        return preferencesHelper.getString(key, defaultValue) ?: defaultValue
    }

    operator fun get(key: String, defaultValue: Int): Int {
        return preferencesHelper.getInt(key, defaultValue)
    }

    operator fun get(key: String, defaultValue: Float): Float {
        return preferencesHelper.getFloat(key, defaultValue)
    }

    operator fun get(key: String, defaultValue: Boolean): Boolean {
        return preferencesHelper.getBoolean(key, defaultValue)
    }

    operator fun get(key: String, defaultValue: Long): Long {
        return preferencesHelper.getLong(key, defaultValue)
    }

    fun <T> getResponse(key: String, clazz: Class<T>): T? {
        return try {
            val response = preferencesHelper.getString(key, "") ?: ""
            Gson().fromJson(response, clazz)
        } catch (e: Exception) {
            null
        }
    }


    fun deleteSavedData(key: String) {
        preferencesHelper.edit().remove(key).apply()
    }

    fun clearSavedData() {
        preferencesHelper.edit().clear().apply()
    }

    companion object {
        private const val PREF_NAME = "chatbot_preferences"
        private const val KEY_USER_NAME = "user_name"
        private const val KEY_USER_PHONE = "user_phone"
        private const val KEY_USER_ID = "user_id"
        private const val KEY_CHAT_ID = "chat_id"
        private const val KEY_LANGUAGE = "language"
        private const val KEY_IS_REGISTERED = "is_registered"
        private const val KEY_LAST_CHAT_TIME = "last_chat_time"
    }

    private val sharedPreferences: SharedPreferences =
        context.getSharedPreferences(PREF_NAME, Context.MODE_PRIVATE)

    fun saveUserInfo(name: String, phone: String, userId: String, chatId: String, language: String = "en") {
        sharedPreferences.edit().apply {
            putString(KEY_USER_NAME, name)
            putString(KEY_USER_PHONE, phone)
            putString(KEY_USER_ID, userId)
            putString(KEY_CHAT_ID, chatId)
            putString(KEY_LANGUAGE, language)
            putBoolean(KEY_IS_REGISTERED, true)
            putLong(KEY_LAST_CHAT_TIME, System.currentTimeMillis())
            apply()
        }
    }

    fun getUserName(): String? = sharedPreferences.getString(KEY_USER_NAME, null)

    fun getUserPhone(): String? = sharedPreferences.getString(KEY_USER_PHONE, null)

    fun getUserId(): String? = sharedPreferences.getString(KEY_USER_ID, null)

    fun getChatId(): String? = sharedPreferences.getString(KEY_CHAT_ID, null)

    fun getLanguage(): String = sharedPreferences.getString(KEY_LANGUAGE, "en") ?: "en"

    fun isUserRegistered(): Boolean = sharedPreferences.getBoolean(KEY_IS_REGISTERED, false)

    fun getLastChatTime(): Long = sharedPreferences.getLong(KEY_LAST_CHAT_TIME, 0L)

    fun setLanguage(language: String) {
        sharedPreferences.edit().putString(KEY_LANGUAGE, language).apply()
    }

    fun clearUserSession() {
        sharedPreferences.edit().clear().apply()
    }

    fun updateChatId(chatId: String) {
        sharedPreferences.edit().putString(KEY_CHAT_ID, chatId).apply()
    }

    fun updateLastChatTime() {
        sharedPreferences.edit().putLong(KEY_LAST_CHAT_TIME, System.currentTimeMillis()).apply()
    }

    fun isSessionExpired(): Boolean {
        val lastChatTime = getLastChatTime()
        val currentTime = System.currentTimeMillis()
        val sessionDuration = 24 * 60 * 60 * 1000L // 24 hours in milliseconds
        return (currentTime - lastChatTime) > sessionDuration
    }


}