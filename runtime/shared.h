// This file will be generated as part of the build, this is an example only

typedef struct {
  float *image;
} kernel_data;

void *cpt_new_kernel(int);
int cpt_dispatch_kernel(void *, kernel_data, int, int, int);
void cpt_free_kernel(void *);
