#include <math.h>
#include <stdio.h>

 typedef float f2 __attribute__((ext_vector_type(8)));

struct vec2 {
    f2 d;
    vec2() {vec2(0,0);}
    vec2(float a, float b) {
        d.x = a;
        d.y = b;
    }
    vec2& operator+=(const vec2& rhs) {
        this->d += rhs.d;
        return *this;
  }
};

int main() {
    vec2 a = vec2(1,1);
    vec2 b = vec2(2,2);
    a += b;
    printf("%f %f %d\n", a.d.x, a.d.y, sizeof(a));
}