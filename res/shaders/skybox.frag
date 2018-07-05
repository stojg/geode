#version 410 core
out vec4 FragColor;

in vec3 TexCoords;

uniform samplerCube x_skybox;

#include "fog.glsl"

const float lowerLimit = 0.0;
const float higherLimit = 0.2;

void main()
{
    vec4 finalColour = texture(x_skybox, TexCoords);
    float factor = (TexCoords.y - lowerLimit) / (higherLimit - lowerLimit);
    factor = clamp(factor, 0.0, 1.0);
    FragColor = mix(vec4(fogColor, 1), finalColour, factor);
}
