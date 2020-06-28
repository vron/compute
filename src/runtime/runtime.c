#include "stdlib.h"
#include "stdio.h"
#include "kernel.h"

void panic(char* msg) {
    printf("panic: %s", msg);
    exit(1);
}


typedef struct {
    int num_threads;
} cpt_kernel;

typedef struct {
    unsigned int gid[3];
    unsigned int lid[3];
} cpt_ctx;

unsigned int __attribute__((overloadable)) __attribute__((const)) get_global_id_(const long _thread_context_, unsigned int a) {
    cpt_ctx* th;
    th = (cpt_ctx*)(_thread_context_);
    return th->gid[a];
}

unsigned int __attribute__((overloadable)) __attribute__((const)) get_local_id_(const long _thread_context_, unsigned int a) {
    cpt_ctx* th;
    th = (cpt_ctx*)(_thread_context_);
    return th->lid[a];
}

void* cpt_new_kernel(int num_threads) {
    // here we should create the threads and sync primitives we will use
    cpt_kernel* k;
    k = (cpt_kernel*)malloc(sizeof(cpt_kernel));
    k->num_threads = num_threads;
    if (num_threads == 1) {
        return k;
    }
    panic("multiple threads not implemented yet");
    return (void*)k;
}

void cpt_free_kernel(void* k) {
    // kill all the threads, deallocate buffers and sync primitives
    free((cpt_kernel*)k);
}


static int dispatch_seq(kernel_data* d,
    int numx, int numy, int numz) {
    cpt_ctx* th;
    // BUG(vron): this implementation will not work with barriers

    th = (cpt_ctx*)(malloc(sizeof(cpt_ctx)));
    for(int gz = 0; gz < numz; ++gz) {
        th->gid[2] = gz;
        for(int gy = 0; gy < numy; ++gy) {
            th->gid[1] = gy;
            for(int gx = 0; gx < numx; ++gx) {
                th->gid[0] = gx;

                // loop over the local work group
                #pragma clang loop unroll(full)
                for(int lz = 0; lz < LOCAL_SIZE_Z; ++lz) {
                    th->lid[2] = lz;
                    #pragma clang loop unroll(full)
                    for(int ly = 0; ly < LOCAL_SIZE_Y; ++ly) {
                        th->lid[1] = ly;
                        #pragma clang loop unroll(full)
                        for(int lx = 0; lx < LOCAL_SIZE_X; ++lx) {
                            th->lid[0] = lx;
                            // TODO(vrom): For small kernels, ensure this call is actually inlined
                            kern((long) th, d->imgData, d->imgWidth);
                        }
                    }
                }
            }
        }
    }
    free(th);
    return 0;
}

int cpt_dispatch_kernel(void* v,
    kernel_data d,
    int numx, int numy, int numz) {
    printf("v %p\n", v);
    cpt_kernel* k = (cpt_kernel*)(v);

    printf("in dispatch\n");
    if (k->num_threads == 1) {
    printf("in if\n");
        return dispatch_seq(&d, numx, numy, numz);
    }

    return 0;
}