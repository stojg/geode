#version 410 core

layout (location = 0) in vec2 aPosition;
layout (location = 1) in vec4 posTrans; // [3]pos, 1 float transp
layout (location = 2) in vec4 scale; // 1 float scale

out float transp;

#include "matrices.glsl"

uniform vec3 x_camUp;
uniform vec3 x_camRight;

void main()
{

    float particleSize = scale.x;

    vec3 vertexPosition_worldspace =
    		posTrans.xyz
    		+ x_camRight * aPosition.x * particleSize
    		+ x_camUp * aPosition.y * particleSize;

    gl_Position = projection * view * vec4(vertexPosition_worldspace, 1.0);
    transp = posTrans.w;
}

