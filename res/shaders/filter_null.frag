#version 410 core

out vec4 FragColor;

in vec2 TexCoords;

uniform sampler2D x_filterTexture;

void main()
{
    FragColor = vec4(texture(x_filterTexture, TexCoords).rgb, 1.0);
}
