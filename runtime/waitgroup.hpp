#pragma once
#include <condition_variable>
#include <mutex>

template <class T> class WaitGroup {
  std::mutex m;
  std::condition_variable cv;
  T counter;

public:
  WaitGroup() : counter(0) {};
  ~WaitGroup(){};

  void add(T n) {
    std::lock_guard<std::mutex> lk(m);
    counter += n;
  }

  void done() {
    bool notify = false;
    {
      std::lock_guard<std::mutex> lk(m);
      counter -= 1;
      assert(counter>=0);
      if (counter == 0)
        notify = true;
    }
    if (notify)
      cv.notify_one();
  }

  void wait() {
    std::unique_lock<std::mutex> lk(m);
    cv.wait(lk, [this] { return this->counter <= 0; });
  }
};
