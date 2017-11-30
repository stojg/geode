#version 410 core

#include "include.frag"

uniform TestStructB light;
uniform TestStructA lights[2];

void main() {
    FragColor = calcLight(light.inner.color);

    for (int i = 0; i < 2; i++) {
        FragColor += lights[i].color;
    }
}
