#version 410 core

#include "light.vert"

uniform SpotLight spotLight;

void main() {
    setOutput(vec4(vec3(spotLight.pointLight.position), 1.0));
}
