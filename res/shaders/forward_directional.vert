#version 410 core

#include "light.vert"

uniform DirectionalLight directionalLight;

void main() {
    setOutput(vec4(vec3(directionalLight.direction), 0.0));
}
