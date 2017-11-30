#version 410 core

#include "include.frag"

uniform TestStructB light;

void main() {
    FragColor = calcLight(light.inner.color);
}
