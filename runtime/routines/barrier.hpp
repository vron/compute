#pragma once
#include <condition_variable>
#include <mutex>

// a barrier which forces multiple threads to wait at that point until a pre-set
// number of threads has reached the barrier. Once activated it is re-set and can
// be used again.
// Should be replaced by std:barrier once that is implemented in clang.
class Barrier {
private:
  std::mutex m;
  std::condition_variable cv;
  int count;
  int gen;
  int original_count;

public:
  Barrier() : count(0), gen(0), original_count(0){};

  void set_count(int c) {
    std::unique_lock<std::mutex> lc(m);
    assert(this->count == this->original_count);
    count = c;
    original_count = c;
  }

  void wait() {
    std::unique_lock<std::mutex> lc(m);
    int tempgen = gen;
    count--;
    if (count == 0) {
      gen++;
      count = original_count;
      lc.unlock();
      cv.notify_all();
    } else {
      cv.wait(lc, [tempgen, this] { return tempgen != this->gen; });
    }
  };
};
