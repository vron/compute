// TODO: Wrap this into no-ops when no debugging.

#include <stdarg.h>
#include <stdio.h>
#include <assert.h>
#include <pthread.h>


#ifdef DEBUG

pthread_mutex_t deb_lock;
''
void init_debug() {
    assert(pthread_mutex_init(&deb_lock, NULL)==0);
}

#define log(f_, ...) { \
    pthread_mutex_lock(&deb_lock); \
    printf("%20s() %12s:%-4d | ", __func__, __FILE__, __LINE__); \
    printf((f_), ##__VA_ARGS__); \
    puts(""); \
    fflush(stdout); \
    pthread_mutex_unlock(&deb_lock); \
};

#define verbose(f_, ...) { \
    pthread_mutex_lock(&deb_lock); \
    printf("%20s() %12s:%-4d | ", __func__, __FILE__, __LINE__); \
    printf((f_), ##__VA_ARGS__); \
    puts(""); \
    fflush(stdout); \
    pthread_mutex_unlock(&deb_lock); \
};

#define trace(f_, ...) { \
    pthread_mutex_lock(&deb_lock); \
    printf("%20s() %12s:%-4d | ", __func__, __FILE__, __LINE__); \
    printf((f_), ##__VA_ARGS__); \
    puts(""); \
    fflush(stdout); \
    pthread_mutex_unlock(&deb_lock); \
};


#else
void init_debug() {}
#define log(f_, ...) {};
#define verbose(f_, ...) {};
#define trace(f_, ...) {};
#endif
