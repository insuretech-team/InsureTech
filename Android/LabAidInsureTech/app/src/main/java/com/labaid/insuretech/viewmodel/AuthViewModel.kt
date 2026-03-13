package com.labaid.insuretech.viewmodel

import androidx.lifecycle.LiveData
import androidx.lifecycle.MutableLiveData
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.labaid.insuretech.repo.AuthRepository
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.launch
import javax.inject.Inject

@HiltViewModel
class AuthViewModel @Inject constructor(
    private val repository: AuthRepository
) : ViewModel(){

    /*private val _loginUser = MutableLiveData<Resource<Data>>()

    val loginUser: LiveData<Resource<Data>> = _loginUser

    fun userLogin(phoneNumber: LoginBody) {
        _loginUser.postValue(Resource.loading())
        viewModelScope.launch(Dispatchers.IO) {
            repository.loginUser(phoneNumber).let {
                _loginUser.postValue(it)
            }
        }
    }

    private val _logoutUser = MutableLiveData<Resource<LogoutResponse>>()

    val logoutUser: LiveData<Resource<LogoutResponse>> = _logoutUser

    fun userLogout() {
        _logoutUser.postValue(Resource.loading())
        viewModelScope.launch(Dispatchers.IO) {
            repository.logoutUser().let {
                _logoutUser.postValue(it)
            }
        }
    }


    private val _verifyOTP = MutableLiveData<Resource<OTPResponse>>()
    val verifyOTP: LiveData<Resource<OTPResponse>> = _verifyOTP

    fun userVerifyOTP(otpBody: OTPBody){
        _verifyOTP.postValue(Resource.loading())
        viewModelScope.launch(Dispatchers.IO) {
            repository.verifyOTP(otpBody).let {
                _verifyOTP.postValue(it)
            }
        }
    }

    fun saveToken(token: String?) = repository.saveToken(token)*/



    fun checkPhone(phone: String): String? {
        if (phone.isEmpty()) return "Please enter a valid phone number"
        if (phone.length != 11) return "Phone must be 11 digit"
        if (phone[0] != '0') return "Phone number must start with 0"
        return null
    }
}