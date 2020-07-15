#pragma once

// #define DEBUG 1
#ifdef DEBUG

#include <stdarg.h>
#include <stdio.h>
#include <pthread.h>

static pthread_mutex_t deb_lock = PTHREAD_MUTEX_INITIALIZER;

#define cpt_log(f_, ...) { \
    pthread_mutex_lock(&deb_lock); \
    printf("%20s() %12s:%-4d | ", __func__, __FILE__, __LINE__); \
    printf((f_), ##__VA_ARGS__); \
    puts(""); \
    fflush(stdout); \
    pthread_mutex_unlock(&deb_lock); \
};

#define cpt_verbose(f_, ...) { \
    pthread_mutex_lock(&deb_lock); \
    printf("%20s() %12s:%-4d | ", __func__, __FILE__, __LINE__); \
    printf((f_), ##__VA_ARGS__); \
    puts(""); \
    fflush(stdout); \
    pthread_mutex_unlock(&deb_lock); \
};

#define cpt_trace(f_, ...) { \
    pthread_mutex_lock(&deb_lock); \
    printf("%20s() %12s:%-4d | ", __func__, __FILE__, __LINE__); \
    printf((f_), ##__VA_ARGS__); \
    puts(""); \
    fflush(stdout); \
    pthread_mutex_unlock(&deb_lock); \
};

#else

#define cpt_log(f_, ...) {};
#define cpt_verbose(f_, ...) {};
#define cpt_trace(f_, ...) {};

#endif
