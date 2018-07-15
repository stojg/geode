#version 410 core

layout (location = 0) in vec2 aPosition;

uniform mat4 model;

void main()
{
    gl_Position = model * vec4(aPosition, 0, 1.0);
}

