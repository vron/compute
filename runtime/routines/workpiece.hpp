#pragma once

#include "routines.hpp"

template <class T> class WorkPiece {
public:
  int no;
  Invocation<T> **threads; // nullptr => end of computation
  WorkPiece() : no(0){};
  WorkPiece(int no, Invocation<T> **threads) : no(no), threads(threads){};
};
