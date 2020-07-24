#pragma once

#include "generated/shared.h"
#include "types/types.hpp"

class WorkPiece {
public:
  uvec3 nwg;
  uvec3 gid;
  cpt_data *data;
  int generation;
  bool quit;
  WorkPiece(uvec3 nwg, uvec3 gid, cpt_data *data, int generation)
      : nwg(nwg), gid(gid), data(data), generation(generation), quit(false) {};
  WorkPiece(bool quit)
      : quit(quit) {};
};
