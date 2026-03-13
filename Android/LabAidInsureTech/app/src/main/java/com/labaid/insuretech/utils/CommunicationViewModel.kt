package com.lifeplus.onetapservice.utils

import androidx.lifecycle.MutableLiveData
import androidx.lifecycle.ViewModel

class CommunicationViewModel:ViewModel() {
    val hideBottomNavBar= MutableLiveData<Boolean>()
    val hideToolbar= MutableLiveData<Boolean>()
    val uncheckBottomNavItem = MutableLiveData<Boolean>()
    val isSwipeRefreshVisible = MutableLiveData<Boolean>()  // to control swipe refresh visibility
    val internetConnectionListen = MutableLiveData<Boolean>()  // to check if internet is reconnected with Connectivity livedata
}