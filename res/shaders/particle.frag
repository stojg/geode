#version 410 core

out vec4 FragColor;

uniform float transparency;

void main()
{
    FragColor = vec4(10,10,8,transparency);
}
