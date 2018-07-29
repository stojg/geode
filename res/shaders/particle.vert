#version 410 core

layout (location = 0) in vec2 aPosition;
layout (location = 1) in vec4 posScale; // [3]pos, 1 float scale
layout (location = 2) in vec4 colourTrans; // [3]colour, 1 float transp

out vec4 colour;

#include "matrices.glsl"

uniform vec3 x_camUp;
uniform vec3 x_camRight;

void main()
{

    float particleSize = posScale.w;

    vec3 vertexPosition_worldspace =
    		posScale.xyz
    		+ x_camRight * aPosition.x * particleSize
    		+ x_camUp * aPosition.y * particleSize;

    gl_Position = projection * view * vec4(vertexPosition_worldspace, 1.0);

   colour = colourTrans;
}

