#version 410 core

out vec4 FragColor;

in float transp;

void main()
{
    FragColor = vec4(10,10,8,transp);
}
