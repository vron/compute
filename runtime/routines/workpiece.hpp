#pragma once

#include "routines.hpp"

template <class T> class WorkPiece {
public:
  int no;
  bool sync_next;
  Invocation<T> **threads; // nullptr => end of computation
  WorkPiece() : no(0), sync_next(true) {};
  WorkPiece(int no, bool sync_next, Invocation<T> **threads) : no(no), sync_next(sync_next), threads(threads){};
};
