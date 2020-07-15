
#define _cpt_WG_SIZE_X 2
#define _cpt_WG_SIZE_Y 2
#define _cpt_WG_SIZE_Z 2

#include <math.h>

class shared_data_t {
public:
};

struct shader {
  uvec3 gl_NumWorkGroups;
  uvec3 gl_WorkGroupSize;
  uvec3 gl_WorkGroupID;
  uvec3 gl_LocalInvocationID;
  uvec3 gl_GlobalInvocationID;
  uint32_t gl_LocalInvocationIndex;
  thread_data *thread;

  float *din;

  float *dout;

  shader(){};

  void main() {
    
  }

  int set_data(cpt_data d) {
#include "setdata.hpp"
  }

  void barrier();

  static shared_data_t *create_shared_data() {
    shared_data_t *sd = new shared_data_t();

    return sd;
  }

  static void free_shared_data(shared_data_t *sd) { delete sd; }

  void set_shared_data(shared_data_t *sd) { (void)sd; }
};
