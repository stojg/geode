#version 410 core

#include "light.vert"

uniform PointLight pointLight;

void main() {
    setOutput(vec4(vec3(pointLight.position), 1.0));
}
