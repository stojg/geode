#version 410 core

layout (location = 0) in vec2 aPosition;
layout (location = 1) in vec4 posTrans; // [3]pos, 1 float transp

out float transp;

#include "matrices.glsl"

uniform vec3 x_camUp;
uniform vec3 x_camRight;

void main()
{

    vec3 vertexPosition_worldspace =
    		posTrans.xyz
    		+ x_camRight * aPosition.x * 0.1
    		+ x_camUp * aPosition.y * 0.1;

    gl_Position = projection * view * vec4(vertexPosition_worldspace, 1.0);
    transp = posTrans.w;
}

