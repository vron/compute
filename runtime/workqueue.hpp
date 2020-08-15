#pragma once

#include <condition_variable>
#include <mutex>
#include <vector>

template <class T> class WorkQueue {
public:
  std::mutex m;
  std::condition_variable cv;
  std::vector<T> queue;

  T receive() {
    std::unique_lock<std::mutex> lk(m);

    cv.wait(lk, [this] { return this->queue.size() > 0; });

    auto e = queue.back();
    queue.pop_back();
    return e;
  };

  void send(T e) {
    {
      std::lock_guard<std::mutex> lk(m);
      queue.push_back(e);
    }
    cv.notify_one();
  };
};
