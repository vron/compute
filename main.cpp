#include <cstdio>

typedef float vec3 __attribute__((ext_vector_type(3)));

typedef float vec3a[4];

int main() {   
    printf("hej\n");

    vec3a b;
    b[0] = 1.0;
    b[1] = 2.0;
    b[2] = 3.0;


    vec3 a = *((vec3*)(b));

    printf("%f %f %f\n", a.x, a.y, a.z);

    return 0;
}