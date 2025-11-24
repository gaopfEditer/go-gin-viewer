#pragma once

#include <iostream>
#include <string>
#include <vector>
#include <memory>
#include "..\..\..\public_include\openssl\evp.h"
#include "..\..\..\public_include\openssl\rand.h"
#include "..\..\..\public_include\openssl\err.h"

unsigned char key[32] = {
        0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77,
        0x88, 0x99, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF,
        0x10, 0x20, 0x30, 0x40, 0x50, 0x60, 0x70, 0x80,
        0x90, 0xA0, 0xB0, 0xC0, 0xD0, 0xE0, 0xF0, 0x00
};  // AES-256 key

unsigned char iv[16] = {
        0x01, 0x23, 0x45, 0x67, 0x89, 0xAB, 0xCD, 0xEF,
        0xFE, 0xDC, 0xBA, 0x98, 0x76, 0x54, 0x32, 0x10
};   // Initialization Vector
       
void handleErrors(const char* context) {
    char errbuf[256];
    ERR_error_string_n(ERR_get_error(), errbuf, sizeof(errbuf));
    throw std::runtime_error(std::string(context) + ": " + errbuf);
}
/**
     * 加密数据
     * @param plaintext 明文数据
     * @param plaintext_len 明文长度
     * @param ciphertext_len 输出参数，返回密文长度
     * @return 密文数据指针，需要调用者释放 (delete[])
     */
unsigned char* encrypt(const unsigned char* plaintext, int plaintext_len, int* ciphertext_len) {
    EVP_CIPHER_CTX* ctx = EVP_CIPHER_CTX_new();
    if (!ctx) handleErrors("Failed to create cipher context");

    if (1 != EVP_EncryptInit_ex(ctx, EVP_aes_256_cbc(), NULL, key, iv)) {
        EVP_CIPHER_CTX_free(ctx);
        handleErrors("Encrypt init failed");
    }

    // 分配输出缓冲区（明文长度 + 一个块大小）
    int max_ciphertext_len = plaintext_len + EVP_MAX_BLOCK_LENGTH;
    unsigned char* ciphertext = new unsigned char[max_ciphertext_len];

    int len1 = 0, len2 = 0;

    try {
        if (1 != EVP_EncryptUpdate(ctx, ciphertext, &len1, plaintext, plaintext_len)) {
            delete[] ciphertext;
            EVP_CIPHER_CTX_free(ctx);
            handleErrors("Encrypt update failed");
        }

        if (1 != EVP_EncryptFinal_ex(ctx, ciphertext + len1, &len2)) {
            delete[] ciphertext;
            EVP_CIPHER_CTX_free(ctx);
            handleErrors("Encrypt final failed");
        }

        EVP_CIPHER_CTX_free(ctx);
        *ciphertext_len = len1 + len2;
        return ciphertext;

    }
    catch (...) {
        delete[] ciphertext;
        EVP_CIPHER_CTX_free(ctx);
        throw;
    }
}

/**
 * 解密数据
 * @param ciphertext 密文数据
 * @param ciphertext_len 密文长度
 * @param plaintext_len 输出参数，返回明文长度
 * @return 明文数据指针，需要调用者释放 (delete[])
 */
unsigned char* decrypt(const unsigned char* ciphertext, int ciphertext_len, int* plaintext_len) {
    EVP_CIPHER_CTX* ctx = EVP_CIPHER_CTX_new();
    if (!ctx) handleErrors("Failed to create cipher context");

    if (1 != EVP_DecryptInit_ex(ctx, EVP_aes_256_cbc(), NULL, key, iv)) {
        EVP_CIPHER_CTX_free(ctx);
        handleErrors("Decrypt init failed");
    }

    // 分配输出缓冲区
    int max_plaintext_len = ciphertext_len + EVP_MAX_BLOCK_LENGTH;
    unsigned char* plaintext = new unsigned char[max_plaintext_len];

    int len1 = 0, len2 = 0;

    try {
        if (1 != EVP_DecryptUpdate(ctx, plaintext, &len1, ciphertext, ciphertext_len)) {
            delete[] plaintext;
            EVP_CIPHER_CTX_free(ctx);
            handleErrors("Decrypt update failed");
        }

        if (1 != EVP_DecryptFinal_ex(ctx, plaintext + len1, &len2)) {
            delete[] plaintext;
            EVP_CIPHER_CTX_free(ctx);
            handleErrors("Decrypt final failed");
        }

        EVP_CIPHER_CTX_free(ctx);
        *plaintext_len = len1 + len2;
        return plaintext;

    }
    catch (...) {
        delete[] plaintext;
        EVP_CIPHER_CTX_free(ctx);
        throw;
    }
}

/**
 * 加密字符串（自动处理字符串长度）
 * @param plaintext 明文字符串
 * @param ciphertext_len 输出参数，返回密文长度
 * @return 密文数据指针，需要调用者释放 (delete[])
 */
unsigned char* encryptString(const char* plaintext, int* ciphertext_len) {
    return encrypt(reinterpret_cast<const unsigned char*>(plaintext),
        strlen(plaintext), ciphertext_len);
}

/**
 * 解密为字符串
 * @param ciphertext 密文数据
 * @param ciphertext_len 密文长度
 * @return 明文字符串指针，需要调用者释放 (delete[])
 */
unsigned char* decryptToString(const unsigned char* ciphertext, int* ciphertext_len) {
    int plaintext_len = 0;
    unsigned char* plaintext = decrypt(ciphertext, *ciphertext_len, &plaintext_len);
    *ciphertext_len = plaintext_len;
    return plaintext;
}

