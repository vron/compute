#include "atemp.hpp"
#include "routines.hpp"
#include <cstdio>

#define WG_SIZE 5
#define NO_DI 3

void coroutine(Invocation<TempArg> *th) {
  auto td = th->get_argument();
  td.shader->thread = th;
  td.shader->main();
  th->exit();
}

int main() {
  auto wg = new WorkGroup<TempArg>(2, 1024 * 1024);

  Invocation<TempArg> *threads[WG_SIZE];
  shader shaders[WG_SIZE];
  for (int i = 0; i < WG_SIZE; i++) {
    threads[i] = wg->create_thread();
  }

  // for each wg dispatch:
  for (int n = 0; n < NO_DI; n++) {
    for (int i = 0; i < WG_SIZE; i++) {
      // do set invocation id's etc
      shaders[i].thread = threads[i];
      shaders[i].id = i + (n + 1) * 100;
      threads[i]->set_function(&coroutine,
                               TempArg(i + (n + 1) * 100, &shaders[i]));
    }

    wg->run(WG_SIZE, threads);
  }
  delete wg;
  return 0;
}

void shader::barrier() { this->thread->barrier(); }
