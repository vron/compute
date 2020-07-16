#pragma once
#include "routines.hpp"
#include <cstdio>

class TempArg;

struct shader {
  Invocation<TempArg> *thread;
  int id;

  shader(){};

  void main() {
    printf("started %d\n", this->id);
    barrier();
    printf("barriered %d\n", this->id);
    barrier();
    printf("done %d\n", this->id);
  }

  void barrier();
};

class TempArg {
public:
  int no;
  shader *shader;
  TempArg() : no(0){};
  TempArg(int no, struct shader *s) : no(no), shader(s){};
};