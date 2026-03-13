package com.lifeplus.onetapservice.utils

import android.app.Activity
import android.content.Context
import android.content.Intent
import android.view.View
import android.view.inputmethod.InputMethodManager
import android.widget.EditText
import android.widget.Toast
import androidx.annotation.StringRes
import androidx.fragment.app.Fragment
import androidx.recyclerview.widget.RecyclerView.ViewHolder
import java.text.SimpleDateFormat
import java.util.Calendar
import java.util.Locale
import android.provider.Settings
import androidx.core.content.ContextCompat
import android.Manifest
import android.content.pm.PackageManager

object Extensions {
    fun View.gone() = run { visibility = View.GONE }
    fun View.visible() = run { visibility = View.VISIBLE }
    fun View.invisible() = run { visibility = View.INVISIBLE }

    val EditText.doubleValue get() = this.text.toString().toDouble()
    val EditText.intValue get() = this.text.toString().toInt()
    val EditText.value get() = this.text?.toString() ?: ""

    fun String.compare(s: String): Boolean {
        return this.isNotEmpty() && s.contains(this)
    }

    fun Activity.startActivity(
        cls: Class<*>,
        finishThis: Boolean = false,
        block: (Intent.() -> Unit)? = null
    ) {
        val intent = Intent(this, cls)
        block?.invoke(intent)
        startActivity(intent)
        if (finishThis) finish()
    }

    fun Fragment.startActivity(
        cls: Class<*>,
        finishThis: Boolean = false,
        block: (Intent.() -> Unit)? = null
    ) {
        val intent = Intent(requireActivity(), cls)
        block?.invoke(intent)
        startActivity(intent)
        if (finishThis) requireActivity().finish()
    }

    fun Fragment.toast(message: String) {
        Toast.makeText(requireContext(), message, Toast.LENGTH_SHORT).show()
    }

    fun Fragment.toast(@StringRes message: Int) {
        Toast.makeText(requireContext(), message, Toast.LENGTH_SHORT).show()
    }

    fun Activity.toast(message: String) {
        Toast.makeText(this, message, Toast.LENGTH_SHORT).show()
    }

    fun Activity.toast(@StringRes message: Int) {
        Toast.makeText(this, message, Toast.LENGTH_SHORT).show()
    }
    fun Fragment.hideKeyboard() {
        view?.let { activity?.hideKeyboard(it) }
    }
    fun Context.hideKeyboard(view: View) {
        val inputMethodManager = getSystemService(Activity.INPUT_METHOD_SERVICE) as InputMethodManager
        inputMethodManager.hideSoftInputFromWindow(view.windowToken, 0)
    }

    fun Activity.hideKeyboard() {
        val inputMethodManager =
            getSystemService(Context.INPUT_METHOD_SERVICE) as InputMethodManager
        currentFocus?.let {
            inputMethodManager.hideSoftInputFromWindow(it.windowToken, 0)
        }

    }
    fun ViewHolder.getString(
        resourceId: Int,
        params1: String,
        params2: String? = null
    ): String {
        return if (params2 == null) itemView.context.getString(resourceId, params1)
        else itemView.context.getString(resourceId, params1, params2)

    }

    fun View.makeGone() {
        visibility = View.GONE
    }

    fun View.makeInvisible() {
        visibility = View.INVISIBLE
    }

    fun View.makeVisible() {
        visibility = View.VISIBLE
    }

    fun Context.getDeviceUniqueId(): String {
        return Settings.Secure.getString(this.contentResolver, Settings.Secure.ANDROID_ID)
    }

    fun getCurrentDateAndDay(): Triple<String, String, String> {
        val calendar = Calendar.getInstance()

        // Get the day of the month
        val date = calendar.get(Calendar.DAY_OF_MONTH).toString()

        // Get the day of the week name
        val dayFormat = SimpleDateFormat("EEE")
        val day = dayFormat.format(calendar.time)

        // Get the current year in two-digit format
        val yearFormat = SimpleDateFormat("yy")
        val year = yearFormat.format(calendar.time)

        return Triple(date, day, year)
    }

    // Extension function for Activity to show a ProgressDialog with custom ProgressBar color


    fun getFormattedDate(scheduleDay: String, scheduleDate: String): String {
        // Get the current year and month
        val calendar = Calendar.getInstance()
        val year = calendar.get(Calendar.YEAR)
        val month = calendar.get(Calendar.MONTH) + 1 // Calendar.MONTH is zero-based

        // Construct the date string in the desired format
        val formattedDate = "$scheduleDate-${if (month < 10) "0$month" else "$month"}-$year"

        // Optional: Validate the day
        val sdf = SimpleDateFormat("dd-MM-yyyy", Locale.getDefault())
        val date = sdf.parse(formattedDate)
        val dayOfWeek = SimpleDateFormat("EEE", Locale.getDefault()).format(date)

        if (!dayOfWeek.equals(scheduleDay, ignoreCase = true)) {
            throw IllegalArgumentException("Provided scheduleDay does not match the actual day of the date.")
        }

        return formattedDate
    }


    fun isLocationPermissionGranted(context: Context) : Boolean =
        ContextCompat.checkSelfPermission(context, Manifest.permission.ACCESS_FINE_LOCATION) ==
                PackageManager.PERMISSION_GRANTED

    fun requestLocationPermission(activity: Activity) {
        activity.requestPermissions(arrayOf(Manifest.permission.ACCESS_FINE_LOCATION), 1000)
    }

}
