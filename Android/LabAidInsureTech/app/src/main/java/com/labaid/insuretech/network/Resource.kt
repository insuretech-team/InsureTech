package com.labaid.insuretech.network

enum class Status {
    SUCCESS,
    ERROR,
    LOADING
}

data class Resource<out T>(
    val status: Status,
    val data: T?,
    val message: String?,
    val extraValue: Int = -1
) {
    companion object {
        fun <T> success(data: T, extraValue: Int = -1): Resource<T> {
            return Resource(Status.SUCCESS, data, null, extraValue)
        }

        fun <T> error(msg: String, extraValue: Int = -1): Resource<T> {
            return Resource(Status.ERROR, null, msg, extraValue)
        }

        fun <T> loading(extraValue: Int = -1): Resource<T> {
            return Resource(Status.LOADING, null, null)
        }
    }
}