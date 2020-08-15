#include "./generated/shared.h"
#include <cstdlib>

int main() {
  for (long long i = 0; i < 1e9; i++) {
    cpt_data d;
    d.data = malloc(64 * 8 * 8 * 8 * 1024 * 4);
    void *s = cpt_new_kernel(2, 1024);
    for (int j = 0; j < 1e9; j++) {
      cpt_dispatch_kernel(s, d, 8, 8, 8);
    }
    free(d.data);
    cpt_free_kernel(s);
  }
}