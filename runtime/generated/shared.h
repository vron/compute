#pragma once
/*
  This header and associated library was generated from a GLSL compute shader
  to be executed on a CPU as static code. The library is safe for threaded use
  as further specified below.
*/
#include "errno.h"
#include "stdint.h"

/*
  error_t represents an error as reported from cpt_dispatch_kernel. The 
  possible errors can mostly be classified as either user errors or underlying
  system errors. In case of underlying errors, such as insufficient resources,
  the .code field will be set to an error code from errno.h. In case of user
  errors, such as providing data with bad alignment, .code will be set to
  EINVAL with further description given in .msg. The data pointed to by .msg is
  only accessible until the next call to cpt_dispatch_kernel or cpt_free_kernel
  for the same kernel reference.
*/
struct error_t {
  int code;
  char* msg;
};

typedef struct {
  float vertices[6];
} cpt_triangle;

typedef struct {
  cpt_triangle triangles[64];
} cpt_polygon;

/*
  cpt_data consists of all the input/output required by the compute kernel. All
  fixed sized fields (including arrays) will be copied internally to ensure
  correct alignment. For all variable sizes fields (type void*) the user must
  ensure sufficient length and data alignment for the relevant use.
*/
typedef struct {
  float transform[9];
  void* polygons;
  void* cogs;
} cpt_data;

/*
  cpt_new_kernel creates a new computational kernel using a maximum of num_t
  threads for kernel calculation, returning a reference to the kernel created.
  If there is insufficient memory available to create a new kernel 0 is
  returned. For all other possible errors a kernel reference is returned and
  the next call to cpt_dispatch_kernel will return the error information.
  cpt_new_kernel is safe for concurrent use from multiple threads.
*/
void *cpt_new_kernel(int32_t num_t);

/*
  cpt_dispatch_kernel issues a calculation of the compute shader using x, y, z
  work groups in x, y, z directions respectively. The kernel reference k passed
  must have been created using cpt_new_kernel and not subsequently deallocated
  using cpt_free_kernel. It is the callers responsibility to ensure that any
  data of non-fixed size in d is properly aligned as required by the kernel and
  of sufficient length for the number of work groups issued. Any error message
  description returned in error_t.msg is only accessible until the next call to
  cpt_dispatch_kernel or cpt_free_kernel for the same kernel reference k.
  cpt_dispatch_kernel is safe for concurrent use by multiple threads for
  different kernel references (k) but must not be called concurrently for the
  same k.
*/
struct error_t cpt_dispatch_kernel(void *k, cpt_data d, int32_t x, int32_t y, int32_t z);

/*
  cpt_free_kernel must be called for any non-null kernel k created to avoid
  leaks. Note that any k for which cpt_free_kernel has been called is unsafe for
  any further use.
*/
void cpt_free_kernel(void *k);
