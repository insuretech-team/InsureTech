package com.labaid.insuretech.data.model.login

data class LoginReq(
    val mobile_number: String?,
    val password: String?,
    val device_id: String?,
    val device_type: String?
)
