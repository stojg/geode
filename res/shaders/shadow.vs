#version 410 core
layout (location = 0) in vec3 aPos;

uniform mat4 lightViewProjection;
uniform mat4 model;

void main()
{
    gl_Position = lightViewProjection * model * vec4(aPos, 1.0);
}
