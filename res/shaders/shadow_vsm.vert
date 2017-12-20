#version 410 core
layout (location = 0) in vec3 aPos;

uniform mat4 LightMVP;

void main()
{
    gl_Position = LightMVP * vec4(aPos, 1.0);
}
