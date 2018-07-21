#version 410 core

layout (location = 0) in vec2 aPosition;
layout (location = 1) in mat4 model;
layout (location = 5) in float transparancy;
//uniform mat4 model;

out float transp;

void main()
{
    gl_Position = model * vec4(aPosition, 0, 1.0);
    transp = transparancy;
}

