#version 410 core
out vec4 FragColor;

in vec3 TexCoords;

uniform samplerCube x_skybox;

void main()
{
    FragColor = texture(x_skybox, TexCoords);
}
