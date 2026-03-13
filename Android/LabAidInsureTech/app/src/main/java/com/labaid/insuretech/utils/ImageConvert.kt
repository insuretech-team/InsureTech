package com.labaid.insuretech.utils

import android.graphics.Bitmap
import android.graphics.BitmapFactory
import android.util.Base64
import java.io.ByteArrayOutputStream
import java.io.IOException
import java.io.InputStream
import java.net.HttpURLConnection
import java.net.URL

class ImageConvert private constructor() {

    companion object {
        @Volatile
        private var instance: ImageConvert? = null

        fun getInstance(): ImageConvert {
            return instance ?: synchronized(this) {
                instance ?: ImageConvert().also { instance = it }
            }
        }
    }

    fun base64ToImage(pic: String?): Bitmap? {
        return pic?.let {
            val decodeString = Base64.decode(it, Base64.DEFAULT)
            BitmapFactory.decodeByteArray(decodeString, 0, decodeString.size)
        }
    }

    fun getEncoded64ImageStringFromBitmap(bitmap: Bitmap?): String? {
        return bitmap?.let {
            val stream = ByteArrayOutputStream()
            it.compress(Bitmap.CompressFormat.JPEG, 70, stream)
            val byteFormat = stream.toByteArray()
            Base64.encodeToString(byteFormat, Base64.NO_WRAP)
        }
    }

    fun getBitmapFromURL(strURL: String): Bitmap? {
        return try {
            val url = URL(strURL)
            val connection = url.openConnection() as HttpURLConnection
            connection.doInput = true
            connection.connect()
            val input: InputStream = connection.inputStream
            BitmapFactory.decodeStream(input)
        } catch (e: IOException) {
            e.printStackTrace()
            null
        }
    }
}